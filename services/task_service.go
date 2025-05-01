package services

import (
	"errors"

	"task-management-system/models"
	"task-management-system/repository"
)

var (
	ErrInvalidTask       = errors.New("invalid task data")
	ErrInvalidPagination = errors.New("invalid pagination parameters")
)

type TaskService interface {
	CreateTask(topic, data string) (*models.Task, error)
	GetTaskByID(id uint) (*models.Task, error)
	GetTasks(page, limit int, filters map[string]string) ([]*models.Task, *models.Pagination, error)
	UpdateTask(id uint, topic, data, status string) (*models.Task, error)
	DeleteTask(id uint) error
}

type TaskServiceImpl struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &TaskServiceImpl{
		repo: repo,
	}
}

func (s *TaskServiceImpl) CreateTask(topic, data string) (*models.Task, error) {
	if topic == "" {
		return nil, ErrInvalidTask
	}

	task := models.NewTask(topic, data)
	err := s.repo.Create(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskServiceImpl) GetTaskByID(id uint) (*models.Task, error) {
	return s.repo.GetByID(id)
}

func (s *TaskServiceImpl) GetTasks(page, limit int, filters map[string]string) ([]*models.Task, *models.Pagination, error) {
	if page < 1 || limit < 1 {
		return nil, nil, ErrInvalidPagination
	}

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

func (s *TaskServiceImpl) UpdateTask(id uint, topic, data, status string) (*models.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if topic != "" {
		task.Topic = topic
	}

	if data != "" {
		task.Data = data
	}

	if status != "" {
		if !models.IsValidStatus(status) {
			return nil, models.ErrInvalidStatus
		}
		task.Status = status
	}

	err = s.repo.Update(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskServiceImpl) DeleteTask(id uint) error {
	return s.repo.Delete(id)
}
