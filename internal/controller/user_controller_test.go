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
	"userapi/internal/infrastructure/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestCreateUser_Success(t *testing.T) {
	// Set up the controller and mocks
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	// Create a test router and include the middleware
	router := gin.New()
	router.Use(middleware.ErrorHandler)
	router.POST("/users", userController.CreateUser)

	user := entity.User{
		Firstname: "Alice",
		Lastname:  "Smith",
		Email:     "alice@example.com",
		Age:       28,
	}

	userJSON, _ := json.Marshal(user)

	// Mock the use case to return no error
	mockUseCase.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(nil)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var response entity.User
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Alice", response.Firstname)
	assert.Equal(t, "Smith", response.Lastname)
	assert.Equal(t, "alice@example.com", response.Email)
	assert.Equal(t, uint(28), response.Age)
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	// Set up the controller and mocks
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	// Create a test router and include the middleware
	router := gin.New()
	router.Use(middleware.ErrorHandler)
	router.POST("/users", userController.CreateUser)

	// Create an invalid JSON payload
	invalidJSON := `{"firstname": "Alice", "lastname": "Smith", "email": "alice@example.com", "age": "twenty"}`
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, 400, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid request", response["error"])
}

func TestCreateUser_ValidationError(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.Use(middleware.ErrorHandler)
	router.POST("/users", userController.CreateUser)

	user := entity.User{
		Firstname: "",
		Lastname:  "Smith",
		Email:     "alice@example.com",
		Age:       28,
	}

	userJSON, _ := json.Marshal(user)

	// Mock the use case to return a validation error
	mockUseCase.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(&appErrors.AppError{Code: 400, Message: "firstname is required"})

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "firstname is required", response["error"])
}

func TestGetUser_Success(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.Use(middleware.ErrorHandler)
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

	req, _ := http.NewRequest("GET", "/user/"+userID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response entity.User
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Alice", response.Firstname)
	assert.Equal(t, "Smith", response.Lastname)
	assert.Equal(t, "alice@example.com", response.Email)
	assert.Equal(t, uint(28), response.Age)
}

func TestGetUser_InvalidUUID(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.Use(middleware.ErrorHandler)
	router.GET("/user/:id", userController.GetUser)

	req, _ := http.NewRequest("GET", "/user/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid uuid", response["error"])
}

func TestGetUser_NotFound(t *testing.T) {
	mockUseCase := new(mocks.UserUseCase)
	userController := NewUserController(mockUseCase)

	router := gin.New()
	router.Use(middleware.ErrorHandler)
	router.GET("/user/:id", userController.GetUser)

	userID := uuid.New()

	// Mock the use case to return a not found error
	mockUseCase.On("GetUser", userID).Return(nil, appErrors.ErrNotFound)

	req, _ := http.NewRequest("GET", "/user/"+userID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "user not found", response["error"])
}

// Continue with the rest of the test functions as previously provided
// (The rest of the test functions are already included in the previous response)
