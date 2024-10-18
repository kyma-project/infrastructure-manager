package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/infrastructure-manager/hack/performance/action"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type OperationType int

const (
	Create OperationType = iota
	Delete
	Unknown
)

func Execute() (OperationType, action.Worker, error) {
	var parsedRuntime imv1.Runtime

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	loadID := createCmd.String("load-id", "", "the identifier (label) of the created load (required)")
	namePrefix := createCmd.String("name-prefix", "", "the prefix used to generate each runtime name (required)")
	kubeconfig := createCmd.String("kubeconfig", "", "the path to the kubeconfig file (required)")
	rtNumber := createCmd.Int("rt-number", 0, "the number of the runtimes to be created (required)")
	templatePath := createCmd.String("rt-template", "", "the path to the yaml file with the runtime template (required)")
	runOnCi := createCmd.String("run-on-ci", "false", "identifies if the load is running on CI")

	loadIDDelete := deleteCmd.String("load-id", "", "the identifier (label) of the created load (required)")
	kubeconfigDelete := deleteCmd.String("kubeconfig", "", "the path to the kubeconfig file (required)")

	if len(os.Args) < 2 {
		fmt.Println("expected 'create' or 'delete' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		if *loadID == "" || *namePrefix == "" || *kubeconfig == "" || *rtNumber == 0 {
			fmt.Println("all flags --load-id, --name-prefix, --kubeconfig and --rt-number are required")
			createCmd.Usage()
			os.Exit(1)
		}

		if *templatePath != "" {
			file, err := os.Open(*templatePath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error opening file:", err)
				return Unknown, nil, err
			}
			defer func(file *os.File) {
				err = file.Close()
				if err != nil {
					fmt.Fprintln(os.Stderr, "error closing file:", err)
				}
			}(file)
			parsedRuntime, err = readFromSource(file)
			if err != nil {
				return Unknown, nil, err
			}
		}

		if *runOnCi == "false" {
			var response string
			fmt.Printf("Do you want to create %d runtimes? [y/n]: ", *rtNumber)
			fmt.Scanln(&response)
			if response != "y" {
				fmt.Println("Operation cancelled.")
				os.Exit(1)
			}
		}
		fmt.Printf("Creating load with ID: %s, Name Prefix: %s, Kubeconfig: %s, Runtime Number: %d\n", *loadID, *namePrefix, *kubeconfig, *rtNumber)
		worker, err := action.NewWorker(*loadID, *namePrefix, *kubeconfig, *rtNumber, parsedRuntime)
		return Create, worker, err
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *loadIDDelete == "" || *kubeconfigDelete == "" {
			fmt.Println("all flags --load-id and --kubeconfig are required")
			deleteCmd.Usage()
			os.Exit(1)
		}
		fmt.Printf("Deleting load with ID: %s, Kubeconfig: %s\n", *loadIDDelete, *kubeconfig)
		worker, err := action.NewWorker(*loadIDDelete, "", *kubeconfigDelete, 0, imv1.Runtime{})
		return Delete, worker, err
	default:
		fmt.Println("expected 'create' or 'delete' subcommands")
		os.Exit(1)
	}
	return Unknown, nil, nil
}

func readFromSource(reader io.Reader) (imv1.Runtime, error) {
	data, err := io.ReadAll(reader)
	runtime, err := parseInputToRuntime(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing input:", err)
		return imv1.Runtime{}, err
	}
	return runtime, nil
}

func parseInputToRuntime(data []byte) (imv1.Runtime, error) {
	runtime := imv1.Runtime{}
	err := yaml.Unmarshal(data, &runtime)
	if err != nil {
		return imv1.Runtime{}, err
	}
	return runtime, nil
}
