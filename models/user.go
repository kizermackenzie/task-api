package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	Password  string         `gorm:"not null" json:"-"`
	FirstName string         `gorm:"not null" json:"first_name" binding:"required"`
	LastName  string         `gorm:"not null" json:"last_name" binding:"required"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	Tasks     []Task         `gorm:"foreignKey:UserID" json:"tasks,omitempty"`
}