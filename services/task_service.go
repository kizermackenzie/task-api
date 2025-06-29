package services

import (
	"task-api/repositories"
)

type TaskService interface {
	CreateTask(userID uint, dto CreateTaskDTO) (*TaskResponseDTO, error)
	GetTaskByID(userID, taskID uint) (*TaskResponseDTO, error)
	GetUserTasks(userID uint, pagination repositories.PaginationParams) (*TaskListResponseDTO, error)
	UpdateTask(userID, taskID uint, dto UpdateTaskDTO) (*TaskResponseDTO, error)
	DeleteTask(userID, taskID uint) error
	GetAllTasks(pagination repositories.PaginationParams) (*TaskListResponseDTO, error)
	CompleteTask(userID, taskID uint) (*TaskResponseDTO, error)
}