package handlers

import (
	"net/http"
	"strconv"

	"task-api/middleware"
	"task-api/repositories"
	"task-api/services"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService services.TaskService
}

func NewTaskHandler(taskService services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID := middleware.RequireUserID(c)
	if userID == 0 {
		return
	}

	var dto services.CreateTaskDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	result, err := h.taskService.CreateTask(userID, dto)
	if err != nil {
		h.handleServiceError(c, err, "Task creation failed")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"data":    result,
	})
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	userID := middleware.RequireUserID(c)
	if userID == 0 {
		return
	}

	taskID, err := h.getTaskIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_task_id",
			"message": "Invalid task ID",
		})
		return
	}

	result, err := h.taskService.GetTaskByID(userID, taskID)
	if err != nil {
		h.handleServiceError(c, err, "Failed to get task")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task retrieved successfully",
		"data":    result,
	})
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID := middleware.RequireUserID(c)
	if userID == 0 {
		return
	}

	taskID, err := h.getTaskIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_task_id",
			"message": "Invalid task ID",
		})
		return
	}

	var dto services.UpdateTaskDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	result, err := h.taskService.UpdateTask(userID, taskID, dto)
	if err != nil {
		h.handleServiceError(c, err, "Task update failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"data":    result,
	})
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID := middleware.RequireUserID(c)
	if userID == 0 {
		return
	}

	taskID, err := h.getTaskIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_task_id",
			"message": "Invalid task ID",
		})
		return
	}

	err = h.taskService.DeleteTask(userID, taskID)
	if err != nil {
		h.handleServiceError(c, err, "Task deletion failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}

func (h *TaskHandler) GetUserTasks(c *gin.Context) {
	userID := middleware.RequireUserID(c)
	if userID == 0 {
		return
	}

	pagination := h.getPaginationParams(c)

	result, err := h.taskService.GetUserTasks(userID, pagination)
	if err != nil {
		h.handleServiceError(c, err, "Failed to get tasks")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tasks retrieved successfully",
		"data":    result,
	})
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	pagination := h.getPaginationParams(c)

	result, err := h.taskService.GetAllTasks(pagination)
	if err != nil {
		h.handleServiceError(c, err, "Failed to get all tasks")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All tasks retrieved successfully",
		"data":    result,
	})
}

func (h *TaskHandler) CompleteTask(c *gin.Context) {
	userID := middleware.RequireUserID(c)
	if userID == 0 {
		return
	}

	taskID, err := h.getTaskIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_task_id",
			"message": "Invalid task ID",
		})
		return
	}

	result, err := h.taskService.CompleteTask(userID, taskID)
	if err != nil {
		h.handleServiceError(c, err, "Task completion failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task completed successfully",
		"data":    result,
	})
}

func (h *TaskHandler) getTaskIDFromParam(c *gin.Context) (uint, error) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (h *TaskHandler) getPaginationParams(c *gin.Context) repositories.PaginationParams {
	page := 1
	pageSize := 10

	if pageParam := c.Query("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeParam := c.Query("page_size"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	return repositories.NewPaginationParams(page, pageSize)
}

func (h *TaskHandler) handleServiceError(c *gin.Context, err error, defaultMessage string) {
	statusCode := http.StatusInternalServerError
	errorType := "internal_error"
	message := defaultMessage

	switch err {
	case services.ErrTaskNotFound:
		statusCode = http.StatusNotFound
		errorType = "task_not_found"
		message = "Task not found"
	case services.ErrUnauthorizedAccess:
		statusCode = http.StatusForbidden
		errorType = "unauthorized_access"
		message = "You don't have permission to access this task"
	case services.ErrTaskAlreadyCompleted:
		statusCode = http.StatusConflict
		errorType = "task_already_completed"
		message = "Task is already completed"
	case services.ErrDueDateInPast:
		statusCode = http.StatusBadRequest
		errorType = "invalid_due_date"
		message = "Due date cannot be in the past"
	case services.ErrInvalidInput:
		statusCode = http.StatusBadRequest
		errorType = "invalid_input"
		message = "Invalid input data"
	default:
		if validationErr, ok := err.(services.ValidationErrors); ok {
			statusCode = http.StatusBadRequest
			errorType = "validation_error"
			message = "Validation failed"
			c.JSON(statusCode, gin.H{
				"error":   errorType,
				"message": message,
				"details": validationErr.Errors,
			})
			return
		}
	}

	c.JSON(statusCode, gin.H{
		"error":   errorType,
		"message": message,
	})
}