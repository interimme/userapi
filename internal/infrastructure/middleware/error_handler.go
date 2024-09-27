package middleware

import (
	"userapi/internal/errors"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a middleware that handles errors uniformly
func ErrorHandler(c *gin.Context) {
	c.Next() // Execute the handlers

	// Check if any errors were set during the request
	if len(c.Errors) > 0 {
		// Retrieve the last error
		err := c.Errors.Last().Err

		// Check if it's an AppError
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}

		// For other errors, return a generic 500 error
		c.JSON(500, gin.H{"error": "Internal server error"})
	}
}
