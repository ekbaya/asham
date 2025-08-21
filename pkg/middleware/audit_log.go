package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuditLogMiddleware provides audit logging functionality for HTTP requests
type AuditLogMiddleware struct {
	auditLogService *services.AuditLogService
}

// NewAuditLogMiddleware creates a new audit log middleware instance
func NewAuditLogMiddleware(auditLogService *services.AuditLogService) *AuditLogMiddleware {
	return &AuditLogMiddleware{
		auditLogService: auditLogService,
	}
}

// responseWriter wraps gin.ResponseWriter to capture response data
type responseWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// AuditLog returns a middleware function that logs HTTP requests
func (m *AuditLogMiddleware) AuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip audit logging for certain paths
		if m.shouldSkipAudit(c.Request.URL.Path, c.Request.Method) {
			c.Next()
			return
		}

		// Generate request ID for correlation
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Capture request data
		startTime := time.Now()
		userID := c.GetString("user_id")
		sessionID := c.GetString("session_id")
		ipAddress := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		// Read and restore request body
		requestBody := m.captureRequestBody(c)

		// Wrap response writer to capture response
		respWriter := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer([]byte{}),
			status:         200,
		}
		c.Writer = respWriter

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Determine if this action should be audited
		auditInfo := m.getAuditInfo(c.Request.Method, c.Request.URL.Path, respWriter.status)
		if auditInfo == nil {
			return
		}

		// Extract resource information from request
		resourceInfo := m.extractResourceInfo(c, requestBody, respWriter.body.String())

		// Create audit log parameters
		var userIDPtr *string
		if userID != "" {
			userIDPtr = &userID
		}

		// Parse metadata
		var metadata map[string]interface{}
		metadataStr := m.buildMetadata(c, requestBody, auditInfo)
		json.Unmarshal([]byte(metadataStr), &metadata)

		// Parse old/new values
		var oldValues, newValues map[string]interface{}
		if resourceInfo.OldValues != "" {
			json.Unmarshal([]byte(resourceInfo.OldValues), &oldValues)
		}
		if resourceInfo.NewValues != "" {
			json.Unmarshal([]byte(resourceInfo.NewValues), &newValues)
		}

		auditParams := services.LogActionParams{
			UserID:        userIDPtr,
			Action:        auditInfo.Action,
			Module:        auditInfo.Module,
			ResourceType:  auditInfo.ResourceType,
			ResourceID:    resourceInfo.ResourceID,
			ResourceTitle: resourceInfo.ResourceTitle,
			Description:   m.generateDescription(auditInfo, resourceInfo, c),
			Metadata:      metadata,
			OldValues:     oldValues,
			NewValues:     newValues,
			IPAddress:     ipAddress,
			UserAgent:     userAgent,
			SessionID:     sessionID,
			RequestID:     requestID,
			Success:       respWriter.status >= 200 && respWriter.status < 400,
			ErrorMessage:  m.extractErrorMessage(respWriter.body.String(), respWriter.status),
			Duration:      duration.Milliseconds(),
		}

		// Log the audit entry asynchronously
		go func() {
			if err := m.auditLogService.LogAction(auditParams); err != nil {
				// Log error but don't fail the request
				fmt.Printf("Failed to log audit entry: %v\n", err)
			}
		}()
	}
}

// AuditInfo contains information about what should be audited
type AuditInfo struct {
	Action       models.ActionType
	Module       models.ModuleType
	ResourceType string
}

// ResourceInfo contains extracted resource information
type ResourceInfo struct {
	ResourceID    *string
	ResourceTitle string
	OldValues     string
	NewValues     string
}

// shouldSkipAudit determines if a request should be skipped from audit logging
func (m *AuditLogMiddleware) shouldSkipAudit(path, method string) bool {
	// Skip health checks, metrics, and other non-business endpoints
	skipPaths := []string{
		"/health",
		"/metrics",
		"/ping",
		"/favicon.ico",
		"/static/",
		"/assets/",
		"/swagger/",
		"/docs/",
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	// Skip GET requests to audit logs (to prevent recursive logging)
	if strings.HasPrefix(path, "/api/audit-logs") && method == "GET" {
		return true
	}

	return false
}

// getAuditInfo determines what should be audited based on the request
func (m *AuditLogMiddleware) getAuditInfo(method, path string, statusCode int) *AuditInfo {
	// Only audit successful operations and some failed ones
	if statusCode >= 500 {
		return nil // Skip server errors
	}

	// Define audit rules based on path patterns
	auditRules := map[string]AuditInfo{
		// User Management
		"POST:/api/users":           {models.ActionUserCreate, models.ModuleUsers, "User"},
		"PUT:/api/users/":           {models.ActionUserUpdate, models.ModuleUsers, "User"},
		"DELETE:/api/users/":        {models.ActionUserDelete, models.ModuleUsers, "User"},
		"POST:/api/auth/login":      {models.ActionUserLogin, models.ModuleUsers, "Session"},
		"POST:/api/auth/logout":     {models.ActionUserLogout, models.ModuleUsers, "Session"},
		"POST:/api/auth/register":   {models.ActionUserCreate, models.ModuleUsers, "User"},

		// Project Management
		"POST:/api/projects":        {models.ActionProjectCreate, models.ModuleProjects, "Project"},
		"PUT:/api/projects/":        {models.ActionProjectUpdate, models.ModuleProjects, "Project"},
		"DELETE:/api/projects/":     {models.ActionProjectDelete, models.ModuleProjects, "Project"},
		"POST:/api/projects/*/submit": {models.ActionProjectPublish, models.ModuleProjects, "Project"},
		"POST:/api/projects/*/approve": {models.ActionProjectApprove, models.ModuleProjects, "Project"},
		"POST:/api/projects/*/reject": {models.ActionProjectReject, models.ModuleProjects, "Project"},

		// Document Management
		"POST:/api/documents":       {models.ActionDocumentUpload, models.ModuleDocuments, "Document"},
		"PUT:/api/documents/":       {models.ActionDocumentUpdate, models.ModuleDocuments, "Document"},
		"DELETE:/api/documents/":    {models.ActionDocumentDelete, models.ModuleDocuments, "Document"},
		"GET:/api/documents/*/download": {models.ActionDocumentDownload, models.ModuleDocuments, "Document"},

		// Ballot Management
		"POST:/api/ballots":         {models.ActionBallotCreate, models.ModuleBalloting, "Ballot"},
		"PUT:/api/ballots/":         {models.ActionBallotUpdate, models.ModuleBalloting, "Ballot"},
		"DELETE:/api/ballots/":      {models.ActionBallotClose, models.ModuleBalloting, "Ballot"},
		"POST:/api/ballots/*/vote":  {models.ActionVoteSubmit, models.ModuleBalloting, "Vote"},
		"POST:/api/ballots/*/submit": {models.ActionBallotSubmit, models.ModuleBalloting, "Ballot"},

		// Meeting Management
		"POST:/api/meetings":        {models.ActionMeetingCreate, models.ModuleMeetings, "Meeting"},
		"PUT:/api/meetings/":        {models.ActionMeetingUpdate, models.ModuleMeetings, "Meeting"},
		"DELETE:/api/meetings/":     {models.ActionMeetingDelete, models.ModuleMeetings, "Meeting"},
		"POST:/api/meetings/*/join":  {models.ActionMeetingStart, models.ModuleMeetings, "Meeting"},

		// Permission Management
		"POST:/api/permissions":     {models.ActionPermissionGrant, models.ModulePermissions, "Permission"},
		"PUT:/api/permissions/":     {models.ActionPermissionGrant, models.ModulePermissions, "Permission"},
		"DELETE:/api/permissions/":  {models.ActionPermissionRevoke, models.ModulePermissions, "Permission"},
		"POST:/api/users/*/roles":   {models.ActionRoleAssign, models.ModulePermissions, "UserRole"},
		"DELETE:/api/users/*/roles": {models.ActionRoleRevoke, models.ModulePermissions, "UserRole"},

		// Report Management
		"POST:/api/reports":         {models.ActionStandardCreate, models.ModuleReports, "Report"},
		"GET:/api/reports/*/export": {models.ActionStandardCreate, models.ModuleReports, "Report"},

		// Notification Management
		"POST:/api/notifications":   {models.ActionNotificationSend, models.ModuleNotifications, "Notification"},
		"PUT:/api/notifications/*/read": {models.ActionNotificationRead, models.ModuleNotifications, "Notification"},

		// Comment Management
		"POST:/api/comments":        {models.ActionCommentCreate, models.ModuleComments, "Comment"},
		"PUT:/api/comments/":        {models.ActionCommentUpdate, models.ModuleComments, "Comment"},
		"DELETE:/api/comments/":     {models.ActionCommentDelete, models.ModuleComments, "Comment"},

		// Public Feedback
		"POST:/api/public-feedback": {models.ActionFeedbackSubmit, models.ModuleComments, "Feedback"},
	}

	// Try exact match first
	key := method + ":" + path
	if info, exists := auditRules[key]; exists {
		return &info
	}

	// Try pattern matching for paths with IDs
	for pattern, info := range auditRules {
		if m.matchesPattern(key, pattern) {
			return &info
		}
	}

	return nil
}

// matchesPattern checks if a request matches a pattern with wildcards
func (m *AuditLogMiddleware) matchesPattern(request, pattern string) bool {
	// Simple pattern matching for paths with * wildcards
	if !strings.Contains(pattern, "*") {
		return request == pattern
	}

	requestParts := strings.Split(request, "/")
	patternParts := strings.Split(pattern, "/")

	if len(requestParts) != len(patternParts) {
		return false
	}

	for i, patternPart := range patternParts {
		if patternPart == "*" {
			continue // Wildcard matches anything
		}
		if requestParts[i] != patternPart {
			return false
		}
	}

	return true
}

// captureRequestBody reads and restores the request body
func (m *AuditLogMiddleware) captureRequestBody(c *gin.Context) string {
	if c.Request.Body == nil {
		return ""
	}

	// Read body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}

	// Restore body for further processing
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Return body as string (limit size for audit)
	body := string(bodyBytes)
	if len(body) > 10000 {
		body = body[:10000] + "... (truncated)"
	}

	return body
}

// extractResourceInfo extracts resource information from request/response
func (m *AuditLogMiddleware) extractResourceInfo(c *gin.Context, requestBody, responseBody string) ResourceInfo {
	info := ResourceInfo{}

	// Extract resource ID from URL path
	pathParts := strings.Split(c.Request.URL.Path, "/")
	for _, part := range pathParts {
		// Look for UUID patterns in path
		if _, err := uuid.Parse(part); err == nil {
			info.ResourceID = &part
			break
		}
	}

	// Extract resource information from request body
	if requestBody != "" {
		var requestData map[string]interface{}
		if err := json.Unmarshal([]byte(requestBody), &requestData); err == nil {
			// Extract title/name fields
			if title, ok := requestData["title"].(string); ok {
				info.ResourceTitle = title
			} else if name, ok := requestData["name"].(string); ok {
				info.ResourceTitle = name
			}

			// For updates, store new values
			if c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
				info.NewValues = requestBody
			}
		}
	}

	// Extract resource information from response body
	if responseBody != "" && (c.Request.Method == "POST" || c.Request.Method == "PUT") {
		var responseData map[string]interface{}
		if err := json.Unmarshal([]byte(responseBody), &responseData); err == nil {
			// Extract ID from response if not found in URL
			if info.ResourceID == nil {
				if data, ok := responseData["data"].(map[string]interface{}); ok {
					if id, ok := data["id"].(string); ok {
						info.ResourceID = &id
					}
				}
			}

			// Extract title if not found in request
			if info.ResourceTitle == "" {
				if data, ok := responseData["data"].(map[string]interface{}); ok {
					if title, ok := data["title"].(string); ok {
						info.ResourceTitle = title
					} else if name, ok := data["name"].(string); ok {
						info.ResourceTitle = name
					}
				}
			}
		}
	}

	return info
}

// generateDescription creates a human-readable description of the action
func (m *AuditLogMiddleware) generateDescription(auditInfo *AuditInfo, resourceInfo ResourceInfo, c *gin.Context) string {
	action := strings.ToLower(string(auditInfo.Action))
	resourceType := strings.ToLower(auditInfo.ResourceType)

	if resourceInfo.ResourceTitle != "" {
		return fmt.Sprintf("%s %s '%s'", action, resourceType, resourceInfo.ResourceTitle)
	}

	if resourceInfo.ResourceID != nil {
		return fmt.Sprintf("%s %s with ID %s", action, resourceType, *resourceInfo.ResourceID)
	}

	return fmt.Sprintf("%s %s", action, resourceType)
}

// buildMetadata creates metadata for the audit log
func (m *AuditLogMiddleware) buildMetadata(c *gin.Context, requestBody string, auditInfo *AuditInfo) string {
	metadata := map[string]interface{}{
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"query":       c.Request.URL.RawQuery,
		"content_type": c.GetHeader("Content-Type"),
		"referer":     c.GetHeader("Referer"),
	}

	// Add request body size
	if requestBody != "" {
		metadata["request_body_size"] = len(requestBody)
	}

	// Add any additional context
	if requestID := c.GetString("request_id"); requestID != "" {
		metadata["request_id"] = requestID
	}

	metadataBytes, _ := json.Marshal(metadata)
	return string(metadataBytes)
}

// extractErrorMessage extracts error message from response
func (m *AuditLogMiddleware) extractErrorMessage(responseBody string, statusCode int) string {
	if statusCode >= 200 && statusCode < 400 {
		return ""
	}

	if responseBody == "" {
		return fmt.Sprintf("HTTP %d", statusCode)
	}

	// Try to extract error message from JSON response
	var responseData map[string]interface{}
	if err := json.Unmarshal([]byte(responseBody), &responseData); err == nil {
		if message, ok := responseData["message"].(string); ok {
			return message
		}
		if error, ok := responseData["error"].(string); ok {
			return error
		}
	}

	// Return truncated response body
	if len(responseBody) > 500 {
		return responseBody[:500] + "... (truncated)"
	}

	return responseBody
}