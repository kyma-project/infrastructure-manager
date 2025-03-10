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

package runtime

import (
	"context"
	"fmt"
	"github.com/kyma-project/infrastructure-manager/internal/controller/customconfig/registrycache"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sync/atomic"
	"time"

	"github.com/go-logr/logr"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/infrastructure-manager/internal/controller/runtime/fsm"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// RuntimeReconciler reconciles a Runtime object
// nolint:revive
type CustomSKRConfigReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Log           logr.Logger
	Cfg           fsm.RCCfg
	EventRecorder record.EventRecorder
	RequestID     atomic.Uint64
}

const fieldManagerName = "customconfigcontroller"

//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=runtimes,verbs=get;list;watch;create;update;patch,namespace=kcp-system
//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=runtimes/status,verbs=get;list;delete;create;update;patch,namespace=kcp-system
//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=runtimes/finalizers,verbs=get;list;delete;create;update;patch,namespace=kcp-system

func (r *CustomSKRConfigReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	r.Log.Info(request.String())

	var secret v1.Secret
	if err := r.Get(ctx, request.NamespacedName, &secret); err != nil {

		if apierrors.IsNotFound(err) {
			r.Log.Info(fmt.Sprintf("Secret %s not found", request.Name))
		} else {
			r.Log.Error(err, fmt.Sprintf("Failed to get secret %s", request.Name))
		}

		return ctrl.Result{
			Requeue: false,
		}, client.IgnoreNotFound(err)
	}

	runtimeID, ok := secret.Labels["kyma-project.io/runtime-id"]
	secretControlledByKIM := ok && secret.Labels["operator.kyma-project.io/managed-by"] == "infrastructure-manager"

	if !secretControlledByKIM {
		r.Log.Info(fmt.Sprintf("Secret doesn't contain kubeconfig %s", request.Name))

		return ctrl.Result{
			Requeue: false,
		}, nil
	}

	var runtime imv1.Runtime
	if err := r.Get(ctx, types.NamespacedName{
		Name:      runtimeID,
		Namespace: request.Namespace,
	}, &runtime); err != nil {
		r.Log.Error(err, fmt.Sprintf("Failed to get runtime %s", runtimeID))

		if apierrors.IsNotFound(err) {
			r.Log.Info(fmt.Sprintf("Runtime %s not found", request.Name))
		}

		return ctrl.Result{
			Requeue: false,
		}, client.IgnoreNotFound(err)
	}

	log := r.Log.WithValues("runtimeID", runtimeID, "shootName", runtime.Spec.Shoot.Name, "requestID", r.RequestID.Add(1))
	log.Info("Reconciling custom configuration", "Name", runtime.Name, "Namespace", runtime.Namespace)

	return r.handleCustomConfig(ctx, runtime, secret)
}

func (r *CustomSKRConfigReconciler) handleCustomConfig(ctx context.Context, runtime imv1.Runtime, kubeconfigSecret v1.Secret) (ctrl.Result, error) {
	getSecretFunc := func() (v1.Secret, error) {
		return kubeconfigSecret, nil
	}
	customConfigExplorer, err := registrycache.NewConfigExplorer(ctx, getSecretFunc)
	if err != nil {
		r.Log.Error(err, "Failed to create custom config explorer")
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: 30 * time.Minute,
		}, err
	}

	exists, err := customConfigExplorer.RegistryCacheConfigExists()
	if err != nil {
		r.Log.Error(err, "Failed to verify custom config explorer")
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Minute,
		}, err
	}

	if exists {
		r.Log.Info(fmt.Sprintf("Custom config exists on runtime %s", runtime.Name))
	} else {
		r.Log.Info(fmt.Sprintf("Custom config doesn't exist on runtime %s", runtime.Name))
	}

	runtime.ManagedFields = nil

	cachingEnabled := runtime.Spec.Caching != nil && runtime.Spec.Caching.Enabled

	if cachingEnabled != exists {
		runtime.Spec.Caching = &imv1.ImageRegistryCache{
			Enabled: exists,
		}

		r.Log.Info(fmt.Sprintf("Updating runtime %s with caching enabled: %t", runtime.Name, exists))

		err := r.Client.Patch(ctx, &runtime, client.Apply, &client.PatchOptions{
			FieldManager: fieldManagerName,
			Force:        ptr.To(true),
		})

		if err != nil {
			r.Log.Error(err, "Failed to patch runtime")
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: time.Minute,
			}, err
		}
	}

	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: 1 * time.Minute,
	}, err
}

func NewCustomSKRConfigReconciler(mgr ctrl.Manager, logger logr.Logger) *CustomSKRConfigReconciler {
	return &CustomSKRConfigReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		EventRecorder: mgr.GetEventRecorderFor("runtime-controller"),
		Log:           logger,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *CustomSKRConfigReconciler) SetupWithManager(mgr ctrl.Manager, channelSource source.Source, numberOfWorkers int) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Secret{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: numberOfWorkers}).
		WithEventFilter(predicate.Or(
			predicate.GenerationChangedPredicate{},
			predicate.LabelChangedPredicate{},
			predicate.AnnotationChangedPredicate{},
		)).
		Named("custom-config-controller").
		WatchesRawSource(channelSource).
		Complete(r)
}
