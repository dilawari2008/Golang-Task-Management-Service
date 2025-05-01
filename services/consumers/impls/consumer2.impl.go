package impls

import (
	"errors"
	"log"
	"task-management-system/constants"
	"task-management-system/models"
	"task-management-system/services/consumers"
)

type Consumer2 struct {
}

func NewConsumer2() consumers.Consumer {
	return &Consumer2{}
}

func (c *Consumer2) Process(task *models.Task) (bool, error) {
	if task == nil {
		log.Println("Consumer2: Received nil task")
		return false, errors.New("Consumer2: Received nil task")
	}

	if task.Topic != string(constants.TopicConsumer2) {
		log.Printf("Consumer2: Received task with incorrect topic: %s", task.Topic)
		return false, errors.New("Consumer2: Received task with incorrect topic")
	}

	log.Printf("Consumer2: Processing task ID: %d with data: %s", task.ID, task.Data)

	// push data into message queue here

	log.Printf("Consumer2: Successfully pushed to queue: task ID: %d", task.ID)
	return false, nil
}
