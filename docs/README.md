# Docs

## Overview

This folder contains documents that relate to the project.

## Development

Run `make test` to see if all tests are passing. 

## Configuration

You can configure the Infrastructure Manager deployment with the following arguments:
1. `gardener-kubeconfig-path` - defines the path to the Gardener project kubeconfig used during API calls
2. `gardener-project` - the name of the Gardener project where the infrastructure operations are performed
3. `minimal-rotation-time` - the ratio determines what is the minimal time that needs to pass to rotate the certificate
4. `kubeconfig-expiration-time` - maximum time after which kubeconfig is rotated. The rotation happens between (`minimal-rotation-time` * `kubeconfig-expiration-time`) and `kubeconfig-expiration-time`.
4. `gardener-request-timeout` - specifies the timeout for requests to Gardener. Default value is `60s`.
5. `runtime-reconciler-enabled` - feature flag responsible for enabling the runtime reconciler. Default value is `false`.


See [manager_gardener_secret_patch.yaml](../config/default/manager_gardener_secret_patch.yaml) for default values.

## Troubleshooting

> TBD: List potential issues and provide tips on how to avoid or solve them. To structure the content, use the following sections:
>
> - **Symptom**
> - **Cause**
> - **Remedy**
