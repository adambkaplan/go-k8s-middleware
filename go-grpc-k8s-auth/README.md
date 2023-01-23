# go-grpc-k8s-auth

gRPC middleware that provides authentication and authorization support for
Kubernetes applications.

## Try the example

To build and run the example, you need the following:

1. A Kubernetes cluster - you can deploy one locally with
  [kind](https://kind.sigs.k8s.io/).
2. [ko](https://github.com/ko-build/ko) to build the application.

Next, run `make deploy` to deploy the server on your cluster.

To connect to the server, first run `kubectl port-forward` in a separate
terminal:

```sh
$ kubectl port-forward -n go-grpc-k8s-auth-example service/greeter-grpc-service 50051:grpc
```

Then run the client with the `-name` flag to change the response of the greeing:

```sh
$ go run example/client/main.go -name Alice
```
