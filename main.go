package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"taskify/config"
	_ "taskify/docs" // Import swagger docs
	"taskify/middleware"
	"taskify/routes"
	"taskify/utils"
)

// @title           Taskify API
// @version         1.0
// @description     A Task Management API with authentication and authorization
// @host           localhost:3000
// @BasePath       /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Initialize validator
	utils.InitValidator()

	// Initialize configuration
	config.LoadConfig()
	config.ConnectDatabase()

	// Initialize Gin
	r := gin.Default()

	// Global middleware
	r.Use(middleware.ErrorHandler()) // Register error handler first

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize Casbin enforcer
	enforcer, err := casbin.NewEnforcer("config/model.conf", "config/policy.csv")
	if err != nil {
		log.Fatal("Failed to initialize Casbin enforcer:", err)
	}

	// Register routes
	routes.RegisterRoutes(r, enforcer)

	// Start server
	serverAddr := fmt.Sprintf("%s:%s", config.AppConfig.ServerAddress, config.AppConfig.ServerPort)
	r.Run(serverAddr)
}
