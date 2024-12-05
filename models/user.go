package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string            `json:"username" bson:"username" binding:"required"`
	Password  string            `json:"-" bson:"password" binding:"required"`  // "-" means this field won't be included in JSON
	Role      string            `json:"role" bson:"role" binding:"required,oneof=admin editor viewer"`
	CreatedAt time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" bson:"updated_at"`
}

// NewUser creates a new user with default values
func NewUser(username, password, role string) *User {
	now := time.Now()
	return &User{
		Username:  username,
		Password:  password,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// HashPassword hashes the user's password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies the provided password against the hashed password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// swagger:model User
type UserResponse struct {
	ID        string    `json:"id" example:"5f7b5e1b9b0b3a1b3c9b4b1a"`
	Username  string    `json:"username" example:"johndoe"`
	Role      string    `json:"role" example:"editor" enum:"admin,editor,viewer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
