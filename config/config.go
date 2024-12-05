package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"taskify/utils"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI      string `validate:"required,url"`
	DatabaseName  string `validate:"required,min=1"`
	ServerPort    string `validate:"required,numeric,min=1,max=65535"`
	ServerAddress string `validate:"required,hostname_port|hostname"`
	Environment   string `validate:"required,oneof=development production test"`
}

var AppConfig Config

// LoadConfig loads configuration from environment variables
func LoadConfig() error {
	// Determine environment
	env := getEnv("GO_ENV", "development")

	// Try to load environment-specific file
	envFile := ".env." + env
	if err := godotenv.Load(envFile); err != nil {
		// If environment-specific file doesn't exist, try default .env
		if err := godotenv.Load(); err != nil {
			log.Printf("No %s or .env file found, using environment variables", envFile)
		}
	} else {
		log.Printf("Loaded configuration from %s", envFile)
	}

	// Set configuration values
	AppConfig = Config{
		Environment:   env,
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName:  getEnv("DB_NAME", "taskify"),
		ServerPort:    getEnv("SERVER_PORT", "3000"),
		ServerAddress: getEnv("SERVER_ADDRESS", "localhost"),
	}

	// Validate configuration
	if err := utils.ValidateStruct(AppConfig); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	log.Printf("Configuration loaded for environment: %s", AppConfig.Environment)
	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// ValidateEnvironment validates if the environment is supported
func ValidateEnvironment(env string) bool {
	validEnvs := []string{"development", "production", "test"}
	env = strings.ToLower(env)
	for _, validEnv := range validEnvs {
		if env == validEnv {
			return true
		}
	}
	return false
}
