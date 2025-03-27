package models

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(getSecretKey()) // Secret key for signing and verifying JWTs, retrieved from environment variable

// Claims represents the custom claims for JWT
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// getSecretKey retrieves the JWT secret from the environment variable
func getSecretKey() string {
	// It's recommended to use an environment variable to store the JWT secret securely
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// If the secret is not found, return a default value (not recommended for production)
		// Replace this with better error handling in production
		return "your-secret-key"
	}
	return secret
}

// GenerateJWT generates a JWT token for a user with the given ID and role
func GenerateJWT(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			// Token expires in 24 hours
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken generates a long-lived refresh token for a user
func GenerateRefreshToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			// Refresh token expires in 7 days
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates a JWT token and returns the claims if valid
func ValidateJWT(tokenStr string) (*Claims, error) {
	// Parse and validate the token
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	// If there is an error or the token is not valid, return an error
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			// You can further check the error type to handle specific cases (like Expired, Invalid)
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token has expired")
			}
			return nil, errors.New("invalid token")
		}
		return nil, err
	}

	// Return the claims from the validated token
	return claims, nil
}

// ValidateRefreshToken validates a refresh token and returns the claims, skipping the expiration check
func ValidateRefreshToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	// If there is an error, ensure it's not due to expiration
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			// Ignore expiration errors for refresh tokens
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return claims, nil // Expired tokens are acceptable for refresh tokens
			}
			return nil, errors.New("invalid token")
		}
		return nil, err
	}

	return claims, nil
}
