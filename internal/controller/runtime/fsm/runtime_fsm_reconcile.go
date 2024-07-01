package fsm

import (
	"context"
	"fmt"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnWaitForShootReconcile(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {

	lastOperation := s.shoot.Status.LastOperation

	switch s.shoot.Status.LastOperation.State {
	case gardener.LastOperationStateProcessing, gardener.LastOperationStatePending:
		msg := fmt.Sprintf("Shoot %s is in %s state, scheduling for retry", s.shoot.Name, lastOperation.State)
		m.log.Info(msg)
		s.instance.UpdateStatePending(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonProcessing,
			"Unknown",
			"Shoot update in progress")

		return updateStatusAndRequeue()

	case gardener.LastOperationStateSucceeded:
		msg := fmt.Sprintf("Shoot %s successfully updated", s.shoot.Name)
		m.log.Info(msg)

		// shoot completed

		if !s.instance.IsStateWithConditionSet(imv1.RuntimeStatePending, imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonProcessing) {
			s.instance.UpdateStatePending(
				imv1.ConditionTypeRuntimeProvisioned,
				imv1.ConditionReasonProcessing,
				"True",
				"Shoot update completed")

			return updateStatusAndRequeue()
		}

		return switchState(sFnProcessShoot)

	case gardener.LastOperationStateFailed:
		var reason ErrReason

		if len(s.shoot.Status.LastErrors) > 0 {
			reason = gardenerErrCodesToErrReason(s.shoot.Status.LastErrors...)
		}

		msg := fmt.Sprintf("error during cluster processing: reconcilation error for shoot %s, reason: %s, scheduling for retry", s.shoot.Name, reason)
		m.log.Info(msg)

		s.instance.UpdateStatePending(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonProcessingErr,
			"False",
			string(reason))

		return updateStatusAndStop()
	}

	m.log.Info("Update did not processed, exiting with no retry")
	return stop()
}
