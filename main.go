package main

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"github.com/baabale/go-taskify/config"
	"github.com/baabale/go-taskify/controllers"
	"github.com/baabale/go-taskify/middleware"
	"github.com/baabale/go-taskify/models"
)

func main() {
	// Initialize configuration
	db := config.InitDB()

	// Auto migrate the schema
	db.AutoMigrate(&models.Task{}, &models.User{})

	// Initialize Casbin enforcer
	enforcer, err := casbin.NewEnforcer("config/model.conf", "config/policy.csv")
	if err != nil {
		log.Fatal("Failed to initialize Casbin enforcer:", err)
	}

	// Initialize controllers
	taskController := controllers.NewTaskController(db)
	authController := controllers.NewAuthController(db)

	// Create Gin router
	r := gin.Default()

	// Public routes
	r.POST("/register", authController.Register)
	r.POST("/login", authController.Login)

	// Protected routes
	api := r.Group("/")
	api.Use(middleware.AuthMiddleware())
	api.Use(middleware.PermissionMiddleware(enforcer))
	{
		api.GET("/tasks", taskController.GetTasks)
		api.GET("/tasks/:id", taskController.GetTask)
		api.POST("/tasks", taskController.CreateTask)
		api.PUT("/tasks/:id", taskController.UpdateTask)
		api.DELETE("/tasks/:id", taskController.DeleteTask)
	}

	r.Run(":8080")
}
