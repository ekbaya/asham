package repository

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ekbaya/asham/pkg/domain/models"
)

// NotificationRepository handles database operations for notifications
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// CreateNotification creates a new notification
func (r *NotificationRepository) CreateNotification(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

// CreateNotifications creates multiple notifications in batch
func (r *NotificationRepository) CreateNotifications(notifications []*models.Notification) error {
	return r.db.CreateInBatches(notifications, 100).Error
}

// GetNotificationByID retrieves a notification by ID
func (r *NotificationRepository) GetNotificationByID(id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.Preload("Recipient").Preload("Project").Preload("Meeting").
		Preload("Document").Preload("Balloting").Preload("CreatedBy").
		First(&notification, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// GetNotificationsByMemberID retrieves notifications for a specific member
func (r *NotificationRepository) GetNotificationsByMemberID(memberID string, limit, offset int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	query := r.db.Where("recipient_id = ?", memberID).
		Preload("Project").Preload("Meeting").Preload("Document").Preload("Balloting").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&notifications).Error
	return notifications, err
}

// GetUnreadNotificationsByMemberID retrieves unread notifications for a member
func (r *NotificationRepository) GetUnreadNotificationsByMemberID(memberID string) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := r.db.Where("recipient_id = ? AND read = ?", memberID, false).
		Preload("Project").Preload("Meeting").Preload("Document").Preload("Balloting").
		Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

// GetCriticalUnreadNotificationsByMemberID retrieves critical unread notifications
func (r *NotificationRepository) GetCriticalUnreadNotificationsByMemberID(memberID string) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := r.db.Where("recipient_id = ? AND read = ? AND priority = ?", 
		memberID, false, models.NotificationPriorityCritical).
		Preload("Project").Preload("Meeting").Preload("Document").Preload("Balloting").
		Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

// SearchNotifications searches notifications based on criteria
func (r *NotificationRepository) SearchNotifications(req *models.NotificationSearchRequest) ([]*models.Notification, int64, error) {
	var notifications []*models.Notification
	var total int64

	query := r.db.Model(&models.Notification{})

	// Apply filters
	if req.MemberID != "" {
		query = query.Where("recipient_id = ?", req.MemberID)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Priority != "" {
		query = query.Where("priority = ?", req.Priority)
	}
	if req.Read != nil {
		query = query.Where("read = ?", *req.Read)
	}
	if req.ProjectID != nil {
		query = query.Where("project_id = ?", *req.ProjectID)
	}
	if req.DateFrom != nil {
		query = query.Where("created_at >= ?", *req.DateFrom)
	}
	if req.DateTo != nil {
		query = query.Where("created_at <= ?", *req.DateTo)
	}
	if req.Keyword != "" {
		keyword := "%" + strings.ToLower(req.Keyword) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(message) LIKE ?", keyword, keyword)
	}

	// Count total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and get results
	query = query.Preload("Project").Preload("Meeting").Preload("Document").Preload("Balloting").
		Order("created_at DESC")

	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}
	if req.Offset > 0 {
		query = query.Offset(req.Offset)
	}

	err = query.Find(&notifications).Error
	return notifications, total, err
}

// MarkNotificationAsRead marks a notification as read
func (r *NotificationRepository) MarkNotificationAsRead(id uuid.UUID, memberID string) error {
	now := time.Now()
	return r.db.Model(&models.Notification{}).
		Where("id = ? AND recipient_id = ?", id, memberID).
		Updates(map[string]interface{}{
			"read":    true,
			"read_at": now,
		}).Error
}

// MarkAllNotificationsAsRead marks all notifications as read for a member
func (r *NotificationRepository) MarkAllNotificationsAsRead(memberID string) error {
	now := time.Now()
	return r.db.Model(&models.Notification{}).
		Where("recipient_id = ? AND read = ?", memberID, false).
		Updates(map[string]interface{}{
			"read":    true,
			"read_at": now,
		}).Error
}

// MarkNotificationEmailSent marks a notification email as sent
func (r *NotificationRepository) MarkNotificationEmailSent(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"email_sent":    true,
			"email_sent_at": now,
		}).Error
}

// GetNotificationCounts returns notification counts for a member
func (r *NotificationRepository) GetNotificationCounts(memberID string) (map[string]int64, error) {
	counts := make(map[string]int64)

	// Total notifications
	var total int64
	err := r.db.Model(&models.Notification{}).Where("recipient_id = ?", memberID).Count(&total).Error
	if err != nil {
		return nil, err
	}
	counts["total"] = total

	// Unread notifications
	var unread int64
	err = r.db.Model(&models.Notification{}).Where("recipient_id = ? AND read = ?", memberID, false).Count(&unread).Error
	if err != nil {
		return nil, err
	}
	counts["unread"] = unread

	// Critical unread notifications
	var critical int64
	err = r.db.Model(&models.Notification{}).Where("recipient_id = ? AND read = ? AND priority = ?", 
		memberID, false, models.NotificationPriorityCritical).Count(&critical).Error
	if err != nil {
		return nil, err
	}
	counts["critical"] = critical

	return counts, nil
}

// GetNotificationsByType returns notifications grouped by type for a member
func (r *NotificationRepository) GetNotificationsByType(memberID string) (map[models.NotificationType]int64, error) {
	type Result struct {
		Type  models.NotificationType `json:"type"`
		Count int64                   `json:"count"`
	}

	var results []Result
	err := r.db.Model(&models.Notification{}).
		Select("type, COUNT(*) as count").
		Where("recipient_id = ?", memberID).
		Group("type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[models.NotificationType]int64)
	for _, result := range results {
		counts[result.Type] = result.Count
	}

	return counts, nil
}

// DeleteNotification deletes a notification
func (r *NotificationRepository) DeleteNotification(id uuid.UUID, memberID string) error {
	return r.db.Where("id = ? AND recipient_id = ?", id, memberID).Delete(&models.Notification{}).Error
}

// DeleteOldNotifications deletes notifications older than specified days
func (r *NotificationRepository) DeleteOldNotifications(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return r.db.Where("created_at < ?", cutoff).Delete(&models.Notification{}).Error
}

// GetPendingEmailNotifications retrieves notifications that need email sending
func (r *NotificationRepository) GetPendingEmailNotifications(limit int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	query := r.db.Where("email_sent = ? AND (channel = ? OR channel = ?)", 
		false, models.NotificationChannelEmail, models.NotificationChannelBoth).
		Preload("Recipient").Preload("Project").Preload("Meeting").
		Preload("Document").Preload("Balloting").
		Order("created_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&notifications).Error
	return notifications, err
}

// NotificationTemplate Repository Methods

// CreateNotificationTemplate creates a new notification template
func (r *NotificationRepository) CreateNotificationTemplate(template *models.NotificationTemplate) error {
	return r.db.Create(template).Error
}

// GetNotificationTemplateByType retrieves a template by type
func (r *NotificationRepository) GetNotificationTemplateByType(notificationType models.NotificationType) (*models.NotificationTemplate, error) {
	var template models.NotificationTemplate
	err := r.db.Where("type = ? AND active = ?", notificationType, true).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetAllNotificationTemplates retrieves all notification templates
func (r *NotificationRepository) GetAllNotificationTemplates() ([]*models.NotificationTemplate, error) {
	var templates []*models.NotificationTemplate
	err := r.db.Order("type ASC").Find(&templates).Error
	return templates, err
}

// UpdateNotificationTemplate updates a notification template
func (r *NotificationRepository) UpdateNotificationTemplate(template *models.NotificationTemplate) error {
	return r.db.Save(template).Error
}

// DeleteNotificationTemplate deletes a notification template
func (r *NotificationRepository) DeleteNotificationTemplate(id uuid.UUID) error {
	return r.db.Delete(&models.NotificationTemplate{}, "id = ?", id).Error
}

// NotificationPreference Repository Methods

// CreateNotificationPreference creates notification preferences for a member
func (r *NotificationRepository) CreateNotificationPreference(preference *models.NotificationPreference) error {
	return r.db.Create(preference).Error
}

// GetNotificationPreferenceByMemberID retrieves preferences for a member
func (r *NotificationRepository) GetNotificationPreferenceByMemberID(memberID string) (*models.NotificationPreference, error) {
	var preference models.NotificationPreference
	err := r.db.Where("member_id = ?", memberID).First(&preference).Error
	if err != nil {
		return nil, err
	}
	return &preference, nil
}

// UpdateNotificationPreference updates notification preferences
func (r *NotificationRepository) UpdateNotificationPreference(preference *models.NotificationPreference) error {
	return r.db.Save(preference).Error
}

// GetOrCreateNotificationPreference gets or creates default preferences for a member
func (r *NotificationRepository) GetOrCreateNotificationPreference(memberID string) (*models.NotificationPreference, error) {
	preference, err := r.GetNotificationPreferenceByMemberID(memberID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create default preferences
			preference = &models.NotificationPreference{
				MemberID:              memberID,
				Language:              models.NotificationLanguageEnglish,
				EmailNotifications:    true,
				InAppNotifications:    true,
				ProjectNotifications:  true,
				BallotNotifications:   true,
				DocumentNotifications: true,
				MeetingNotifications:  true,
				CommentNotifications:  true,
				SystemNotifications:   true,
				DeadlineReminders:     true,
			}
			err = r.CreateNotificationPreference(preference)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return preference, nil
}

// NotificationHistory Repository Methods

// CreateNotificationHistory creates a notification history record
func (r *NotificationRepository) CreateNotificationHistory(history *models.NotificationHistory) error {
	return r.db.Create(history).Error
}

// GetNotificationHistory retrieves notification history for a member
func (r *NotificationRepository) GetNotificationHistory(memberID string, limit, offset int) ([]*models.NotificationHistory, error) {
	var history []*models.NotificationHistory
	query := r.db.Where("recipient_id = ?", memberID).Order("sent_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&history).Error
	return history, err
}

// GetNotificationDashboard retrieves dashboard data for a member
func (r *NotificationRepository) GetNotificationDashboard(memberID string) (*models.NotificationDashboard, error) {
	dashboard := &models.NotificationDashboard{}

	// Get counts
	counts, err := r.GetNotificationCounts(memberID)
	if err != nil {
		return nil, err
	}
	dashboard.UnreadCount = counts["unread"]
	dashboard.CriticalCount = counts["critical"]

	// Get recent notifications (last 10)
	recentNotifications, err := r.GetNotificationsByMemberID(memberID, 10, 0)
	if err != nil {
		return nil, err
	}
	// Convert to slice of models.Notification
	dashboard.RecentNotifications = make([]models.Notification, len(recentNotifications))
	for i, n := range recentNotifications {
		dashboard.RecentNotifications[i] = *n
	}

	// Get notifications by type
	notificationsByType, err := r.GetNotificationsByType(memberID)
	if err != nil {
		return nil, err
	}
	dashboard.NotificationsByType = notificationsByType

	return dashboard, nil
}

// Helper method to convert data map to JSON string
func (r *NotificationRepository) dataToJSON(data map[string]interface{}) string {
	if data == nil {
		return ""
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

// Helper method to convert JSON string to data map
func (r *NotificationRepository) jsonToData(jsonStr string) map[string]interface{} {
	if jsonStr == "" {
		return nil
	}
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil
	}
	return data
}