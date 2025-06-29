package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrTokenMalformed   = errors.New("token is malformed")
	ErrMissingSecretKey = errors.New("JWT secret key not found")
)

const (
	DefaultAccessTokenExpiry  = 15 * time.Minute
	DefaultRefreshTokenExpiry = 7 * 24 * time.Hour
)

func getJWTSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, ErrMissingSecretKey
	}
	return []byte(secret), nil
}

func getTokenExpiry(envKey string, defaultDuration time.Duration) time.Duration {
	if envValue := os.Getenv(envKey); envValue != "" {
		if minutes, err := strconv.Atoi(envValue); err == nil {
			return time.Duration(minutes) * time.Minute
		}
	}
	return defaultDuration
}

func GenerateTokenPair(userID uint, email, firstName, lastName string) (*TokenPair, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return nil, err
	}

	accessTokenExpiry := getTokenExpiry("JWT_ACCESS_TOKEN_EXPIRY", DefaultAccessTokenExpiry)
	refreshTokenExpiry := getTokenExpiry("JWT_REFRESH_TOKEN_EXPIRY", DefaultRefreshTokenExpiry)

	now := time.Now()
	
	accessClaims := JWTClaims{
		UserID:    userID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "task-api",
			Subject:   strconv.Itoa(int(userID)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(secret)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(refreshTokenExpiry)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    "task-api",
		Subject:   strconv.Itoa(int(userID)),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(secret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(accessTokenExpiry.Seconds()),
	}, nil
}

func ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenMalformed
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func ValidateRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenMalformed
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func RefreshAccessToken(refreshTokenString string, userID uint, email, firstName, lastName string) (*TokenPair, error) {
	claims, err := ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	userIDFromToken, err := strconv.Atoi(claims.Subject)
	if err != nil || uint(userIDFromToken) != userID {
		return nil, ErrInvalidToken
	}

	return GenerateTokenPair(userID, email, firstName, lastName)
}

func ExtractUserIDFromToken(tokenString string) (uint, error) {
	claims, err := ValidateAccessToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}