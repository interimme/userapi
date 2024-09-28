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

// Predefined error instances
var (
	ErrBadRequest          = &AppError{Code: http.StatusBadRequest, Message: "bad request"}
	ErrUnauthorized        = &AppError{Code: http.StatusUnauthorized, Message: "unauthorized"}
	ErrForbidden           = &AppError{Code: http.StatusForbidden, Message: "forbidden"}
	ErrNotFound            = &AppError{Code: http.StatusNotFound, Message: "user not found"}
	ErrConflict            = &AppError{Code: http.StatusConflict, Message: "email already exists"}
	ErrInternalServerError = &AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
)
