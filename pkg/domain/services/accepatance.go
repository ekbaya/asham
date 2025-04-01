package services

import (
	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type AcceptanceService struct {
	repo *repository.AcceptanceRepository
}

func NewAcceptanceService(repo *repository.AcceptanceRepository) *AcceptanceService {
	return &AcceptanceService{repo: repo}
}

func (service *AcceptanceService) CreateNSBResponse(response *models.NSBResponse) error {
	response.ID = uuid.New()
	return service.repo.CreateNSBResponse(response)
}

func (service *AcceptanceService) GetNSBResponse(id string) (*models.NSBResponse, error) {
	return service.repo.GetNSBResponse(id)
}

func (service *AcceptanceService) GetNSBResponsesByProjectID(projectID string) ([]models.NSBResponse, error) {
	return service.repo.GetNSBResponsesByProjectID(projectID)
}

func (service *AcceptanceService) UpdateNSBResponse(response *models.NSBResponse) error {
	return service.repo.UpdateNSBResponse(response)
}

func (service *AcceptanceService) DeleteNSBResponse(id string) error {
	return service.repo.DeleteNSBResponse(id)
}

func (service *AcceptanceService) GetAcceptance(id string) (*models.Acceptance, error) {
	return service.repo.GetAcceptance(id)
}

func (service *AcceptanceService) GetAcceptanceByProjectID(projectID string) (*models.Acceptance, error) {
	return service.repo.GetAcceptanceByProjectID(projectID)
}

func (service *AcceptanceService) GetAcceptances() (*[]models.Acceptance, error) {
	return service.repo.GetAcceptances()
}

func (service *AcceptanceService) UpdateAcceptance(acceptance *models.Acceptance) error {
	return service.repo.UpdateAcceptance(acceptance)
}

func (service *AcceptanceService) GetAcceptanceWithResponses(id string) (*models.Acceptance, error) {
	return service.repo.GetAcceptanceWithResponses(id)
}

func (service *AcceptanceService) CountNSBResponsesByType(projectID string) (map[models.Response]int, error) {
	return service.repo.CountNSBResponsesByType(projectID)
}

func (service *AcceptanceService) CalculateNSBResponseStats(projectID string) error {
	return service.repo.CalculateNSBResponseStats(projectID)
}

func (service *AcceptanceService) SetAcceptanceApproval(id string, approved bool) error {
	return service.repo.SetAcceptanceApproval(id, approved)
}

func (service *AcceptanceService) GetAcceptanceResults(id string) (*models.AcceptanceResults, error) {
	return service.repo.GetAcceptanceResults(id)
}
