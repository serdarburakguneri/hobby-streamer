package errors

import (
	"errors"
	"fmt"
)

type ErrorType string

const (
	ErrorTypeValidation     ErrorType = "validation"
	ErrorTypeNotFound       ErrorType = "not_found"
	ErrorTypeUnauthorized   ErrorType = "unauthorized"
	ErrorTypeForbidden      ErrorType = "forbidden"
	ErrorTypeConflict       ErrorType = "conflict"
	ErrorTypeInternal       ErrorType = "internal"
	ErrorTypeExternal       ErrorType = "external"
	ErrorTypeTransient      ErrorType = "transient"
	ErrorTypeTimeout        ErrorType = "timeout"
	ErrorTypeCircuitBreaker ErrorType = "circuit_breaker"
)

type AppError struct {
	Type    ErrorType
	Message string
	Cause   error
	Context map[string]interface{}
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func (e *AppError) IsTransient() bool {
	return e.Type == ErrorTypeTransient || e.Type == ErrorTypeTimeout
}

func (e *AppError) IsExternal() bool {
	return e.Type == ErrorTypeExternal || e.Type == ErrorTypeCircuitBreaker
}

func (e *AppError) HTTPStatus() int {
	switch e.Type {
	case ErrorTypeValidation:
		return 400
	case ErrorTypeNotFound:
		return 404
	case ErrorTypeUnauthorized:
		return 401
	case ErrorTypeForbidden:
		return 403
	case ErrorTypeConflict:
		return 409
	case ErrorTypeTransient:
		return 503
	case ErrorTypeTimeout:
		return 504
	case ErrorTypeCircuitBreaker:
		return 503
	case ErrorTypeExternal:
		return 502
	default:
		return 500
	}
}

func NewValidationError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Cause:   cause,
	}
}

func NewNotFoundError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Cause:   cause,
	}
}

func NewUnauthorizedError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
		Cause:   cause,
	}
}

func NewForbiddenError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeForbidden,
		Message: message,
		Cause:   cause,
	}
}

func NewConflictError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
		Cause:   cause,
	}
}

func NewInternalError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Cause:   cause,
	}
}

func NewExternalError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeExternal,
		Message: message,
		Cause:   cause,
	}
}

func NewTransientError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeTransient,
		Message: message,
		Cause:   cause,
	}
}

func NewTimeoutError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeTimeout,
		Message: message,
		Cause:   cause,
	}
}

func NewCircuitBreakerError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeCircuitBreaker,
		Message: message,
		Cause:   cause,
	}
}

func WithContext(err error, context map[string]interface{}) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		if appErr.Context == nil {
			appErr.Context = make(map[string]interface{})
		}
		for k, v := range context {
			appErr.Context[k] = v
		}
		return appErr
	}

	return &AppError{
		Type:    ErrorTypeInternal,
		Message: err.Error(),
		Cause:   err,
		Context: context,
	}
}

func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

func GetErrorType(err error) ErrorType {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type
	}
	return ErrorTypeInternal
}

func IsTransient(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.IsTransient()
	}
	return false
}

func IsExternal(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.IsExternal()
	}
	return false
}

func IsValidationError(err error) bool {
	return GetErrorType(err) == ErrorTypeValidation
}

func IsNotFoundError(err error) bool {
	return GetErrorType(err) == ErrorTypeNotFound
}

func IsConflictError(err error) bool {
	return GetErrorType(err) == ErrorTypeConflict
}

func WrapWithContext(err error, context string) *AppError {
	errorType := GetErrorType(err)

	switch errorType {
	case ErrorTypeValidation:
		return NewValidationError(context, err)
	case ErrorTypeNotFound:
		return NewNotFoundError(context, err)
	case ErrorTypeTransient:
		return NewTransientError(context, err)
	default:
		return NewInternalError(context, err)
	}
}
