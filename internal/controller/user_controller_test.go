package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"userapi/internal/controller/mocks"
	"userapi/internal/entity"
	appErrors "userapi/internal/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestCreateUser_Success(t *testing.T) {
	// Set up the controller and mocks
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	// Create a test router
	router := gin.New()
	router.POST("/users", userController.CreateUser)

	user := entity.User{
		Firstname: "Alice",
		Lastname:  "Smith",
		Email:     "alice@example.com",
		Age:       28,
	}

	userJSON, err := json.Marshal(user)
	require.NoError(t, err, "Failed to marshal user to JSON")

	// Mock the use case to return no error
	mockUseCase.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(nil)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	require.NoError(t, err, "Failed to create HTTP request")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response entity.User
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response JSON")
	assert.Equal(t, "Alice", response.Firstname)
	assert.Equal(t, "Smith", response.Lastname)
	assert.Equal(t, "alice@example.com", response.Email)
	assert.Equal(t, uint(28), response.Age)
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	// Set up the controller and mocks
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	// Create a test router
	router := gin.New()
	router.POST("/users", userController.CreateUser)

	// Create an invalid JSON payload
	invalidJSON := `{"firstname": "Alice", "lastname": "Smith", "email": "alice@example.com", "age": "twenty"}`
	req, err := http.NewRequest("POST", "/users", strings.NewReader(invalidJSON))
	require.NoError(t, err, "Failed to create HTTP request")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response JSON")
	assert.Equal(t, "invalid request", response["error"])
}

func TestCreateUser_ValidationError(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.POST("/users", userController.CreateUser)

	user := entity.User{
		Firstname: "",
		Lastname:  "Smith",
		Email:     "alice@example.com",
		Age:       28,
	}

	userJSON, err := json.Marshal(user)
	require.NoError(t, err, "Failed to marshal user to JSON")

	// Mock the use case to return a validation error
	mockUseCase.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(&appErrors.AppError{Code: http.StatusBadRequest, Message: "firstname is required"})

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	require.NoError(t, err, "Failed to create HTTP request")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response JSON")
	assert.Equal(t, "firstname is required", response["error"])
}

func TestGetUser_Success(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.GET("/user/:id", userController.GetUser)

	userID := uuid.New()
	user := &entity.User{
		ID:        userID,
		Firstname: "Alice",
		Lastname:  "Smith",
		Email:     "alice@example.com",
		Age:       28,
	}

	// Mock the use case to return the user
	mockUseCase.On("GetUser", userID).Return(user, nil)

	req, err := http.NewRequest("GET", "/user/"+userID.String(), nil)
	require.NoError(t, err, "Failed to create HTTP request")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response entity.User
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response JSON")
	assert.Equal(t, "Alice", response.Firstname)
	assert.Equal(t, "Smith", response.Lastname)
	assert.Equal(t, "alice@example.com", response.Email)
	assert.Equal(t, uint(28), response.Age)
}

func TestGetUser_InvalidUUID(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.GET("/user/:id", userController.GetUser)

	req, err := http.NewRequest("GET", "/user/invalid-uuid", nil)
	require.NoError(t, err, "Failed to create HTTP request")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response JSON")
	assert.Equal(t, "invalid uuid", response["error"])
}

func TestGetUser_NotFound(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.GET("/user/:id", userController.GetUser)

	userID := uuid.New()

	// Mock the use case to return ErrNotFound
	mockUseCase.On("GetUser", userID).Return(nil, appErrors.ErrNotFound)

	req, err := http.NewRequest("GET", "/user/"+userID.String(), nil)
	require.NoError(t, err, "Failed to create HTTP request")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response JSON")
	assert.Equal(t, "user not found", response["error"])
}

func TestUpdateUser_Success(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.PATCH("/user/:id", userController.UpdateUser)

	userID := uuid.New()
	user := entity.User{
		Firstname: "Alice",
		Lastname:  "Johnson",
		Email:     "alice.johnson@example.com",
		Age:       29,
	}

	userJSON, err := json.Marshal(user)
	require.NoError(t, err, "Failed to marshal user to JSON")

	// Mock the use case to return no error
	mockUseCase.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(nil)

	req, err := http.NewRequest("PATCH", "/user/"+userID.String(), bytes.NewBuffer(userJSON))
	require.NoError(t, err, "Failed to create HTTP request")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response entity.User
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response JSON")
	assert.Equal(t, "Alice", response.Firstname)
	assert.Equal(t, "Johnson", response.Lastname)
	assert.Equal(t, "alice.johnson@example.com", response.Email)
	assert.Equal(t, uint(29), response.Age)
}

func TestUpdateUser_InvalidUUID(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.PATCH("/user/:id", userController.UpdateUser)

	user := entity.User{
		Firstname: "Alice",
		Lastname:  "Johnson",
		Email:     "alice.johnson@example.com",
		Age:       29,
	}

	userJSON, err := json.Marshal(user)
	require.NoError(t, err, "Failed to marshal user to JSON")

	req, err := http.NewRequest("PATCH", "/user/invalid-uuid", bytes.NewBuffer(userJSON))
	require.NoError(t, err, "Failed to create HTTP request")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response JSON")
	assert.Equal(t, "invalid uuid", response["error"])
}

func TestUpdateUser_ValidationError(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.PATCH("/user/:id", userController.UpdateUser)

	userID := uuid.New()
	user := entity.User{
		Firstname: "",
		Lastname:  "Johnson",
		Email:     "alice.johnson@example.com",
		Age:       29,
	}

	userJSON, err := json.Marshal(user)
	require.NoError(t, err, "Failed to marshal user to JSON")

	// Mock the use case to return a validation error
	mockUseCase.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(&appErrors.AppError{Code: http.StatusBadRequest, Message: "firstname is required"})

	req, err := http.NewRequest("PATCH", "/user/"+userID.String(), bytes.NewBuffer(userJSON))
	require.NoError(t, err, "Failed to create HTTP request")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response JSON")
	assert.Equal(t, "firstname is required", response["error"])
}

func TestUpdateUser_NotFound(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.PATCH("/user/:id", userController.UpdateUser)

	userID := uuid.New()
	user := entity.User{
		Firstname: "Alice",
		Lastname:  "Johnson",
		Email:     "alice.johnson@example.com",
		Age:       29,
	}

	userJSON, err := json.Marshal(user)
	require.NoError(t, err, "Failed to marshal user to JSON")

	// Mock the use case to return ErrNotFound
	mockUseCase.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(appErrors.ErrNotFound)

	req, err := http.NewRequest("PATCH", "/user/"+userID.String(), bytes.NewBuffer(userJSON))
	require.NoError(t, err, "Failed to create HTTP request")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response JSON")
	assert.Equal(t, "user not found", response["error"])
}

func TestDeleteUser_Success(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.DELETE("/user/:id", userController.DeleteUser)

	userID := uuid.New()

	// Mock the use case to return no error
	mockUseCase.On("DeleteUser", userID).Return(nil)

	req, err := http.NewRequest("DELETE", "/user/"+userID.String(), nil)
	require.NoError(t, err, "Failed to create HTTP request")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response JSON")
	assert.Equal(t, "User deleted successfully", response["message"])
}

func TestDeleteUser_InvalidUUID(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.DELETE("/user/:id", userController.DeleteUser)

	req, err := http.NewRequest("DELETE", "/user/invalid-uuid", nil)
	require.NoError(t, err, "Failed to create HTTP request")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response JSON")
	assert.Equal(t, "invalid uuid", response["error"])
}

func TestDeleteUser_NotFound(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.DELETE("/user/:id", userController.DeleteUser)

	userID := uuid.New()

	// Mock the use case to return ErrNotFound
	mockUseCase.On("DeleteUser", userID).Return(appErrors.ErrNotFound)

	req, err := http.NewRequest("DELETE", "/user/"+userID.String(), nil)
	require.NoError(t, err, "Failed to create HTTP request")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response JSON")
	assert.Equal(t, "user not found", response["error"])
}
