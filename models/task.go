package models

import (
	"time"
	"taskify/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task represents a task in the system
type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string            `json:"title" bson:"title" binding:"required,min=3,max=100"`
	Description string            `json:"description,omitempty" bson:"description" binding:"omitempty,max=500"`
	Status      string            `json:"status,omitempty" bson:"status" binding:"omitempty,oneof=pending in_progress completed"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
}

// NewTask creates a new task with default values
func NewTask(title string) *Task {
	now := time.Now()
	return &Task{
		Title:     title,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update updates task fields with non-empty values
func (t *Task) Update(title, description, status string) error {
	if title != "" {
		t.Title = title
	}
	if description != "" {
		t.Description = description
	}
	if status != "" {
		t.Status = status
	}
	t.UpdatedAt = time.Now()
	return utils.ValidateStruct(t)
}

// swagger:model Task
type TaskResponse struct {
	ID          string    `json:"id" example:"5f7b5e1b9b0b3a1b3c9b4b1a"`
	Title       string    `json:"title" example:"Complete project documentation" minLength:"3" maxLength:"100"`
	Description string    `json:"description" example:"Write comprehensive documentation for the Taskify project" maxLength:"500"`
	Status      string    `json:"status" example:"pending" enum:"pending,in_progress,completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
