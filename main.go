package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"task-management-system/api"
	"task-management-system/api/handlers"
	"task-management-system/config"
	"task-management-system/services"
	"task-management-system/repository"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup database
	db, err := config.SetupDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create repository
	taskRepo := repository.NewGormTaskRepository(db)
	
	// Create service
	taskService := services.NewTaskService(taskRepo)
	
	// Create handler
	taskHandler := handlers.NewTaskHandler(taskService)
	
	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	api.SetupRoutes(router, taskHandler)
	
	// Configure server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}
	
	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s\n", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	
	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	
	log.Println("Server stopped")
}
