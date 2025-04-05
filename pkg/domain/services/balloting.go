package services

import (
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type BallotingService struct {
	repo *repository.BallotingRepository
}

func NewBallotingService(repo *repository.BallotingRepository) *BallotingService {
	return &BallotingService{repo: repo}
}

func (service *BallotingService) CreateVote(vote *models.Vote) error {
	return service.repo.CreateVote(vote)
}

func (service *BallotingService) FindVoteByID(id uuid.UUID) (*models.Vote, error) {
	return service.repo.FindVoteByID(id)
}

func (service *BallotingService) FindVotesByBallotingID(ballotingID uuid.UUID) ([]models.Vote, error) {
	return service.repo.FindVotesByBallotingID(ballotingID)
}

func (service *BallotingService) FindByProjectID(projectID string) ([]models.Vote, error) {
	return service.repo.FindByProjectID(projectID)
}

func (service *BallotingService) FindVotesByMemberID(memberID string) ([]models.Vote, error) {
	return service.repo.FindVotesByMemberID(memberID)
}

func (service *BallotingService) UpdateVote(vote *models.Vote) error {
	return service.repo.UpdateVote(vote)
}

func (service *BallotingService) DeleteVote(id uuid.UUID) error {
	return service.repo.DeleteVote(id)
}

func (service *BallotingService) FindVotesWithAssociations() ([]models.Vote, error) {
	return service.repo.FindVotesWithAssociations()
}

func (service *BallotingService) CountVotesByBalloting(ballotingID uuid.UUID) (int64, error) {
	return service.repo.CountVotesByBalloting(ballotingID)
}

func (service *BallotingService) CheckAcceptanceCriteria(projectID string) (*models.AcceptanceCriteriaResult, error) {
	return service.repo.CheckAcceptanceCriteria(projectID)
}

func (service *BallotingService) FindBallotingByID(id uuid.UUID) (*models.Balloting, error) {
	return service.repo.FindBallotingByID(id)
}

func (service *BallotingService) CreateBalloting(balloting *models.Balloting) error {
	return service.repo.CreateBalloting(balloting)
}

func (service *BallotingService) UpdateBalloting(balloting *models.Balloting) error {
	return service.repo.UpdateBalloting(balloting)
}

func (service *BallotingService) DeleteBalloting(id uuid.UUID) error {
	return service.repo.DeleteBalloting(id)
}

func (service *BallotingService) FindAll() ([]models.Balloting, error) {
	return service.repo.FindAll()
}

func (service *BallotingService) FindActiveBallotingWithVotes() ([]models.Balloting, error) {
	return service.repo.FindActiveBallotingWithVotes()
}

func (service *BallotingService) FindBallotingByPeriod(startDate, endDate time.Time) ([]models.Balloting, error) {
	return service.repo.FindBallotingByPeriod(startDate, endDate)
}
