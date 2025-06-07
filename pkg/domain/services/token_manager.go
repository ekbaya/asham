package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ekbaya/asham/pkg/config"
	redisClient "github.com/ekbaya/asham/pkg/db/redis"
	"github.com/redis/go-redis/v9"
)

type MSAzureConfig struct {
	TenantID     string
	ClientID     string
	ClientSecret string
}

type TokenManager struct {
	redisClient  *redis.Client
	msConfig     *MSAzureConfig
	emailService *EmailService
}

const tokenKey = "microsoft_graph_access_token"
const delegateTokenKey = "microsoft_graph_delegate_access_token"

func NewTokenManager(msConfig *MSAzureConfig, emailService *EmailService) *TokenManager {
	return &TokenManager{redisClient: redisClient.GetRedis(), msConfig: msConfig, emailService: emailService}
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

func (tm *TokenManager) RetrieveDelegateToken(ctx context.Context) (string, error) {
	fmt.Println("[TokenManager] Attempting to retrieve delegate token from Redis...")
	// Try getting the token from Redis
	token, err := tm.redisClient.Get(ctx, delegateTokenKey).Result()
	if err == nil {
		fmt.Println("[TokenManager] Delegate Token found in Redis.")
		return token, nil // Found in Redis
	}
	fmt.Printf("[TokenManager] Delegate Token not found in Redis or error occurred: %v\n", err)

	// Token not found or expired, fetch a new one
	fmt.Println("[TokenManager] Fetching new delegate token from Microsoft Graph...")
	token, err = tm.GetDelegatedAccessTokenViaDeviceCode(tm.msConfig.ClientID, tm.msConfig.ClientSecret, tm.msConfig.TenantID)
	if err != nil {
		fmt.Printf("[TokenManager] Failed to fetch new delegate token: %v\n", err)
		return "", fmt.Errorf("failed to fetch new delegate token: %w", err)
	}
	fmt.Println("[TokenManager] Successfully fetched new delegate token.")

	// Store in Redis with expiry slightly less than 1 hour (to be safe)
	fmt.Println("[TokenManager] Storing new delegate token in Redis...")
	err = tm.redisClient.Set(ctx, delegateTokenKey, token, 55*time.Minute).Err()
	if err != nil {
		fmt.Printf("[TokenManager] Failed to store delegate token in Redis: %v\n", err)
		return "", fmt.Errorf("failed to store delegate token in Redis: %w", err)
	}
	fmt.Println("[TokenManager] Delegate Token successfully stored in Redis.")

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

func (tm *TokenManager) GetDelegatedAccessTokenViaDeviceCode(clientID, clientSecret, tenantID string) (string, error) {
	deviceCodeURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/devicecode", tenantID)
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)

	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("scope", "https://graph.microsoft.com/User.Read Files.ReadWrite.All offline_access")

	// Step 1: Get device code
	resp, err := http.PostForm(deviceCodeURL, form)
	if err != nil {
		return "", fmt.Errorf("failed to get device code: %w", err)
	}
	defer resp.Body.Close()

	var deviceResp struct {
		UserCode        string `json:"user_code"`
		DeviceCode      string `json:"device_code"`
		VerificationURL string `json:"verification_uri"`
		ExpiresIn       int    `json:"expires_in"`
		Interval        int    `json:"interval"`
		Message         string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&deviceResp); err != nil {
		return "", fmt.Errorf("failed to parse device code response: %w", err)
	}

	fmt.Println(deviceResp.Message)

	go func() {
		email := config.GetConfig().AZURE_USER_EMAIL
		recipients := []RecipientEmail{
			{
				To:    email,
				Title: "Action Required: Consent to Grant Access",
				Body: fmt.Sprintf(`
				<!DOCTYPE html>
				<html>
				<body style="font-family: Arial, sans-serif; line-height: 1.6;">
					<h2>Hello Doreen,</h2>
					<p>To proceed, please grant access to the application by following the instructions below:</p>
					<p><strong>Step 1:</strong> Visit the verification page:</p>
					<p><a href="%s" target="_blank" style="color: #1a73e8;">%s</a></p>
					<p><strong>Step 2:</strong> Enter the code shown below:</p>
					<p style="font-size: 20px; font-weight: bold; color: #333;">%s</p>
					<p>This code will expire in approximately %d minutes. Please complete the authorization promptly.</p>
					<br>
					<p>Thank you,<br>ASHAM Dev Team</p>
				</body>
				</html>
			`, deviceResp.VerificationURL, deviceResp.VerificationURL, deviceResp.UserCode, deviceResp.ExpiresIn/60),
			},
		}

		if err := tm.emailService.SendCustomEmails(recipients); err != nil {
			log.Fatalf("error sending consent email: %v", err)
		}
	}()

	// Step 2: Poll for token
	start := time.Now()
	for {
		if time.Since(start) > time.Duration(deviceResp.ExpiresIn)*time.Second {
			return "", fmt.Errorf("device code expired before authorization")
		}

		time.Sleep(time.Duration(deviceResp.Interval) * time.Second)

		tokenForm := url.Values{}
		tokenForm.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
		tokenForm.Set("client_id", clientID)
		tokenForm.Set("device_code", deviceResp.DeviceCode)
		tokenForm.Set("client_secret", clientSecret)

		tokenResp, err := http.PostForm(tokenURL, tokenForm)
		if err != nil {
			return "", fmt.Errorf("failed to request token: %w", err)
		}
		defer tokenResp.Body.Close()

		body, _ := io.ReadAll(tokenResp.Body)
		if tokenResp.StatusCode == 200 {
			var token TokenResponse
			if err := json.Unmarshal(body, &token); err != nil {
				return "", fmt.Errorf("failed to parse token: %w", err)
			}
			return token.AccessToken, nil
		} else if strings.Contains(string(body), "authorization_pending") {
			continue
		} else {
			return "", fmt.Errorf("token request error: %s", string(body))
		}
	}
}
