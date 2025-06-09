package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type AcceptanceService struct {
	repo           *repository.AcceptanceRepository
	docService     *DocumentService
	projectService *ProjectService
}

func NewAcceptanceService(repo *repository.AcceptanceRepository, docService *DocumentService, projectService *ProjectService) *AcceptanceService {
	return &AcceptanceService{repo: repo, docService: docService, projectService: projectService}
}

func (service *AcceptanceService) CreateNSBResponse(response *models.NSBResponse) error {
	response.ID = uuid.New()
	response.ResponseDate = time.Now()
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

func (service *AcceptanceService) SetAcceptanceApproval(acceptance models.Acceptance) error {
	err := service.repo.SetAcceptanceApproval(acceptance)
	if err == nil {
		projectUUID, err := uuid.Parse(acceptance.ProjectID)
		if err != nil {
			return err
		}
		project, err := service.projectService.GetProjectByID(projectUUID)
		if err != nil {
			return err
		}
		fileName := fmt.Sprintf("%s.docx", strings.ReplaceAll(project.Reference, "/", "-"))
		doc, errr := service.docService.CopyOneDriveFile(context.Background(), *project.SharepointDocID, fileName, project.Number)
		if errr != nil {
			return fmt.Errorf("failed to copy OneDrive file: %w", errr)
		}
		project.SharepointDocID = &doc.ID
		err = service.projectService.UpdateProject(project)
		if err != nil {
			return fmt.Errorf("failed to update project after WD review: %w", err)
		}
	}
	return err
}

func (service *AcceptanceService) GetAcceptanceResults(id string) (*models.AcceptanceResults, error) {
	return service.repo.GetAcceptanceResults(id)
}
