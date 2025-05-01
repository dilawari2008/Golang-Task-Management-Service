package models

import (
	"errors"
	"time"
)

type Task struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Status    string     `gorm:"not null;default:'pending'" json:"status"`
	Topic     string     `gorm:"not null" json:"topic"`
	Data      string     `json:"data"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

type TaskStatus string
type TaskTopic string

const (
	StatusPending   TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
	StatusExpired   TaskStatus = "expired"
)

const (
	TopicConsumer1 TaskTopic = "consumer1"
	TopicConsumer2 TaskTopic = "consumer2"
	TopicConsumer3 TaskTopic = "consumer3"
)

func NewTask(topic string, data string) *Task {
	return &Task{
		Status:      string(StatusPending),
		Topic:       string(topic),
		Data:        data,
	}
}

func IsValidStatus(status string) bool {
	return status == string(StatusPending) ||
		status == string(StatusInProgress) ||
		status == string(StatusCompleted) ||
		status == string(StatusFailed) ||
		status == string(StatusExpired)
}

var ErrInvalidStatus = errors.New("invalid status value")
