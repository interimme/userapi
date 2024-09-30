package usecase

import (
	"github.com/interimme/userapi/internal/entity"
	"github.com/interimme/userapi/internal/usecase/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userUseCase := NewUserUseCase(mockRepo)

	user := &entity.User{
		Firstname: "Alice",
		Lastname:  "Smith",
		Email:     "alice@example.com",
		Age:       28,
	}

	// Mock GetByEmail to return nil, indicating email does not exist
	mockRepo.On("GetByEmail", "alice@example.com").Return(nil, nil)
	// Mock Create to return nil, indicating successful creation
	mockRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(nil)

	err := userUseCase.CreateUser(user)

	require.NoError(t, err, "Expected no error when creating user")
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_ValidationError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userUseCase := NewUserUseCase(mockRepo)

	user := &entity.User{
		Firstname: "",
		Lastname:  "Smith",
		Email:     "alice@example.com",
		Age:       28,
	}

	err := userUseCase.CreateUser(user)

	require.Error(t, err, "Expected an error due to validation failure")
	assert.Equal(t, "firstname is required", err.Error())
}

func TestCreateUser_EmailExists(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userUseCase := NewUserUseCase(mockRepo)

	user := &entity.User{
		Firstname: "Alice",
		Lastname:  "Smith",
		Email:     "alice@example.com",
		Age:       28,
	}

	// Mock existing user with the same email
	mockRepo.On("GetByEmail", "alice@example.com").Return(&entity.User{}, nil)

	err := userUseCase.CreateUser(user)

	require.Error(t, err, "Expected an error due to existing email")
	assert.Equal(t, "email already exists", err.Error())
}

func TestGetUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userUseCase := NewUserUseCase(mockRepo)

	userID := uuid.New()
	user := &entity.User{
		ID:        userID,
		Firstname: "Alice",
		Lastname:  "Smith",
		Email:     "alice@example.com",
		Age:       28,
	}

	// Mock GetByID to return the user
	mockRepo.On("GetByID", userID).Return(user, nil)

	result, err := userUseCase.GetUser(userID)

	require.NoError(t, err, "Expected no error when getting user")
	assert.Equal(t, user, result)
}

func TestGetUser_NotFound(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userUseCase := NewUserUseCase(mockRepo)

	userID := uuid.New()

	// Mock GetByID to return ErrRecordNotFound
	mockRepo.On("GetByID", userID).Return(nil, gorm.ErrRecordNotFound)

	result, err := userUseCase.GetUser(userID)

	require.Error(t, err, "Expected an error due to user not found")
	require.Nil(t, result)
	assert.Equal(t, "user not found", err.Error())
}

func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userUseCase := NewUserUseCase(mockRepo)

	userID := uuid.New()
	user := &entity.User{
		ID:        userID,
		Firstname: "Alice",
		Lastname:  "Johnson",
		Email:     "alice.johnson@example.com",
		Age:       29,
	}

	existingUser := &entity.User{
		ID:        userID,
		Firstname: "Alice",
		Lastname:  "Smith",
		Email:     "alice@example.com",
		Age:       28,
	}

	// Mock GetByID to return the existing user
	mockRepo.On("GetByID", userID).Return(existingUser, nil)
	// Mock Update to return nil
	mockRepo.On("Update", existingUser).Return(nil)

	err := userUseCase.UpdateUser(user)

	require.NoError(t, err, "Expected no error when updating user")
	assert.Equal(t, "Alice", existingUser.Firstname)
	assert.Equal(t, "Johnson", existingUser.Lastname)
	assert.Equal(t, "alice.johnson@example.com", existingUser.Email)
	assert.Equal(t, uint(29), existingUser.Age)
}

func TestUpdateUser_ValidationError(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userUseCase := NewUserUseCase(mockRepo)

	userID := uuid.New()
	user := &entity.User{
		ID:        userID,
		Firstname: "",
		Lastname:  "Johnson",
		Email:     "alice.johnson@example.com",
		Age:       29,
	}

	err := userUseCase.UpdateUser(user)

	require.Error(t, err, "Expected an error due to validation failure")
	assert.Equal(t, "firstname is required", err.Error())
}

func TestUpdateUser_NotFound(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userUseCase := NewUserUseCase(mockRepo)

	userID := uuid.New()
	user := &entity.User{
		ID:        userID,
		Firstname: "Alice",
		Lastname:  "Johnson",
		Email:     "alice.johnson@example.com",
		Age:       29,
	}

	// Mock GetByID to return ErrRecordNotFound
	mockRepo.On("GetByID", userID).Return(nil, gorm.ErrRecordNotFound)

	err := userUseCase.UpdateUser(user)

	require.Error(t, err, "Expected an error due to user not found")
	assert.Equal(t, "user not found", err.Error())
}

func TestDeleteUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userUseCase := NewUserUseCase(mockRepo)

	userID := uuid.New()
	user := &entity.User{
		ID: userID,
	}

	// Mock GetByID to return the user
	mockRepo.On("GetByID", userID).Return(user, nil)
	// Mock Delete to return nil
	mockRepo.On("Delete", user).Return(nil)

	err := userUseCase.DeleteUser(userID)

	require.NoError(t, err, "Expected no error when deleting user")
}

func TestDeleteUser_NotFound(t *testing.T) {
	mockRepo := new(mocks.UserRepository)
	userUseCase := NewUserUseCase(mockRepo)

	userID := uuid.New()

	// Mock GetByID to return ErrRecordNotFound
	mockRepo.On("GetByID", userID).Return(nil, gorm.ErrRecordNotFound)

	err := userUseCase.DeleteUser(userID)

	require.Error(t, err, "Expected an error due to user not found")
	assert.Equal(t, "user not found", err.Error())
}
