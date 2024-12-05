package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"taskify/config"
	_ "taskify/docs" // This is required for swagger
	"taskify/middleware"
	"taskify/routes"
	"taskify/utils"
)

// @title           Taskify API
// @version         1.0
// @description     A task management RESTful API implementation with MongoDB
// @host           localhost:3000
// @BasePath       /api/v1

func main() {
	// Initialize validator
	utils.InitValidator()

	// Load configuration from environment variables
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Configuration error:", err)
	}

	// Initialize MongoDB connection
	config.ConnectDatabase()

	// Set Gin mode based on environment
	if config.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Add middleware
	r.Use(middleware.ErrorHandler())

	// Register routes
	routes.RegisterRoutes(r)

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server with configured host and port
	serverAddr := fmt.Sprintf("%s:%s", config.AppConfig.ServerAddress, config.AppConfig.ServerPort)
	log.Printf("Server starting on %s in %s mode", serverAddr, config.AppConfig.Environment)
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
