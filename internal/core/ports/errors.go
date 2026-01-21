package ports

import "fmt"

type ErrorType string

const (
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeInvalidInput ErrorType = "INVALID_INPUT"
	ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
)

// AppError is a custom error type that allows distinguishing between
// user-facing messages and internal logs, as well as error types.
type AppError struct {
	Type    ErrorType
	Message string // User facing message
	Err     error  // Original error (for logging)
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
}

// Helper constructors
func NewNotFoundError(msg string, err error) *AppError {
	return &AppError{Type: ErrorTypeNotFound, Message: msg, Err: err}
}

func NewInvalidInputError(msg string, err error) *AppError {
	return &AppError{Type: ErrorTypeInvalidInput, Message: msg, Err: err}
}

func NewInternalError(msg string, err error) *AppError {
	return &AppError{Type: ErrorTypeInternal, Message: msg, Err: err}
}
