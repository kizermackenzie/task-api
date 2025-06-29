package repositories

import (
	"task-api/models"
)

type TaskRepository interface {
	Create(task *models.Task) error
	GetByID(id uint) (*models.Task, error)
	GetByUserID(userID uint, pagination PaginationParams) ([]models.Task, PaginationResult, error)
	Update(task *models.Task) error
	Delete(id uint) error
	List(pagination PaginationParams) ([]models.Task, PaginationResult, error)
}