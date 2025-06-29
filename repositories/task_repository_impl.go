package repositories

import (
	"task-api/models"

	"gorm.io/gorm"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (r *taskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) GetByID(id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.Preload("User").First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) GetByUserID(userID uint, pagination PaginationParams) ([]models.Task, PaginationResult, error) {
	var tasks []models.Task
	var total int64

	query := r.db.Model(&models.Task{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, PaginationResult{}, err
	}

	err := query.Preload("User").
		Offset(pagination.GetOffset()).
		Limit(pagination.PageSize).
		Order("created_at DESC").
		Find(&tasks).Error

	if err != nil {
		return nil, PaginationResult{}, err
	}

	paginationResult := NewPaginationResult(pagination.Page, pagination.PageSize, total)
	return tasks, paginationResult, nil
}

func (r *taskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(id uint) error {
	return r.db.Delete(&models.Task{}, id).Error
}

func (r *taskRepository) List(pagination PaginationParams) ([]models.Task, PaginationResult, error) {
	var tasks []models.Task
	var total int64

	query := r.db.Model(&models.Task{})

	if err := query.Count(&total).Error; err != nil {
		return nil, PaginationResult{}, err
	}

	err := query.Preload("User").
		Offset(pagination.GetOffset()).
		Limit(pagination.PageSize).
		Order("created_at DESC").
		Find(&tasks).Error

	if err != nil {
		return nil, PaginationResult{}, err
	}

	paginationResult := NewPaginationResult(pagination.Page, pagination.PageSize, total)
	return tasks, paginationResult, nil
}