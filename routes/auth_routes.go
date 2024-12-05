package routes

import (
	"taskify/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers all authentication related routes
// These are public endpoints that don't require authentication
func RegisterAuthRoutes(r gin.IRouter) {
	// Public authentication routes
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", controllers.Register) // Public endpoint for user registration
		auth.POST("/login", controllers.Login)       // Public endpoint for user login
	}
}
