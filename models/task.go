package models

import (
	"errors"
	"time"
)

// Task represents a task in the system with GORM tags
type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"not null" json:"title" binding:"required"`
	Description string     `json:"description"`
	Status      string     `gorm:"not null;default:'pending'" json:"status"`
	DueDate     time.Time  `json:"due_date"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
}

// TaskStatus defines the possible statuses for a task
type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
)

// NewTask creates a new task with default values
func NewTask(title, description string, dueDate time.Time) *Task {
	return &Task{
		Title:       title,
		Description: description,
		Status:      string(StatusPending),
		DueDate:     dueDate,
	}
}

// IsValidStatus checks if a status is valid
func IsValidStatus(status string) bool {
	return status == string(StatusPending) ||
		status == string(StatusInProgress) ||
		status == string(StatusCompleted)
}

// ErrInvalidStatus is returned when an invalid status is provided
var ErrInvalidStatus = errors.New("invalid status value")
