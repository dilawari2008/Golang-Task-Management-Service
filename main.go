package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"task-management-system/api"
	"task-management-system/api/handlers/grpc"
	"task-management-system/api/handlers/rest"
	"task-management-system/config"
	"task-management-system/repository"
	"task-management-system/services"
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
	
	// Create REST handler
	taskHandler := handlers.NewTaskHandler(taskService)
	
	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	api.SetupRoutes(router, taskHandler)
	
	// Configure HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}
	
	// Start REST server in a goroutine
	go func() {
		log.Printf("Starting REST server on port %s\n", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("REST server error: %v", err)
		}
	}()
	
	// Start gRPC server in a goroutine
	go func() {
		grpcPort, err := strconv.Atoi(cfg.GrpcServer.Port)
		if err != nil {
			log.Fatalf("Invalid gRPC port: %v", err)
		}
		
		if err := grpc.StartServer(taskService, grpcPort); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()
	
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down servers...")
	
	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	
	log.Println("Servers stopped")
}