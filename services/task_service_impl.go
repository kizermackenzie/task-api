package services

import (
	"errors"
	"strings"
	"time"

	"task-api/models"
	"task-api/repositories"

	"gorm.io/gorm"
)

type taskService struct {
	taskRepo repositories.TaskRepository
}

func NewTaskService(taskRepo repositories.TaskRepository) TaskService {
	return &taskService{
		taskRepo: taskRepo,
	}
}

func (s *taskService) CreateTask(userID uint, dto CreateTaskDTO) (*TaskResponseDTO, error) {
	if err := s.validateCreateTask(dto); err != nil {
		return nil, err
	}

	task := dto.ToModel(userID)
	
	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	response := TaskToResponseDTO(task)
	return &response, nil
}

func (s *taskService) GetTaskByID(userID, taskID uint) (*TaskResponseDTO, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	if task.UserID != userID {
		return nil, ErrUnauthorizedAccess
	}

	response := TaskToResponseDTO(task)
	return &response, nil
}

func (s *taskService) GetUserTasks(userID uint, pagination repositories.PaginationParams) (*TaskListResponseDTO, error) {
	tasks, paginationResult, err := s.taskRepo.GetByUserID(userID, pagination)
	if err != nil {
		return nil, err
	}

	return &TaskListResponseDTO{
		Tasks:      TasksToResponseDTO(tasks),
		Pagination: paginationResult,
	}, nil
}

func (s *taskService) UpdateTask(userID, taskID uint, dto UpdateTaskDTO) (*TaskResponseDTO, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	if task.UserID != userID {
		return nil, ErrUnauthorizedAccess
	}

	if err := s.validateUpdateTask(dto, task); err != nil {
		return nil, err
	}

	s.applyUpdates(task, dto)

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	response := TaskToResponseDTO(task)
	return &response, nil
}

func (s *taskService) DeleteTask(userID, taskID uint) error {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	if task.UserID != userID {
		return ErrUnauthorizedAccess
	}

	return s.taskRepo.Delete(taskID)
}

func (s *taskService) GetAllTasks(pagination repositories.PaginationParams) (*TaskListResponseDTO, error) {
	tasks, paginationResult, err := s.taskRepo.List(pagination)
	if err != nil {
		return nil, err
	}

	return &TaskListResponseDTO{
		Tasks:      TasksToResponseDTO(tasks),
		Pagination: paginationResult,
	}, nil
}

func (s *taskService) CompleteTask(userID, taskID uint) (*TaskResponseDTO, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	if task.UserID != userID {
		return nil, ErrUnauthorizedAccess
	}

	if task.Status == models.TaskStatusCompleted {
		return nil, ErrTaskAlreadyCompleted
	}

	now := time.Now()
	task.Status = models.TaskStatusCompleted
	task.CompletedAt = &now

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	response := TaskToResponseDTO(task)
	return &response, nil
}

func (s *taskService) validateCreateTask(dto CreateTaskDTO) error {
	var validationErrors ValidationErrors

	if strings.TrimSpace(dto.Title) == "" {
		validationErrors.AddError("title", "title is required")
	}

	if dto.DueDate != nil && dto.DueDate.Before(time.Now()) {
		validationErrors.AddError("due_date", "due date cannot be in the past")
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}

func (s *taskService) validateUpdateTask(dto UpdateTaskDTO, currentTask *models.Task) error {
	var validationErrors ValidationErrors

	if dto.Title != nil && strings.TrimSpace(*dto.Title) == "" {
		validationErrors.AddError("title", "title cannot be empty")
	}

	if dto.DueDate != nil && dto.DueDate.Before(time.Now()) {
		validationErrors.AddError("due_date", "due date cannot be in the past")
	}

	if dto.Status != nil && currentTask.Status == models.TaskStatusCompleted && *dto.Status != models.TaskStatusCompleted {
		validationErrors.AddError("status", "cannot change status of completed task")
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}

func (s *taskService) applyUpdates(task *models.Task, dto UpdateTaskDTO) {
	if dto.Title != nil {
		task.Title = *dto.Title
	}
	if dto.Description != nil {
		task.Description = *dto.Description
	}
	if dto.Status != nil {
		task.Status = *dto.Status
		if *dto.Status == models.TaskStatusCompleted && task.CompletedAt == nil {
			now := time.Now()
			task.CompletedAt = &now
		} else if *dto.Status != models.TaskStatusCompleted {
			task.CompletedAt = nil
		}
	}
	if dto.Priority != nil {
		task.Priority = *dto.Priority
	}
	if dto.DueDate != nil {
		task.DueDate = dto.DueDate
	}
}