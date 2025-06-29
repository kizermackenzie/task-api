package handlers

import (
	"net/http"

	"task-api/middleware"
	"task-api/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var dto services.RegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	result, err := h.authService.Register(dto)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorType := "internal_error"
		message := "Registration failed"

		switch err {
		case services.ErrEmailAlreadyExists:
			statusCode = http.StatusConflict
			errorType = "email_exists"
			message = "Email address is already registered"
		case services.ErrWeakPassword:
			statusCode = http.StatusBadRequest
			errorType = "weak_password"
			message = "Password must be at least 8 characters with uppercase, lowercase, and digit"
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
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"data":    result,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var dto services.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	result, err := h.authService.Login(dto)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorType := "internal_error"
		message := "Login failed"

		switch err {
		case services.ErrInvalidCredentials:
			statusCode = http.StatusUnauthorized
			errorType = "invalid_credentials"
			message = "Invalid email or password"
		case services.ErrUserInactive:
			statusCode = http.StatusForbidden
			errorType = "account_inactive"
			message = "Account is inactive"
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
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    result,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var dto services.RefreshTokenDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	result, err := h.authService.RefreshToken(dto)
	if err != nil {
		statusCode := http.StatusUnauthorized
		errorType := "invalid_token"
		message := "Token refresh failed"

		switch err {
		case services.ErrUserNotFound:
			statusCode = http.StatusNotFound
			errorType = "user_not_found"
			message = "User not found"
		case services.ErrUserInactive:
			statusCode = http.StatusForbidden
			errorType = "account_inactive"
			message = "Account is inactive"
		}

		c.JSON(statusCode, gin.H{
			"error":   errorType,
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"data":    result,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.RequireUserID(c)
	if userID == 0 {
		return
	}

	result, err := h.authService.GetUserProfile(userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorType := "internal_error"
		message := "Failed to get user profile"

		if err == services.ErrUserNotFound {
			statusCode = http.StatusNotFound
			errorType = "user_not_found"
			message = "User not found"
		}

		c.JSON(statusCode, gin.H{
			"error":   errorType,
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    result,
	})
}