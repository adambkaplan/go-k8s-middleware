package gink8sauth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	authnv1 "k8s.io/api/authentication/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"

	"github.com/gin-gonic/gin"
)

func TestK8SAuthenticator(t *testing.T) {
	router := gin.New()
	fake := fake.NewSimpleClientset()
	fake.PrependReactor("create", "tokenreviews", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
		tr := action.(ktesting.CreateActionImpl).Object.(*authnv1.TokenReview)
		tr.Status = authnv1.TokenReviewStatus{
			Authenticated: true,
			User: authnv1.UserInfo{
				Username: "kubeadmin",
				UID:      "1234",
			},
		}
		return true, tr, nil
	})
	router.Use(K8SAuthenticator(fake))
	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world!")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/hello", nil)
	req.Header.Add("authorization", "Bearer: SGVsbyB3b3JsZCEK")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
	body := w.Body.String()
	if body != "hello world!" {
		t.Errorf("expected response %q, got %q", "hello world!", body)
	}

}
