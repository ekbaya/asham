package services

import (
	"fmt"
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type BallotingService struct {
	repo         *repository.BallotingRepository
	auditService *AuditLogService
}

func NewBallotingService(repo *repository.BallotingRepository, auditService *AuditLogService) *BallotingService {
	return &BallotingService{repo: repo, auditService: auditService}
}

func (service *BallotingService) CreateVote(vote *models.Vote, userID, ipAddress, userAgent, sessionID, requestID string) error {
	startTime := time.Now()
	vote.ID = uuid.New()
	
	err := service.repo.CreateVote(vote)
	
	// Log the action
	metadata := map[string]interface{}{
		"acceptance":        vote.Acceptance,
		"project_id":        vote.ProjectID,
		"balloting_id":      vote.BallotingID.String(),
		"comment":           vote.Comment,
		"execution_time_ms": time.Since(startTime).Milliseconds(),
	}
	
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	
	service.auditService.LogBallotAction(
		&userID, models.ActionVoteSubmit, vote.ID.String(), fmt.Sprintf("Vote for project %s", vote.ProjectID),
		metadata, err == nil, errorMsg, time.Since(startTime).Milliseconds(),
		ipAddress, userAgent, sessionID, requestID,
	)
	
	return err
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

func (service *BallotingService) UpdateVote(vote *models.Vote, userID, ipAddress, userAgent, sessionID, requestID string) error {
	startTime := time.Now()
	
	err := service.repo.UpdateVote(vote)
	
	// Log the action
	metadata := map[string]interface{}{
		"acceptance":        vote.Acceptance,
		"project_id":        vote.ProjectID,
		"balloting_id":      vote.BallotingID.String(),
		"comment":           vote.Comment,
		"execution_time_ms": time.Since(startTime).Milliseconds(),
	}
	
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	
	service.auditService.LogBallotAction(
		&userID, models.ActionVoteUpdate, vote.ID.String(), fmt.Sprintf("Updated vote for project %s", vote.ProjectID),
		metadata, err == nil, errorMsg, time.Since(startTime).Milliseconds(),
		ipAddress, userAgent, sessionID, requestID,
	)
	
	return err
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

func (service *BallotingService) CreateBalloting(balloting *models.Balloting, userID, ipAddress, userAgent, sessionID, requestID string) error {
	startTime := time.Now()
	balloting.ID = uuid.New()
	balloting.CreatedAt = time.Now()
	
	err := service.repo.CreateBalloting(balloting)
	
	// Log the action
	metadata := map[string]interface{}{
		"project_id":        balloting.ProjectID,
		"start_date":        balloting.StartDate,
		"end_date":          balloting.EndDate,
		"active":            balloting.Active,
		"execution_time_ms": time.Since(startTime).Milliseconds(),
	}
	
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	
	service.auditService.LogBallotAction(
		&userID, models.ActionBallotCreate, balloting.ID.String(), fmt.Sprintf("Balloting for project %s", balloting.ProjectID),
		metadata, err == nil, errorMsg, time.Since(startTime).Milliseconds(),
		ipAddress, userAgent, sessionID, requestID,
	)
	
	return err
}

func (service *BallotingService) UpdateBalloting(balloting *models.Balloting, userID, ipAddress, userAgent, sessionID, requestID string) error {
	startTime := time.Now()
	
	err := service.repo.UpdateBalloting(balloting)
	
	// Log the action
	metadata := map[string]interface{}{
		"project_id":        balloting.ProjectID,
		"start_date":        balloting.StartDate,
		"end_date":          balloting.EndDate,
		"active":            balloting.Active,
		"execution_time_ms": time.Since(startTime).Milliseconds(),
	}
	
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	
	service.auditService.LogBallotAction(
		&userID, models.ActionBallotUpdate, balloting.ID.String(), fmt.Sprintf("Updated balloting for project %s", balloting.ProjectID),
		metadata, err == nil, errorMsg, time.Since(startTime).Milliseconds(),
		ipAddress, userAgent, sessionID, requestID,
	)
	
	return err
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

func (service *BallotingService) RecommendFDARS(memberId, projectId string, recommended bool) error {
	return service.repo.RecommendFDARS(memberId, projectId, recommended)
}

func (service *BallotingService) VerifyFDARSRecommendation(memberId, projectId string) error {
	return service.repo.VerifyFDARSRecommendation(memberId, projectId)
}
