package grpc

import (
	"fmt"
	"log"
	"net"

	"task-management-system/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// StartServer initializes and starts the gRPC server
func StartServer(taskService services.TaskService, port int) error {
	// Create a listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	
	// Create a gRPC server
	grpcServer := grpc.NewServer()
	
	// Register our service
	RegisterServer(grpcServer, taskService)
	
	// Register reflection service for grpcurl and other tools
	reflection.Register(grpcServer)
	
	// Start serving
	log.Printf("gRPC server listening on port %d", port)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	
	return nil
}