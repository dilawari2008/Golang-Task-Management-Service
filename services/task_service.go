package services

import (
	"errors"
	"fmt"

	"task-management-system/constants"
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
	repo      repository.TaskRepository
	taskQueue *TaskQueue
}

func NewTaskService(repo repository.TaskRepository, taskQueue *TaskQueue) TaskService {
	return &TaskServiceImpl{
		repo:      repo,
		taskQueue: taskQueue,
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

	s.taskQueue.Push(task)

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

func (s *TaskServiceImpl) UpdateTask(id uint, topic string, data string, status string) (*models.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if task.Status != string(constants.StatusPending) && (topic != "" || data != "") {
		return nil, errors.New("can only update topic and data when task status is pending")
	}

	if topic != "" && task.Status == string(constants.StatusPending) {
		task.Topic = topic
	}

	if data != "" && task.Status == string(constants.StatusPending) {
		task.Data = data
	}

	if status != "" {
		if !models.IsValidStatus(status) {
			return nil, models.ErrInvalidStatus
		}

		if !isValidStatusTransition(task.Status, status) {
			return nil, fmt.Errorf("invalid status transition: cannot change from %s to %s", task.Status, status)
		}

		task.Status = status
	}

	err = s.repo.Update(task)
	if err != nil {
		return nil, err
	}

	s.taskQueue.Push(task)

	return task, nil
}

func isValidStatusTransition(currentStatus, newStatus string) bool {
	hierarchyLevel := map[string]int{
		string(constants.StatusPending):    1,
		string(constants.StatusInProgress): 2,
		string(constants.StatusCompleted):  3,
		string(constants.StatusFailed):     3,
		string(constants.StatusExpired):    3,
	}

	currentLevel, _ := hierarchyLevel[currentStatus]
	newLevel, _ := hierarchyLevel[newStatus]

	return newLevel >= currentLevel
}

func (s *TaskServiceImpl) DeleteTask(id uint) error {
	return s.repo.Delete(id)
}
