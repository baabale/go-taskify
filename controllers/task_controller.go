package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"taskify/config"
	"taskify/errors"
	"taskify/models"
)

// GetTasks godoc
// @Summary Get all tasks
// @Description Get a list of all tasks with optional filtering, pagination, and sorting
// @Tags tasks
// @Accept json
// @Produce json
// @Param status query string false "Filter by status (pending/in_progress/completed)"
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Param sort query string false "Sort field (created_at/-created_at)"
// @Success 200 {array} models.TaskResponse
// @Failure 400 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /tasks [get]
func GetTasks(c *gin.Context) {
	collection := config.DB.Collection("tasks")
	ctx := context.Background()

	// Build filter
	filter := bson.M{}
	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}

	// Pagination
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	// get the limit from the query parameter or default to 10
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	skip := (page - 1) * limit

	// Build options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	// Sort
	if sort := c.Query("sort"); sort != "" {
		order := 1
		if sort[0] == '-' {
			order = -1
			sort = sort[1:]
		}
		findOptions.SetSort(bson.D{{Key: sort, Value: order}})
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		_ = c.Error(errors.NewDatabaseError(err))
		return
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		_ = c.Error(errors.NewDatabaseError(err))
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task with the provided information
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.CreateTaskDTO true "Task object"
// @Success 201 {object} models.TaskResponse
// @Failure 400 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /tasks [post]
func CreateTask(c *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required,min=3,max=100"`
		Description string `json:"description,omitempty" binding:"omitempty,max=500"`
		Status      string `json:"status,omitempty" binding:"omitempty,oneof=pending in_progress completed"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		_ = c.Error(errors.NewInvalidInput(err.Error()))
		return
	}

	task := models.NewTask(input.Title)
	if input.Description != "" {
		task.Description = input.Description
	}
	if input.Status != "" {
		task.Status = input.Status
	}

	collection := config.DB.Collection("tasks")
	ctx := context.Background()

	result, err := collection.InsertOne(ctx, task)
	if err != nil {
		_ = c.Error(errors.NewDatabaseError(err))
		return
	}

	task.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, task)
}

// GetTask godoc
// @Summary Get a task by ID
// @Description Get details of a specific task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} models.TaskResponse
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /tasks/{id} [get]
func GetTask(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(errors.NewInvalidInput("Invalid task ID format"))
		return
	}

	collection := config.DB.Collection("tasks")
	ctx := context.Background()

	var task models.Task
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_ = c.Error(errors.NewNotFound("Task"))
			return
		}
		_ = c.Error(errors.NewDatabaseError(err))
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update a task's information
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body models.Task true "Task object"
// @Success 200 {object} models.TaskResponse
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		_ = c.Error(errors.NewInvalidInput("invalid task ID"))
		return
	}

	var input struct {
		Title       string `json:"title,omitempty" binding:"omitempty,min=3,max=100"`
		Description string `json:"description,omitempty" binding:"omitempty,max=500"`
		Status      string `json:"status,omitempty" binding:"omitempty,oneof=pending in_progress completed"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		_ = c.Error(errors.NewInvalidInput(err.Error()))
		return
	}

	collection := config.DB.Collection("tasks")
	ctx := context.Background()

	var task models.Task
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_ = c.Error(errors.NewNotFound("task not found"))
			return
		}
		_ = c.Error(errors.NewDatabaseError(err))
		return
	}

	if err := task.Update(input.Title, input.Description, input.Status); err != nil {
		_ = c.Error(errors.NewInvalidInput(err.Error()))
		return
	}

	_, err = collection.ReplaceOne(ctx, bson.M{"_id": objectID}, task)
	if err != nil {
		_ = c.Error(errors.NewDatabaseError(err))
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 204 "No Content"
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(errors.NewInvalidInput("Invalid task ID format"))
		return
	}

	collection := config.DB.Collection("tasks")
	ctx := context.Background()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		_ = c.Error(errors.NewDatabaseError(err))
		return
	}

	if result.DeletedCount == 0 {
		_ = c.Error(errors.NewNotFound("Task"))
		return
	}

	c.Status(http.StatusNoContent)
}
