package lib

import (
	"errors"
	"strings"
)

func extractBearerToken(authHeader string) (string, error) {
	s := strings.SplitN(authHeader, " ", 2)
	if len(s) < 2 {
		return "", errors.New("unknown auth token format")
	}
	token := s[1]
	if token == "" {
		return "", errors.New("empty auth token")
	}
	return token, nil
}
