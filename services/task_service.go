package services

import (
	"errors"
	"time"

	"task-management-system/models"
	"task-management-system/repository"
)

var (
	ErrInvalidTask       = errors.New("invalid task data")
	ErrInvalidPagination = errors.New("invalid pagination parameters")
)

// TaskService defines the interface for task business logic
type TaskService interface {
	CreateTask(title, description string, dueDate time.Time) (*models.Task, error)
	GetTaskByID(id uint) (*models.Task, error)
	GetTasks(page, limit int, filters map[string]string) ([]*models.Task, *models.Pagination, error)
	UpdateTask(id uint, title, description, status string, dueDate time.Time) (*models.Task, error)
	DeleteTask(id uint) error
}

// TaskServiceImpl implements TaskService
type TaskServiceImpl struct {
	repo repository.TaskRepository
}

// NewTaskService creates a new task service
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &TaskServiceImpl{
		repo: repo,
	}
}

// CreateTask creates a new task
func (s *TaskServiceImpl) CreateTask(title, description string, dueDate time.Time) (*models.Task, error) {
	if title == "" {
		return nil, ErrInvalidTask
	}

	task := models.NewTask(title, description, dueDate)
	err := s.repo.Create(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetTaskByID retrieves a task by its ID
func (s *TaskServiceImpl) GetTaskByID(id uint) (*models.Task, error) {
	return s.repo.GetByID(id)
}

// GetTasks retrieves tasks with pagination and filtering
func (s *TaskServiceImpl) GetTasks(page, limit int, filters map[string]string) ([]*models.Task, *models.Pagination, error) {
	if page < 1 || limit < 1 {
		return nil, nil, ErrInvalidPagination
	}

	// Validate status filter if present
	if status, ok := filters["status"]; ok && status != "" {
		if !models.IsValidStatus(status) {
			return nil, nil, models.ErrInvalidStatus
		}
	}

	tasks, totalCount, err := s.repo.GetAll(page, limit, filters)
	if err != nil {
		return nil, nil, err
	}

	pagination := models.NewPagination(totalCount, page, limit)
	return tasks, pagination, nil
}

// UpdateTask updates an existing task
func (s *TaskServiceImpl) UpdateTask(id uint, title, description, status string, dueDate time.Time) (*models.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		task.Title = title
	}

	if description != "" {
		task.Description = description
	}

	if status != "" {
		if !models.IsValidStatus(status) {
			return nil, models.ErrInvalidStatus
		}
		task.Status = status
	}

	if !dueDate.IsZero() {
		task.DueDate = dueDate
	}

	err = s.repo.Update(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteTask deletes a task
func (s *TaskServiceImpl) DeleteTask(id uint) error {
	return s.repo.Delete(id)
}