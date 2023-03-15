package lib

import "testing"

func TestExtractBearerToken(t *testing.T) {
	cases := []struct {
		name          string
		token         string
		expectedToken string
		expectErr     bool
	}{
		{
			name:          "bearer",
			token:         "Bearer av2323908va3e3==",
			expectedToken: "av2323908va3e3==",
		},
		{
			name:      "no token",
			token:     "Bearer ",
			expectErr: true,
		},
		{
			name:      "no space",
			token:     "Basic",
			expectErr: true,
		},
		{
			name:      "empty",
			token:     "",
			expectErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := extractBearerToken(tc.token)
			if tc.expectErr && err == nil {
				t.Error("expected error, got none")
			}
			if err != nil && !tc.expectErr {
				t.Errorf("unexpected error: %v", err)
			}
			if token != tc.expectedToken {
				t.Errorf("expected token %q, got %q", tc.expectedToken, token)
			}
		})
	}
}
