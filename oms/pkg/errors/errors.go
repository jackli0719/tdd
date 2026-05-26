package errors

import "fmt"

// Error codes matching the API contract
const (
	ErrCodeSuccess      = 0
	ErrCodeParam        = 400
	ErrCodeUnauthorized = 401
	ErrCodeForbidden    = 403
	ErrCodeNotFound     = 404
	ErrCodeInternal     = 500
)

// AppError represents an application error with code and message
type AppError struct {
	Code    int    // Error code matching HTTP status
	Message string // Error message
	Err     error  // Original error (optional)
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%d: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error for errors.Is/As
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// NewWithError creates a new AppError with an underlying error
func NewWithError(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// Param creates a parameter validation error (400)
func Param(message string) *AppError {
	return New(ErrCodeParam, message)
}

// Unauthorized creates an unauthorized error (401)
func Unauthorized(message string) *AppError {
	return New(ErrCodeUnauthorized, message)
}

// Forbidden creates a forbidden error (403)
func Forbidden(message string) *AppError {
	return New(ErrCodeForbidden, message)
}

// NotFound creates a not found error (404)
func NotFound(message string) *AppError {
	return New(ErrCodeNotFound, message)
}

// Internal creates an internal server error (500)
func Internal(message string) *AppError {
	return New(ErrCodeInternal, message)
}

// InternalWithError creates an internal server error with underlying error
func InternalWithError(message string, err error) *AppError {
	return NewWithError(ErrCodeInternal, message, err)
}
