package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"task-management-system/models"
	"task-management-system/repository"
	"task-management-system/services"
)

// TaskHandler handles HTTP requests for tasks
type TaskHandler struct {
	service services.TaskService
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(service services.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
}

// UpdateTaskRequest represents the request body for updating a task
type UpdateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
	DueDate     time.Time `json:"due_date" time_format:"2006-01-02T15:04:05Z07:00"`
}

// Example JSON payload:
// {
//   "title": "Update project documentation",
//   "description": "Review and update all API documentation for the new release",
//   "status": "in_progress",
//   "due_date": "2023-12-31T23:59:59Z"
// }

// Example JSON payload for UpdateTaskRequest:
// {
//   "title": "Update project documentation",
//   "description": "Review and update all API documentation for the new release",
//   "status": "in_progress",
//   "due_date": "2023-12-31T23:59:59Z"
// }

// Response represents the standard API response
type Response struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
	Error      string      `json:"error,omitempty"`
}

// CreateTask handles the creation of a new task
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "Invalid request body"})
		return
	}

	task, err := h.service.CreateTask(req.Title, req.Description, req.DueDate)
	if err != nil {
		if err == services.ErrInvalidTask {
			c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, Response{Success: true, Data: task})
}

// GetTask handles retrieving a task by ID
func (h *TaskHandler) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "Invalid task ID"})
		return
	}

	task, err := h.service.GetTaskByID(uint(id))
	if err != nil {
		if err == repository.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, Response{Success: false, Error: "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, Response{Success: true, Data: task})
}

// GetTasks handles retrieving all tasks with pagination and filtering
func (h *TaskHandler) GetTasks(c *gin.Context) {
	// Get page and limit parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Get filters
	filters := make(map[string]string)
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	tasks, pagination, err := h.service.GetTasks(page, limit, filters)
	if err != nil {
		if err == models.ErrInvalidStatus {
			c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success:    true,
		Data:       tasks,
		Pagination: pagination,
	})
}

// UpdateTask handles updating a task
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "Invalid task ID"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "Invalid request body"})
		return
	}

	task, err := h.service.UpdateTask(uint(id), req.Title, req.Description, req.Status, req.DueDate)
	if err != nil {
		switch err {
		case repository.ErrTaskNotFound:
			c.JSON(http.StatusNotFound, Response{Success: false, Error: "Task not found"})
		case models.ErrInvalidStatus:
			c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, Response{Success: true, Data: task})
}

// DeleteTask handles deleting a task
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Success: false, Error: "Invalid task ID"})
		return
	}

	err = h.service.DeleteTask(uint(id))
	if err != nil {
		if err == repository.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, Response{Success: false, Error: "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, Response{Success: true, Data: map[string]string{"message": "Task deleted successfully"}})
}
