package lib

import (
	"context"
	"fmt"
	"testing"

	authnv1 "k8s.io/api/authentication/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"
)

func TestAuthenticateFromHeader(t *testing.T) {
	ctx := context.Background()
	bearerToken := "SGVsbyB3b3JsZCEK"
	authHeader := fmt.Sprintf("Bearer %s", bearerToken)
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
	authn := NewTokenAuthenticator(fake)
	_, user, err := authn.AuthenticateFromHeader(ctx, authHeader)
	if err != nil {
		t.Errorf("failed to authenticate: %v", err)
	}
	if user.Username != "kubeadmin" {
		t.Errorf("expected username %q, got %q", "kubeadmin", user.Username)
	}
	if user.UID != "1234" {
		t.Errorf("expected UID %q, got %q", "1234", user.UID)
	}
}
