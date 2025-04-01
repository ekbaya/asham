package services

import (
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type ProposalService struct {
	repo *repository.ProposalRepository
}

func NewProposalService(repo *repository.ProposalRepository) *ProposalService {
	return &ProposalService{repo: repo}
}

func (service *ProposalService) Create(proposal *models.Proposal) error {
	proposal.ID = uuid.New()
	proposal.CreatedAt = time.Now()
	return service.repo.Create(proposal)
}

func (service *ProposalService) GetByID(id uuid.UUID) (*models.Proposal, error) {
	return service.repo.GetByID(id)
}

func (service *ProposalService) GetByProjectID(projectID uuid.UUID) (*models.Proposal, error) {
	return service.repo.GetByProjectID(projectID)
}

func (service *ProposalService) GetByProposingNSB(nsbID uuid.UUID, limit, offset int) ([]models.Proposal, int64, error) {
	return service.repo.GetByProposingNSB(nsbID, limit, offset)
}

func (service *ProposalService) GetByCreator(memberID uuid.UUID, limit, offset int) ([]models.Proposal, int64, error) {
	return service.repo.GetByCreator(memberID, limit, offset)
}

func (service *ProposalService) Update(proposal *models.Proposal) error {
	return service.repo.Update(proposal)
}

func (service *ProposalService) UpdatePartial(id uuid.UUID, updates map[string]interface{}) error {
	return service.repo.UpdatePartial(id, updates)
}

func (service *ProposalService) Delete(id uuid.UUID) error {
	return service.repo.Delete(id)
}

func (service *ProposalService) List(limit, offset int) ([]models.Proposal, int64, error) {
	return service.repo.List(limit, offset)
}

func (service *ProposalService) Search(query string, limit, offset int) ([]models.Proposal, int64, error) {
	return service.repo.Search(query, limit, offset)
}

func (service *ProposalService) Exists(projectID string) (bool, error) {
	return service.repo.Exists(projectID)
}

func (service *ProposalService) AddReferencedStandard(proposalID, documentID uuid.UUID) error {
	return service.repo.AddReferencedStandard(proposalID, documentID)
}

func (service *ProposalService) RemoveReferencedStandard(proposalID, documentID uuid.UUID) error {
	return service.repo.RemoveReferencedStandard(proposalID, documentID)
}

func (service *ProposalService) GetProposalCountByNSB() (map[uuid.UUID]int64, error) {
	return service.repo.GetProposalCountByNSB()
}

func (service *ProposalService) Transfer(proposalID uuid.UUID, newProjectID string) error {
	return service.repo.Transfer(proposalID, newProjectID)
}
