package fsm

import (
	"context"
	"fmt"
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ErrReason string

func sFnSelectClusterProcessing(_ context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	if s.shoot == nil {
		m.log.Info("Gardener shoot does not exist, creating new one")
		return switchState(sFnCreateShoot)
	}

	if s.shoot.Spec.DNS == nil || s.shoot.Spec.DNS.Domain == nil {
		msg := fmt.Sprintf("DNS Domain is not set yet for shoot: %s, scheduling for retry", s.shoot.Name)
		m.log.Info(msg)
		return updateStatusAndRequeueAfter(gardenerRequeueDuration)
	}

	lastOperation := s.shoot.Status.LastOperation
	if lastOperation == nil {
		msg := fmt.Sprintf("Last operation is nil for shoot: %s, scheduling for retry", s.shoot.Name)
		m.log.Info(msg)
		return updateStatusAndRequeueAfter(gardenerRequeueDuration)
	}

	if lastOperation.Type == gardener.LastOperationTypeCreate {
		return switchState(sFnWaitForShootCreation)
	}

	if lastOperation.Type == gardener.LastOperationTypeReconcile {
		return switchState(sFnWaitForShootReconcile)
	}

	m.log.Info("Unknown shoot operation type, exiting with no retry")
	return stop()
}
