package lib

import (
	"context"

	authnv1 "k8s.io/api/authentication/v1"
	authzv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	authzclient "k8s.io/client-go/kubernetes/typed/authorization/v1"
)

type SubjectAuthorizer struct {
	client authzclient.AuthorizationV1Interface
}

// NewSubjectAuthorizer creates a new authorizer with the provided Kubernetes client.
func NewSubjectAuthorizer(client kubernetes.Interface) *SubjectAuthorizer {
	return &SubjectAuthorizer{
		client: client.AuthorizationV1(),
	}
}

// AuthorizeSubject checks if the user is allowed to perform the provided action
func (s *SubjectAuthorizer) AuthorizeSubject(ctx context.Context, user authnv1.UserInfo, namespace, group, resource, verb string) (bool, error) {
	sar := &authzv1.SubjectAccessReview{
		Spec: authzv1.SubjectAccessReviewSpec{
			User:   user.Username,
			Groups: user.Groups,
			ResourceAttributes: &authzv1.ResourceAttributes{
				Namespace: namespace,
				Group:     group,
				Resource:  resource,
				Verb:      verb,
			},
		},
	}
	authzExtras := map[string]authzv1.ExtraValue{}
	for key, authnExtra := range user.Extra {
		authzExtras[key] = authzv1.ExtraValue(authnExtra)
	}

	returnSar, err := s.client.SubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	return returnSar.Status.Allowed, nil
}
