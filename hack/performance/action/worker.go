package action

import (
	"context"
	"math/rand"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Worker interface {
	Create() error
	Delete() error
}

type WorkerData struct {
	loadID     string
	namePrefix string
	rtNumber   int
	k8sClient  client.Client
	rtTemplate imv1.Runtime
}

func NewWorker(loadID, namePrefix, kubeconfigPath string, rtNumber int, rtTemplate imv1.Runtime) (Worker, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	k8sClient, err := client.New(config, client.Options{})
	if err != nil {
		return nil, err
	}

	err = imv1.AddToScheme(k8sClient.Scheme())
	if err != nil {
		return nil, err
	}

	return &WorkerData{
		loadID:     loadID,
		namePrefix: namePrefix,
		rtNumber:   rtNumber,
		k8sClient:  k8sClient,
		rtTemplate: rtTemplate,
	}, nil
}

func (w WorkerData) Create() error {
	runtimes := w.prepareRuntimeBatch()

	for i := 0; i < w.rtNumber; i++ {
		err := w.k8sClient.Create(context.Background(), &runtimes.Items[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (w WorkerData) Delete() error {
	runtimes, err := w.deleteRuntimeBatch()
	if err != nil {
		return err
	}
	for _, item := range runtimes.Items {
		err = w.k8sClient.Delete(context.Background(), &item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w WorkerData) prepareRuntimeBatch() imv1.RuntimeList {

	baseRuntime := w.rtTemplate.DeepCopy()
	baseRuntime.Name = ""
	baseRuntime.GenerateName = w.namePrefix + "-"
	baseRuntime.Labels["kim.performance.loadId"] = w.loadID

	if baseRuntime.Spec.Shoot.Name == "" {
		baseRuntime.Spec.Shoot.Name = generateRandomName(7) + "-" + w.loadID
	}

	runtimeBatch := imv1.RuntimeList{}

	for i := 0; i < w.rtNumber; i++ {
		runtimeBatch.Items = append(runtimeBatch.Items, *baseRuntime)
	}

	return runtimeBatch
}

func (w WorkerData) deleteRuntimeBatch() (imv1.RuntimeList, error) {
	runtimesToDelete := imv1.RuntimeList{}
	err := w.k8sClient.List(context.Background(), &runtimesToDelete, client.MatchingLabels{"kim.performance.loadId": w.loadID})
	if err != nil {
		return imv1.RuntimeList{}, err
	}
	return runtimesToDelete, nil
}

func generateRandomName(count int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	runes := make([]rune, count)
	for i := range runes {
		runes[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(runes)
}
