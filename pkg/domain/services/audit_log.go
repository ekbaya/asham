package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

// AuditLogService provides business logic for audit logging
type AuditLogService struct {
	repo       *repository.AuditLogRepository
	memberRepo *repository.MemberRepository
}

// NewAuditLogService creates a new audit log service instance
func NewAuditLogService(repo *repository.AuditLogRepository, memberRepo *repository.MemberRepository) *AuditLogService {
	return &AuditLogService{
		repo:       repo,
		memberRepo: memberRepo,
	}
}

// LogAction creates a new audit log entry for an action
func (s *AuditLogService) LogAction(params LogActionParams) error {
	auditLog := &models.AuditLog{
		ID:            uuid.New(),
		UserID:        params.UserID,
		Action:        params.Action,
		Module:        params.Module,
		ResourceType:  params.ResourceType,
		ResourceID:    params.ResourceID,
		ResourceTitle: params.ResourceTitle,
		Description:   params.Description,
		IPAddress:     params.IPAddress,
		UserAgent:     params.UserAgent,
		SessionID:     params.SessionID,
		RequestID:     params.RequestID,
		Success:       params.Success,
		ErrorMessage:  params.ErrorMessage,
		Duration:      params.Duration,
		CreatedAt:     time.Now(),
	}

	// Serialize metadata if provided
	if params.Metadata != nil {
		metadataJSON, err := json.Marshal(params.Metadata)
		if err == nil {
			auditLog.Metadata = string(metadataJSON)
		}
	}

	// Serialize old values if provided
	if params.OldValues != nil {
		oldValuesJSON, err := json.Marshal(params.OldValues)
		if err == nil {
			auditLog.OldValues = string(oldValuesJSON)
		}
	}

	// Serialize new values if provided
	if params.NewValues != nil {
		newValuesJSON, err := json.Marshal(params.NewValues)
		if err == nil {
			auditLog.NewValues = string(newValuesJSON)
		}
	}

	return s.repo.Create(auditLog)
}

// LogActionParams represents parameters for logging an action
type LogActionParams struct {
	UserID        *string                `json:"user_id"`
	Action        models.ActionType      `json:"action"`
	Module        models.ModuleType      `json:"module"`
	ResourceType  string                 `json:"resource_type"`
	ResourceID    *string                `json:"resource_id"`
	ResourceTitle string                 `json:"resource_title"`
	Description   string                 `json:"description"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	OldValues     map[string]interface{} `json:"old_values,omitempty"`
	NewValues     map[string]interface{} `json:"new_values,omitempty"`
	IPAddress     string                 `json:"ip_address"`
	UserAgent     string                 `json:"user_agent"`
	SessionID     string                 `json:"session_id"`
	RequestID     string                 `json:"request_id"`
	Success       bool                   `json:"success"`
	ErrorMessage  string                 `json:"error_message,omitempty"`
	Duration      int64                  `json:"duration"`
}

// LogBatchActions creates multiple audit log entries in a single transaction
func (s *AuditLogService) LogBatchActions(params []LogActionParams) error {
	auditLogs := make([]models.AuditLog, len(params))

	for i, param := range params {
		auditLog := models.AuditLog{
			ID:            uuid.New(),
			UserID:        param.UserID,
			Action:        param.Action,
			Module:        param.Module,
			ResourceType:  param.ResourceType,
			ResourceID:    param.ResourceID,
			ResourceTitle: param.ResourceTitle,
			Description:   param.Description,
			IPAddress:     param.IPAddress,
			UserAgent:     param.UserAgent,
			SessionID:     param.SessionID,
			RequestID:     param.RequestID,
			Success:       param.Success,
			ErrorMessage:  param.ErrorMessage,
			Duration:      param.Duration,
			CreatedAt:     time.Now(),
		}

		// Serialize metadata if provided
		if param.Metadata != nil {
			metadataJSON, err := json.Marshal(param.Metadata)
			if err == nil {
				auditLog.Metadata = string(metadataJSON)
			}
		}

		// Serialize old values if provided
		if param.OldValues != nil {
			oldValuesJSON, err := json.Marshal(param.OldValues)
			if err == nil {
				auditLog.OldValues = string(oldValuesJSON)
			}
		}

		// Serialize new values if provided
		if param.NewValues != nil {
			newValuesJSON, err := json.Marshal(param.NewValues)
			if err == nil {
				auditLog.NewValues = string(newValuesJSON)
			}
		}

		auditLogs[i] = auditLog
	}

	return s.repo.CreateBatch(auditLogs)
}

// GetAuditLogs retrieves audit logs with filtering and pagination
func (s *AuditLogService) GetAuditLogs(filter models.AuditLogFilter) ([]models.AuditLog, int64, error) {
	return s.repo.List(filter)
}

// GetAuditLogByID retrieves a specific audit log by ID
func (s *AuditLogService) GetAuditLogByID(id uuid.UUID) (*models.AuditLog, error) {
	return s.repo.GetByID(id)
}

// GetResourceAuditTrail retrieves the audit trail for a specific resource
func (s *AuditLogService) GetResourceAuditTrail(resourceType, resourceID string, limit int) ([]models.AuditLog, error) {
	return s.repo.GetByResourceID(resourceType, resourceID, limit)
}

// GetUserActivity retrieves audit logs for a specific user
func (s *AuditLogService) GetUserActivity(userID string, limit int) ([]models.AuditLog, error) {
	return s.repo.GetByUserID(userID, limit)
}

// GetAuditSummary generates audit log summary statistics
func (s *AuditLogService) GetAuditSummary(filter models.AuditLogFilter) ([]models.AuditLogSummary, error) {
	return s.repo.GetSummary(filter)
}

// GetActivityTimeline retrieves audit logs grouped by time periods
func (s *AuditLogService) GetActivityTimeline(filter models.AuditLogFilter, interval string) (map[string]int64, error) {
	return s.repo.GetActivityTimeline(filter, interval)
}

// GetFailedActions retrieves audit logs for failed actions
func (s *AuditLogService) GetFailedActions(filter models.AuditLogFilter) ([]models.AuditLog, error) {
	return s.repo.GetFailedActions(filter)
}

// SearchAuditLogs performs a text search across audit logs
func (s *AuditLogService) SearchAuditLogs(searchTerm string, filter models.AuditLogFilter) ([]models.AuditLog, int64, error) {
	return s.repo.SearchLogs(searchTerm, filter)
}

// GetComplianceReport generates a comprehensive compliance report
func (s *AuditLogService) GetComplianceReport(filter models.AuditLogFilter) (map[string]interface{}, error) {
	return s.repo.GetComplianceReport(filter)
}

// GetUserActivityStats retrieves user activity statistics
func (s *AuditLogService) GetUserActivityStats(filter models.AuditLogFilter) (map[string]int64, error) {
	return s.repo.GetUserActivity(filter)
}

// GetModuleActivityStats retrieves module activity statistics
func (s *AuditLogService) GetModuleActivityStats(filter models.AuditLogFilter) (map[string]int64, error) {
	return s.repo.GetModuleActivity(filter)
}

// ValidateUserAccess checks if a user has permission to access audit logs
func (s *AuditLogService) ValidateUserAccess(userID string, requestedFilter models.AuditLogFilter) error {
	// Get user information
	user, err := s.memberRepo.GetMemberByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user has audit access roles
	hasAuditAccess := false
	for _, role := range user.Roles {
		if s.isAuditRole(role.Title) {
			hasAuditAccess = true
			break
		}
	}

	if !hasAuditAccess {
		return fmt.Errorf("user does not have permission to access audit logs")
	}

	// Additional access controls based on role
	for _, role := range user.Roles {
		switch role.Title {
		case "ARSO_SECRETARIAT", "SMC_MEMBER":
			// Full access - no additional restrictions
			return nil
		case "TC_SECRETARY", "TC_MEMBER":
			// Limited access - can only view logs related to their committee
			if requestedFilter.Module != nil && *requestedFilter.Module != models.ModuleProjects {
				return fmt.Errorf("insufficient permissions to access logs from this module")
			}
		case "AUDITOR":
			// Auditor role - full read access but no modification
			return nil
		default:
			// Regular users can only view their own activity
			if requestedFilter.UserID == nil || *requestedFilter.UserID != userID {
				return fmt.Errorf("insufficient permissions to access other users' audit logs")
			}
		}
	}

	return nil
}

// isAuditRole checks if a role has audit access permissions
func (s *AuditLogService) isAuditRole(roleTitle string) bool {
	auditRoles := []string{
		"ARSO_SECRETARIAT",
		"SMC_MEMBER",
		"TC_SECRETARY",
		"TC_MEMBER",
		"AUDITOR",
		"ADMIN",
	}

	for _, role := range auditRoles {
		if strings.EqualFold(roleTitle, role) {
			return true
		}
	}
	return false
}

// CleanupOldLogs removes audit logs older than the specified retention period
func (s *AuditLogService) CleanupOldLogs(retentionPeriod time.Duration) (int64, error) {
	return s.repo.DeleteOldLogs(retentionPeriod)
}

// LogProjectAction logs project-related actions
func (s *AuditLogService) LogProjectAction(userID *string, action models.ActionType, projectID, projectTitle string, metadata map[string]interface{}, success bool, errorMsg string, duration int64, ipAddress, userAgent, sessionID, requestID string) error {
	description := s.generateProjectActionDescription(action, projectTitle)

	return s.LogAction(LogActionParams{
		UserID:        userID,
		Action:        action,
		Module:        models.ModuleProjects,
		ResourceType:  "Project",
		ResourceID:    &projectID,
		ResourceTitle: projectTitle,
		Description:   description,
		Metadata:      metadata,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		SessionID:     sessionID,
		RequestID:     requestID,
		Success:       success,
		ErrorMessage:  errorMsg,
		Duration:      duration,
	})
}

// LogBallotAction logs ballot-related actions
func (s *AuditLogService) LogBallotAction(userID *string, action models.ActionType, ballotID, ballotTitle string, metadata map[string]interface{}, success bool, errorMsg string, duration int64, ipAddress, userAgent, sessionID, requestID string) error {
	description := s.generateBallotActionDescription(action, ballotTitle)

	return s.LogAction(LogActionParams{
		UserID:        userID,
		Action:        action,
		Module:        models.ModuleBalloting,
		ResourceType:  "Ballot",
		ResourceID:    &ballotID,
		ResourceTitle: ballotTitle,
		Description:   description,
		Metadata:      metadata,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		SessionID:     sessionID,
		RequestID:     requestID,
		Success:       success,
		ErrorMessage:  errorMsg,
		Duration:      duration,
	})
}

// LogDocumentAction logs document-related actions
func (s *AuditLogService) LogDocumentAction(userID *string, action models.ActionType, documentID, documentTitle string, metadata map[string]interface{}, success bool, errorMsg string, duration int64, ipAddress, userAgent, sessionID, requestID string) error {
	description := s.generateDocumentActionDescription(action, documentTitle)

	return s.LogAction(LogActionParams{
		UserID:        userID,
		Action:        action,
		Module:        models.ModuleDocuments,
		ResourceType:  "Document",
		ResourceID:    &documentID,
		ResourceTitle: documentTitle,
		Description:   description,
		Metadata:      metadata,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		SessionID:     sessionID,
		RequestID:     requestID,
		Success:       success,
		ErrorMessage:  errorMsg,
		Duration:      duration,
	})
}

// LogUserAction logs user-related actions
func (s *AuditLogService) LogUserAction(userID *string, action models.ActionType, targetUserID, targetUserName string, metadata map[string]interface{}, success bool, errorMsg string, duration int64, ipAddress, userAgent, sessionID, requestID string) error {
	description := s.generateUserActionDescription(action, targetUserName)

	return s.LogAction(LogActionParams{
		UserID:        userID,
		Action:        action,
		Module:        models.ModuleUsers,
		ResourceType:  "User",
		ResourceID:    &targetUserID,
		ResourceTitle: targetUserName,
		Description:   description,
		Metadata:      metadata,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		SessionID:     sessionID,
		RequestID:     requestID,
		Success:       success,
		ErrorMessage:  errorMsg,
		Duration:      duration,
	})
}

// Helper methods to generate action descriptions
func (s *AuditLogService) generateProjectActionDescription(action models.ActionType, projectTitle string) string {
	switch action {
	case models.ActionProjectCreate:
		return fmt.Sprintf("Created project: %s", projectTitle)
	case models.ActionProjectUpdate:
		return fmt.Sprintf("Updated project: %s", projectTitle)
	case models.ActionProjectDelete:
		return fmt.Sprintf("Deleted project: %s", projectTitle)
	case models.ActionProjectStageChange:
		return fmt.Sprintf("Changed stage for project: %s", projectTitle)
	case models.ActionProjectApprove:
		return fmt.Sprintf("Approved project: %s", projectTitle)
	case models.ActionProjectReject:
		return fmt.Sprintf("Rejected project: %s", projectTitle)
	case models.ActionProjectCancel:
		return fmt.Sprintf("Cancelled project: %s", projectTitle)
	case models.ActionProjectPublish:
		return fmt.Sprintf("Published project: %s", projectTitle)
	default:
		return fmt.Sprintf("Performed action %s on project: %s", action, projectTitle)
	}
}

func (s *AuditLogService) generateBallotActionDescription(action models.ActionType, ballotTitle string) string {
	switch action {
	case models.ActionBallotCreate:
		return fmt.Sprintf("Created ballot: %s", ballotTitle)
	case models.ActionBallotUpdate:
		return fmt.Sprintf("Updated ballot: %s", ballotTitle)
	case models.ActionBallotSubmit:
		return fmt.Sprintf("Submitted ballot: %s", ballotTitle)
	case models.ActionBallotClose:
		return fmt.Sprintf("Closed ballot: %s", ballotTitle)
	case models.ActionBallotResults:
		return fmt.Sprintf("Generated results for ballot: %s", ballotTitle)
	case models.ActionVoteSubmit:
		return fmt.Sprintf("Submitted vote for ballot: %s", ballotTitle)
	case models.ActionVoteUpdate:
		return fmt.Sprintf("Updated vote for ballot: %s", ballotTitle)
	default:
		return fmt.Sprintf("Performed action %s on ballot: %s", action, ballotTitle)
	}
}

func (s *AuditLogService) generateDocumentActionDescription(action models.ActionType, documentTitle string) string {
	switch action {
	case models.ActionDocumentCreate:
		return fmt.Sprintf("Created document: %s", documentTitle)
	case models.ActionDocumentUpdate:
		return fmt.Sprintf("Updated document: %s", documentTitle)
	case models.ActionDocumentDelete:
		return fmt.Sprintf("Deleted document: %s", documentTitle)
	case models.ActionDocumentDownload:
		return fmt.Sprintf("Downloaded document: %s", documentTitle)
	case models.ActionDocumentUpload:
		return fmt.Sprintf("Uploaded document: %s", documentTitle)
	default:
		return fmt.Sprintf("Performed action %s on document: %s", action, documentTitle)
	}
}

func (s *AuditLogService) generateUserActionDescription(action models.ActionType, userName string) string {
	switch action {
	case models.ActionUserCreate:
		return fmt.Sprintf("Created user: %s", userName)
	case models.ActionUserUpdate:
		return fmt.Sprintf("Updated user: %s", userName)
	case models.ActionUserDelete:
		return fmt.Sprintf("Deleted user: %s", userName)
	case models.ActionUserLogin:
		return fmt.Sprintf("User logged in: %s", userName)
	case models.ActionUserLogout:
		return fmt.Sprintf("User logged out: %s", userName)
	case models.ActionRoleAssign:
		return fmt.Sprintf("Assigned role to user: %s", userName)
	case models.ActionRoleRevoke:
		return fmt.Sprintf("Revoked role from user: %s", userName)
	case models.ActionPermissionGrant:
		return fmt.Sprintf("Granted permission to user: %s", userName)
	case models.ActionPermissionRevoke:
		return fmt.Sprintf("Revoked permission from user: %s", userName)
	default:
		return fmt.Sprintf("Performed action %s on user: %s", action, userName)
	}
}