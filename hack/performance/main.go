package main

import (
	"log/slog"

	"github.com/kyma-project/infrastructure-manager/hack/performance/cmd"
)

func main() {

	log := slog.Default()
	operation, worker, err := cmd.Execute()

	if err != nil {
		log.Error("error during initialization", err)
		return
	}

	switch operation {
	case cmd.Create:
		err = worker.Create()
		if err != nil {
			log.Error("error in creation operation", err)
		}
	case cmd.Delete:
		err = worker.Delete()
		if err != nil {
			log.Error("error in deletion operation", err)
		}
	default:
		log.Info("unknown operation specified, exiting")
	}
}
