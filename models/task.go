package models

import (
	"errors"
	"time"

	"task-management-system/constants"
)

type Task struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Status       string     `gorm:"not null;default:'pending'" json:"status"`
	Topic        string     `gorm:"not null" json:"topic"`
	Data         string     `json:"data"`
	ErrorMessage string     `json:"error_message,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at"`
}

type TaskStatus string
type TaskTopic string

func NewTask(topic string, data string) *Task {
	return &Task{
		Status: string(constants.StatusPending),
		Topic:  string(topic),
		Data:   data,
	}
}

func IsValidStatus(status string) bool {
	return status == string(constants.StatusPending) ||
		status == string(constants.StatusInProgress) ||
		status == string(constants.StatusCompleted) ||
		status == string(constants.StatusFailed) ||
		status == string(constants.StatusExpired)
}

var ErrInvalidStatus = errors.New("invalid status value")
