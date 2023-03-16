# gin-k8s-auth

Kubernetes auth middlewares for the gin web framework.
There are separate middlewares for authentication and authorization.

## K8SAuthenticator

`K8SAuthenticator` uses Kubernetes TokenAccessReviews to authenticate an incoming connection.
This middleware can be used for authentication against incoming requests to all endpoints,
or a subset of endpoints.

To use the authenticator, it must be provided a Kubernetes client interface object.
See the [server example](examples/README.md) for details of how to create one from a
local `kubeconfig` file.

Once the Kubernetes interface has been created, you can use the `K8SAuthenticator` just like
any other gin middleware:

```go
    r := gin.Default()
    // Use the authenticator for all requests
    r.Use(gink8sauth.K8SAuthenticator(client))
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world!",
		})
	})
```
