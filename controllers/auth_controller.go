package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"taskify/config"
	"taskify/errors"
	"taskify/middleware"
	"taskify/models"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
	Role     string `json:"role" binding:"required,oneof=admin editor viewer" example:"editor"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// @Summary Register a new user
// @Description Register a new user with the provided credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration details"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewInvalidInput(err.Error()))
		return
	}

	// Create new user instance
	user := models.NewUser(req.Username, req.Password, req.Role)
	if err := user.HashPassword(); err != nil {
		c.Error(errors.NewInternalError(err))
		return
	}

	// Check if username already exists
	collection := config.DB.Collection("users")
	ctx := context.Background()

	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err != mongo.ErrNoDocuments {
		if err == nil {
			c.Error(errors.NewInvalidInput("Username already exists"))
			return
		}
		c.Error(errors.NewDatabaseError(err))
		return
	}

	// Insert new user
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		c.Error(errors.NewDatabaseError(err))
		return
	}

	// Get the inserted user
	var createdUser models.User
	err = collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdUser)
	if err != nil {
		c.Error(errors.NewDatabaseError(err))
		return
	}

	userResponse := models.UserResponse{
		ID:        createdUser.ID.Hex(),
		Username:  createdUser.Username,
		Role:      createdUser.Role,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": userResponse,
	})
}

// Login handles user login
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewInvalidInput(err.Error()))
		return
	}

	fmt.Println("Login request received:")

	collection := config.DB.Collection("users")
	ctx := context.Background()

	// Find user by username
	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.Error(errors.NewInvalidInput("Invalid username or password"))
		return
	}
	if err != nil {
		c.Error(errors.NewInternalError(err))
		return
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		c.Error(errors.NewInvalidInput("Invalid username or password"))
		return
	}

	// Generate token
	token, err := middleware.GenerateToken(user.Username, user.Role)
	if err != nil {
		c.Error(errors.NewInternalError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
