package lib

import "net/http"

var (
	ErrBadRequest      = AuthErr(http.StatusBadRequest, "bad request: invalid token format")
	ErrUnauthenticated = AuthErr(http.StatusUnauthorized, "bad request: unauthenticated")
	ErrInternalAuthN   = AuthErr(http.StatusInternalServerError, "internal error: failed to authenticate user")
	ErrInternalAuthZ   = AuthErr(http.StatusInternalServerError, "internal error: failed to authorize user")
)

type authErr struct {
	message    string
	statusCode int
}

func AuthErr(statusCode int, message string) authErr {
	return authErr{
		statusCode: statusCode,
		message:    message,
	}
}

func (a authErr) Error() string {
	return a.message
}

func (a authErr) Is(target error) bool {
	return a.message == target.Error()
}
