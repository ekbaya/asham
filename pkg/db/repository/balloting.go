package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BallotingRepository struct {
	db *gorm.DB
}

func NewBallotingRepository(db *gorm.DB) *BallotingRepository {
	return &BallotingRepository{db: db}
}

func (r *BallotingRepository) CreateVote(vote *models.Vote) error {
	isEligible, err := r.IsEligibleToVote(vote.MemberID, vote.ProjectID)
	if err != nil {
		return err
	}

	if !isEligible {
		return errors.New("member is not eligible to vote")
	}
	return r.db.Create(vote).Error
}

func (r *BallotingRepository) IsEligibleToVote(memberID string, projectID string) (bool, error) {
	// Check if member has already voted
	var voteCount int64
	if err := r.db.Model(&models.Vote{}).
		Where("member_id = ? AND project_id = ?", memberID, projectID).
		Count(&voteCount).Error; err != nil {
		return false, err
	}
	if voteCount > 0 {
		return false, nil
	}

	// Check if the balloting is active
	var balloting models.Balloting
	if err := r.db.Where("project_id = ?", projectID).
		First(&balloting).Error; err != nil {
		return false, err
	}
	if time.Now().After(balloting.EndDate) {
		return false, nil
	}

	// Load member to get NSB ID
	var member models.Member
	if err := r.db.Where("id = ?", memberID).
		First(&member).Error; err != nil {
		return false, err
	}

	// Check involvement in committee stage
	var committeeCount int64
	if err := r.db.Model(&models.CommentObservation{}).
		Joins("JOIN members ON members.id = comment_observations.national_secretary_id").
		Joins("JOIN national_standard_bodies ON national_standard_bodies.id = members.national_standard_body_id").
		Where("national_standard_bodies.id = ?", member.NationalStandardBodyID).
		Count(&committeeCount).Error; err != nil {
		return false, err
	}

	// Check involvement in enquiry stage
	var enquiryCount int64
	if err := r.db.Model(&models.NationalConsultation{}).
		Joins("JOIN members ON members.id = national_consultations.national_secretary_id").
		Joins("JOIN national_standard_bodies ON national_standard_bodies.id = members.national_standard_body_id").
		Where("national_standard_bodies.id = ?", member.NationalStandardBodyID).
		Count(&enquiryCount).Error; err != nil {
		return false, err
	}

	if committeeCount == 0 && enquiryCount == 0 {
		return false, nil
	}

	return true, nil
}

func (r *BallotingRepository) FindVoteByID(id uuid.UUID) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.Where("id = ?", id).First(&vote).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

func (r *BallotingRepository) FindVotesByBallotingID(ballotingID uuid.UUID) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.Where("balloting_id = ?", ballotingID).Preload("Member").Preload("NationalStandardBody").Find(&votes).Error
	if err != nil {
		return nil, err
	}
	return votes, nil
}

func (r *BallotingRepository) FindByProjectID(projectID string) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.Where("project_id = ?", projectID).Preload("Member").Preload("NationalStandardBody").Find(&votes).Error
	if err != nil {
		return nil, err
	}
	return votes, nil
}

func (r *BallotingRepository) FindVotesByMemberID(memberID string) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.Where("member_id = ?", memberID).Find(&votes).Error
	if err != nil {
		return nil, err
	}
	return votes, nil
}

func (r *BallotingRepository) UpdateVote(vote *models.Vote) error {
	return r.db.Save(vote).Error
}

func (r *BallotingRepository) DeleteVote(id uuid.UUID) error {
	return r.db.Delete(&models.Vote{}, id).Error
}

func (r *BallotingRepository) FindVotesWithAssociations() ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.Preload("Project").Preload("Member").Find(&votes).Error
	if err != nil {
		return nil, err
	}
	return votes, nil
}

func (r *BallotingRepository) CountVotesByBalloting(ballotingID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Vote{}).Where("balloting_id = ?", ballotingID).Count(&count).Error
	return count, err
}

// CheckAcceptanceCriteria determines if a project has met the 75% acceptance threshold
func (r *BallotingRepository) CheckAcceptanceCriteria(projectID string) (*models.AcceptanceCriteriaResult, error) {
	const requiredRate float64 = 0.75 // 75% acceptance required

	result := &models.AcceptanceCriteriaResult{
		RequiredRate: requiredRate,
	}

	// Count total votes for this project
	err := r.db.Model(&models.Vote{}).Where("project_id = ?", projectID).Count(&result.TotalVotes).Error
	if err != nil {
		return nil, err
	}

	// If no votes, criteria not met
	if result.TotalVotes == 0 {
		result.Message = "No votes recorded for this project."
		return result, nil
	}

	// Count accepted votes for this project
	err = r.db.Model(&models.Vote{}).Where("project_id = ? AND acceptance = ?", projectID, true).Count(&result.AcceptedVotes).Error
	if err != nil {
		return nil, err
	}

	// Calculate acceptance rate
	result.AcceptanceRate = float64(result.AcceptedVotes) / float64(result.TotalVotes)

	// Check if acceptance threshold is met
	result.CriteriaMet = result.AcceptanceRate >= requiredRate

	// Generate appropriate message
	if result.CriteriaMet {
		result.Message = fmt.Sprintf("Project accepted with %.1f%% approval (required: %.1f%%)",
			result.AcceptanceRate*100, requiredRate*100)
	} else {
		result.Message = fmt.Sprintf("Project not accepted. Current approval: %.1f%% (required: %.1f%%)",
			result.AcceptanceRate*100, requiredRate*100)
	}

	return result, nil
}

func (r *BallotingRepository) FindBallotingByID(id uuid.UUID) (*models.Balloting, error) {
	var balloting models.Balloting
	err := r.db.Where("id = ?", id).First(&balloting).Error
	if err != nil {
		return nil, err
	}
	return &balloting, nil
}

func (r *BallotingRepository) CreateBalloting(balloting *models.Balloting) error {
	// Generate UUID if not provided
	if balloting.ID == uuid.Nil {
		balloting.ID = uuid.New()
	}
	return r.db.Create(balloting).Error
}

func (r *BallotingRepository) UpdateBalloting(balloting *models.Balloting) error {
	return r.db.Save(balloting).Error
}

func (r *BallotingRepository) DeleteBalloting(id uuid.UUID) error {
	return r.db.Delete(&models.Balloting{}, id).Error
}

func (r *BallotingRepository) FindAll() ([]models.Balloting, error) {
	var ballotings []models.Balloting
	err := r.db.Find(&ballotings).Error
	if err != nil {
		return nil, err
	}
	return ballotings, nil
}

func (r *BallotingRepository) FindActiveBallotingWithVotes() ([]models.Balloting, error) {
	var ballotings []models.Balloting
	err := r.db.Where("active = ?", true).Preload("Votes").Find(&ballotings).Error
	if err != nil {
		return nil, err
	}
	return ballotings, nil
}

func (r *BallotingRepository) FindBallotingByPeriod(startDate, endDate time.Time) ([]models.Balloting, error) {
	var ballotings []models.Balloting
	err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&ballotings).Error
	if err != nil {
		return nil, err
	}
	return ballotings, nil
}
