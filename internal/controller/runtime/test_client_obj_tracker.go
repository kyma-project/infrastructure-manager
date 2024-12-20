package runtime

import (
	"fmt"
	"sync"

	gardener_api "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	apimachinery "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clienttesting "k8s.io/client-go/testing"
)

const shootType = "shoots"

// CustomTracker implements ObjectTracker with a sequence of Shoot objects
// it will be updated with a different shoot sequence for each test case
type CustomTracker struct {
	clienttesting.ObjectTracker
	shootSequence []*gardener_api.Shoot
	shootCallCnt  int
	mu            sync.Mutex
}

func NewCustomTracker(tracker clienttesting.ObjectTracker, shoots []*gardener_api.Shoot) *CustomTracker {
	return &CustomTracker{
		ObjectTracker: tracker,
		shootSequence: shoots,
	}
}

func (t *CustomTracker) IsSequenceFullyUsed() bool {
	return (t.shootCallCnt == len(t.shootSequence) && len(t.shootSequence) > 0)
}

func (t *CustomTracker) Get(gvr schema.GroupVersionResource, ns, name string, opts ...apimachinery.GetOptions) (runtime.Object, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if gvr.Resource == shootType {
		return getNextObject(t.shootSequence, &t.shootCallCnt)
	}
	return t.ObjectTracker.Get(gvr, ns, name, opts...)
}

func getNextObject[T any](sequence []*T, counter *int) (*T, error) {
	if *counter < len(sequence) {
		obj := sequence[*counter]
		*counter++

		if obj == nil {
			return nil, k8serrors.NewNotFound(schema.GroupResource{}, "")
		}
		return obj, nil
	}
	return nil, fmt.Errorf("no more objects in sequence")
}

func (t *CustomTracker) Update(gvr schema.GroupVersionResource, obj runtime.Object, ns string, opts ...apimachinery.UpdateOptions) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if gvr.Resource == shootType {
		shoot, ok := obj.(*gardener_api.Shoot)
		if !ok {
			return fmt.Errorf("object is not of type Gardener Shoot")
		}
		for index, existingShoot := range t.shootSequence {
			if existingShoot != nil && existingShoot.Name == shoot.Name {
				t.shootSequence[index] = shoot
				return nil
			}
		}
		return k8serrors.NewNotFound(schema.GroupResource{}, shoot.Name)
	}
	return t.ObjectTracker.Update(gvr, obj, ns, opts...)
}

func (t *CustomTracker) Delete(gvr schema.GroupVersionResource, ns, name string, opts ...apimachinery.DeleteOptions) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if gvr.Resource == shootType {
		for index, shoot := range t.shootSequence {
			if shoot != nil && shoot.Name == name {
				t.shootSequence[index] = nil
				return nil
			}
		}
		return k8serrors.NewNotFound(schema.GroupResource{}, "")
	}
	return t.ObjectTracker.Delete(gvr, ns, name, opts...)
}
