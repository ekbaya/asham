package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationService services.NotificationService
}

func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// GetNotificationDashboard returns dashboard data for the authenticated user
func (h *NotificationHandler) GetNotificationDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)
	dashboard, err := h.notificationService.GetNotificationDashboard(userIDStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dashboard})
}

// GetNotifications returns paginated notifications for the authenticated user
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Parse filter parameters
	unreadOnly := c.Query("unread_only") == "true"
	criticalOnly := c.Query("critical_only") == "true"
	notificationType := c.Query("type")

	notifications, total, err := h.notificationService.GetNotificationsByMemberID(
		userIDStr, page, limit, unreadOnly, criticalOnly, notificationType,
	)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  notifications,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// SearchNotifications searches notifications for the authenticated user
func (h *NotificationHandler) SearchNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)

	// Parse search parameters
	query := c.Query("q")
	notificationType := c.Query("type")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Parse dates
	var startDate, endDate *time.Time
	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		}
	}

	// Create search request
	searchReq := &models.NotificationSearchRequest{
		MemberID: userIDStr,
		Keyword:  query,
		DateFrom: startDate,
		DateTo:   endDate,
		Limit:    limit,
		Offset:   (page - 1) * limit,
	}

	if notificationType != "" {
		searchReq.Type = models.NotificationType(notificationType)
	}

	notifications, total, err := h.notificationService.SearchNotifications(searchReq)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  notifications,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// MarkNotificationAsRead marks a notification as read
func (h *NotificationHandler) MarkNotificationAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)
	notificationIDStr := c.Param("id")

	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	err = h.notificationService.MarkNotificationAsRead(notificationID, userIDStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// MarkAllNotificationsAsRead marks all notifications as read for the authenticated user
func (h *NotificationHandler) MarkAllNotificationsAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)

	err := h.notificationService.MarkAllNotificationsAsRead(userIDStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

// GetNotificationPreferences returns notification preferences for the authenticated user
func (h *NotificationHandler) GetNotificationPreferences(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)
	preferences, err := h.notificationService.GetNotificationPreferences(userIDStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": preferences})
}

// UpdateNotificationPreferences updates notification preferences for the authenticated user
func (h *NotificationHandler) UpdateNotificationPreferences(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)

	var payload models.NotificationPreference
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.ShowError(c, http.StatusBadRequest, formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	payload.MemberID = userIDStr

	err := h.notificationService.UpdateNotificationPreferences(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification preferences updated successfully"})
}

// SendAdminAnnouncement sends an admin announcement (admin only)
func (h *NotificationHandler) SendAdminAnnouncement(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)

	type AnnouncementPayload struct {
		Title       string   `json:"title" binding:"required"`
		Message     string   `json:"message" binding:"required"`
		Type        string   `json:"type" binding:"required"`
		IsCritical  bool     `json:"is_critical"`
		TargetRoles []string `json:"target_roles"`
	}

	var payload AnnouncementPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.ShowError(c, http.StatusBadRequest, formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Determine priority based on critical flag
	priority := models.NotificationPriorityMedium
	if payload.IsCritical {
		priority = models.NotificationPriorityCritical
	}

	err := h.notificationService.SendAdminAnnouncement(
		payload.Title,
		payload.Message,
		payload.TargetRoles,
		priority,
		userIDStr,
	)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Announcement sent successfully"})
}

// GetNotificationHistory returns notification history for the authenticated user
func (h *NotificationHandler) GetNotificationHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	history, total, err := h.notificationService.GetNotificationHistory(userIDStr, page, limit)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  history,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// DeleteNotification deletes a notification (soft delete)
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)
	notificationIDStr := c.Param("id")

	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	err = h.notificationService.DeleteNotification(notificationID, userIDStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}