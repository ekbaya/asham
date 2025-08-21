package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ActionType represents the type of action performed
type ActionType string

const (
	// Document actions
	ActionDocumentCreate   ActionType = "DOCUMENT_CREATE"
	ActionDocumentUpdate   ActionType = "DOCUMENT_UPDATE"
	ActionDocumentDelete   ActionType = "DOCUMENT_DELETE"
	ActionDocumentDownload ActionType = "DOCUMENT_DOWNLOAD"
	ActionDocumentUpload   ActionType = "DOCUMENT_UPLOAD"

	// Project actions
	ActionProjectCreate      ActionType = "PROJECT_CREATE"
	ActionProjectUpdate      ActionType = "PROJECT_UPDATE"
	ActionProjectDelete      ActionType = "PROJECT_DELETE"
	ActionProjectStageChange ActionType = "PROJECT_STAGE_CHANGE"
	ActionProjectApprove     ActionType = "PROJECT_APPROVE"
	ActionProjectReject      ActionType = "PROJECT_REJECT"
	ActionProjectCancel      ActionType = "PROJECT_CANCEL"
	ActionProjectPublish     ActionType = "PROJECT_PUBLISH"

	// Ballot actions
	ActionBallotCreate   ActionType = "BALLOT_CREATE"
	ActionBallotUpdate   ActionType = "BALLOT_UPDATE"
	ActionBallotSubmit   ActionType = "BALLOT_SUBMIT"
	ActionBallotClose    ActionType = "BALLOT_CLOSE"
	ActionBallotResults  ActionType = "BALLOT_RESULTS"
	ActionVoteSubmit     ActionType = "VOTE_SUBMIT"
	ActionVoteUpdate     ActionType = "VOTE_UPDATE"

	// User and role actions
	ActionUserCreate       ActionType = "USER_CREATE"
	ActionUserUpdate       ActionType = "USER_UPDATE"
	ActionUserDelete       ActionType = "USER_DELETE"
	ActionUserLogin        ActionType = "USER_LOGIN"
	ActionUserLogout       ActionType = "USER_LOGOUT"
	ActionRoleAssign       ActionType = "ROLE_ASSIGN"
	ActionRoleRevoke       ActionType = "ROLE_REVOKE"
	ActionPermissionGrant  ActionType = "PERMISSION_GRANT"
	ActionPermissionRevoke ActionType = "PERMISSION_REVOKE"

	// Meeting actions
	ActionMeetingCreate ActionType = "MEETING_CREATE"
	ActionMeetingUpdate ActionType = "MEETING_UPDATE"
	ActionMeetingDelete ActionType = "MEETING_DELETE"
	ActionMeetingStart  ActionType = "MEETING_START"
	ActionMeetingEnd    ActionType = "MEETING_END"

	// Comment and feedback actions
	ActionCommentCreate ActionType = "COMMENT_CREATE"
	ActionCommentUpdate ActionType = "COMMENT_UPDATE"
	ActionCommentDelete ActionType = "COMMENT_DELETE"
	ActionFeedbackSubmit ActionType = "FEEDBACK_SUBMIT"

	// Notification actions
	ActionNotificationSend ActionType = "NOTIFICATION_SEND"
	ActionNotificationRead ActionType = "NOTIFICATION_READ"

	// Standard actions
	ActionStandardCreate ActionType = "STANDARD_CREATE"
	ActionStandardUpdate ActionType = "STANDARD_UPDATE"
	ActionStandardDelete ActionType = "STANDARD_DELETE"

	// System actions
	ActionSystemBackup  ActionType = "SYSTEM_BACKUP"
	ActionSystemRestore ActionType = "SYSTEM_RESTORE"
	ActionConfigUpdate  ActionType = "CONFIG_UPDATE"

	// Workflow actions
	ActionWorkflowStart      ActionType = "WORKFLOW_START"
	ActionWorkflowTransition ActionType = "WORKFLOW_TRANSITION"
	ActionWorkflowComplete   ActionType = "WORKFLOW_COMPLETE"
)

// ModuleType represents the module where the action occurred
type ModuleType string

const (
	ModuleProjects      ModuleType = "PROJECTS"
	ModuleBalloting     ModuleType = "BALLOTING"
	ModuleDocuments     ModuleType = "DOCUMENTS"
	ModuleUsers         ModuleType = "USERS"
	ModulePermissions   ModuleType = "PERMISSIONS"
	ModuleMeetings      ModuleType = "MEETINGS"
	ModuleComments      ModuleType = "COMMENTS"
	ModuleNotifications ModuleType = "NOTIFICATIONS"
	ModuleStandards     ModuleType = "STANDARDS"
	ModuleReports       ModuleType = "REPORTS"
	ModuleSystem        ModuleType = "SYSTEM"
	ModuleWorkflow      ModuleType = "WORKFLOW"
	ModuleLibrary       ModuleType = "LIBRARY"
	ModuleOrganization  ModuleType = "ORGANIZATION"
)

// AuditLog represents a comprehensive audit log entry
type AuditLog struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID           *string        `json:"user_id" gorm:"index"` // Can be null for system actions
	User             *Member        `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Action           ActionType     `json:"action" gorm:"index;not null"`
	Module           ModuleType     `json:"module" gorm:"index;not null"`
	ResourceType     string         `json:"resource_type" gorm:"index"` // e.g., "Project", "Document", "User"
	ResourceID       *string        `json:"resource_id" gorm:"index"`   // ID of the affected resource
	ResourceTitle    string         `json:"resource_title"`              // Human-readable title/name of the resource
	Description      string         `json:"description"`                 // Human-readable description of the action
	Metadata         string         `json:"metadata" gorm:"type:jsonb"`  // Additional metadata as JSON
	OldValues        string         `json:"old_values" gorm:"type:jsonb"` // Previous values for update operations
	NewValues        string         `json:"new_values" gorm:"type:jsonb"` // New values for update operations
	IPAddress        string         `json:"ip_address"`                  // IP address of the user
	UserAgent        string         `json:"user_agent"`                  // User agent string
	SessionID        string         `json:"session_id"`                  // Session identifier
	RequestID        string         `json:"request_id"`                  // Request correlation ID
	Success          bool           `json:"success" gorm:"default:true"` // Whether the action was successful
	ErrorMessage     string         `json:"error_message"`               // Error message if action failed
	Duration         int64          `json:"duration"`                    // Duration of the action in milliseconds
	CreatedAt        time.Time      `json:"created_at" gorm:"index"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete for data retention
}

// AuditLogFilter represents filters for querying audit logs
type AuditLogFilter struct {
	UserID       *string     `json:"user_id"`
	Action       *ActionType `json:"action"`
	Module       *ModuleType `json:"module"`
	ResourceType *string     `json:"resource_type"`
	ResourceID   *string     `json:"resource_id"`
	Success      *bool       `json:"success"`
	DateFrom     *time.Time  `json:"date_from"`
	DateTo       *time.Time  `json:"date_to"`
	IPAddress    *string     `json:"ip_address"`
	Limit        int         `json:"limit"`
	Offset       int         `json:"offset"`
	OrderBy      string      `json:"order_by"` // Default: "created_at DESC"
}

// AuditLogSummary represents aggregated audit log data
type AuditLogSummary struct {
	Module      ModuleType `json:"module"`
	Action      ActionType `json:"action"`
	Count       int64      `json:"count"`
	SuccessRate float64    `json:"success_rate"`
	LastOccurred time.Time `json:"last_occurred"`
}

// AuditLogExportRequest represents a request to export audit logs
type AuditLogExportRequest struct {
	Filter AuditLogFilter `json:"filter"`
	Format string         `json:"format"` // "csv", "excel", "pdf"
	Fields []string       `json:"fields"` // Fields to include in export
}

// BeforeCreate sets default values before creating an audit log
func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}
	return nil
}

// TableName returns the table name for the AuditLog model
func (AuditLog) TableName() string {
	return "audit_logs"
}