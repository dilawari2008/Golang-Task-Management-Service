package grpc

import (
	"fmt"
	"log"
	"net"

	"task-management-system/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartServer(taskService services.TaskService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	
	grpcServer := grpc.NewServer()
	
	RegisterServer(grpcServer, taskService)
	
	reflection.Register(grpcServer)
	
	log.Printf("gRPC server listening on port %d", port)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	
	return nil
}