package lib

import (
	"context"
	"errors"
	"testing"

	authnv1 "k8s.io/api/authentication/v1"
	authzv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"
)

func TestAuthorizeSubject(t *testing.T) {
	cases := []struct {
		name             string
		namespace        string
		apiGroup         string
		apiResource      string
		verb             string
		thrownError      string
		expectAuthorized bool
	}{
		{
			name:             "allowed",
			namespace:        "default",
			apiGroup:         "",
			apiResource:      "pods",
			verb:             "get",
			expectAuthorized: true,
		},
		{
			name:             "denied",
			namespace:        "default",
			apiGroup:         "",
			apiResource:      "pods",
			verb:             "create",
			expectAuthorized: false,
		},
		{
			name:             "error",
			namespace:        "default",
			apiGroup:         "",
			apiResource:      "pods",
			verb:             "get",
			thrownError:      "apiserver throttle, too many requests",
			expectAuthorized: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			userInfo := authnv1.UserInfo{
				Username: "reader",
			}
			fake := fake.NewSimpleClientset()
			fake.PrependReactor("create", "subjectaccessreviews", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
				if tc.thrownError != "" {
					return true, nil, errors.New(tc.thrownError)
				}
				sar := action.(ktesting.CreateActionImpl).Object.(*authzv1.SubjectAccessReview)
				resourceAttr := sar.Spec.ResourceAttributes
				allowed := true
				if resourceAttr.Namespace != "default" {
					allowed = false
				}
				if resourceAttr.Group != "" && resourceAttr.Group != "core" {
					allowed = false
				}
				if resourceAttr.Resource != "pods" {
					allowed = false
				}
				if resourceAttr.Verb != "get" && resourceAttr.Verb != "list" {
					allowed = false
				}
				sar.Status = authzv1.SubjectAccessReviewStatus{
					Allowed: allowed,
				}
				return true, sar, nil
			})
			authz := NewSubjectAuthorizer(fake)
			authorized, err := authz.AuthorizeSubject(ctx, userInfo, tc.namespace, "", tc.apiResource, tc.verb)
			if tc.expectAuthorized != authorized {
				t.Errorf("expected subject authorization %t, got %t", tc.expectAuthorized, authorized)
			}
			if err != nil && tc.thrownError == "" {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
