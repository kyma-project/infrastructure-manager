package fsm

import (
	"context"
	"fmt"
	"io"
	"os"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/yaml"
)

func getWriterForFilesystem(filePath string) (io.Writer, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to create file: %w", err)
	}
	return file, nil
}

func persist(path string, s *gardener.Shoot, saveFunc writerGetter) error {
	writer, err := saveFunc(path)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}

	b, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("unable to marshal shoot: %w", err)
	}

	if _, err = writer.Write(b); err != nil {
		return fmt.Errorf("unable to write to file: %w", err)
	}
	return nil
}

func sFnPersistShoot(_ context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	path := fmt.Sprintf("%s/%s-%s.yaml", m.PVCPath, s.shoot.Namespace, s.shoot.Name)
	if err := persist(path, s.shoot, m.writerProvider); err != nil {
		return updateStatusAndStopWithError(err)
	}
	return updateStatusAndRequeueAfter(gardenerRequeueDuration)
}
