package fsm

import (
	"context"
	"encoding/json"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnDeleteShoot(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	if !isGardenerCloudDelConfirmationSet(s.shoot.Annotations) {
		m.log.Info("patching shoot with del-confirmation")
		// workaround for Gardener client
		s.shoot.Kind = "Shoot"
		s.shoot.APIVersion = "core.gardener.cloud/v1beta1"
		s.shoot.Annotations = addGardenerCloudDelConfirmation(s.shoot.Annotations)
		s.shoot.ManagedFields = nil
		// attempt to marshall patched instance
		shootData, err := json.Marshal(&s.shoot)
		if err != nil {
			// unrecoverable error
			s.instance.UpdateStateDeletion(
				imv1.ConditionTypeRuntimeProvisioned,
				imv1.ConditionReasonSerializationError,
				"False",
				err.Error())
			return updateStatusAndStop()
		}
		// see: https://gardener.cloud/docs/gardener/projects/#four-eyes-principle-for-resource-deletion
		if s.shoot, err = m.ShootClient.Patch(ctx, s.shoot.Name, types.ApplyPatchType, shootData,
			metav1.PatchOptions{
				FieldManager: "kim",
				Force:        ptr.To(true),
			}); err != nil {
			m.log.Error(err, "unable to patch shoot:")
			return requeue()
		}
	}

	if !s.instance.IsStateWithConditionSet(imv1.RuntimeStateTerminating, imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonDeletion) {
		m.log.Info("setting state to in deletion")
		s.instance.UpdateStateDeletion(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonDeletion,
			"Unknown",
			"Runtime deletion initialised",
		)
		return updateStatusAndRequeue()
	}

	m.log.Info("deleting shoot")
	err := m.ShootClient.Delete(ctx, s.instance.Name, metav1.DeleteOptions{})
	if err != nil {
		m.log.Error(err, "Failed to delete gardener Shoot")

		s.instance.UpdateStateDeletion(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonGardenerError,
			"False",
			"Gardener API delete error",
		)
		return updateStatusAndRequeue()
	}

	return updateStatusAndRequeue()
}

func isGardenerCloudDelConfirmationSet(a map[string]string) bool {
	if len(a) == 0 {
		return false
	}
	val, found := a[imv1.AnnotationGardenerCloudDelConfirmation]
	return found && (val == "true")
}

func addGardenerCloudDelConfirmation(a map[string]string) map[string]string {
	if len(a) == 0 {
		a = map[string]string{}
	}
	a[imv1.AnnotationGardenerCloudDelConfirmation] = "true"
	return a
}
