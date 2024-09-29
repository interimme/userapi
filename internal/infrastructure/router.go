package infrastructure

import (
	"userapi/internal/controller"

	"github.com/gin-gonic/gin"
)

// NewRouter initializes the Gin router with routes and handlers
func NewRouter(userController *controller.UserController) *gin.Engine {
	router := gin.Default()

	// Define the routes and handlers
	router.POST("/users", userController.CreateUser)
	router.GET("/user/:id", userController.GetUser)
	router.PATCH("/user/:id", userController.UpdateUser)
	router.DELETE("/user/:id", userController.DeleteUser)

	return router
}
