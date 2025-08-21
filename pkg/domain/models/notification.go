package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	// Project related notifications
	NotificationProjectCreated   NotificationType = "PROJECT_CREATED"
	NotificationProjectAssigned  NotificationType = "PROJECT_ASSIGNED"
	NotificationProjectUpdated   NotificationType = "PROJECT_UPDATED"

	// Ballot related notifications
	NotificationBallotOpened     NotificationType = "BALLOT_OPENED"
	NotificationBallotReminder   NotificationType = "BALLOT_REMINDER"
	NotificationBallotClosing    NotificationType = "BALLOT_CLOSING"
	NotificationBallotClosed     NotificationType = "BALLOT_CLOSED"

	// Document related notifications
	NotificationDocumentUploaded NotificationType = "DOCUMENT_UPLOADED"
	NotificationDocumentUpdated  NotificationType = "DOCUMENT_UPDATED"
	NotificationDocumentVersion  NotificationType = "DOCUMENT_VERSION"

	// Meeting related notifications
	NotificationMeetingInvitation NotificationType = "MEETING_INVITATION"
	NotificationMeetingChanged    NotificationType = "MEETING_CHANGED"
	NotificationMeetingReminder   NotificationType = "MEETING_REMINDER"
	NotificationMeetingCancelled  NotificationType = "MEETING_CANCELLED"

	// Comment related notifications
	NotificationCommentWindowOpened NotificationType = "COMMENT_WINDOW_OPENED"
	NotificationCommentWindowClosed NotificationType = "COMMENT_WINDOW_CLOSED"
	NotificationCommentReceived     NotificationType = "COMMENT_RECEIVED"

	// System notifications
	NotificationSystemUpdate      NotificationType = "SYSTEM_UPDATE"
	NotificationPolicyChange      NotificationType = "POLICY_CHANGE"
	NotificationMaintenance       NotificationType = "MAINTENANCE"
	NotificationTraining          NotificationType = "TRAINING"
	NotificationDeadlineReminder  NotificationType = "DEADLINE_REMINDER"
	NotificationTaskEscalation    NotificationType = "TASK_ESCALATION"
	NotificationAnnouncement      NotificationType = "ANNOUNCEMENT"
)

// NotificationPriority represents the priority level of a notification
type NotificationPriority string

const (
	NotificationPriorityLow      NotificationPriority = "LOW"
	NotificationPriorityMedium   NotificationPriority = "MEDIUM"
	NotificationPriorityHigh     NotificationPriority = "HIGH"
	NotificationPriorityCritical NotificationPriority = "CRITICAL"
)

// NotificationChannel represents how the notification should be delivered
type NotificationChannel string

const (
	NotificationChannelInApp NotificationChannel = "IN_APP"
	NotificationChannelEmail NotificationChannel = "EMAIL"
	NotificationChannelBoth  NotificationChannel = "BOTH"
)

// NotificationLanguage represents the language preference
type NotificationLanguage string

const (
	NotificationLanguageEnglish NotificationLanguage = "EN"
	NotificationLanguageFrench  NotificationLanguage = "FR"
)

// Notification represents a notification sent to a user
type Notification struct {
	ID          uuid.UUID             `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	RecipientID string                `json:"recipient_id" gorm:"type:uuid;not null;index"`
	Recipient   *Member               `json:"recipient,omitempty" gorm:"foreignKey:RecipientID;references:ID"`
	Type        NotificationType      `json:"type" gorm:"not null;index"`
	Priority    NotificationPriority  `json:"priority" gorm:"default:MEDIUM"`
	Channel     NotificationChannel   `json:"channel" gorm:"default:BOTH"`
	Language    NotificationLanguage  `json:"language" gorm:"default:EN"`
	Title       string                `json:"title" gorm:"not null"`
	Message     string                `json:"message" gorm:"type:text;not null"`
	Data        string                `json:"data,omitempty" gorm:"type:jsonb"` // Additional structured data
	Read        bool                  `json:"read" gorm:"default:false;index"`
	ReadAt      *time.Time            `json:"read_at,omitempty"`
	EmailSent   bool                  `json:"email_sent" gorm:"default:false"`
	EmailSentAt *time.Time            `json:"email_sent_at,omitempty"`

	// Related entities (optional foreign keys)
	ProjectID   *string   `json:"project_id,omitempty" gorm:"type:uuid"`
	Project     *Project  `json:"project,omitempty" gorm:"foreignKey:ProjectID;references:ID"`
	MeetingID   *string   `json:"meeting_id,omitempty" gorm:"type:uuid"`
	Meeting     *Meeting  `json:"meeting,omitempty" gorm:"foreignKey:MeetingID;references:ID"`
	DocumentID  *string   `json:"document_id,omitempty" gorm:"type:uuid"`
	Document    *Document `json:"document,omitempty" gorm:"foreignKey:DocumentID;references:ID"`
	BallotingID *string   `json:"balloting_id,omitempty" gorm:"type:uuid"`
	Balloting   *Balloting `json:"balloting,omitempty" gorm:"foreignKey:BallotingID;references:ID"`

	// Metadata
	CreatedByID *string   `json:"created_by_id,omitempty" gorm:"type:uuid"`
	CreatedBy   *Member   `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID;references:ID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NotificationTemplate represents a template for generating notifications
type NotificationTemplate struct {
	ID          uuid.UUID             `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Type        NotificationType      `json:"type" gorm:"unique;not null"`
	Language    NotificationLanguage  `json:"language" gorm:"not null"`
	TitleEN     string                `json:"title_en" gorm:"not null"`
	TitleFR     string                `json:"title_fr" gorm:"not null"`
	MessageEN   string                `json:"message_en" gorm:"type:text;not null"`
	MessageFR   string                `json:"message_fr" gorm:"type:text;not null"`
	EmailEN     string                `json:"email_en,omitempty" gorm:"type:text"`
	EmailFR     string                `json:"email_fr,omitempty" gorm:"type:text"`
	Priority    NotificationPriority  `json:"priority" gorm:"default:MEDIUM"`
	Channel     NotificationChannel   `json:"channel" gorm:"default:BOTH"`
	Active      bool                  `json:"active" gorm:"default:true"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// NotificationPreference represents user preferences for notifications
type NotificationPreference struct {
	ID                    uuid.UUID             `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MemberID              string                `json:"member_id" gorm:"type:uuid;not null;index"`
	Member                *Member               `json:"member,omitempty" gorm:"foreignKey:MemberID;references:ID"`
	Language              NotificationLanguage  `json:"language" gorm:"default:EN"`
	EmailNotifications    bool                  `json:"email_notifications" gorm:"default:true"`
	InAppNotifications    bool                  `json:"in_app_notifications" gorm:"default:true"`

	// Specific notification type preferences
	ProjectNotifications  bool `json:"project_notifications" gorm:"default:true"`
	BallotNotifications   bool `json:"ballot_notifications" gorm:"default:true"`
	DocumentNotifications bool `json:"document_notifications" gorm:"default:true"`
	MeetingNotifications  bool `json:"meeting_notifications" gorm:"default:true"`
	CommentNotifications  bool `json:"comment_notifications" gorm:"default:true"`
	SystemNotifications   bool `json:"system_notifications" gorm:"default:true"`
	DeadlineReminders     bool `json:"deadline_reminders" gorm:"default:true"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NotificationHistory represents a log of all sent notifications for audit purposes
type NotificationHistory struct {
	ID             uuid.UUID            `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	NotificationID uuid.UUID            `json:"notification_id" gorm:"type:uuid;not null;index"`
	RecipientID    string               `json:"recipient_id" gorm:"type:uuid;not null;index"`
	Type           NotificationType     `json:"type" gorm:"not null;index"`
	Channel        NotificationChannel  `json:"channel" gorm:"not null"`
	Status         string               `json:"status" gorm:"not null"` // sent, failed, delivered, read
	Error          string               `json:"error,omitempty"`
	SentAt         time.Time            `json:"sent_at"`
	CreatedAt      time.Time            `json:"created_at"`
}

// BeforeCreate hook for Notification
func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	if n.ID == uuid.Nil {
		n.ID, err = uuid.NewRandom()
		if err != nil {
			return err
		}
	}
	return nil
}

// BeforeCreate hook for NotificationTemplate
func (nt *NotificationTemplate) BeforeCreate(tx *gorm.DB) (err error) {
	if nt.ID == uuid.Nil {
		nt.ID, err = uuid.NewRandom()
		if err != nil {
			return err
		}
	}
	return nil
}

// BeforeCreate hook for NotificationPreference
func (np *NotificationPreference) BeforeCreate(tx *gorm.DB) (err error) {
	if np.ID == uuid.Nil {
		np.ID, err = uuid.NewRandom()
		if err != nil {
			return err
		}
	}
	return nil
}

// BeforeCreate hook for NotificationHistory
func (nh *NotificationHistory) BeforeCreate(tx *gorm.DB) (err error) {
	if nh.ID == uuid.Nil {
		nh.ID, err = uuid.NewRandom()
		if err != nil {
			return err
		}
	}
	return nil
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	n.Read = true
	now := time.Now()
	n.ReadAt = &now
}

// MarkEmailSent marks the notification email as sent
func (n *Notification) MarkEmailSent() {
	n.EmailSent = true
	now := time.Now()
	n.EmailSentAt = &now
}

// GetTitle returns the title in the specified language
func (nt *NotificationTemplate) GetTitle(language NotificationLanguage) string {
	if language == NotificationLanguageFrench {
		return nt.TitleFR
	}
	return nt.TitleEN
}

// GetMessage returns the message in the specified language
func (nt *NotificationTemplate) GetMessage(language NotificationLanguage) string {
	if language == NotificationLanguageFrench {
		return nt.MessageFR
	}
	return nt.MessageEN
}

// GetEmailTemplate returns the email template in the specified language
func (nt *NotificationTemplate) GetEmailTemplate(language NotificationLanguage) string {
	if language == NotificationLanguageFrench {
		return nt.EmailFR
	}
	return nt.EmailEN
}

// NotificationRequest represents a request to create a notification
type NotificationRequest struct {
	RecipientIDs []string              `json:"recipient_ids" binding:"required"`
	Type         NotificationType      `json:"type" binding:"required"`
	Priority     NotificationPriority  `json:"priority"`
	Channel      NotificationChannel   `json:"channel"`
	Title        string                `json:"title"`
	Message      string                `json:"message" binding:"required"`
	Data         map[string]interface{} `json:"data,omitempty"`
	ProjectID    *string               `json:"project_id,omitempty"`
	MeetingID    *string               `json:"meeting_id,omitempty"`
	DocumentID   *string               `json:"document_id,omitempty"`
	BallotingID  *string               `json:"balloting_id,omitempty"`
}

// NotificationSearchRequest represents search parameters for notifications
type NotificationSearchRequest struct {
	MemberID   string                `json:"member_id,omitempty"`
	Type       NotificationType      `json:"type,omitempty"`
	Priority   NotificationPriority  `json:"priority,omitempty"`
	Read       *bool                 `json:"read,omitempty"`
	Keyword    string                `json:"keyword,omitempty"`
	DateFrom   *time.Time            `json:"date_from,omitempty"`
	DateTo     *time.Time            `json:"date_to,omitempty"`
	ProjectID  *string               `json:"project_id,omitempty"`
	Limit      int                   `json:"limit,omitempty"`
	Offset     int                   `json:"offset,omitempty"`
}

// NotificationDashboard represents the notification dashboard data
type NotificationDashboard struct {
	UnreadCount       int64          `json:"unread_count"`
	CriticalCount     int64          `json:"critical_count"`
	RecentNotifications []Notification `json:"recent_notifications"`
	NotificationsByType map[NotificationType]int64 `json:"notifications_by_type"`
}