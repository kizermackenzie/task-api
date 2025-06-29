package middleware

import (
	"net/http"
	"strings"

	"task-api/utils"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey    = "user_id"
	UserEmailKey = "user_email"
	UserClaimsKey = "user_claims"
	AuthorizationHeader = "Authorization"
	BearerPrefix = "Bearer "
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func AuthRequired() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		token := extractTokenFromHeader(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "Authorization token is required",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateAccessToken(token)
		if err != nil {
			var message string
			switch err {
			case utils.ErrExpiredToken:
				message = "Token has expired"
			case utils.ErrTokenMalformed:
				message = "Token is malformed"
			case utils.ErrInvalidToken:
				message = "Invalid token"
			case utils.ErrMissingSecretKey:
				message = "Authentication service unavailable"
			default:
				message = "Authentication failed"
			}

			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: message,
			})
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Set(UserClaimsKey, claims)

		c.Next()
	})
}

func AuthOptional() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		token := extractTokenFromHeader(c)
		if token == "" {
			c.Next()
			return
		}

		claims, err := utils.ValidateAccessToken(token)
		if err == nil {
			c.Set(UserIDKey, claims.UserID)
			c.Set(UserEmailKey, claims.Email)
			c.Set(UserClaimsKey, claims)
		}

		c.Next()
	})
}

func extractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader(AuthorizationHeader)
	if authHeader == "" {
		return ""
	}

	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return ""
	}

	return strings.TrimPrefix(authHeader, BearerPrefix)
}

func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	return id, ok
}

func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(UserEmailKey)
	if !exists {
		return "", false
	}

	emailStr, ok := email.(string)
	return emailStr, ok
}

func GetUserClaims(c *gin.Context) (*utils.JWTClaims, bool) {
	claims, exists := c.Get(UserClaimsKey)
	if !exists {
		return nil, false
	}

	userClaims, ok := claims.(*utils.JWTClaims)
	return userClaims, ok
}

func RequireUserID(c *gin.Context) uint {
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		c.Abort()
		return 0
	}
	return userID
}