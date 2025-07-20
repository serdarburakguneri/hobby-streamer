package constants

const (
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"

	StatusUnauthorized = "Unauthorized"
	StatusForbidden    = "Forbidden"
	StatusNotFound     = "Not Found"
	StatusBadRequest   = "Bad Request"
	StatusOK           = "OK"

	ErrInvalidToken        = "Invalid token"
	ErrTokenExpired        = "Token expired"
	ErrInsufficientPerm    = "Insufficient permissions"
	ErrAuthorizationHeader = "Authorization header required"
	ErrInvalidRequestBody  = "Invalid request body"
	ErrUserNotFound        = "User not found in context"

	BearerPrefix = "Bearer "
)
