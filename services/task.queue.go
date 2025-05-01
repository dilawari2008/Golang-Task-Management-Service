package services

import (
	"log"
	"sync"
	"task-management-system/models"
)

type TaskQueue struct {
	tasks []*models.Task
	mutex sync.Mutex
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		tasks: make([]*models.Task, 0),
	}
}

func (q *TaskQueue) Push(task *models.Task) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.tasks = append(q.tasks, task)
	log.Printf("Task pushed to queue. Queue size: %d", len(q.tasks))
}

func (q *TaskQueue) Pop() *models.Task {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.tasks) == 0 {
		log.Println("Attempted to pop from empty queue")
		return nil
	}

	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	log.Printf("Task popped from queue. Queue size: %d", len(q.tasks))
	return task
}

func (q *TaskQueue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.tasks)
}
