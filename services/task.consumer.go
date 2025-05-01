package services

import (
	"log"
	"time"

	"task-management-system/constants"
	"task-management-system/repository"
	"task-management-system/services/consumers"
	"task-management-system/services/consumers/impls"
)

type TaskConsumer struct {
	queue     *TaskQueue
	isRunning bool
	consumers map[string]consumers.Consumer
	repo      repository.TaskRepository
}

func NewTaskConsumer(queue *TaskQueue, repo repository.TaskRepository) *TaskConsumer {
	return &TaskConsumer{
		queue: queue,
		consumers: map[string]consumers.Consumer{
			string(constants.TopicConsumer1): impls.NewConsumer1(),
			string(constants.TopicConsumer2): impls.NewConsumer2(),
		},
		repo: repo,
	}
}

func (c *TaskConsumer) Start() {
	if c.isRunning {
		return
	}

	c.isRunning = true
	go c.consumeTasks()
}

func (c *TaskConsumer) Stop() {
	if !c.isRunning {
		return
	}

	c.isRunning = false
}

func (c *TaskConsumer) consumeTasks() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for c.isRunning {
		<-ticker.C
		c.processNextTask()
	}

	log.Println("Task consumer stopped")
}
func (c *TaskConsumer) processNextTask() {
	// 1. start: pre-process
	task := c.queue.Pop()
	if task == nil {
		return
	}

	log.Printf("Processing task ID: %d, Topic: %s", task.ID, task.Topic)

	consumer, exists := c.consumers[task.Topic]
	if !exists {
		log.Printf("Unknown topic: %s", task.Topic)
		return
	}
	// 1. end: pre-process

	// 2. start: process
	success, err := consumer.Process(task)
	// 2. end: process

	// 3. start: post-process
	if success && err == nil {
		task.Status = string(constants.StatusCompleted)
		log.Printf("Successfully completed task ID: %d, status updated to completed", task.ID)
	} else if !success && err == nil {
		task.Status = string(constants.StatusInProgress)
		log.Printf("Task ID %d needs more processing, status updated to in progress", task.ID)
	} else if !success && err != nil {
		task.Status = string(constants.StatusFailed)
		task.ErrorMessage = err.Error()
		log.Printf("Error processing task ID %d: %v, status updated to failed", task.ID, err)
	}

	if err := c.repo.Update(task); err != nil {
		log.Printf("Failed to update task status in repository: %v", err)
	}
	// 3. end: post-process
}
