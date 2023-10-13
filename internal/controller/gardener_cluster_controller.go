/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	lastKubeconfigSyncAnnotation = "operator.kyma-project.io/last-sync"
	clusterCRNameLabel           = "operator.kyma-project.io/cluster-name"
	defaultRequeuInSeconds       = 60 * 1000
)

// GardenerClusterController reconciles a GardenerCluster object
type GardenerClusterController struct {
	client.Client
	Scheme             *runtime.Scheme
	KubeconfigProvider KubeconfigProvider
	log                logr.Logger
}

func NewGardenerClusterController(mgr ctrl.Manager, kubeconfigProvider KubeconfigProvider, logger logr.Logger) *GardenerClusterController {
	return &GardenerClusterController{
		Client:             mgr.GetClient(),
		Scheme:             mgr.GetScheme(),
		KubeconfigProvider: kubeconfigProvider,
		log:                logger,
	}
}

//go:generate mockery --name=KubeconfigProvider
type KubeconfigProvider interface {
	Fetch(shootName string) (string, error)
}

//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=gardenerclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=gardenerclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=gardenerclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the GardenerCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (controller *GardenerClusterController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) { //nolint:revive
	controller.log.Info("Starting reconciliation loop")

	var cluster imv1.GardenerCluster

	err := controller.Client.Get(ctx, req.NamespacedName, &cluster)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			err = controller.deleteSecret(req.NamespacedName.Name)
			if err != nil {
				controller.log.Error(err, "failed to delete secret")
			}
		}

		return controller.resultWithoutRequeue(), err
	}

	cluster.UpdateConditionForProcessingState(imv1.ConditionTypeKubeconfigManagement, imv1.ConditionReasonKubeconfigReadingSecret, metav1.ConditionTrue)

	secret, err := controller.getSecret(cluster.Spec.Shoot.Name)
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			controller.log.Error(err, "failed to get the Secret for "+cluster.Spec.Shoot.Name)
			cluster.UpdateConditionForErrorState(imv1.ConditionTypeKubeconfigManagement, imv1.ConditionReasonFailedToCheckSecret, metav1.ConditionTrue, err)
			if err := controller.persistStatusChange(ctx, &cluster); err != nil {
				controller.log.Error(err, "Failed to set state for GardenerCluster", req.NamespacedName)
			}
			return controller.resultWithoutRequeue(), err
		}
	}

	if secret == nil {
		err = controller.createSecret(ctx, cluster)

		if err != nil {
			cluster.UpdateConditionForErrorState(imv1.ConditionTypeKubeconfigManagement, imv1.ConditionReasonFailedToCreateSecret, metav1.ConditionTrue, err)
			if err := controller.persistStatusChange(ctx, &cluster); err != nil {
				controller.log.Error(err, "Failed to set state for GardenerCluster", req.NamespacedName)
			}
			return controller.resultWithoutRequeue(), err
		}
	}

	cluster.UpdateConditionForReadyState(imv1.ConditionTypeKubeconfigManagement, imv1.ConditionReasonKubeconfigSecretReady, metav1.ConditionTrue)

	return controller.resultWithRequeue(), controller.persistStatusChange(ctx, &cluster)
}

func (controller *GardenerClusterController) resultWithRequeue() ctrl.Result {
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: defaultRequeuInSeconds,
	}
}

func (controller *GardenerClusterController) resultWithoutRequeue() ctrl.Result {
	return ctrl.Result{
		Requeue: false,
	}
}

func (controller *GardenerClusterController) persistStatusChange(ctx context.Context, cluster *imv1.GardenerCluster) error {
	return controller.Client.Status().Update(ctx, cluster)
}

func (controller *GardenerClusterController) deleteSecret(clusterCRName string) error {
	selector := client.MatchingLabels(map[string]string{
		clusterCRNameLabel: clusterCRName,
	})

	var secretList corev1.SecretList
	err := controller.Client.List(context.TODO(), &secretList, selector)
	if err != nil {
		return err
	}

	if len(secretList.Items) != 1 {
		return errors.Errorf("unexpected numer of secrets found for cluster CR `%s`", clusterCRName)
	}

	return controller.Client.Delete(context.TODO(), &secretList.Items[0])
}

func (controller *GardenerClusterController) getSecret(shootName string) (*corev1.Secret, error) {
	var secretList corev1.SecretList

	shootNameSelector := client.MatchingLabels(map[string]string{
		"kyma-project.io/shoot-name": shootName,
	})

	err := controller.Client.List(context.Background(), &secretList, shootNameSelector)
	if err != nil {
		return nil, err
	}

	size := len(secretList.Items)

	if size == 0 {
		return nil, k8serrors.NewNotFound(schema.GroupResource{}, "") //nolint:all TODO: update with valid GroupResource
	}

	if size > 1 {
		return nil, fmt.Errorf("unexpected numer of secrets found for shoot `%s`", shootName)
	}

	return &secretList.Items[0], nil
}

func (controller *GardenerClusterController) createSecret(ctx context.Context, cluster imv1.GardenerCluster) error {
	secret, err := controller.newSecret(cluster)
	if err != nil {
		return err
	}

	return controller.Client.Create(ctx, &secret)
}

func (controller *GardenerClusterController) newSecret(cluster imv1.GardenerCluster) (corev1.Secret, error) {
	labels := map[string]string{}

	for key, val := range cluster.Labels {
		labels[key] = val
	}
	labels["operator.kyma-project.io/managed-by"] = "infrastructure-manager"
	labels[clusterCRNameLabel] = cluster.Name

	kubeconfig, err := controller.KubeconfigProvider.Fetch(cluster.Spec.Shoot.Name)
	if err != nil {
		return corev1.Secret{}, err
	}

	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        cluster.Spec.Kubeconfig.Secret.Name,
			Namespace:   cluster.Spec.Kubeconfig.Secret.Namespace,
			Labels:      labels,
			Annotations: map[string]string{lastKubeconfigSyncAnnotation: time.Now().UTC().String()},
		},
		StringData: map[string]string{cluster.Spec.Kubeconfig.Secret.Key: kubeconfig},
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (controller *GardenerClusterController) SetupWithManager(mgr ctrl.Manager) error {
	omitStatusChanged := predicate.Or(
		predicate.GenerationChangedPredicate{},
		predicate.LabelChangedPredicate{},
		predicate.AnnotationChangedPredicate{},
	)

	return ctrl.NewControllerManagedBy(mgr).
		For(&imv1.GardenerCluster{}, builder.WithPredicates(
			predicate.And(
				predicate.ResourceVersionChangedPredicate{},
				omitStatusChanged),
		)).
		Complete(controller)
}
