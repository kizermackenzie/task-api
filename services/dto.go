package services

import (
	"task-api/models"
	"task-api/repositories"
	"time"
)

type CreateTaskDTO struct {
	Title       string                `json:"title" binding:"required,min=1,max=255"`
	Description string                `json:"description" binding:"max=1000"`
	Priority    models.TaskPriority   `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time            `json:"due_date,omitempty"`
}

type UpdateTaskDTO struct {
	Title       *string               `json:"title,omitempty" binding:"omitempty,min=1,max=255"`
	Description *string               `json:"description,omitempty" binding:"omitempty,max=1000"`
	Status      *models.TaskStatus    `json:"status,omitempty" binding:"omitempty,oneof=pending in_progress completed cancelled"`
	Priority    *models.TaskPriority  `json:"priority,omitempty" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time            `json:"due_date,omitempty"`
}

type TaskResponseDTO struct {
	ID          uint                 `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Status      models.TaskStatus    `json:"status"`
	Priority    models.TaskPriority  `json:"priority"`
	DueDate     *time.Time           `json:"due_date,omitempty"`
	CompletedAt *time.Time           `json:"completed_at,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	UserID      uint                 `json:"user_id"`
}

type TaskListResponseDTO struct {
	Tasks      []TaskResponseDTO           `json:"tasks"`
	Pagination repositories.PaginationResult `json:"pagination"`
}

func (dto CreateTaskDTO) ToModel(userID uint) *models.Task {
	task := &models.Task{
		Title:       dto.Title,
		Description: dto.Description,
		Status:      models.TaskStatusPending,
		Priority:    dto.Priority,
		DueDate:     dto.DueDate,
		UserID:      userID,
	}

	if task.Priority == "" {
		task.Priority = models.TaskPriorityMedium
	}

	return task
}

func TaskToResponseDTO(task *models.Task) TaskResponseDTO {
	return TaskResponseDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		CompletedAt: task.CompletedAt,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		UserID:      task.UserID,
	}
}

func TasksToResponseDTO(tasks []models.Task) []TaskResponseDTO {
	result := make([]TaskResponseDTO, len(tasks))
	for i, task := range tasks {
		result[i] = TaskToResponseDTO(&task)
	}
	return result
}