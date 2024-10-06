package grpcserver

import (
	"context"

	"github.com/google/uuid"
	"github.com/interimme/userapi/internal/apperrors"
	"github.com/interimme/userapi/internal/entity"
	"github.com/interimme/userapi/internal/usecase"
	userapi "github.com/interimme/userapi/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server implements the UserServiceServer interface generated from Proto.
type Server struct {
	userapi.UnimplementedUserServiceServer
	UserUseCase *usecase.UserUseCase
}

// NewServer creates a new UserService with the provided UserUseCase.
func NewServer(userUseCase *usecase.UserUseCase) *Server {
	return &Server{
		UserUseCase: userUseCase,
	}
}

// CreateUser implements the CreateUser RPC method.
func (s *Server) CreateUser(ctx context.Context, req *userapi.CreateUserRequest) (*userapi.CreateUserResponse, error) {
	// Validate the request.
	if req.GetUser() == nil {
		return nil, status.Error(codes.InvalidArgument, "user data is required")
	}

	// Convert Proto User to Entity User.
	user := &entity.User{
		Firstname: req.GetUser().GetFirstname(),
		Lastname:  req.GetUser().GetLastname(),
		Email:     req.GetUser().GetEmail(),
		Age:       uint(req.GetUser().GetAge()), // Safe conversion as age is small.
		// ID and Created are set by the usecase.
	}

	// Call the usecase to create the user.
	if err := s.UserUseCase.CreateUser(user); err != nil {
		// Map internal errors to gRPC status codes using a type switch.
		switch e := err.(type) {
		case *apperrors.AppError:
			return nil, status.Error(e.Code, e.Message)
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	// Convert Entity User back to Proto User.
	respUser := &userapi.User{
		Id:        user.ID.String(),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Age:       uint32(user.Age),              // Convert back to uint32 for Proto.
		Created:   timestamppb.New(user.Created), // Proper Timestamp handling.
	}

	return &userapi.CreateUserResponse{User: respUser}, nil
}

// GetUser implements the GetUser RPC method.
func (s *Server) GetUser(ctx context.Context, req *userapi.GetUserRequest) (*userapi.GetUserResponse, error) {
	// Validate the request.
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Parse the UUID.
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	// Call the usecase to retrieve the user.
	user, err := s.UserUseCase.GetUser(userID)
	if err != nil {
		// Map internal errors to gRPC status codes using a type switch.
		switch e := err.(type) {
		case *apperrors.AppError:
			return nil, status.Error(e.Code, e.Message)
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	// Convert Entity User back to Proto User.
	respUser := &userapi.User{
		Id:        user.ID.String(),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Age:       uint32(user.Age),              // Convert back to uint32 for Proto.
		Created:   timestamppb.New(user.Created), // Proper Timestamp handling.
	}

	return &userapi.GetUserResponse{User: respUser}, nil
}

// UpdateUser implements the UpdateUser RPC method.
func (s *Server) UpdateUser(ctx context.Context, req *userapi.UpdateUserRequest) (*userapi.UpdateUserResponse, error) {
	// Validate the request.
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	if req.GetUser() == nil {
		return nil, status.Error(codes.InvalidArgument, "user data is required")
	}

	// Parse the UUID.
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	// Convert Proto User to Entity User.
	user := &entity.User{
		ID:        userID, // Ensure the ID is set from the request.
		Firstname: req.GetUser().GetFirstname(),
		Lastname:  req.GetUser().GetLastname(),
		Email:     req.GetUser().GetEmail(),
		Age:       uint(req.GetUser().GetAge()), // Safe conversion as age is small.
		// Created remains unchanged.
	}

	// Call the usecase to update the user.
	if err := s.UserUseCase.UpdateUser(user); err != nil {
		// Map internal errors to gRPC status codes using a type switch.
		switch e := err.(type) {
		case *apperrors.AppError:
			return nil, status.Error(e.Code, e.Message)
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	// Retrieve the updated user to include the 'created' timestamp.
	updatedUser, err := s.UserUseCase.GetUser(userID)
	if err != nil {
		switch e := err.(type) {
		case *apperrors.AppError:
			return nil, status.Error(e.Code, e.Message)
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	// Convert Entity User back to Proto User.
	respUser := &userapi.User{
		Id:        updatedUser.ID.String(),
		Firstname: updatedUser.Firstname,
		Lastname:  updatedUser.Lastname,
		Email:     updatedUser.Email,
		Age:       uint32(updatedUser.Age),              // Convert back to uint32 for Proto.
		Created:   timestamppb.New(updatedUser.Created), // Proper Timestamp handling.
	}

	return &userapi.UpdateUserResponse{User: respUser}, nil
}

// DeleteUser implements the DeleteUser RPC method.
func (s *Server) DeleteUser(ctx context.Context, req *userapi.DeleteUserRequest) (*userapi.DeleteUserResponse, error) {
	// Validate the request.
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Parse the UUID.
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	// Call the usecase to delete the user.
	if err := s.UserUseCase.DeleteUser(userID); err != nil {
		// Map internal errors to gRPC status codes using a type switch.
		switch e := err.(type) {
		case *apperrors.AppError:
			return nil, status.Error(e.Code, e.Message)
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &userapi.DeleteUserResponse{Message: "User deleted successfully"}, nil
}
