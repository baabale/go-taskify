package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Common errors
var (
	ErrNotFound          = errors.New("resource not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrDatabaseOperation = errors.New("database operation failed")
	ErrInternal         = errors.New("internal server error")
)

// AppError represents an application error
type AppError struct {
	Err        error
	Message    string
	StatusCode int
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap implements the unwrap interface
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewNotFound creates a new not found error
func NewNotFound(resource string) *AppError {
	return &AppError{
		Err:        ErrNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

// NewInvalidInput creates a new invalid input error
func NewInvalidInput(message string) *AppError {
	return &AppError{
		Err:        ErrInvalidInput,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewDatabaseError creates a new database error
func NewDatabaseError(err error) *AppError {
	return &AppError{
		Err:        fmt.Errorf("%w: %v", ErrDatabaseOperation, err),
		Message:    "Database operation failed",
		StatusCode: http.StatusInternalServerError,
	}
}

// NewInternalError creates a new internal error
func NewInternalError(err error) *AppError {
	return &AppError{
		Err:        fmt.Errorf("%w: %v", ErrInternal, err),
		Message:    "Internal server error",
		StatusCode: http.StatusInternalServerError,
	}
}

// AsAppError converts an error to an AppError
func AsAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	// Handle validation errors specifically
	if errors.Is(err, ErrInvalidInput) {
		return NewInvalidInput(err.Error())
	}

	// Default to internal error
	return NewInternalError(err)
}
