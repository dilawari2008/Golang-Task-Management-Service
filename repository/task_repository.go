package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"task-management-system/models"
)

var (
	// ErrTaskNotFound is returned when a task is not found
	ErrTaskNotFound = errors.New("task not found")
)

// TaskRepository defines the interface for task data access
type TaskRepository interface {
	Create(task *models.Task) error
	GetByID(id uint) (*models.Task, error)
	GetAll(page, limit int, filters map[string]string) ([]*models.Task, int64, error)
	Update(task *models.Task) error
	Delete(id uint) error
}

// GormTaskRepository implements TaskRepository with GORM and PostgreSQL
type GormTaskRepository struct {
	db *gorm.DB
}

// NewGormTaskRepository creates a new GORM-based task repository
func NewGormTaskRepository(db *gorm.DB) *GormTaskRepository {
	return &GormTaskRepository{
		db: db,
	}
}

// Create adds a new task to the repository
func (r *GormTaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

// GetByID retrieves a task by its ID
func (r *GormTaskRepository) GetByID(id uint) (*models.Task, error) {
	var task models.Task
	result := r.db.Where("deleted_at IS NULL").First(&task, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, result.Error
	}
	return &task, nil
}

// GetAll retrieves all tasks with pagination and filtering
func (r *GormTaskRepository) GetAll(page, limit int, filters map[string]string) ([]*models.Task, int64, error) {
	var tasks []*models.Task
	var totalCount int64

	offset := (page - 1) * limit

	query := r.db.Model(&models.Task{}).Where("deleted_at IS NULL")

	// Apply filters
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, totalCount, nil
}

// Update updates an existing task
func (r *GormTaskRepository) Update(task *models.Task) error {
	result := r.db.Where("deleted_at IS NULL").Save(task)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// Delete marks a task as deleted
func (r *GormTaskRepository) Delete(id uint) error {
	result := r.db.Model(&models.Task{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}
