package impls

import (
	"errors"
	"fmt"
	"log"
	"task-management-system/constants"
	"task-management-system/models"
	"task-management-system/services/consumers"
)

type Consumer1 struct {
}

func NewConsumer1() consumers.Consumer {
	return &Consumer1{}
}

func (c *Consumer1) Process(task *models.Task) (bool, error) {
	if task == nil {
		log.Println("Consumer1: Received nil task")
		return false, errors.New("Consumer1: Received nil task")
	}

	if task.Topic != string(constants.TopicConsumer1) {
		log.Printf("Consumer1: Received task with incorrect topic: %s", task.Topic)
		return false, errors.New("Consumer1: Received task with incorrect topic")
	}

	log.Printf("Consumer1: Processing task ID: %d with data: %s", task.ID, task.Data)

	// add actual rest api call here
	
	// Simulate API call with 1/3 failure rate
	if task.ID%3 == 0 {
		log.Printf("Consumer1: API call failed for task ID: %d", task.ID)
		return false, fmt.Errorf("API call failed")
	}
	
	log.Printf("Consumer1: API call successful for task ID: %d", task.ID)
	return true, nil
}


