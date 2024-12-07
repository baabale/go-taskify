package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"taskify/middleware"
	"net/http"
	"time"
)

var startTime = time.Now()

// RegisterRoutes registers all application routes
func RegisterRoutes(r *gin.Engine, enforcer *casbin.Enforcer) {
	// Health check route
	r.GET("/health", healthCheck)

	// Public routes
	RegisterAuthRoutes(r)

	// Protected API routes
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	api.Use(middleware.PermissionMiddleware(enforcer))

	// Register protected routes under /api/v1
	RegisterTaskRoutes(api)
}

// Health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"uptime": time.Since(startTime).String(),
	})
}
