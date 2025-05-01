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

type Config struct {
	Server     ServerConfig
	GrpcServer GrpcServerConfig
	Database   DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type GrpcServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL string
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("GRPC_PORT", "9090")
	viper.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/taskdb?sslmode=disable")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = viper.GetString("DATABASE_URL")
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = viper.GetString("GRPC_PORT")
	}

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
		GrpcServer: GrpcServerConfig{
			Port: grpcPort,
		},
		Database: DatabaseConfig{
			URL: dbURL,
		},
	}

	return config, nil
}

func SetupDatabase(cfg *DatabaseConfig) (*gorm.DB, error) {
	fmt.Println("Connecting to database:", cfg.URL)
	db, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Database connection established successfully")

	fmt.Println("Running database migrations...")
	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Database migrations completed successfully")

	return db, nil
}