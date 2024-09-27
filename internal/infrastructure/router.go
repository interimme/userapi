package infrastructure

import (
	"userapi/internal/controller"
	"userapi/internal/infrastructure/middleware"
	"userapi/internal/infrastructure/persistence"
	"userapi/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewRouter initializes the Gin router with routes and handlers
func NewRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Apply the error handling middleware
	router.Use(middleware.ErrorHandler)

	// Initialize the repository
	userRepo := persistence.NewUserRepository(db)

	// Initialize the use case
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Initialize the controller
	userController := controller.NewUserController(userUseCase)

	// Define the routes and handlers
	router.POST("/users", userController.CreateUser)
	router.GET("/user/:id", userController.GetUser)
	router.PATCH("/user/:id", userController.UpdateUser)
	router.DELETE("/user/:id", userController.DeleteUser)

	return router
}
