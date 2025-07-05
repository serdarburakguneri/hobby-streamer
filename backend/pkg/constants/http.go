package constants

const (
	// HTTP Headers
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"

	// HTTP Status Messages
	StatusUnauthorized = "Unauthorized"
	StatusForbidden    = "Forbidden"
	StatusNotFound     = "Not Found"
	StatusBadRequest   = "Bad Request"
	StatusOK           = "OK"

	// Common Error Messages
	ErrInvalidToken        = "Invalid token"
	ErrTokenExpired        = "Token expired"
	ErrInsufficientPerm    = "Insufficient permissions"
	ErrAuthorizationHeader = "Authorization header required"
	ErrInvalidRequestBody  = "Invalid request body"
	ErrUserNotFound        = "User not found in context"

	// Common Prefixes
	BearerPrefix = "Bearer "
)
