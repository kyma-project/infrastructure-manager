# Runtime Load Generator

## Overview

The `rt-load` program is used to generate performance test loads by creating or deleting Runtime resources in a Kubernetes cluster. All of the Runtime Custom Resources created by this program will be linked to the same Gardener cluster.

## Building the Binary

To build the `rt-load` binary from the `main.go` file, follow these steps:

1. Open a terminal and navigate to the directory containing the `main.go` file.
2. Run the following command to build the binary:

    ```sh
    go build -o rt-load main.go
    ```
   
## Usage

```sh
rt-load <command> [options]
```

## Commands

### `create`

Creates a specified number of Runtime resources.

#### Options

- `--load-id <STRING>`: The identifier (label) of the created load (**required**).
- `--name-prefix <STRING>`: The prefix used to generate each runtime name (**required**).
- `--kubeconfig <STRING>`: The path to the Kubeconfig file (**required**).
- `--rt-number <INT>`: The number of the runtimes to be created (**required**).
- `--rt-template <STRING>`: The absolute path to the YAML file with the runtime template (**required**).
- `--run-on-ci <STRING>`: Identifies if the load is running on CI (**optional**, default is `false`).

#### Example

```sh
./rt-load create --load-id my-load --name-prefix my-runtime --kubeconfig /path/to/kubeconfig --rt-number 10 --rt-template /path/to/template.yaml
```

### `delete`

Deletes Runtime resources based on the specified load ID.

#### Options

- `--load-id <LOAD-ID>`: The identifier of the created load (**required**).
- `--kubeconfig <FILE>`: The path to the Kubeconfig file (**required**).

#### Example

```sh
./rt-load delete --load-id my-load --kubeconfig /path/to/kubeconfig
```

## Running on CI

When running the `rt-load` program in a Continuous Integration (CI) environment, you can use the `--run-on-ci` option to bypass interactive prompts. This is useful for automated CI/CD pipelines where user interaction is not possible.

### Example

To create runtime resources in a CI environment, use the following command:

```sh
./rt-load create --load-id my-load --name-prefix my-runtime --kubeconfig /path/to/kubeconfig --rt-number 10 --rt-template /path/to/template.yaml --run-on-ci true
```

In this example, the `--run-on-ci` option is set to `true`, which ensures that the program runs without requiring any user input.

## Notes

- Ensure that the `kubeconfig` file points to the correct Kubernetes cluster where the runtime resources are to be created or deleted.
- The `--load-id ` will be included as a value of the label `kim.performance.loadId` on the Runtime CR that was created by this program.
- The `--run-on-ci` option is useful for automated CI/CD pipelines to bypass interactive prompts.