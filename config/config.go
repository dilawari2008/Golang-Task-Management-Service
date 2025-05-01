package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"task-management-system/models"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port string
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	URL string
}

// LoadConfig loads configuration from environment variables or config file
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	// Default values
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/taskdb?sslmode=disable")

	// Try to read from config file
	if err := viper.ReadInConfig(); err != nil {
		// Just use environment variables if config file is not found
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Prioritize DATABASE_URL from environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = viper.GetString("DATABASE_URL")
	}

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			URL: dbURL,
		},
	}

	return config, nil
}

// SetupDatabase initializes the database connection
func SetupDatabase(cfg *DatabaseConfig) (*gorm.DB, error) {
	fmt.Println("Connecting to database:", cfg.URL)
	db, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Database connection established successfully")

	// Run migrations
	fmt.Println("Running database migrations...")
	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Database migrations completed successfully")

	return db, nil
}