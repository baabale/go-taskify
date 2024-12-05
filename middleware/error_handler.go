package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"taskify/errors"
)

// ErrorHandler is a middleware that handles errors globally
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			appErr := errors.AsAppError(err)

			// Log error details (in production, you might want to use a proper logger)
			log.Printf("Error: %v", appErr.Err)

			// Send error response
			c.JSON(appErr.StatusCode, gin.H{
				"error": appErr.Message,
			})

			// Stop processing
			c.Abort()
		}
	}
}
