package client

import (
	"context"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetShootClient returns a Kubernetes client for the shoot cluster associated with the given runtime.
//nolint:gochecknoglobals
var GetShootClient = func(ctx context.Context, cnt client.Client, runtime imv1.Runtime) (client.Client, error) {
	runtimeID := runtime.Labels[imv1.LabelKymaRuntimeID]

	secret, err := GetKubeconfigSecret(ctx, cnt, runtimeID, runtime.Namespace)
	if err != nil {
		return nil, err
	}

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(secret.Data[KubeconfigSecretKey])
	if err != nil {
		return nil, err
	}

	shootClientWithAdmin, err := client.New(restConfig, client.Options{})
	if err != nil {
		return nil, err
	}

	return shootClientWithAdmin, nil
}