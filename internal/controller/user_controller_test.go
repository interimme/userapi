package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"userapi/internal/controller/mocks"
	"userapi/internal/entity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_Success(t *testing.T) {
	mockUseCase := new(mocks.UserUseCaseMock)
	userController := NewUserController(mockUseCase)

	user := &entity.User{
		Firstname: "Alice",
		Lastname:  "Smith",
		Email:     "alice.smith@example.com",
		Age:       25,
	}

	mockUseCase.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonValue, _ := json.Marshal(user)
	c.Request, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	userController.CreateUser(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	mockUseCase := new(mocks.UserUseCaseMock)
	userController := NewUserController(mockUseCase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	invalidJSON := []byte(`{"firstname": "Alice", "lastname": "Smith",`) // Malformed JSON
	c.Request, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	userController.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUseCase.AssertNotCalled(t, "CreateUser", mock.Anything)
}

func TestCreateUser_ValidationError(t *testing.T) {
	mockUseCase := new(mocks.UserUseCaseMock)
	userController := NewUserController(mockUseCase)

	user := &entity.User{
		Firstname: "", // Missing firstname
		Lastname:  "Smith",
		Email:     "invalid-email",
		Age:       200, // Invalid age
	}

	validationError := errors.New("firstname is required")
	mockUseCase.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(validationError)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonValue, _ := json.Marshal(user)
	c.Request, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	userController.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestGetUser_Success(t *testing.T) {
	mockUseCase := new(mocks.UserUseCaseMock)
	userController := NewUserController(mockUseCase)

	userID := uuid.New()
	user := &entity.User{
		ID:        userID,
		Firstname: "Bob",
		Lastname:  "Jones",
		Email:     "bob.jones@example.com",
		Age:       40,
		Created:   time.Now(),
	}

	mockUseCase.On("GetUser", userID).Return(user, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: userID.String()}}
	c.Request, _ = http.NewRequest("GET", "/user/"+userID.String(), nil)

	userController.GetUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestGetUser_InvalidUUID(t *testing.T) {
	mockUseCase := new(mocks.UserUseCaseMock)
	userController := NewUserController(mockUseCase)

	invalidID := "invalid-uuid"

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: invalidID}}
	c.Request, _ = http.NewRequest("GET", "/user/"+invalidID, nil)

	userController.GetUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUseCase.AssertNotCalled(t, "GetUser", mock.Anything)
}

func TestGetUser_NotFound(t *testing.T) {
	mockUseCase := new(mocks.UserUseCaseMock)
	userController := NewUserController(mockUseCase)

	userID := uuid.New()

	mockUseCase.On("GetUser", userID).Return(nil, errors.New("user not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: userID.String()}}
	c.Request, _ = http.NewRequest("GET", "/user/"+userID.String(), nil)

	userController.GetUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	mockUseCase := new(mocks.UserUseCaseMock)
	userController := NewUserController(mockUseCase)

	userID := uuid.New()
	user := &entity.User{
		Firstname: "Charlie",
		Lastname:  "Brown",
		Email:     "charlie.brown@example.com",
		Age:       35,
	}

	mockUseCase.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: userID.String()}}

	jsonValue, _ := json.Marshal(user)
	c.Request, _ = http.NewRequest("PATCH", "/user/"+userID.String(), bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	userController.UpdateUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestUpdateUser_InvalidUUID(t *testing.T) {
	mockUseCase := new(mocks.UserUseCaseMock)
	userController := NewUserController(mockUseCase)

	invalidID := "invalid-uuid"

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: invalidID}}
	c.Request, _ = http.NewRequest("PATCH", "/user/"+invalidID, nil)

	userController.UpdateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUseCase.AssertNotCalled(t, "UpdateUser", mock.Anything)
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	mockUseCase := new(mocks.UserUseCaseMock)
	userController := NewUserController(mockUseCase)

	userID := uuid.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: userID.String()}}

	invalidJSON := []byte(`{"firstname": "Charlie", "lastname": "Brown",`) // Malformed JSON
	c.Request, _ = http.NewRequest("PATCH", "/user/"+userID.String(), bytes.NewBuffer(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	userController.UpdateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUseCase.AssertNotCalled(t, "UpdateUser", mock.Anything)
}

func TestUpdateUser_ValidationError(t *testing.T) {
	mockUseCase := new(mocks.UserUseCaseMock)
	userController := NewUserController(mockUseCase)

	userID := uuid.New()
	user := &entity.User{
		Firstname: "", // Missing firstname
		Lastname:  "Brown",
		Email:     "invalid-email",
		Age:       200, // Invalid age
	}

	validationError := errors.New("firstname is required")
	mockUseCase.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(validationError)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: userID.String()}}

	jsonValue, _ := json.Marshal(user)
	c.Request, _ = http.NewRequest("PATCH", "/user/"+userID.String(), bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	userController.UpdateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	mockUseCase := new(mocks.UserUseCaseMock)
	userController := NewUserController(mockUseCase)

	userID := uuid.New()

	mockUseCase.On("DeleteUser", userID).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: userID.String()}}
	c.Request, _ = http.NewRequest("DELETE", "/user/"+userID.String(), nil)

	userController.DeleteUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestDeleteUser_InvalidUUID(t *testing.T) {
	mockUseCase := new(mocks.UserUseCaseMock)
	userController := NewUserController(mockUseCase)

	invalidID := "invalid-uuid"

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: invalidID}}
	c.Request, _ = http.NewRequest("DELETE", "/user/"+invalidID, nil)

	userController.DeleteUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUseCase.AssertNotCalled(t, "DeleteUser", mock.Anything)
}

func TestDeleteUser_NotFound(t *testing.T) {
	mockUseCase := new(mocks.UserUseCaseMock)
	userController := NewUserController(mockUseCase)

	userID := uuid.New()

	mockUseCase.On("DeleteUser", userID).Return(errors.New("user not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: userID.String()}}
	c.Request, _ = http.NewRequest("DELETE", "/user/"+userID.String(), nil)

	userController.DeleteUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUseCase.AssertExpectations(t)
}
