package events

import (
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
)

// NotificationEventHandler handles notification events triggered by system actions
type NotificationEventHandler struct {
	notificationService *services.NotificationService
}

// NewNotificationEventHandler creates a new notification event handler
func NewNotificationEventHandler(notificationService *services.NotificationService) *NotificationEventHandler {
	return &NotificationEventHandler{
		notificationService: notificationService,
	}
}

// OnProjectCreated handles project creation events
func (h *NotificationEventHandler) OnProjectCreated(project *models.Project, createdByID string) error {
	return h.notificationService.NotifyProjectCreated(project, createdByID)
}

// OnProjectAssigned handles project assignment events
func (h *NotificationEventHandler) OnProjectAssigned(project *models.Project, assignedToID string, assignedByID string) error {
	return h.notificationService.NotifyProjectAssigned(project, assignedToID, assignedByID)
}

// OnBallotOpened handles ballot opening events
func (h *NotificationEventHandler) OnBallotOpened(project *models.Project, balloting *models.Balloting) error {
	return h.notificationService.NotifyBallotOpened(balloting, project)
}

// OnBallotReminder handles ballot reminder events
func (h *NotificationEventHandler) OnBallotReminder(project *models.Project, balloting *models.Balloting, daysLeft int) error {
	return h.notificationService.NotifyBallotReminder(balloting, project)
}

// OnBallotClosed handles ballot closing events
func (h *NotificationEventHandler) OnBallotClosed(project *models.Project, balloting *models.Balloting) error {
	return h.notificationService.NotifyBallotClosed(project, balloting)
}

// OnDocumentUploaded handles document upload events
func (h *NotificationEventHandler) OnDocumentUploaded(project *models.Project, document *models.Document, uploadedByID string) error {
	return h.notificationService.NotifyDocumentUploaded(document, project, uploadedByID)
}

// OnDocumentUpdated handles document update events
func (h *NotificationEventHandler) OnDocumentUpdated(project *models.Project, document *models.Document, updatedByID string) error {
	return h.notificationService.NotifyDocumentUpdated(project, document, updatedByID)
}

// OnMeetingInvitation handles meeting invitation events
func (h *NotificationEventHandler) OnMeetingInvitation(meeting *models.Meeting) error {
	return h.notificationService.NotifyMeetingInvitation(meeting)
}

// OnMeetingChanged handles meeting change events
func (h *NotificationEventHandler) OnMeetingChanged(meeting *models.Meeting, changedByID string) error {
	return h.notificationService.NotifyMeetingChanged(meeting, changedByID)
}

// OnMeetingCancelled handles meeting cancellation events
func (h *NotificationEventHandler) OnMeetingCancelled(meeting *models.Meeting, cancelledByID string) error {
	return h.notificationService.NotifyMeetingCancelled(meeting, cancelledByID)
}

// OnCommentWindowOpened handles comment window opening events
func (h *NotificationEventHandler) OnCommentWindowOpened(project *models.Project, endDate string) error {
	return h.notificationService.NotifyCommentWindowOpened(project, endDate)
}

// OnCommentWindowClosed handles comment window closing events
func (h *NotificationEventHandler) OnCommentWindowClosed(project *models.Project) error {
	return h.notificationService.NotifyCommentWindowClosed(project)
}

// OnCommentReceived handles new comment events
func (h *NotificationEventHandler) OnCommentReceived(project *models.Project, commentByID string) error {
	return h.notificationService.NotifyCommentReceived(project, commentByID)
}

// OnDeadlineReminder handles deadline reminder events
func (h *NotificationEventHandler) OnDeadlineReminder(project *models.Project, deadlineType string, daysLeft int) error {
	return h.notificationService.NotifyDeadlineReminder(project, deadlineType, daysLeft)
}

// OnTaskEscalation handles task escalation events
func (h *NotificationEventHandler) OnTaskEscalation(project *models.Project, task string, escalatedToRole string, originalAssigneeID string) error {
	return h.notificationService.NotifyTaskEscalation(project, task, escalatedToRole, originalAssigneeID)
}

// OnSystemMaintenance handles system maintenance notification events
func (h *NotificationEventHandler) OnSystemMaintenance(startTime string, endTime string, description string) error {
	return h.notificationService.NotifySystemMaintenance(startTime, endTime, description)
}

// OnPolicyChange handles policy change notification events
func (h *NotificationEventHandler) OnPolicyChange(policyName string, effectiveDate string, description string) error {
	return h.notificationService.NotifyPolicyChange(policyName, effectiveDate, description)
}

// OnTrainingOpportunity handles training opportunity notification events
func (h *NotificationEventHandler) OnTrainingOpportunity(trainingTitle string, startDate string, registrationDeadline string, targetRoles []string) error {
	return h.notificationService.NotifyTrainingOpportunity(trainingTitle, startDate, registrationDeadline, targetRoles)
}

// OnAdminAnnouncement handles admin announcement events
func (h *NotificationEventHandler) OnAdminAnnouncement(title, message string, targetRoles []string, priority models.NotificationPriority, createdByID string) error {
	return h.notificationService.SendAdminAnnouncement(title, message, targetRoles, priority, createdByID)
}