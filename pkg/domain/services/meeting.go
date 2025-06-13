package services

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type MeetingService struct {
	repo *repository.MeetingRepository
}

func NewMeetingService(repo *repository.MeetingRepository) *MeetingService {
	return &MeetingService{repo: repo}
}

func (service *MeetingService) CreateMeeting(meeting *models.Meeting) error {
	meeting.ID = uuid.New()
	meeting.CreatedAt = time.Now()
	return service.repo.CreateMeeting(meeting)
}

func (service *MeetingService) GetMeetingByID(id string) (*models.Meeting, error) {
	return service.repo.GetMeetingByID(id)
}

func (service *MeetingService) GetAllMeetings(page, pageSize int) (*[]models.Meeting, error) {
	return service.repo.GetAllMeetings(page, pageSize)
}

func (service *MeetingService) GetMeetingsByCommittee(committeeID string) (*[]models.Meeting, error) {
	return service.repo.GetMeetingsByCommittee(committeeID)
}

func (service *MeetingService) GetMeetingsByProject(projectID string) (*[]models.Meeting, error) {
	return service.repo.GetMeetingsByProject(projectID)
}

func (service *MeetingService) GetMeetingsByType(meetingType models.MeetingType) (*[]models.Meeting, error) {
	return service.repo.GetMeetingsByType(meetingType)
}

func (service *MeetingService) GetMeetingsByStatus(status models.MeetingStatus) (*[]models.Meeting, error) {
	return service.repo.GetMeetingsByStatus(status)
}

func (service *MeetingService) GetMeetingsByHostOrganization(hostOrg string) (*[]models.Meeting, error) {
	return service.repo.GetMeetingsByHostOrganization(hostOrg)
}

func (service *MeetingService) GetUpcomingMeetings() (*[]models.Meeting, error) {
	return service.repo.GetUpcomingMeetings()
}

func (service *MeetingService) GetMeetingsCreatedByMember(memberID string) (*[]models.Meeting, error) {
	return service.repo.GetMeetingsCreatedByMember(memberID)
}

func (service *MeetingService) UpdateMeeting(meeting *models.Meeting) error {
	return service.repo.UpdateMeeting(meeting)
}

func (service *MeetingService) DeleteMeeting(meetingID string) error {
	return service.repo.DeleteMeeting(meetingID)
}

func (service *MeetingService) AddAttendeeToMeeting(meetingID string, memberID string) error {
	return service.repo.AddAttendeeToMeeting(meetingID, memberID)
}

func (service *MeetingService) RemoveAttendeeFromMeeting(meetingID string, memberID string) error {
	return service.repo.RemoveAttendeeFromMeeting(meetingID, memberID)
}

func (service *MeetingService) AddRelatedDocumentToMeeting(meetingID string, documentID string) error {
	return service.repo.AddRelatedDocumentToMeeting(meetingID, documentID)
}

func (service *MeetingService) RemoveRelatedDocumentFromMeeting(meetingID string, documentID string) error {
	return service.repo.RemoveRelatedDocumentFromMeeting(meetingID, documentID)
}

func (service *MeetingService) AddCommitteeDraftToMeeting(meetingID string, documentID string) error {
	return service.repo.AddCommitteeDraftToMeeting(meetingID, documentID)
}

func (service *MeetingService) UpdateMeetingStatus(meetingID string, status models.MeetingStatus, reason string) error {
	return service.repo.UpdateMeetingStatus(meetingID, status, reason)
}

func (service *MeetingService) CheckQuorum(meetingID string) (bool, error) {
	return service.repo.CheckQuorum(meetingID)
}

// SendMeetingInvitations sends email invitations to meeting attendees
func (service *MeetingService) SendMeetingInvitations(meetingID string, emailService *EmailService) error {
	// Get meeting details
	meeting, err := service.GetMeetingByID(meetingID)
	if err != nil {
		return fmt.Errorf("failed to get meeting: %w", err)
	}

	// Verify it's an electronic meeting
	if meeting.Format != models.MeetingFormatElectronic {
		return errors.New("meeting invitations can only be sent for electronic meetings")
	}

	// Verify video conference link exists
	if meeting.VideoConferenceLink == "" {
		return errors.New("video conference link is required for electronic meetings")
	}

	if meeting.Attendees == nil || len(*meeting.Attendees) == 0 {
		return errors.New("no attendees found for the meeting")
	}

	// Prepare email content
	meetingDate := meeting.Date.Format("Monday, January 2, 2006")
	meetingTime := fmt.Sprintf("%s - %s", meeting.StartTime, meeting.EndTime)

	// Read the email template from file
	tmpl, err := os.ReadFile("templates/meeting_invitation.html")
	if err != nil {
		return fmt.Errorf("failed to read email template: %w", err)
	}

	// Parse the template
	t, err := template.New("meetingInvitation").Parse(string(tmpl))
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Prepare template data
	data := struct {
		Title               string
		CommitteeName       string
		Date                string
		Time                string
		Venue               string
		Format              string
		Language            string
		Agenda              string
		VideoConferenceLink string
	}{
		Title:               meeting.Title,
		CommitteeName:       meeting.CommitteeName,
		Date:                meetingDate,
		Time:                meetingTime,
		Venue:               meeting.Venue,
		Format:              string(meeting.Format),
		Language:            meeting.Language,
		Agenda:              meeting.Agenda,
		VideoConferenceLink: meeting.VideoConferenceLink,
	}

	// Execute template with data
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Send emails to all attendees
	var emails []RecipientEmail
	for _, attendee := range *meeting.Attendees {
		emails = append(emails, RecipientEmail{
			To:    attendee.Email,
			Title: fmt.Sprintf("Meeting Invitation: %s", meeting.Title),
			Body:  body.String(),
		})
	}

	return emailService.SendCustomEmails(emails)
}
