package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	redisClient "github.com/ekbaya/asham/pkg/db/redis"
	"github.com/redis/go-redis/v9"
)

type MSAzureConfig struct {
	TenantID     string
	ClientID     string
	ClientSecret string
}

type TokenManager struct {
	redisClient *redis.Client
	msConfig    *MSAzureConfig
}

const tokenKey = "microsoft_graph_access_token"

func NewTokenManager(msConfig *MSAzureConfig) *TokenManager {
	return &TokenManager{redisClient: redisClient.GetRedis(), msConfig: msConfig}
}

func (tm *TokenManager) RetrieveToken(ctx context.Context) (string, error) {
	fmt.Println("[TokenManager] Attempting to retrieve token from Redis...")
	// Try getting the token from Redis
	token, err := tm.redisClient.Get(ctx, tokenKey).Result()
	if err == nil {
		fmt.Println("[TokenManager] Token found in Redis.")
		return token, nil // Found in Redis
	}
	fmt.Printf("[TokenManager] Token not found in Redis or error occurred: %v\n", err)

	// Token not found or expired, fetch a new one
	fmt.Println("[TokenManager] Fetching new token from Microsoft Graph...")
	token, err = tm.GetMicrosoftGraphAccessToken(tm.msConfig.TenantID, tm.msConfig.ClientID, tm.msConfig.ClientSecret)
	if err != nil {
		fmt.Printf("[TokenManager] Failed to fetch new token: %v\n", err)
		return "", fmt.Errorf("failed to fetch new token: %w", err)
	}
	fmt.Println("[TokenManager] Successfully fetched new token.")

	// Store in Redis with expiry slightly less than 1 hour (to be safe)
	fmt.Println("[TokenManager] Storing new token in Redis...")
	err = tm.redisClient.Set(ctx, tokenKey, token, 55*time.Minute).Err()
	if err != nil {
		fmt.Printf("[TokenManager] Failed to store token in Redis: %v\n", err)
		return "", fmt.Errorf("failed to store token in Redis: %w", err)
	}
	fmt.Println("[TokenManager] Token successfully stored in Redis.")

	return token, nil
}

func (tm *TokenManager) GetMicrosoftGraphAccessToken(tenantID, clientID, clientSecret string) (string, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("scope", "https://graph.microsoft.com/.default")
	form.Set("client_secret", clientSecret)
	form.Set("grant_type", "client_credentials")

	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed: %s", string(body))
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	return token.AccessToken, nil
}
