package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"userapi/internal/entity"
	appErrors "userapi/internal/errors"
	"userapi/internal/usecase"
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
		c.Error(&appErrors.AppError{Code: 400, Message: err.Error()})
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

func (ctrl *UserController) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.Error(&appErrors.AppError{Code: 400, Message: "Invalid UUID"})
		return
	}

	user, err := ctrl.UserUseCase.GetUser(userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, user)
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.Error(&appErrors.AppError{Code: 400, Message: "Invalid UUID"})
		return
	}

	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(&appErrors.AppError{Code: 400, Message: err.Error()})
		return
	}
	user.ID = userID

	if err := ctrl.UserUseCase.UpdateUser(&user); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, user)
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.Error(&appErrors.AppError{Code: 400, Message: "Invalid UUID"})
		return
	}

	if err := ctrl.UserUseCase.DeleteUser(userID); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}
