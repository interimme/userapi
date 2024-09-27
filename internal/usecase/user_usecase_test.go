package usecase

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
	"userapi/internal/entity"
	"userapi/internal/usecase/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	userUseCase := NewUserUseCase(mockRepo)

	user := &entity.User{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
	}

	mockRepo.On("GetByEmail", user.Email).Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*entity.User")).Return(nil)

	err := userUseCase.CreateUser(user)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.WithinDuration(t, time.Now(), user.Created, time.Second)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_EmailExists(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	userUseCase := NewUserUseCase(mockRepo)

	existingUser := &entity.User{
		ID:        uuid.New(),
		Firstname: "Jane",
		Lastname:  "Doe",
		Email:     "jane.doe@example.com",
		Age:       28,
	}

	mockRepo.On("GetByEmail", existingUser.Email).Return(existingUser, nil)

	err := userUseCase.CreateUser(existingUser)

	assert.Error(t, err)
	assert.Equal(t, "email already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_InvalidInput(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	userUseCase := NewUserUseCase(mockRepo)

	user := &entity.User{
		Firstname: "", // Missing firstname
		Lastname:  "Doe",
		Email:     "invalid-email",
		Age:       200, // Invalid age
	}

	err := userUseCase.CreateUser(user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "firstname is required")
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestGetUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	userUseCase := NewUserUseCase(mockRepo)

	userID := uuid.New()
	user := &entity.User{
		ID:        userID,
		Firstname: "Jane",
		Lastname:  "Doe",
		Email:     "jane.doe@example.com",
		Age:       28,
		Created:   time.Now(),
	}

	mockRepo.On("GetByID", userID).Return(user, nil)

	result, err := userUseCase.GetUser(userID)

	assert.NoError(t, err)
	assert.Equal(t, user, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	userUseCase := NewUserUseCase(mockRepo)

	userID := uuid.New()

	mockRepo.On("GetByID", userID).Return(nil, errors.New("user not found"))

	result, err := userUseCase.GetUser(userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	userUseCase := NewUserUseCase(mockRepo)

	userID := uuid.New()
	existingUser := &entity.User{
		ID:        userID,
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
		Created:   time.Now(),
	}

	updatedUser := &entity.User{
		ID:        userID,
		Firstname: "John",
		Lastname:  "Smith",
		Email:     "john.smith@example.com",
		Age:       31,
	}

	mockRepo.On("GetByID", userID).Return(existingUser, nil)
	mockRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(nil)

	err := userUseCase.UpdateUser(updatedUser)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_NotFound(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	userUseCase := NewUserUseCase(mockRepo)

	userID := uuid.New()
	updatedUser := &entity.User{
		ID:        userID,
		Firstname: "John",
		Lastname:  "Smith",
		Email:     "john.smith@example.com",
		Age:       31,
	}

	mockRepo.On("GetByID", userID).Return(nil, errors.New("user not found"))

	err := userUseCase.UpdateUser(updatedUser)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_InvalidInput(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	userUseCase := NewUserUseCase(mockRepo)

	userID := uuid.New()
	updatedUser := &entity.User{
		ID:        userID,
		Firstname: "", // Missing firstname
		Lastname:  "Smith",
		Email:     "invalid-email",
		Age:       200, // Invalid age
	}

	err := userUseCase.UpdateUser(updatedUser)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "firstname is required")
	mockRepo.AssertNotCalled(t, "Update", mock.Anything)
}

func TestDeleteUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	userUseCase := NewUserUseCase(mockRepo)

	userID := uuid.New()
	user := &entity.User{
		ID:        userID,
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
		Created:   time.Now(),
	}

	mockRepo.On("GetByID", userID).Return(user, nil)
	mockRepo.On("Delete", user).Return(nil)

	err := userUseCase.DeleteUser(userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_NotFound(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	userUseCase := NewUserUseCase(mockRepo)

	userID := uuid.New()

	mockRepo.On("GetByID", userID).Return(nil, errors.New("user not found"))

	err := userUseCase.DeleteUser(userID)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}
