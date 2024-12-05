package middleware

import (
	"net/http"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func PermissionMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Get the request path and method
		path := filepath.Clean(c.Request.URL.Path)
		method := c.Request.Method

		// Check if the user has permission
		allowed, err := e.Enforce(role, path, method)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking permissions"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
