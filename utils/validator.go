package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Custom validation tags
const (
	StatusEnum = "status_enum"
)

// Use a single instance of Validator, it caches struct info
var validate *validator.Validate

// InitValidator initializes the validator with custom validations
func InitValidator() {
	validate = validator.New()

	// Register custom validation tags
	validate.RegisterValidation(StatusEnum, validateTaskStatus)

	// Get validator from Gin's validator engine
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom tag name function
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// Register custom validations
		v.RegisterValidation(StatusEnum, validateTaskStatus)
	}
}

// ValidateStruct validates a struct using tags
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, formatValidationError(err))
		}
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}
	return nil
}

// Custom validation function for task status
func validateTaskStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	validStatuses := []string{"pending", "in_progress", "completed"}
	for _, s := range validStatuses {
		if status == s {
			return true
		}
	}
	return false
}

// formatValidationError formats validation errors in a user-friendly way
func formatValidationError(err validator.FieldError) string {
	field := err.Field()
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, err.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case StatusEnum:
		return fmt.Sprintf("%s must be one of: pending, in_progress, completed", field)
	default:
		return fmt.Sprintf("%s failed %s validation", field, err.Tag())
	}
}
