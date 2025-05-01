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

type TaskHandler struct {
	service services.TaskService
}

func NewTaskHandler(service services.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
	DueDate     time.Time `json:"due_date" time_format:"2006-01-02T15:04:05Z07:00"`
}

type Response struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
	Error      string      `json:"error,omitempty"`
}

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

func (h *TaskHandler) GetTasks(c *gin.Context) {
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
