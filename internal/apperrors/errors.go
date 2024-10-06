package apperrors

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

// AppError represents a custom application error tailored for gRPC.
type AppError struct {
	Code    codes.Code // gRPC status code
	Message string     // Error message
}

// Error implements the error interface for AppError.
func (e *AppError) Error() string {
	return fmt.Sprintf("code: %v, message: %s", e.Code, e.Message)
}

// NewAppError creates a new AppError with the specified gRPC code and message.
func NewAppError(code codes.Code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Predefined error instances using gRPC codes.
var (
	ErrBadRequest          = NewAppError(codes.InvalidArgument, "bad request")
	ErrUnauthorized        = NewAppError(codes.Unauthenticated, "unauthorized")
	ErrForbidden           = NewAppError(codes.PermissionDenied, "forbidden")
	ErrNotFound            = NewAppError(codes.NotFound, "user not found")
	ErrConflict            = NewAppError(codes.AlreadyExists, "email already exists")
	ErrInternalServerError = NewAppError(codes.Internal, "internal server error")
)
