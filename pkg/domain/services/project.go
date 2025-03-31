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

func (service *ProjectService) Exists(projectID uuid.UUID) (bool, error) {
	return service.repo.Exists(projectID)
}

/*
	For new WDs, shall indicate WD/TC NN/XXX/YYYY, where NN is the TC code, XXX is the serial
	number allocated to the Working Draft by the TC Secretariat and YYYY is the year of circulation.
	For WDs on revision of the standard, shall indicate WD/XXX: YYYY where XXX is the ARS
	number of the current standard and YYYY is the year of circulation. For example, when revising
	ARS 461:2021 in 2024, the corresponding drafts shall be numbered as WD/461:2024. This kind
	of numbering WDs applies also for various stages (CD, DARS and FDARS).
*/

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

func (service *ProjectService) GetNextAvailableNumber() (int64, error) {
	return service.repo.GetNextAvailableNumber()
}

func (service *ProjectService) UpdateProjectStage(projectID uuid.UUID, newStageID uuid.UUID, notes string) error {
	return service.repo.UpdateProjectStage(projectID, newStageID, notes)
}

func (service *ProjectService) GetProjectWithStageHistory(projectID uuid.UUID) (*models.Project, error) {
	return service.repo.GetProjectWithStageHistory(projectID)
}

func (service *ProjectService) GetProjectStageHistory(projectID uuid.UUID) ([]models.ProjectStageHistory, error) {
	return service.repo.GetProjectStageHistory(projectID)
}

func (service *ProjectService) FindProjectsByStage(stageID uuid.UUID, limit, offset int) ([]models.Project, int64, error) {
	return service.repo.FindProjectsByStage(stageID, limit, offset)
}

func (service *ProjectService) FindProjectsByStageTimeline(stageID uuid.UUID, startDate, endDate time.Time) ([]models.Project, error) {
	return service.repo.FindProjectsByStageTimeline(stageID, startDate, endDate)
}

func (service *ProjectService) GetProjectByID(projectID uuid.UUID) (*models.Project, error) {
	return service.repo.GetProjectByID(projectID)
}

func (service *ProjectService) UpdateProject(project *models.Project) error {
	return service.repo.UpdateProject(project)
}

func (service *ProjectService) DeleteProject(projectID uuid.UUID) error {
	return service.repo.DeleteProject(projectID)
}

func (service *ProjectService) FindProjects(params map[string]interface{}, limit, offset int) ([]models.Project, int64, error) {
	return service.repo.FindProjects(params, limit, offset)
}

func (service *ProjectService) GetProjectsByTimeframe(startDate, endDate time.Time) ([]models.Project, error) {
	return service.repo.GetProjectsByTimeframe(startDate, endDate)
}

func (service *ProjectService) GetProjectCountByType() (map[models.ProjectType]int64, error) {
	return service.repo.GetProjectCountByType()
}

func (service *ProjectService) GetProjectsWithStageTransitions(fromStageID, toStageID uuid.UUID) ([]models.Project, error) {
	return service.repo.GetProjectsWithStageTransitions(fromStageID, toStageID)
}

func (service *ProjectService) GetProjectsByReferenceBase(referenceBase string) ([]models.Project, error) {
	return service.repo.GetProjectsByReferenceBase(referenceBase)
}

func (service *ProjectService) CreateProjectRevision(baseProjectID uuid.UUID) (*models.Project, error) {
	return service.repo.CreateProjectRevision(baseProjectID)
}

func (service *ProjectService) GetProjectsApproachingDeadline(daysThreshold int) ([]models.Project, error) {
	return service.repo.GetProjectsApproachingDeadline(daysThreshold)
}

func (service *ProjectService) GetProjectsInStageForTooLong(stageID uuid.UUID, dayThreshold int) ([]models.Project, error) {
	return service.repo.GetProjectsInStageForTooLong(stageID, dayThreshold)
}

func (service *ProjectService) GetRelatedProjects(projectID uuid.UUID) ([]models.Project, error) {
	return service.repo.GetRelatedProjects(projectID)
}
