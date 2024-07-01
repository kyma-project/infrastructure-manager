package fsm

import (
	"context"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// There is a decision made to not rely on state of the Runtime CR we have already set
// All the states we set in the operator are about to be read only by the external clients

func sFnInitialize(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	instanceIsNotBeingDeleted := s.instance.GetDeletionTimestamp().IsZero()
	instanceHasFinalizer := controllerutil.ContainsFinalizer(&s.instance, m.Finalizer)
	provisioningCondition := meta.FindStatusCondition(s.instance.Status.Conditions, string(imv1.ConditionTypeRuntimeProvisioned))

	if instanceIsNotBeingDeleted && !instanceHasFinalizer {
		return addFinalizerAndRequeue(ctx, m, s)
	}

	if instanceIsNotBeingDeleted && s.shoot == nil && provisioningCondition == nil {
		m.log.Info("Update Runtime state to Pending - initialised")
		s.instance.UpdateStatePending(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonInitialized,
			"Unknown",
			"Runtime initialized",
		)
		return updateStatusAndRequeue()
	}

	if instanceIsNotBeingDeleted {
		m.log.Info("Gardener shoot exists, processing")
		return switchState(sFnSelectClusterProcessing)
	}

	if !instanceIsNotBeingDeleted && instanceHasFinalizer && s.shoot != nil {
		m.log.Info("Instance is being deleted")
		return switchState(sFnDeleteShoot)
	}

	if !instanceIsNotBeingDeleted && instanceHasFinalizer && s.shoot == nil {
		return removeFinalizerAndStop(ctx, m, s)
	}

	return stop()
}

func addFinalizerAndRequeue(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("adding finalizer")
	controllerutil.AddFinalizer(&s.instance, m.Finalizer)

	err := m.Update(ctx, &s.instance)
	if err != nil {
		return updateStatusAndStopWithError(err)
	}
	return requeue()
}

func removeFinalizerAndStop(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("removing finalizer")
	controllerutil.RemoveFinalizer(&s.instance, m.Finalizer)

	err := m.Update(ctx, &s.instance)
	if err != nil {
		return updateStatusAndStopWithError(err)
	}
	return stop()
}
