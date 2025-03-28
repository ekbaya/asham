package services

import (
	"fmt"
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (service *ProjectService) CreateProject(project *models.Project) error {
	// Generate project ID
	project.ID = uuid.New()
	project.CreatedAt = time.Now()

	// Generate reference number
	project.Reference = generateProjectReference(project)

	// Save project in the repository
	return service.repo.CreateProject(project)
}

// generateProjectReference generates the reference number for a project
func generateProjectReference(project *models.Project) string {
	currentYear := time.Now().Year()

	if project.PartNo > 0 && project.TechnicalCommittee != nil {
		// New Working Drafts: WD/TC NN/XXX/YYYY
		return fmt.Sprintf("WD/TC %d/%d/%d", project.TechnicalCommittee.Code, project.Number, currentYear)
	}

	// Revisions: WD/XXX:YYYY
	return fmt.Sprintf("WD/%d:%d", project.Number, currentYear)
}
