package controller

import (
	"errors"
	"net/http"
	"userapi/internal/entity"
	appErrors "userapi/internal/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserController handles HTTP requests related to users
type UserController struct {
	UserUseCase UserUseCase // Use the interface defined in this package
}

// NewUserController creates a new UserController instance
func NewUserController(uc UserUseCase) *UserController {
	return &UserController{
		UserUseCase: uc,
	}
}

// CreateUser handles the creation of a new user
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var user entity.User

	// Bind JSON input to the user entity
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Call the use case to create the user
	if err := ctrl.UserUseCase.CreateUser(&user); err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Respond with the created user
	c.JSON(http.StatusCreated, user)
}

// GetUser handles fetching a user by ID
func (ctrl *UserController) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	user, err := ctrl.UserUseCase.GetUser(userID)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser handles updating an existing user
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user.ID = userID

	if err := ctrl.UserUseCase.UpdateUser(&user); err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if appErr, ok := err.(*appErrors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles deleting a user by ID
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	if err := ctrl.UserUseCase.DeleteUser(userID); err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
