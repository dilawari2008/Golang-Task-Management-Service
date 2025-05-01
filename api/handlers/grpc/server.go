package grpc

import (
	"context"

	"task-management-system/api/proto"
	"task-management-system/services"
	"google.golang.org/grpc"
)

type TaskServer struct {
	proto.UnimplementedTaskServiceServer
	service services.TaskService
}

func NewTaskServer(service services.TaskService) *TaskServer {
	return &TaskServer{
		service: service,
	}
}

func RegisterServer(s *grpc.Server, service services.TaskService) {
	proto.RegisterTaskServiceServer(s, NewTaskServer(service))
}

func (s *TaskServer) GetAllTasks(ctx context.Context, req *proto.GetAllTasksRequest) (*proto.GetAllTasksResponse, error) {
	page := int(req.Page)
	if page < 1 {
		page = 1
	}

	limit := int(req.Limit)
	if limit < 1 {
		limit = 10
	}

	filters := make(map[string]string)
	for k, v := range req.Filters {
		filters[k] = v
	}

	tasks, pagination, err := s.service.GetTasks(page, limit, filters)
	if err != nil {
		return &proto.GetAllTasksResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	protoTasks := make([]*proto.Task, len(tasks))
	for i, task := range tasks {
		protoTasks[i] = &proto.Task{
			Id:          int64(task.ID),
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt.Unix(),
			UpdatedAt:   task.UpdatedAt.Unix(),
		}
	}

	return &proto.GetAllTasksResponse{
		Success: true,
		Tasks:   protoTasks,
		Pagination: &proto.Pagination{
			CurrentPage: int32(pagination.Page),
			TotalPages:  int32(pagination.TotalPages),
			PageSize:    int32(pagination.Limit),
			TotalCount:  int32(pagination.TotalCount),
		},
	}, nil
}