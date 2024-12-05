package routes

import (
	"github.com/gin-gonic/gin"
	"taskify/controllers"
)

// RegisterTaskRoutes registers all task related routes
func RegisterTaskRoutes(rg *gin.RouterGroup) {
	tasks := rg.Group("/tasks")
	{
		tasks.GET("", controllers.GetTasks)
		tasks.POST("", controllers.CreateTask)
		tasks.GET("/:id", controllers.GetTask)
		tasks.PUT("/:id", controllers.UpdateTask)
		tasks.DELETE("/:id", controllers.DeleteTask)
	}
}
