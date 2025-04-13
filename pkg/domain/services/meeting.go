package services

import (
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
