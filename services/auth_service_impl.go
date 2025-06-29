package services

import (
	"errors"
	"strings"

	"task-api/repositories"
	"task-api/utils"

	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrWeakPassword       = errors.New("password does not meet requirements")
)

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(dto RegisterDTO) (*AuthResponseDTO, error) {
	if err := s.validateRegistration(dto); err != nil {
		return nil, err
	}

	exists, err := s.userRepo.EmailExists(dto.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	user := dto.ToModel()
	user.Password = hashedPassword

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	tokens, err := utils.GenerateTokenPair(user.ID, user.Email, user.FirstName, user.LastName)
	if err != nil {
		return nil, err
	}

	return &AuthResponseDTO{
		User:   UserToResponseDTO(user),
		Tokens: *tokens,
	}, nil
}

func (s *authService) Login(dto LoginDTO) (*AuthResponseDTO, error) {
	if err := s.validateLogin(dto); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByEmail(dto.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if !utils.CheckPasswordHash(dto.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	tokens, err := utils.GenerateTokenPair(user.ID, user.Email, user.FirstName, user.LastName)
	if err != nil {
		return nil, err
	}

	return &AuthResponseDTO{
		User:   UserToResponseDTO(user),
		Tokens: *tokens,
	}, nil
}

func (s *authService) RefreshToken(dto RefreshTokenDTO) (*AuthResponseDTO, error) {
	claims, err := utils.ValidateRefreshToken(dto.RefreshToken)
	if err != nil {
		return nil, err
	}

	userID := uint(0)
	if claims.Subject != "" {
		if id, parseErr := utils.ExtractUserIDFromToken(dto.RefreshToken); parseErr == nil {
			userID = id
		}
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	tokens, err := utils.RefreshAccessToken(dto.RefreshToken, user.ID, user.Email, user.FirstName, user.LastName)
	if err != nil {
		return nil, err
	}

	return &AuthResponseDTO{
		User:   UserToResponseDTO(user),
		Tokens: *tokens,
	}, nil
}

func (s *authService) GetUserProfile(userID uint) (*UserResponseDTO, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	response := UserToResponseDTO(user)
	return &response, nil
}

func (s *authService) validateRegistration(dto RegisterDTO) error {
	var validationErrors ValidationErrors

	if strings.TrimSpace(dto.Email) == "" {
		validationErrors.AddError("email", "email is required")
	}

	if strings.TrimSpace(dto.FirstName) == "" {
		validationErrors.AddError("first_name", "first name is required")
	}

	if strings.TrimSpace(dto.LastName) == "" {
		validationErrors.AddError("last_name", "last name is required")
	}

	if err := s.validatePassword(dto.Password); err != nil {
		validationErrors.AddError("password", err.Error())
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}

func (s *authService) validateLogin(dto LoginDTO) error {
	var validationErrors ValidationErrors

	if strings.TrimSpace(dto.Email) == "" {
		validationErrors.AddError("email", "email is required")
	}

	if strings.TrimSpace(dto.Password) == "" {
		validationErrors.AddError("password", "password is required")
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}

func (s *authService) validatePassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return ErrWeakPassword
	}

	return nil
}