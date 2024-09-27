package controller

import (
	"net/http"
	"userapi/internal/entity"
	"userapi/internal/usecase"
	"userapi/internal/utils"

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
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Call the use case to create the user
	err := ctrl.UserUseCase.CreateUser(&user)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Respond with the created user
	utils.RespondJSON(c, http.StatusCreated, user)
}

// GetUser retrieves a user by ID
func (ctrl *UserController) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid UUID")
		return
	}

	user, err := ctrl.UserUseCase.GetUser(userID)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondJSON(c, http.StatusOK, user)
}

// UpdateUser updates an existing user's information
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid UUID")
		return
	}

	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	user.ID = userID

	err = ctrl.UserUseCase.UpdateUser(&user)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, user)
}

// DeleteUser removes a user by ID
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid UUID")
		return
	}

	err = ctrl.UserUseCase.DeleteUser(userID)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}
