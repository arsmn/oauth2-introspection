package introspection

import "errors"

var (
	ErrMalformedToken = errors.New("Missing or malformed token")
	ErrUnauthorized   = errors.New("Access credentials are invalid")
	ErrForbidden      = errors.New("Access credentials are not sufficient to access this resource")
)
