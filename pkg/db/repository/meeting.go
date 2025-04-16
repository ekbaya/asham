package repository

import (
	"errors"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MeetingRepository struct {
	db *gorm.DB
}

func NewMeetingRepository(db *gorm.DB) *MeetingRepository {
	return &MeetingRepository{db: db}
}

func (r *MeetingRepository) CreateMeeting(meeting *models.Meeting) error {
	return r.db.Create(meeting).Error
}

func (r *MeetingRepository) GetMeetingByID(id string) (*models.Meeting, error) {
	var meeting models.Meeting
	result := r.db.Where("id = ?", id).Preload(clause.Associations).
		First(&meeting)
	if result.Error != nil {
		return nil, result.Error
	}
	return &meeting, nil
}

func (r *MeetingRepository) GetAllMeetings(page, pageSize int) (*[]models.Meeting, error) {
	var meetings []models.Meeting
	offset := (page - 1) * pageSize
	result := r.db.Preload("Project").
		Limit(pageSize).
		Offset(offset).
		Find(&meetings)
	if result.Error != nil {
		return nil, result.Error
	}
	return &meetings, nil
}

func (r *MeetingRepository) GetMeetingsByCommittee(committeeID string) (*[]models.Meeting, error) {
	var meetings []models.Meeting
	result := r.db.Where("committee_id = ?", committeeID).
		Preload("Project").
		Find(&meetings)
	if result.Error != nil {
		return nil, result.Error
	}
	return &meetings, nil
}

func (r *MeetingRepository) GetMeetingsByProject(projectID string) (*[]models.Meeting, error) {
	var meetings []models.Meeting
	result := r.db.Where("project_id = ?", projectID).
		Preload("Project").
		Find(&meetings)
	if result.Error != nil {
		return nil, result.Error
	}
	return &meetings, nil
}

func (r *MeetingRepository) GetMeetingsByType(meetingType models.MeetingType) (*[]models.Meeting, error) {
	var meetings []models.Meeting
	result := r.db.Where("meeting_type = ?", meetingType).
		Preload("Project").
		Find(&meetings)
	if result.Error != nil {
		return nil, result.Error
	}
	return &meetings, nil
}

func (r *MeetingRepository) GetMeetingsByStatus(status models.MeetingStatus) (*[]models.Meeting, error) {
	var meetings []models.Meeting
	result := r.db.Where("status = ?", status).
		Preload("Project").
		Find(&meetings)
	if result.Error != nil {
		return nil, result.Error
	}
	return &meetings, nil
}

func (r *MeetingRepository) GetMeetingsByHostOrganization(hostOrg string) (*[]models.Meeting, error) {
	var meetings []models.Meeting
	result := r.db.Where("host_organization_id = ?", hostOrg).
		Preload("Project").
		Find(&meetings)
	if result.Error != nil {
		return nil, result.Error
	}
	return &meetings, nil
}

func (r *MeetingRepository) GetUpcomingMeetings() (*[]models.Meeting, error) {
	var meetings []models.Meeting
	result := r.db.Where("date > NOW() AND status != ?", models.MeetingStatusCancelled).
		Preload("Project").
		Order("date ASC").
		Find(&meetings)
	if result.Error != nil {
		return nil, result.Error
	}
	return &meetings, nil
}

func (r *MeetingRepository) GetMeetingsCreatedByMember(memberID string) (*[]models.Meeting, error) {
	var meetings []models.Meeting
	result := r.db.Where("created_by = ?", memberID).
		Preload("Project").
		Find(&meetings)
	if result.Error != nil {
		return nil, result.Error
	}
	return &meetings, nil
}

func (r *MeetingRepository) UpdateMeeting(meeting *models.Meeting) error {
	return r.db.Save(meeting).Error
}

func (r *MeetingRepository) DeleteMeeting(meetingID string) error {
	return r.db.Delete(&models.Meeting{}, "id = ?", meetingID).Error
}

func (r *MeetingRepository) AddAttendeeToMeeting(meetingID string, memberID string) error {
	var meeting models.Meeting
	if err := r.db.First(&meeting, "id = ?", meetingID).Error; err != nil {
		return err
	}

	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&meeting).Association("Attendees").Append(&member)
}

func (r *MeetingRepository) RemoveAttendeeFromMeeting(meetingID string, memberID string) error {
	var meeting models.Meeting
	if err := r.db.First(&meeting, "id = ?", meetingID).Error; err != nil {
		return err
	}

	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&meeting).Association("Attendees").Delete(&member)
}

func (r *MeetingRepository) AddRelatedDocumentToMeeting(meetingID string, documentID string) error {
	var meeting models.Meeting
	if err := r.db.First(&meeting, "id = ?", meetingID).Error; err != nil {
		return err
	}

	var document models.Document
	if err := r.db.First(&document, "id = ?", documentID).Error; err != nil {
		return err
	}

	return r.db.Model(&meeting).Association("RelatedDocuments").Append(&document)
}

func (r *MeetingRepository) RemoveRelatedDocumentFromMeeting(meetingID string, documentID string) error {
	var meeting models.Meeting
	if err := r.db.First(&meeting, "id = ?", meetingID).Error; err != nil {
		return err
	}

	var document models.Document
	if err := r.db.First(&document, "id = ?", documentID).Error; err != nil {
		return err
	}

	return r.db.Model(&meeting).Association("RelatedDocuments").Delete(&document)
}

func (r *MeetingRepository) AddCommitteeDraftToMeeting(meetingID string, documentID string) error {
	var meeting models.Meeting
	if err := r.db.First(&meeting, "id = ?", meetingID).Error; err != nil {
		return err
	}

	var document models.Document
	if err := r.db.First(&document, "id = ?", documentID).Error; err != nil {
		return err
	}

	return r.db.Model(&meeting).Association("CommitteeDrafts").Append(&document)
}

func (r *MeetingRepository) UpdateMeetingStatus(meetingID string, status models.MeetingStatus, reason string) error {
	var meeting models.Meeting
	if err := r.db.First(&meeting, "id = ?", meetingID).Error; err != nil {
		return err
	}

	// Check if the meeting is being cancelled and enforce the two-week rule
	if status == models.MeetingStatusCancelled {
		// Calculate the time difference between now and the meeting date
		timeUntilMeeting := time.Until(meeting.Date)
		twoWeeks := 14 * 24 * time.Hour

		// If the meeting is less than two weeks away, prevent cancellation
		if timeUntilMeeting < twoWeeks {
			return errors.New("meetings cannot be cancelled within two weeks of the scheduled date")
		}
	}

	meeting.Status = status
	if status == models.MeetingStatusCancelled || status == models.MeetingStatusPostponed {
		meeting.CancellationReason = reason
	}

	return r.db.Save(&meeting).Error
}

func (r *MeetingRepository) CheckQuorum(meetingID string) (bool, error) {
	var meeting models.Meeting
	if err := r.db.First(&meeting, "id = ?", meetingID).Error; err != nil {
		return false, err
	}

	// According to the document, quorum requires at least 50% of P-members plus one
	hasQuorum := meeting.PresentPMembers >= (meeting.TotalPMembers/2 + 1)

	// Update the meeting record with quorum status
	meeting.HasQuorum = hasQuorum
	if err := r.db.Save(&meeting).Error; err != nil {
		return hasQuorum, err
	}

	return hasQuorum, nil
}
