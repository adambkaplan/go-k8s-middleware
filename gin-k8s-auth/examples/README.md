# Gin K8S Auth Example

This directory has an example server which utilizes the k8s auth middleware,
and a client which connects to the server with credentials provided by a
`kubeconfig` file.

## Prerequisites

To run the example, you must have:

- Go 1.19 or later
- Access to a Kubernetes cluster

## Running the example

- Launch the server in a terminal shell:

  ```sh
  $ go run ./server/main.go
  ```

- In a separate terminal shell, run the client:

  ```sh
  $ go run ./client/main.go
  ```
