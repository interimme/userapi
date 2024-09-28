package controller

import (
	"userapi/internal/entity"
	appErrors "userapi/internal/errors"
	"userapi/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserController handles HTTP requests related to users
type UserController struct {
	UserUseCase usecase.UserUseCase
}

// NewUserController creates a new UserController instance
func NewUserController(uc usecase.UserUseCase) *UserController {
	return &UserController{
		UserUseCase: uc,
	}
}

// CreateUser handles the creation of a new user
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var user entity.User

	// Bind JSON input to the user entity
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(&appErrors.AppError{Code: 400, Message: "invalid request"})
		return
	}

	// Call the use case to create the user
	if err := ctrl.UserUseCase.CreateUser(&user); err != nil {
		c.Error(err)
		return
	}

	// Respond with the created user
	c.JSON(201, user)
}

// GetUser handles fetching a user by ID
func (ctrl *UserController) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.Error(&appErrors.AppError{Code: 400, Message: "invalid uuid"})
		return
	}

	user, err := ctrl.UserUseCase.GetUser(userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, user)
}

// UpdateUser handles updating an existing user
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.Error(&appErrors.AppError{Code: 400, Message: "invalid uuid"})
		return
	}

	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(&appErrors.AppError{Code: 400, Message: "invalid request"})
		return
	}

	user.ID = userID

	if err := ctrl.UserUseCase.UpdateUser(&user); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, user)
}

// DeleteUser handles deleting a user by ID
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.Error(&appErrors.AppError{Code: 400, Message: "invalid uuid"})
		return
	}

	if err := ctrl.UserUseCase.DeleteUser(userID); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}
