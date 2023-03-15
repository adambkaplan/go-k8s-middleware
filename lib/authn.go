package lib

import (
	"context"

	authnv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	authnclient "k8s.io/client-go/kubernetes/typed/authentication/v1"
)

type TokenAuthenticator struct {
	client authnclient.AuthenticationV1Interface
}

func NewTokenAuthenticator(client kubernetes.Interface) *TokenAuthenticator {
	return &TokenAuthenticator{
		client: client.AuthenticationV1(),
	}
}

func (t *TokenAuthenticator) AuthenticateFromHeader(ctx context.Context, authHeader string) (authnv1.UserInfo, error) {
	bearerToken, err := extractBearerToken(authHeader)
	if err != nil || bearerToken == "" {
		return authnv1.UserInfo{}, ErrBadRequest
	}
	tokenReview, err := t.client.TokenReviews().Create(ctx, &authnv1.TokenReview{
		Spec: authnv1.TokenReviewSpec{
			Token: bearerToken,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return authnv1.UserInfo{}, ErrInternalAuthN
	}
	if !tokenReview.Status.Authenticated {
		return tokenReview.Status.User, ErrUnauthenticated
	}

	return tokenReview.Status.User, nil
}
