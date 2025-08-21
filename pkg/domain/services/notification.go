package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
)

// NotificationService handles notification business logic
type NotificationService struct {
	notificationRepo *repository.NotificationRepository
	memberRepo       *repository.MemberRepository
	rbacRepo         *repository.RbacRepository
	ballotingRepo    *repository.BallotingRepository
	emailService     *EmailService
	db               *gorm.DB
}

// NewNotificationService creates a new notification service
func NewNotificationService(
	notificationRepo *repository.NotificationRepository,
	memberRepo *repository.MemberRepository,
	rbacRepo *repository.RbacRepository,
	ballotingRepo *repository.BallotingRepository,
	emailService *EmailService,
	db *gorm.DB,
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		memberRepo:       memberRepo,
		rbacRepo:         rbacRepo,
		ballotingRepo:    ballotingRepo,
		emailService:     emailService,
		db:               db,
	}
}

// CreateNotification creates a single notification
func (s *NotificationService) CreateNotification(req *models.NotificationRequest, recipients []string) error {
	if len(recipients) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	// Get notification template if title is not provided
	var template *models.NotificationTemplate
	var err error
	if req.Title == "" {
		template, err = s.notificationRepo.GetNotificationTemplateByType(req.Type)
		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to get notification template: %w", err)
		}
	}

	// Create notifications for each recipient
	var notifications []*models.Notification
	for _, recipientID := range recipients {
		// Get recipient preferences
		preferences, err := s.notificationRepo.GetOrCreateNotificationPreference(recipientID)
		if err != nil {
			return fmt.Errorf("failed to get notification preferences for %s: %w", recipientID, err)
		}

		// Check if user wants this type of notification
		if !s.shouldSendNotification(req.Type, preferences) {
			continue
		}

		// Determine title and message
		title := req.Title
		message := req.Message
		if template != nil {
			if title == "" {
				title = template.GetTitle(preferences.Language)
			}
			message = s.processTemplate(template.GetMessage(preferences.Language), req.Data)
		}

		// Determine channel based on preferences
		channel := s.determineChannel(req.Channel, preferences)

		// Create notification
		notification := &models.Notification{
			RecipientID: recipientID,
			Type:        req.Type,
			Priority:    req.Priority,
			Channel:     channel,
			Language:    preferences.Language,
			Title:       title,
			Message:     message,
			ProjectID:   req.ProjectID,
			MeetingID:   req.MeetingID,
			DocumentID:  req.DocumentID,
			BallotingID: req.BallotingID,
			CreatedByID: nil,
		}

		// Add structured data if provided
		if req.Data != nil {
			dataJSON, _ := json.Marshal(req.Data)
			notification.Data = string(dataJSON)
		}

		notifications = append(notifications, notification)
	}

	if len(notifications) == 0 {
		return fmt.Errorf("no notifications to send after filtering preferences")
	}

	// Save notifications to database
	err = s.notificationRepo.CreateNotifications(notifications)
	if err != nil {
		return fmt.Errorf("failed to create notifications: %w", err)
	}

	// Send email notifications asynchronously
	go s.processEmailNotifications(notifications)

	// Create notification history records
	go s.createNotificationHistory(notifications)

	return nil
}

// NotifyProjectAssigned sends notification when a project is assigned
func (s *NotificationService) NotifyProjectAssigned(project *models.Project, assignedToID string, assignedByID string) error {
	notificationReq := &models.NotificationRequest{
		Type:     models.NotificationProjectAssigned,
		Priority: models.NotificationPriorityMedium,
		Channel:  models.NotificationChannelBoth,
		Data: map[string]interface{}{
			"project_title": project.Title,
			"project_id":    project.ID.String(),
			"assigned_by":   assignedByID,
		},
		ProjectID: func() *string { s := project.ID.String(); return &s }(),
	}

	return s.CreateNotification(notificationReq, []string{assignedToID})
}

// NotifyBallotClosed sends notification when a ballot is closed
func (s *NotificationService) NotifyBallotClosed(project *models.Project, balloting *models.Balloting) error {
	// Get eligible voters
	members, _, err := s.memberRepo.GetAllMembers(1, 1000)
	if err != nil {
		return err
	}

	var recipients []string
	for _, member := range *members {
		if eligible, _ := s.ballotingRepo.IsEligibleToVote(member.ID.String(), balloting.ID.String()); eligible {
			recipients = append(recipients, member.ID.String())
		}
	}

	notificationReq := &models.NotificationRequest{
		Type:     models.NotificationBallotClosed,
		Priority: models.NotificationPriorityHigh,
		Channel:  models.NotificationChannelBoth,
		Data: map[string]interface{}{
			"project_title": project.Title,
			"project_id":    project.ID.String(),
			"end_date":      balloting.EndDate,
		},
		ProjectID:   func() *string { s := project.ID.String(); return &s }(),
		BallotingID: func() *string { s := balloting.ID.String(); return &s }(),
	}

	return s.CreateNotification(notificationReq, recipients)
}

// NotifyDocumentUpdated sends notification when a document is updated
func (s *NotificationService) NotifyDocumentUpdated(project *models.Project, document *models.Document, updatedByID string) error {
	// Get project stakeholders
	members, _, err := s.memberRepo.GetAllMembers(1, 1000)
	if err != nil {
		return err
	}

	var recipients []string
	for _, member := range *members {
		recipients = append(recipients, member.ID.String())
	}

	notificationReq := &models.NotificationRequest{
		Type:     models.NotificationDocumentUpdated,
		Priority: models.NotificationPriorityMedium,
		Channel:  models.NotificationChannelBoth,
		Data: map[string]interface{}{
			"project_title":  project.Title,
			"project_id":     project.ID.String(),
			"document_title": document.Title,
			"updated_by":     updatedByID,
		},
		ProjectID:  func() *string { s := project.ID.String(); return &s }(),
		DocumentID: func() *string { s := document.ID.String(); return &s }(),
	}

	return s.CreateNotification(notificationReq, recipients)
}

// NotifyMeetingChanged sends notification when a meeting is changed
func (s *NotificationService) NotifyMeetingChanged(meeting *models.Meeting, changedByID string) error {
	var recipients []string

	// Add attendees if available
	if meeting.Attendees != nil {
		for _, attendee := range *meeting.Attendees {
			recipients = append(recipients, attendee.ID.String())
		}
	}

	notificationReq := &models.NotificationRequest{
		Type:     models.NotificationMeetingChanged,
		Priority: models.NotificationPriorityHigh,
		Channel:  models.NotificationChannelBoth,
		Data: map[string]interface{}{
			"meeting_title": meeting.Title,
			"meeting_date":  meeting.Date,
			"venue":         meeting.Venue,
			"changed_by":    changedByID,
		},
		MeetingID: func() *string { s := meeting.ID.String(); return &s }(),
	}

	return s.CreateNotification(notificationReq, recipients)
}

// NotifyMeetingCancelled sends notification when a meeting is cancelled
func (s *NotificationService) NotifyMeetingCancelled(meeting *models.Meeting, cancelledByID string) error {
	var recipients []string

	// Add attendees if available
	if meeting.Attendees != nil {
		for _, attendee := range *meeting.Attendees {
			recipients = append(recipients, attendee.ID.String())
		}
	}

	notificationReq := &models.NotificationRequest{
		Type:     models.NotificationMeetingCancelled,
		Priority: models.NotificationPriorityCritical,
		Channel:  models.NotificationChannelBoth,
		Data: map[string]interface{}{
			"meeting_title": meeting.Title,
			"meeting_date":  meeting.Date,
			"venue":         meeting.Venue,
			"cancelled_by":  cancelledByID,
		},
		MeetingID: func() *string { s := meeting.ID.String(); return &s }(),
	}

	return s.CreateNotification(notificationReq, recipients)
}

// NotifyCommentWindowOpened sends notification when a comment window opens
func (s *NotificationService) NotifyCommentWindowOpened(project *models.Project, endDate string) error {
	// Get all members for public comment notifications
	members, _, err := s.memberRepo.GetAllMembers(1, 1000)
	if err != nil {
		return err
	}

	var recipients []string
	for _, member := range *members {
		recipients = append(recipients, member.ID.String())
	}

	notificationReq := &models.NotificationRequest{
		Type:     models.NotificationCommentWindowOpened,
		Priority: models.NotificationPriorityMedium,
		Channel:  models.NotificationChannelBoth,
		Data: map[string]interface{}{
			"project_title": project.Title,
			"project_id":    project.ID.String(),
			"end_date":      endDate,
		},
		ProjectID: func() *string { s := project.ID.String(); return &s }(),
	}

	return s.CreateNotification(notificationReq, recipients)
}

// NotifyCommentWindowClosed sends notification when a comment window closes
func (s *NotificationService) NotifyCommentWindowClosed(project *models.Project) error {
	// Get all members for public comment notifications
	members, _, err := s.memberRepo.GetAllMembers(1, 1000)
	if err != nil {
		return err
	}

	var recipients []string
	for _, member := range *members {
		recipients = append(recipients, member.ID.String())
	}

	notificationReq := &models.NotificationRequest{
		Type:     models.NotificationCommentWindowClosed,
		Priority: models.NotificationPriorityMedium,
		Channel:  models.NotificationChannelBoth,
		Data: map[string]interface{}{
			"project_title": project.Title,
			"project_id":    project.ID.String(),
		},
		ProjectID: func() *string { s := project.ID.String(); return &s }(),
	}

	return s.CreateNotification(notificationReq, recipients)
}

// NotifyCommentReceived sends notification when a new comment is received
func (s *NotificationService) NotifyCommentReceived(project *models.Project, commentByID string) error {
	// Notify project stakeholders about new comment
	members, _, err := s.memberRepo.GetAllMembers(1, 1000)
	if err != nil {
		return err
	}

	var recipients []string
	for _, member := range *members {
		// Don't notify the commenter
		if member.ID.String() != commentByID {
			recipients = append(recipients, member.ID.String())
		}
	}

	notificationReq := &models.NotificationRequest{
		Type:     models.NotificationCommentReceived,
		Priority: models.NotificationPriorityMedium,
		Channel:  models.NotificationChannelBoth,
		Data: map[string]interface{}{
			"project_title": project.Title,
			"project_id":    project.ID.String(),
			"comment_by":    commentByID,
		},
		ProjectID: func() *string { s := project.ID.String(); return &s }(),
	}

	return s.CreateNotification(notificationReq, recipients)
}

// NotifyDeadlineReminder sends deadline reminder notifications
func (s *NotificationService) NotifyDeadlineReminder(project *models.Project, deadlineType string, daysLeft int) error {
	// Get relevant members based on deadline type
	members, _, err := s.memberRepo.GetAllMembers(1, 1000)
	if err != nil {
		return err
	}

	var recipients []string
	for _, member := range *members {
		recipients = append(recipients, member.ID.String())
	}

	notificationReq := &models.NotificationRequest{
		Type:     models.NotificationDeadlineReminder,
		Priority: models.NotificationPriorityHigh,
		Channel:  models.NotificationChannelBoth,
		Data: map[string]interface{}{
			"project_title": project.Title,
			"project_id":    project.ID.String(),
			"deadline_type": deadlineType,
			"days_left":     daysLeft,
		},
		ProjectID: func() *string { s := project.ID.String(); return &s }(),
	}

	return s.CreateNotification(notificationReq, recipients)
}

// NotifyTaskEscalation sends task escalation notifications
func (s *NotificationService) NotifyTaskEscalation(project *models.Project, task string, escalatedToRole string, originalAssigneeID string) error {
	// Get members with the escalated role
	members, _, err := s.memberRepo.GetAllMembers(1, 1000)
	if err != nil {
		return err
	}

	var recipients []string
	for _, member := range *members {
		// TODO: Filter by role when GetMembersByRole is implemented
		recipients = append(recipients, member.ID.String())
	}

	notificationReq := &models.NotificationRequest{
		Type:     models.NotificationTaskEscalation,
		Priority: models.NotificationPriorityCritical,
		Channel:  models.NotificationChannelBoth,
		Data: map[string]interface{}{
			"project_title":        project.Title,
			"project_id":           project.ID.String(),
			"task":                 task,
			"escalated_to_role":    escalatedToRole,
			"original_assignee_id": originalAssigneeID,
		},
		ProjectID: func() *string { s := project.ID.String(); return &s }(),
	}

	return s.CreateNotification(notificationReq, recipients)
}

// NotifySystemMaintenance sends system maintenance notifications
func (s *NotificationService) NotifySystemMaintenance(startTime string, endTime string, description string) error {
	// Get all members for system-wide notifications
	members, _, err := s.memberRepo.GetAllMembers(1, 1000)
	if err != nil {
		return err
	}

	var recipients []string
	for _, member := range *members {
		recipients = append(recipients, member.ID.String())
	}

	notificationReq := &models.NotificationRequest{
		Type:     models.NotificationMaintenance,
		Priority: models.NotificationPriorityHigh,
		Channel:  models.NotificationChannelBoth,
		Data: map[string]interface{}{
			"start_time":  startTime,
			"end_time":    endTime,
			"description": description,
		},
	}

	return s.CreateNotification(notificationReq, recipients)
}

// NotifyPolicyChange sends policy change notifications
func (s *NotificationService) NotifyPolicyChange(policyName string, effectiveDate string, description string) error {
	// Get all members for policy change notifications
	members, _, err := s.memberRepo.GetAllMembers(1, 1000)
	if err != nil {
		return err
	}

	var recipients []string
	for _, member := range *members {
		recipients = append(recipients, member.ID.String())
	}

	notificationReq := &models.NotificationRequest{
		Type:     models.NotificationPolicyChange,
		Priority: models.NotificationPriorityMedium,
		Channel:  models.NotificationChannelBoth,
		Data: map[string]interface{}{
			"policy_name":    policyName,
			"effective_date": effectiveDate,
			"description":    description,
		},
	}

	return s.CreateNotification(notificationReq, recipients)
}

// NotifyTrainingOpportunity sends training opportunity notifications
func (s *NotificationService) NotifyTrainingOpportunity(trainingTitle string, startDate string, registrationDeadline string, targetRoles []string) error {
	// Get members based on target roles
	members, _, err := s.memberRepo.GetAllMembers(1, 1000)
	if err != nil {
		return err
	}

	var recipients []string
	for _, member := range *members {
		// TODO: Filter by target roles when GetMembersByRole is implemented
		recipients = append(recipients, member.ID.String())
	}

	notificationReq := &models.NotificationRequest{
		Type:     models.NotificationTraining,
		Priority: models.NotificationPriorityMedium,
		Channel:  models.NotificationChannelBoth,
		Data: map[string]interface{}{
			"training_title":        trainingTitle,
			"start_date":            startDate,
			"registration_deadline": registrationDeadline,
			"target_roles":          targetRoles,
		},
	}

	return s.CreateNotification(notificationReq, recipients)
}

// NotifyProjectCreated sends notifications for new project creation
func (s *NotificationService) NotifyProjectCreated(project *models.Project, createdByID string) error {
	// Get relevant members (TC members, SMC members, etc.)
	recipients, err := s.getProjectRelatedMembers(project)
	if err != nil {
		return fmt.Errorf("failed to get project-related members: %w", err)
	}

	req := &models.NotificationRequest{
		Type:     models.NotificationProjectCreated,
		Priority: models.NotificationPriorityMedium,
		Channel:  models.NotificationChannelBoth,
		Message:  fmt.Sprintf("New project '%s' has been created", project.Title),
		Data: map[string]interface{}{
			"project_id":    project.ID,
			"project_title": project.Title,
			"created_by":    createdByID,
		},
		ProjectID: func() *string { s := project.ID.String(); return &s }(),
	}

	return s.CreateNotification(req, recipients)
}

// NotifyBallotOpened sends notifications when a ballot is opened
func (s *NotificationService) NotifyBallotOpened(balloting *models.Balloting, project *models.Project) error {
	// Get eligible voters for this ballot
	recipients, err := s.getBallotEligibleMembers(balloting, project)
	if err != nil {
		return fmt.Errorf("failed to get ballot-eligible members: %w", err)
	}

	req := &models.NotificationRequest{
		Type:     models.NotificationBallotOpened,
		Priority: models.NotificationPriorityHigh,
		Channel:  models.NotificationChannelBoth,
		Message:  fmt.Sprintf("Ballot opened for project '%s'. Please cast your vote.", project.Title),
		Data: map[string]interface{}{
			"project_id":    project.ID,
			"project_title": project.Title,
			"balloting_id":  balloting.ID,
			"start_date":    balloting.StartDate,
			"end_date":      balloting.EndDate,
		},
		ProjectID:   func() *string { s := project.ID.String(); return &s }(),
		BallotingID: func() *string { s := balloting.ID.String(); return &s }(),
	}

	return s.CreateNotification(req, recipients)
}

// NotifyBallotReminder sends ballot reminder notifications
func (s *NotificationService) NotifyBallotReminder(balloting *models.Balloting, project *models.Project) error {
	// Get members who haven't voted yet
	recipients, err := s.getNonVotingMembers(balloting, project)
	if err != nil {
		return fmt.Errorf("failed to get non-voting members: %w", err)
	}

	if len(recipients) == 0 {
		return nil // Everyone has voted
	}

	daysLeft := int(time.Until(balloting.EndDate).Hours() / 24)
	req := &models.NotificationRequest{
		Type:     models.NotificationBallotReminder,
		Priority: models.NotificationPriorityHigh,
		Channel:  models.NotificationChannelBoth,
		Message:  fmt.Sprintf("Reminder: Ballot for project '%s' closes in %d days. Please cast your vote.", project.Title, daysLeft),
		Data: map[string]interface{}{
			"project_id":    project.ID,
			"project_title": project.Title,
			"balloting_id":  balloting.ID,
			"end_date":      balloting.EndDate,
			"days_left":     daysLeft,
		},
		ProjectID:   func() *string { s := project.ID.String(); return &s }(),
		BallotingID: func() *string { s := balloting.ID.String(); return &s }(),
	}

	return s.CreateNotification(req, recipients)
}

// NotifyDocumentUploaded sends notifications for document uploads
func (s *NotificationService) NotifyDocumentUploaded(document *models.Document, project *models.Project, uploadedByID string) error {
	recipients, err := s.getProjectRelatedMembers(project)
	if err != nil {
		return fmt.Errorf("failed to get project-related members: %w", err)
	}

	req := &models.NotificationRequest{
		Type:     models.NotificationDocumentUploaded,
		Priority: models.NotificationPriorityMedium,
		Channel:  models.NotificationChannelBoth,
		Message:  fmt.Sprintf("New document '%s' uploaded for project '%s'", document.Title, project.Title),
		Data: map[string]interface{}{
			"document_id":    document.ID,
			"document_title": document.Title,
			"project_id":     project.ID,
			"project_title":  project.Title,
			"uploaded_by":    uploadedByID,
		},
		ProjectID:  func() *string { s := project.ID.String(); return &s }(),
		DocumentID: func() *string { s := document.ID.String(); return &s }(),
	}

	return s.CreateNotification(req, recipients)
}

// NotifyMeetingInvitation sends meeting invitation notifications
func (s *NotificationService) NotifyMeetingInvitation(meeting *models.Meeting) error {
	// Get meeting attendees
	var recipients []string
	if meeting.Attendees != nil {
		for _, attendee := range *meeting.Attendees {
			recipients = append(recipients, attendee.ID.String())
		}
	}

	req := &models.NotificationRequest{
		Type:     models.NotificationMeetingInvitation,
		Priority: models.NotificationPriorityHigh,
		Channel:  models.NotificationChannelBoth,
		Message:  fmt.Sprintf("You are invited to meeting '%s' on %s", meeting.Title, meeting.Date.Format("2006-01-02")),
		Data: map[string]interface{}{
			"meeting_id":    meeting.ID,
			"meeting_title": meeting.Title,
			"start_time":    meeting.StartTime,
			"end_time":      meeting.EndTime,
			"venue":         meeting.Venue,
		},
		MeetingID: func() *string { s := meeting.ID.String(); return &s }(),
	}

	return s.CreateNotification(req, recipients)
}

// SendAdminAnnouncement sends platform-wide or role-targeted announcements
func (s *NotificationService) SendAdminAnnouncement(title, message string, targetRoles []string, priority models.NotificationPriority, createdByID string) error {
	var recipients []string
	var err error

	if len(targetRoles) == 0 {
		// Platform-wide announcement - get all active members
		members, _, err := s.memberRepo.GetAllMembers(1000, 0) // Get first 1000 members
		if err != nil {
			return fmt.Errorf("failed to get all members: %w", err)
		}
		for _, member := range *members {
			recipients = append(recipients, member.ID.String())
		}
	} else {
		// Role-targeted announcement
		recipients, err = s.getMembersByRoles(targetRoles)
		if err != nil {
			return fmt.Errorf("failed to get members by roles: %w", err)
		}
	}

	req := &models.NotificationRequest{
		Type:     models.NotificationAnnouncement,
		Priority: priority,
		Channel:  models.NotificationChannelBoth,
		Title:    title,
		Message:  message,
		Data: map[string]interface{}{
			"announcement": true,
			"target_roles": targetRoles,
		},
	}

	return s.CreateNotification(req, recipients)
}

// GetNotificationDashboard retrieves dashboard data for a member
func (s *NotificationService) GetNotificationDashboard(memberID string) (*models.NotificationDashboard, error) {
	return s.notificationRepo.GetNotificationDashboard(memberID)
}

// GetNotifications retrieves notifications for a member with pagination
func (s *NotificationService) GetNotifications(memberID string, limit, offset int) ([]*models.Notification, error) {
	return s.notificationRepo.GetNotificationsByMemberID(memberID, offset/limit+1, limit)
}

// GetNotificationsByMemberID gets notifications for a member with filters
func (s *NotificationService) GetNotificationsByMemberID(memberID string, page, limit int, unreadOnly, criticalOnly bool, notificationType string) ([]*models.Notification, int64, error) {
	offset := (page - 1) * limit
	notifications, err := s.notificationRepo.GetNotificationsByMemberID(memberID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Apply filters (this is a simplified implementation)
	var filtered []*models.Notification
	for _, notification := range notifications {
		if unreadOnly && notification.Read {
			continue
		}
		if criticalOnly && notification.Priority != models.NotificationPriorityCritical {
			continue
		}
		if notificationType != "" && string(notification.Type) != notificationType {
			continue
		}
		filtered = append(filtered, notification)
	}

	// Get total count from notification counts
	counts, err := s.notificationRepo.GetNotificationCounts(memberID)
	if err != nil {
		return nil, 0, err
	}
	total := counts["total"]

	return filtered, total, nil
}

// SearchNotifications searches notifications based on criteria
func (s *NotificationService) SearchNotifications(req *models.NotificationSearchRequest) ([]*models.Notification, int64, error) {
	return s.notificationRepo.SearchNotifications(req)
}

// GetNotificationHistory gets notification history for a member
func (s *NotificationService) GetNotificationHistory(memberID string, page, limit int) ([]*models.NotificationHistory, int64, error) {
	offset := (page - 1) * limit
	history, err := s.notificationRepo.GetNotificationHistory(memberID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	counts, err := s.notificationRepo.GetNotificationCounts(memberID)
	if err != nil {
		return nil, 0, err
	}
	total := counts["total"]

	return history, total, nil
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(notificationID uuid.UUID, memberID string) error {
	return s.notificationRepo.DeleteNotification(notificationID, memberID)
}

// MarkNotificationAsRead marks a notification as read
func (s *NotificationService) MarkNotificationAsRead(notificationID uuid.UUID, memberID string) error {
	return s.notificationRepo.MarkNotificationAsRead(notificationID, memberID)
}

// MarkAllNotificationsAsRead marks all notifications as read for a member
func (s *NotificationService) MarkAllNotificationsAsRead(memberID string) error {
	return s.notificationRepo.MarkAllNotificationsAsRead(memberID)
}

// GetNotificationPreferences retrieves notification preferences for a member
func (s *NotificationService) GetNotificationPreferences(memberID string) (*models.NotificationPreference, error) {
	return s.notificationRepo.GetOrCreateNotificationPreference(memberID)
}

// UpdateNotificationPreferences updates notification preferences for a member
func (s *NotificationService) UpdateNotificationPreferences(preferences *models.NotificationPreference) error {
	return s.notificationRepo.UpdateNotificationPreference(preferences)
}

// ProcessPendingEmailNotifications processes pending email notifications
func (s *NotificationService) ProcessPendingEmailNotifications() error {
	notifications, err := s.notificationRepo.GetPendingEmailNotifications(50)
	if err != nil {
		return fmt.Errorf("failed to get pending email notifications: %w", err)
	}

	for _, notification := range notifications {
		err := s.sendEmailNotification(notification)
		if err != nil {
			// Log error but continue processing other notifications
			fmt.Printf("Failed to send email notification %s: %v\n", notification.ID, err)
			continue
		}

		// Mark as sent
		err = s.notificationRepo.MarkNotificationEmailSent(notification.ID)
		if err != nil {
			fmt.Printf("Failed to mark email as sent for notification %s: %v\n", notification.ID, err)
		}
	}

	return nil
}

// Helper methods

// shouldSendNotification checks if a notification should be sent based on user preferences
func (s *NotificationService) shouldSendNotification(notificationType models.NotificationType, preferences *models.NotificationPreference) bool {
	switch notificationType {
	case models.NotificationProjectCreated, models.NotificationProjectAssigned, models.NotificationProjectUpdated:
		return preferences.ProjectNotifications
	case models.NotificationBallotOpened, models.NotificationBallotReminder, models.NotificationBallotClosing, models.NotificationBallotClosed:
		return preferences.BallotNotifications
	case models.NotificationDocumentUploaded, models.NotificationDocumentUpdated, models.NotificationDocumentVersion:
		return preferences.DocumentNotifications
	case models.NotificationMeetingInvitation, models.NotificationMeetingChanged, models.NotificationMeetingReminder, models.NotificationMeetingCancelled:
		return preferences.MeetingNotifications
	case models.NotificationCommentWindowOpened, models.NotificationCommentWindowClosed, models.NotificationCommentReceived:
		return preferences.CommentNotifications
	case models.NotificationSystemUpdate, models.NotificationPolicyChange, models.NotificationMaintenance, models.NotificationTraining, models.NotificationAnnouncement:
		return preferences.SystemNotifications
	case models.NotificationDeadlineReminder, models.NotificationTaskEscalation:
		return preferences.DeadlineReminders
	default:
		return true
	}
}

// determineChannel determines the notification channel based on request and preferences
func (s *NotificationService) determineChannel(requestedChannel models.NotificationChannel, preferences *models.NotificationPreference) models.NotificationChannel {
	if requestedChannel == models.NotificationChannelEmail && preferences.EmailNotifications {
		return models.NotificationChannelEmail
	}
	if requestedChannel == models.NotificationChannelInApp && preferences.InAppNotifications {
		return models.NotificationChannelInApp
	}
	if requestedChannel == models.NotificationChannelBoth {
		if preferences.EmailNotifications && preferences.InAppNotifications {
			return models.NotificationChannelBoth
		} else if preferences.EmailNotifications {
			return models.NotificationChannelEmail
		} else if preferences.InAppNotifications {
			return models.NotificationChannelInApp
		}
	}

	// Default to in-app if preferences allow
	if preferences.InAppNotifications {
		return models.NotificationChannelInApp
	}

	return models.NotificationChannelEmail // Fallback
}

// processTemplate processes notification template with data
func (s *NotificationService) processTemplate(template string, data map[string]interface{}) string {
	if data == nil {
		return template
	}

	result := template
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}

// getProjectRelatedMembers gets members related to a project (TC members, SMC members, etc.)
func (s *NotificationService) getProjectRelatedMembers(project *models.Project) ([]string, error) {
	// This would need to be implemented based on your specific business logic
	// For now, returning a placeholder implementation
	var recipients []string

	// Get TC members for the project's committee
	tcMembers, err := s.getMembersByRoles([]string{"TC_MEMBER", "TC_SECRETARIAT"})
	if err != nil {
		return nil, err
	}
	recipients = append(recipients, tcMembers...)

	// Get SMC members
	smcMembers, err := s.getMembersByRoles([]string{"SMC_MEMBER"})
	if err != nil {
		return nil, err
	}
	recipients = append(recipients, smcMembers...)

	return s.removeDuplicates(recipients), nil
}

// getBallotEligibleMembers gets members eligible to vote on a ballot
func (s *NotificationService) getBallotEligibleMembers(balloting *models.Balloting, project *models.Project) ([]string, error) {
	// This would need to be implemented based on your balloting eligibility logic
	// For now, returning TC members as eligible voters
	return s.getMembersByRoles([]string{"TC_MEMBER", "TC_SECRETARIAT"})
}

// getNonVotingMembers gets members who haven't voted yet
func (s *NotificationService) getNonVotingMembers(balloting *models.Balloting, project *models.Project) ([]string, error) {
	// Get all eligible members
	eligibleMembers, err := s.getBallotEligibleMembers(balloting, project)
	if err != nil {
		return nil, err
	}

	// Get members who have already voted
	// This would need to query the votes table
	// For now, returning all eligible members as a placeholder
	return eligibleMembers, nil
}

// getMembersByRoles gets all members with specific roles
func (s *NotificationService) getMembersByRoles(roleNames []string) ([]string, error) {
	var memberIDs []string

	for _, roleName := range roleNames {
		// Get role by name
		roles, err := s.rbacRepo.ListRoles()
		if err != nil {
			return nil, err
		}

		var roleID string
		for _, role := range roles {
			if role.Slug == roleName {
				roleID = role.ID.String()
				break
			}
		}

		if roleID == "" {
			continue
		}

		// Get members with this role - this would need to be implemented
		// For now, skipping as the method doesn't exist
		// TODO: Implement GetMembersByRole in MemberRepository
	}

	return s.removeDuplicates(memberIDs), nil
}

// removeDuplicates removes duplicate member IDs
func (s *NotificationService) removeDuplicates(memberIDs []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, id := range memberIDs {
		if !keys[id] {
			keys[id] = true
			result = append(result, id)
		}
	}

	return result
}

// processEmailNotifications processes email notifications asynchronously
func (s *NotificationService) processEmailNotifications(notifications []*models.Notification) {
	for _, notification := range notifications {
		if notification.Channel == models.NotificationChannelEmail || notification.Channel == models.NotificationChannelBoth {
			err := s.sendEmailNotification(notification)
			if err != nil {
				fmt.Printf("Failed to send email notification %s: %v\n", notification.ID, err)
				continue
			}

			// Mark as sent
			err = s.notificationRepo.MarkNotificationEmailSent(notification.ID)
			if err != nil {
				fmt.Printf("Failed to mark email as sent for notification %s: %v\n", notification.ID, err)
			}
		}
	}
}

// sendEmailNotification sends an email notification
func (s *NotificationService) sendEmailNotification(notification *models.Notification) error {
	if notification.Recipient == nil {
		return fmt.Errorf("recipient not loaded")
	}

	// Get email template if available
	template, err := s.notificationRepo.GetNotificationTemplateByType(notification.Type)
	var emailContent string
	if err == nil && template != nil {
		emailContent = template.GetEmailTemplate(notification.Language)
		if emailContent != "" {
			// Process template with notification data
			var data map[string]interface{}
			if notification.Data != "" {
				json.Unmarshal([]byte(notification.Data), &data)
			}
			emailContent = s.processTemplate(emailContent, data)
		} else {
			emailContent = notification.Message
		}
	} else {
		emailContent = notification.Message
	}

	// Use SendCustomEmails method which exists
	emails := []RecipientEmail{
		{
			To:    notification.Recipient.Email,
			Body:  emailContent,
			Title: notification.Title,
		},
	}
	return s.emailService.SendCustomEmails(emails)
}

// createNotificationHistory creates history records for notifications
func (s *NotificationService) createNotificationHistory(notifications []*models.Notification) {
	for _, notification := range notifications {
		history := &models.NotificationHistory{
			NotificationID: notification.ID,
			RecipientID:    notification.RecipientID,
			Type:           notification.Type,
			Channel:        notification.Channel,
			Status:         "sent",
			SentAt:         time.Now(),
		}

		err := s.notificationRepo.CreateNotificationHistory(history)
		if err != nil {
			fmt.Printf("Failed to create notification history for %s: %v\n", notification.ID, err)
		}
	}
}

// CleanupOldNotifications removes old notifications based on retention policy
func (s *NotificationService) CleanupOldNotifications(retentionDays int) error {
	return s.notificationRepo.DeleteOldNotifications(retentionDays)
}

// EscalateOverdueTasks escalates overdue tasks to higher roles
func (s *NotificationService) EscalateOverdueTasks() error {
	// This would need to be implemented based on your specific escalation logic
	// For now, this is a placeholder
	return nil
}

// SendDeadlineReminders sends deadline reminder notifications
func (s *NotificationService) SendDeadlineReminders() error {
	// This would need to be implemented based on your specific deadline logic
	// For now, this is a placeholder
	return nil
}
