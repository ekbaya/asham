package models

import (
	"context"
	"errors"
	"os"
	"time"

	redisClient "github.com/ekbaya/asham/pkg/db/redis"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

var jwtSecret = []byte(getSecretKey()) // Secret key for signing and verifying JWTs

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

// GenerateJWT generates a JWT token for a user with the given ID
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
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
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

	// Check if token has been invalidated
	ctx := context.Background()
	val, err := redisClient.GetRedis().Get(ctx, "blacklist:"+tokenStr).Result()
	if err != redis.Nil {
		// Token exists in blacklist or there was an error
		if err == nil && val == "1" {
			return nil, errors.New("token has been invalidated")
		}
		// If it's any other error, log it and continue (option to be strict here)
	}

	// Return the claims from the validated token
	if token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// ValidateRefreshToken validates a refresh token and returns the claims
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
				// Even for expired tokens, check if they're blacklisted
				ctx := context.Background()
				val, redisErr := redisClient.GetRedis().Get(ctx, "blacklist:"+tokenStr).Result()
				if redisErr != redis.Nil {
					if redisErr == nil && val == "1" {
						return nil, errors.New("refresh token has been invalidated")
					}
				}
				return claims, nil // Expired tokens are acceptable for refresh tokens if not blacklisted
			}
			return nil, errors.New("invalid token")
		}
		return nil, err
	}

	// Check if token has been invalidated
	ctx := context.Background()
	val, redisErr := redisClient.GetRedis().Get(ctx, "blacklist:"+tokenStr).Result()
	if redisErr != redis.Nil {
		if redisErr == nil && val == "1" {
			return nil, errors.New("refresh token has been invalidated")
		}
	}

	return claims, nil
}

// Logout invalidates both access token and refresh token by adding them to Redis blacklist
func Logout(accessToken string, refreshToken string) error {
	ctx := context.Background()

	// Parse tokens to get their expiration time
	accessClaims := &Claims{}
	_, err := jwt.ParseWithClaims(accessToken, accessClaims, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return err
	}

	refreshClaims := &Claims{}
	_, err = jwt.ParseWithClaims(refreshToken, refreshClaims, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return err
	}

	// Calculate TTL for Redis blacklist entries (use token's remaining lifetime)
	accessTTL := time.Until(accessClaims.ExpiresAt.Time)
	if accessTTL < 0 {
		accessTTL = time.Hour // Default 1 hour if already expired
	}

	refreshTTL := time.Until(refreshClaims.ExpiresAt.Time)
	if refreshTTL < 0 {
		refreshTTL = time.Hour // Default 1 hour if already expired
	}

	// Add tokens to blacklist with their respective TTLs
	if err := redisClient.GetRedis().Set(ctx, "blacklist:"+accessToken, "1", accessTTL).Err(); err != nil {
		return err
	}

	if err := redisClient.GetRedis().Set(ctx, "blacklist:"+refreshToken, "1", refreshTTL).Err(); err != nil {
		return err
	}

	return nil
}

// LogoutUser invalidates all tokens for a specific user
func LogoutUser(userID string) error {
	// This is a more advanced implementation
	// It requires storing user->token mapping when tokens are created
	// We won't implement the full logic here, but this shows the pattern

	ctx := context.Background()

	// Get all tokens for user (from a hypothetical storage)
	userTokensKey := "user:tokens:" + userID
	tokens, err := redisClient.GetRedis().SMembers(ctx, userTokensKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	// Add all tokens to blacklist
	for _, token := range tokens {
		// Add to blacklist with a long TTL (e.g., 30 days)
		if err := redisClient.GetRedis().Set(ctx, "blacklist:"+token, "1", 30*24*time.Hour).Err(); err != nil {
			return err
		}
	}

	// Remove the user's token set
	if err := redisClient.GetRedis().Del(ctx, userTokensKey).Err(); err != nil {
		return err
	}

	return nil
}
