package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null" json:"title" binding:"required"`
	Description string         `gorm:"type:text" json:"description"`
	Status      TaskStatus     `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Priority    TaskPriority   `gorm:"type:varchar(20);default:'medium'" json:"priority"`
	DueDate     *time.Time     `json:"due_date,omitempty"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}