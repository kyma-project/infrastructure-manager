package fsm

import (
	"context"
	"fmt"
	gardenerhelper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"
	"strings"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnWaitForShootCreation(_ context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {

	switch s.shoot.Status.LastOperation.State {

	case gardener.LastOperationStateProcessing, gardener.LastOperationStatePending:
		msg := fmt.Sprintf("Shoot %s is in %s state, scheduling for retry", s.shoot.Name, s.shoot.Status.LastOperation.State)
		m.log.Info(msg)

		s.instance.UpdateStatePending(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonShootCreationPending,
			"Unknown",
			"Shoot creation in progress")

		return updateStatusAndRequeueAfter(gardenerRequeueDuration)

	case gardener.LastOperationStateSucceeded:
		msg := fmt.Sprintf("Shoot %s successfully created", s.shoot.Name)
		m.log.Info(msg)

		if !s.instance.IsStateWithConditionSet(imv1.RuntimeStatePending, imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonShootCreationCompleted) {
			s.instance.UpdateStatePending(
				imv1.ConditionTypeRuntimeProvisioned,
				imv1.ConditionReasonShootCreationCompleted,
				"True",
				"Shoot creation completed")

			return updateStatusAndRequeue()
		}

		return switchState(sFnProcessShoot)

	case gardener.LastOperationStateFailed:
		if gardenerhelper.HasErrorCode(s.shoot.Status.LastErrors, gardener.ErrorInfraRateLimitsExceeded) {
			msg := fmt.Sprintf("Error during cluster provisioning: Rate limits exceeded for Shoot %s, scheduling for retry", s.shoot.Name)
			m.log.Info(msg)
			return updateStatusAndRequeueAfter(gardenerRequeueDuration)
		}

		// also handle other retryable errors here
		// ErrorRetryableConfigurationProblem
		// ErrorRetryableInfraDependencies

		msg := fmt.Sprintf("Provisioning failed for shoot: %s ! Last state: %s, Description: %s", s.shoot.Name, s.shoot.Status.LastOperation.State, s.shoot.Status.LastOperation.Description)
		m.log.Info(msg)

		s.instance.UpdateStatePending(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonCreationError,
			"False",
			"Shoot creation failed")

		return updateStatusAndStop()

	default:
		m.log.Info("Unknown shoot operation state, exiting with no retry")
		return stop()
	}
}

func gardenerErrCodesToErrReason(lastErrors ...gardener.LastError) ErrReason {
	var codes []gardener.ErrorCode
	var vals []string

	for _, e := range lastErrors {
		if len(e.Codes) > 0 {
			codes = append(codes, e.Codes...)
		}
	}

	for _, code := range codes {
		vals = append(vals, string(code))
	}

	return ErrReason(strings.Join(vals, ", "))
}
