package controller

import (
	"context"
	"time"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Gardener Cluster controller", func() {
	Context("Secret with kubeconfig doesn't exist", func() {
		It("Should create secret, and set Ready status on CR", func() {
			Skip("Skipping for now, do not merge with this!")
			kymaName := "kymaname1"
			secretName := "secret-name1"
			shootName := "shootName1"
			namespace := "default"

			By("Create GardenerCluster CR")

			gardenerClusterCR := fixGardenerClusterCR(kymaName, namespace, shootName, secretName)
			Expect(k8sClient.Create(context.Background(), &gardenerClusterCR)).To(Succeed())

			By("Wait for secret creation")
			var kubeconfigSecret corev1.Secret
			secretKey := types.NamespacedName{Name: secretName, Namespace: namespace}

			Eventually(func() bool {
				return k8sClient.Get(context.Background(), secretKey, &kubeconfigSecret) == nil
			}, time.Second*30, time.Second*3).Should(BeTrue())

			gardenerClusterKey := types.NamespacedName{Name: gardenerClusterCR.Name, Namespace: gardenerClusterCR.Namespace}
			var newGardenerCluster imv1.GardenerCluster
			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), gardenerClusterKey, &newGardenerCluster)
				if err != nil {
					return false
				}

				return newGardenerCluster.Status.State == imv1.ReadyState
			}, time.Second*30, time.Second*3).Should(BeTrue())

			err := k8sClient.Get(context.Background(), secretKey, &kubeconfigSecret)
			Expect(err).To(BeNil())
			expectedSecret := fixNewSecret(secretName, namespace, kymaName, shootName, "kubeconfig1", "")
			Expect(kubeconfigSecret.Labels).To(Equal(expectedSecret.Labels))
			Expect(kubeconfigSecret.Data).To(Equal(expectedSecret.Data))
			lastSyncTime := kubeconfigSecret.Annotations[lastKubeconfigSyncAnnotation]
			Expect(lastSyncTime).ToNot(BeEmpty())

		})

		It("Should delete secret", func() {
			Skip("Skipping for now, do not merge with this!")
			kymaName := "kymaname2"
			secretName := "secret-name2"
			shootName := "shootName2"
			namespace := "default"

			By("Create GardenerCluster CR")

			gardenerClusterCR := fixGardenerClusterCR(kymaName, namespace, shootName, secretName)
			Expect(k8sClient.Create(context.Background(), &gardenerClusterCR)).To(Succeed())

			By("Wait for secret creation")
			var kubeconfigSecret corev1.Secret
			secretKey := types.NamespacedName{Name: secretName, Namespace: namespace}

			Eventually(func() bool {
				return k8sClient.Get(context.Background(), secretKey, &kubeconfigSecret) == nil
			}, time.Second*30, time.Second*3).Should(BeTrue())

			By("Delete Cluster CR")
			Expect(k8sClient.Delete(context.Background(), &gardenerClusterCR)).To(Succeed())

			By("Wait for secret deletion")
			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), secretKey, &kubeconfigSecret)
				return err != nil && k8serrors.IsNotFound(err)
			}, time.Second*30, time.Second*3).Should(BeTrue())
		})

		It("Should set Error status on CR if failed to fetch kubeconfig", func() {
			Skip("Skipping for now, do not merge with this!")
			kymaName := "kymaname3"
			secretName := "secret-name3"
			shootName := "shootName3"
			namespace := "default"

			gardenerClusterCR := fixGardenerClusterCR(kymaName, namespace, shootName, secretName)
			Expect(k8sClient.Create(context.Background(), &gardenerClusterCR)).To(Succeed())

			gardenerClusterKey := types.NamespacedName{Name: gardenerClusterCR.Name, Namespace: gardenerClusterCR.Namespace}
			var newGardenerCluster imv1.GardenerCluster
			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), gardenerClusterKey, &newGardenerCluster)
				if err != nil {
					return false
				}

				return newGardenerCluster.Status.State == imv1.ErrorState
			}, time.Second*30, time.Second*3).Should(BeTrue())
		})
	})

	Context("Secret with kubeconfig exists", func() {
		namespace := "default"
		DescribeTable("Should update secret", func(gardenerClusterCR imv1.GardenerCluster, secret corev1.Secret, expectedKubeconfig string) {
			Skip("Skipping for now, do not merge with this!")
			By("Create kubeconfig secret")
			Expect(k8sClient.Create(context.Background(), &secret)).To(Succeed())

			previousTimestamp := secret.Annotations[lastKubeconfigSyncAnnotation]

			By("Create Cluster CR")
			Expect(k8sClient.Create(context.Background(), &gardenerClusterCR)).To(Succeed())

			var kubeconfigSecret corev1.Secret
			secretKey := types.NamespacedName{Name: secret.Name, Namespace: namespace}

			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), secretKey, &kubeconfigSecret)
				if err != nil {
					return false
				}

				timestampAnnotation := kubeconfigSecret.Annotations[lastKubeconfigSyncAnnotation]

				return timestampAnnotation != previousTimestamp
			}, time.Second*30, time.Second*3).Should(BeTrue())

			gardenerClusterKey := types.NamespacedName{Name: gardenerClusterCR.Name, Namespace: gardenerClusterCR.Namespace}
			var newGardenerCluster imv1.GardenerCluster

			Eventually(func() bool {
				err := k8sClient.Get(context.Background(), gardenerClusterKey, &newGardenerCluster)
				if err != nil {
					return false
				}

				readyState := newGardenerCluster.Status.State == imv1.ReadyState
				_, forceRotationAnnotationFound := newGardenerCluster.GetAnnotations()[forceKubeconfigRotationAnnotation]

				return readyState && !forceRotationAnnotationFound
			}, time.Second*45, time.Second*3).Should(BeTrue())

			err := k8sClient.Get(context.Background(), secretKey, &kubeconfigSecret)
			Expect(err).To(BeNil())
			Expect(string(kubeconfigSecret.Data["config"])).To(Equal(expectedKubeconfig))
			lastSyncTime := kubeconfigSecret.Annotations[lastKubeconfigSyncAnnotation]
			Expect(lastSyncTime).ToNot(BeEmpty())

		},
			Entry("Rotate kubeconfig when rotation time passed",
				fixGardenerClusterCR("kymaname4", namespace, "shootName4", "secret-name4"),
				fixNewSecret("secret-name4", namespace, "kymaname4", "shootName4", "kubeconfig4", "2023-10-09T23:00:00Z"),
				"kubeconfig4"),
			Entry("Force rotation",
				fixGardenerClusterCRWithForceRotationAnnotation("kymaname5", namespace, "shootName5", "secret-name5"),
				fixNewSecret("secret-name5", namespace, "kymaname5", "shootName5", "kubeconfig5", time.Now().UTC().Format(time.RFC3339)),
				"kubeconfig5"),
		)

		It("Should skip rotation", func() {
			By("Create kubeconfig secret")
			secret := fixNewSecret("secret-name6", namespace, "kymaname6", "shootName6", "kubeconfig6", time.Now().UTC().Format(time.RFC3339))
			Expect(k8sClient.Create(context.Background(), &secret)).To(Succeed())

			previousTimestamp := secret.Annotations[lastKubeconfigSyncAnnotation]

			By("Create Cluster CR")
			gardenerClusterCR := fixGardenerClusterCR("kymaname6", namespace, "shootName6", "secret-name6")
			Expect(k8sClient.Create(context.Background(), &gardenerClusterCR)).To(Succeed())

			var kubeconfigSecret corev1.Secret
			secretKey := types.NamespacedName{Name: secret.Name, Namespace: namespace}

			Consistently(func() bool {
				err := k8sClient.Get(context.Background(), secretKey, &kubeconfigSecret)
				if err != nil {
					return false
				}

				timestampAnnotation := kubeconfigSecret.Annotations[lastKubeconfigSyncAnnotation]

				return timestampAnnotation == previousTimestamp
			}, time.Second*45, time.Second*3).Should(BeTrue())
		})
	})
})

func fixNewSecret(name, namespace, kymaName, shootName, data string, lastSyncTime string) corev1.Secret {
	labels := fixSecretLabels(kymaName, shootName)
	annotations := map[string]string{lastKubeconfigSyncAnnotation: lastSyncTime}

	builder := newTestSecret(name, namespace)
	return builder.WithLabels(labels).WithAnnotations(annotations).WithData(data).ToSecret()
}

func (sb *TestSecret) WithAnnotations(annotations map[string]string) *TestSecret {
	sb.secret.Annotations = annotations

	return sb
}

func (sb *TestSecret) WithLabels(labels map[string]string) *TestSecret {
	sb.secret.Labels = labels

	return sb
}

func (sb *TestSecret) WithData(data string) *TestSecret {
	sb.secret.Data = map[string][]byte{"config": []byte(data)}

	return sb
}

func (sb *TestSecret) ToSecret() corev1.Secret {
	return sb.secret
}

func newTestSecret(name, namespace string) *TestSecret {
	return &TestSecret{
		secret: corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}
}

type TestSecret struct {
	secret corev1.Secret
}

func fixSecretLabels(kymaName, shootName string) map[string]string {
	labels := fixGardenerClusterLabels(kymaName, shootName)
	labels["operator.kyma-project.io/managed-by"] = "infrastructure-manager"
	labels["operator.kyma-project.io/cluster-name"] = kymaName
	return labels
}

func fixGardenerClusterCR(kymaName, namespace, shootName, secretName string) imv1.GardenerCluster {
	return newTestGardenerClusterCR(kymaName, namespace, shootName, secretName).
		WithLabels(fixGardenerClusterLabels(kymaName, shootName)).ToCluster()
}

func fixGardenerClusterCRWithForceRotationAnnotation(kymaName, namespace, shootName, secretName string) imv1.GardenerCluster {
	annotations := map[string]string{forceKubeconfigRotationAnnotation: "true"}

	return newTestGardenerClusterCR(kymaName, namespace, shootName, secretName).
		WithLabels(fixGardenerClusterLabels(kymaName, shootName)).
		WithAnnotations(annotations).
		ToCluster()
}

func newTestGardenerClusterCR(name, namespace, shootName, secretName string) *TestGardenerClusterCR {
	return &TestGardenerClusterCR{
		gardenerCluster: imv1.GardenerCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: imv1.GardenerClusterSpec{
				Shoot: imv1.Shoot{
					Name: shootName,
				},
				Kubeconfig: imv1.Kubeconfig{
					Secret: imv1.Secret{
						Name:      secretName,
						Namespace: namespace,
						Key:       "config", //nolint:all TODO: fill it up with the actual data
					},
				},
			},
		},
	}
}

func (sb *TestGardenerClusterCR) WithLabels(labels map[string]string) *TestGardenerClusterCR {
	sb.gardenerCluster.Labels = labels

	return sb
}

func (sb *TestGardenerClusterCR) WithAnnotations(annotations map[string]string) *TestGardenerClusterCR {
	sb.gardenerCluster.Annotations = annotations

	return sb
}

func (sb *TestGardenerClusterCR) ToCluster() imv1.GardenerCluster {
	return sb.gardenerCluster
}

type TestGardenerClusterCR struct {
	gardenerCluster imv1.GardenerCluster
}

func fixGardenerClusterLabels(kymaName, shootName string) map[string]string {
	labels := map[string]string{}

	labels["kyma-project.io/instance-id"] = "instanceID"
	labels["kyma-project.io/runtime-id"] = "runtimeID"
	labels["kyma-project.io/broker-plan-id"] = "planID"
	labels["kyma-project.io/broker-plan-name"] = "planName"
	labels["kyma-project.io/global-account-id"] = "globalAccountID"
	labels["kyma-project.io/subaccount-id"] = "subAccountID"
	labels["kyma-project.io/shoot-name"] = shootName
	labels["kyma-project.io/region"] = "region"
	labels["operator.kyma-project.io/kyma-name"] = kymaName

	return labels
}
