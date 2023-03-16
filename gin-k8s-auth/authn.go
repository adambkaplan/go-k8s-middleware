package gink8sauth

import (
	"net/http"

	"github.com/adambkaplan/go-k8s-middleware/lib"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

// K8SAuthenticator returns gin middleware that authenticates a request against a Kubernetes
// apiserver. If the authentication is successful, the middleware provides the k8s user info via
// the "k8s-user-info" context key.
func K8SAuthenticator(client kubernetes.Interface) gin.HandlerFunc {
	authenticator := lib.NewTokenAuthenticator(client)
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			// TODO: Send WWW-Authenticate response header?
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		authenticated, user, err := authenticator.AuthenticateFromHeader(ctx.Request.Context(), authHeader)
		if !authenticated {
			// TODO: Send WWW-Authenticate response header?
			ctx.AbortWithError(http.StatusUnauthorized, err)
		}
		ctx.Set("k8s-user-info", user)
		ctx.Next()
	}

}
