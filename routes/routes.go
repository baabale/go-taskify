package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"taskify/middleware"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(r *gin.Engine, enforcer *casbin.Enforcer) {
	// Public routes
	RegisterAuthRoutes(r)

	// Protected API routes
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	api.Use(middleware.PermissionMiddleware(enforcer))

	// Register protected routes under /api/v1
	RegisterTaskRoutes(api)
}
