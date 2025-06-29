package services

import (
	"task-api/models"
	"task-api/utils"
	"time"
)

type RegisterDTO struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=100"`
	FirstName string `json:"first_name" binding:"required,min=1,max=50"`
	LastName  string `json:"last_name" binding:"required,min=1,max=50"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponseDTO struct {
	User   UserResponseDTO  `json:"user"`
	Tokens utils.TokenPair `json:"tokens"`
}

type UserResponseDTO struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (dto RegisterDTO) ToModel() *models.User {
	return &models.User{
		Email:     dto.Email,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		IsActive:  true,
	}
}

func UserToResponseDTO(user *models.User) UserResponseDTO {
	return UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}