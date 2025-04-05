package repository

import (
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

// Create adds a new vote to the database
func (r *BallotingRepository) Create(vote *models.Vote) error {
	// Generate UUID if not provided
	if vote.ID == uuid.Nil {
		vote.ID = uuid.New()
	}
	return r.db.Create(vote).Error
}

// FindByID retrieves a vote by its ID
func (r *BallotingRepository) FindByID(id uuid.UUID) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.Where("id = ?", id).First(&vote).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

// FindByBallotingID retrieves all votes for a specific balloting
func (r *BallotingRepository) FindByBallotingID(ballotingID uuid.UUID) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.Where("balloting_id = ?", ballotingID).Find(&votes).Error
	if err != nil {
		return nil, err
	}
	return votes, nil
}

// FindByProjectID retrieves all votes for a specific project
func (r *BallotingRepository) FindByProjectID(projectID string) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.Where("project_id = ?", projectID).Find(&votes).Error
	if err != nil {
		return nil, err
	}
	return votes, nil
}

// FindByMemberID retrieves all votes cast by a specific member
func (r *BallotingRepository) FindByMemberID(memberID string) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.Where("member_id = ?", memberID).Find(&votes).Error
	if err != nil {
		return nil, err
	}
	return votes, nil
}

// Update modifies an existing vote
func (r *BallotingRepository) Update(vote *models.Vote) error {
	return r.db.Save(vote).Error
}

// Delete removes a vote from the database
func (r *BallotingRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Vote{}, id).Error
}

// FindVotesWithAssociations retrieves votes with their Project and Member data loaded
func (r *BallotingRepository) FindVotesWithAssociations() ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.Preload("Project").Preload("Member").Find(&votes).Error
	if err != nil {
		return nil, err
	}
	return votes, nil
}

// CountVotesByBalloting counts votes for a specific balloting
func (r *BallotingRepository) CountVotesByBalloting(ballotingID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Vote{}).Where("balloting_id = ?", ballotingID).Count(&count).Error
	return count, err
}

// GetAcceptanceRate calculates the acceptance rate for a project
func (r *BallotingRepository) GetAcceptanceRate(projectID string) (float64, error) {
	var acceptedCount, totalCount int64

	err := r.db.Model(&models.Vote{}).Where("project_id = ?", projectID).Count(&totalCount).Error
	if err != nil {
		return 0, err
	}

	if totalCount == 0 {
		return 0, nil
	}

	err = r.db.Model(&models.Vote{}).Where("project_id = ? AND acceptance = ?", projectID, true).Count(&acceptedCount).Error
	if err != nil {
		return 0, err
	}

	return float64(acceptedCount) / float64(totalCount), nil
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
