package errors

import "net/http"

// AppError represents a custom application error
type AppError struct {
	Code    int    // HTTP status code
	Message string // Error message
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrBadRequest          = &AppError{Code: http.StatusBadRequest, Message: "Bad request"}
	ErrUnauthorized        = &AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrForbidden           = &AppError{Code: http.StatusForbidden, Message: "Forbidden"}
	ErrNotFound            = &AppError{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrConflict            = &AppError{Code: http.StatusConflict, Message: "Resource conflict"}
	ErrInternalServerError = &AppError{Code: http.StatusInternalServerError, Message: "Internal server error"}
)
