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

type ProjectService struct {
	repo            *repository.ProjectRepository
	docService      *DocumentService
	auditLogService *AuditLogService
}

func NewProjectService(repo *repository.ProjectRepository, docService *DocumentService, auditLogService *AuditLogService) *ProjectService {
	return &ProjectService{
		repo:            repo,
		docService:      docService,
		auditLogService: auditLogService,
	}
}

func (service *ProjectService) CreateProject(project *models.Project, userID *string, ipAddress, userAgent, sessionID, requestID string) error {
	startTime := time.Now()
	
	// Generate project ID
	project.ID = uuid.New()
	project.CreatedAt = time.Now()
	project.ProposalApproved = false

	number, err := service.repo.GetNextAvailableNumber()
	if err != nil {
		// Log failed action
		if service.auditLogService != nil {
			service.auditLogService.LogProjectAction(
				userID, models.ActionProjectCreate, project.ID.String(), project.Title,
				map[string]interface{}{"error": "failed to get next available number"},
				false, err.Error(), time.Since(startTime).Milliseconds(),
				ipAddress, userAgent, sessionID, requestID,
			)
		}
		return err
	}

	tc, err := service.repo.GetTCByID(project.TechnicalCommitteeID)
	if err != nil {
		// Log failed action
		if service.auditLogService != nil {
			service.auditLogService.LogProjectAction(
				userID, models.ActionProjectCreate, project.ID.String(), project.Title,
				map[string]interface{}{"error": "failed to get technical committee"},
				false, err.Error(), time.Since(startTime).Milliseconds(),
				ipAddress, userAgent, sessionID, requestID,
			)
		}
		return err
	}

	project.Number = number + 1

	// Generate reference number
	project.Reference = generateProjectReference(project, tc.Code)

	stage, err := service.repo.GetStageByNumber(0)
	if err != nil {
		// Log failed action
		if service.auditLogService != nil {
			service.auditLogService.LogProjectAction(
				userID, models.ActionProjectCreate, project.ID.String(), project.Title,
				map[string]interface{}{"error": "failed to get initial stage"},
				false, err.Error(), time.Since(startTime).Milliseconds(),
				ipAddress, userAgent, sessionID, requestID,
			)
		}
		return err
	}

	// project is at stage 0
	project.StageID = stage.ID.String()

	// Save project in the repository
	err = service.repo.CreateProject(project)
	
	// Log the action
	if service.auditLogService != nil {
		metadata := map[string]interface{}{
			"project_number":     project.Number,
			"project_reference":  project.Reference,
			"technical_committee": tc.Code,
			"initial_stage":      stage.Name,
		}
		
		service.auditLogService.LogProjectAction(
			userID, models.ActionProjectCreate, project.ID.String(), project.Title,
			metadata, err == nil, 
			func() string { if err != nil { return err.Error() } else { return "" } }(),
			time.Since(startTime).Milliseconds(),
			ipAddress, userAgent, sessionID, requestID,
		)
	}
	
	return err
}

func (service *ProjectService) Exists(projectID string) (bool, error) {
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
func generateProjectReference(project *models.Project, code string) string {
	currentYear := time.Now().Year()
	return fmt.Sprintf("PWI/TC %s/%03d/%d", code, project.Number, currentYear)
}

func (service *ProjectService) GetNextAvailableNumber() (int64, error) {
	return service.repo.GetNextAvailableNumber()
}

func (service *ProjectService) UpdateProjectStage(projectID uuid.UUID, newStageID uuid.UUID, notes string) error {
	return service.repo.UpdateProjectStage(projectID, newStageID, notes)
}

func (service *ProjectService) ApproveProject(projectID string, approved bool, comment, approvedBy string) error {
	var sharepointDocID string = ""
	if approved {
		project, err := service.repo.GetProjectByID(uuid.MustParse(projectID))

		if err != nil {
			return fmt.Errorf("failed to get project by ID: %w", err)
		}

		fileName := fmt.Sprintf("PROJECT_%d/%s.docx", project.Number, strings.ReplaceAll(project.Reference, "/", "-"))

		doc, err := service.docService.UploadFileToOneDriveFolder(context.Background(), fileName)

		if err != nil {
			return fmt.Errorf("failed to upload project template: %w", err)
		}
		sharepointDocID = doc.ID
	}
	return service.repo.ApproveProject(projectID, approved, comment, approvedBy, sharepointDocID)
}

func (service *ProjectService) ApproveProjectProposal(projectID string, approved bool, comment, approvedBy, procedure string) error {
	return service.repo.ApproveProjectProposal(projectID, approved, comment, approvedBy, procedure)
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

func (service *ProjectService) UpdateProject(project *models.Project, userID *string, ipAddress, userAgent, sessionID, requestID string) error {
	startTime := time.Now()
	
	err := service.repo.UpdateProject(project)
	
	// Log the action
	if service.auditLogService != nil {
		metadata := map[string]interface{}{
			"project_number":    project.Number,
			"project_reference": project.Reference,
		}
		
		service.auditLogService.LogProjectAction(
			userID, models.ActionProjectUpdate, project.ID.String(), project.Title,
			metadata, err == nil,
			func() string { if err != nil { return err.Error() } else { return "" } }(),
			time.Since(startTime).Milliseconds(),
			ipAddress, userAgent, sessionID, requestID,
		)
	}
	
	return err
}

func (service *ProjectService) DeleteProject(projectID uuid.UUID, userID *string, ipAddress, userAgent, sessionID, requestID string) error {
	startTime := time.Now()
	
	// Get project data before deletion for audit trail
	project, err := service.repo.GetProjectByID(projectID)
	if err != nil {
		if service.auditLogService != nil {
			service.auditLogService.LogProjectAction(
				userID, models.ActionProjectDelete, projectID.String(), "Unknown Project",
				map[string]interface{}{"error": "failed to get project data before deletion"},
				false, err.Error(), time.Since(startTime).Milliseconds(),
				ipAddress, userAgent, sessionID, requestID,
			)
		}
		return err
	}
	
	err = service.repo.DeleteProject(projectID)
	
	// Log the action
	if service.auditLogService != nil {
		metadata := map[string]interface{}{
			"project_number":    project.Number,
			"project_reference": project.Reference,
			"deleted_project":   project,
		}
		
		service.auditLogService.LogProjectAction(
			userID, models.ActionProjectDelete, projectID.String(), project.Title,
			metadata, err == nil,
			func() string { if err != nil { return err.Error() } else { return "" } }(),
			time.Since(startTime).Milliseconds(),
			ipAddress, userAgent, sessionID, requestID,
		)
	}
	
	return err
}

func (service *ProjectService) FindProjects(params map[string]interface{}, limit, offset int) ([]models.Project, int64, error) {
	return service.repo.FindProjects(params, limit, offset, true)
}

func (service *ProjectService) FindProjectRequests(params map[string]interface{}, limit, offset int) ([]models.Project, int64, error) {
	return service.repo.FindProjects(params, limit, offset, false)
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

func (service *ProjectService) FetchStages() (*[]models.Stage, error) {
	return service.repo.FetchStages()
}

func (service *ProjectService) ReviewWD(secretary, projectID, comment string, status models.WorkingDraftStatus, userID *string, ipAddress, userAgent, sessionID, requestID string) error {
	err := service.repo.ReviewWD(secretary, projectID, comment, status)
	if err == nil && status == models.ACCEPTED {
		projectUUID, err := uuid.Parse(projectID)
		if err != nil {
			return err
		}
		project, err := service.GetProjectByID(projectUUID)
		if err != nil {
			return err
		}
		fileName := fmt.Sprintf("%s.docx", strings.ReplaceAll(project.Reference, "/", "-"))
		doc, errr := service.docService.CopyOneDriveFile(context.Background(), *project.SharepointDocID, fileName, project.Number)
		if errr != nil {
			return fmt.Errorf("failed to copy OneDrive file: %w", errr)
		}

		project.SharepointDocID = &doc.ID
		err = service.UpdateProject(project, userID, ipAddress, userAgent, sessionID, requestID)
		if err != nil {
			return fmt.Errorf("failed to update project after WD review: %w", err)
		}

		err = service.docService.UpdateProjectDoc(projectID, "CD", fileName, project.MemberID)
		if err != nil {
			return fmt.Errorf("failed to create CD: %w", err)
		}
	}
	return err
}

func (service *ProjectService) ReviewCD(secretary, projectId string, isConsensusReached bool, action models.ProposalAction, meetingRequired bool, userID *string, ipAddress, userAgent, sessionID, requestID string) error {
	err := service.repo.ReviewCD(secretary, projectId, isConsensusReached, action, meetingRequired)
	if err == nil && isConsensusReached {
		projectUUID, err := uuid.Parse(projectId)
		if err != nil {
			return err
		}
		project, err := service.GetProjectByID(projectUUID)
		if err != nil {
			return err
		}
		fileName := fmt.Sprintf("%s.docx", strings.ReplaceAll(project.Reference, "/", "-"))
		doc, errr := service.docService.CopyOneDriveFile(context.Background(), *project.SharepointDocID, fileName, project.Number)
		if errr != nil {
			return fmt.Errorf("failed to copy OneDrive file: %w", errr)
		}

		project.SharepointDocID = &doc.ID
		err = service.UpdateProject(project, userID, ipAddress, userAgent, sessionID, requestID)
		if err != nil {
			return fmt.Errorf("failed to update project after CD review: %w", err)
		}

		err = service.docService.UpdateProjectDoc(projectId, "DARS", fileName, project.MemberID)
		if err != nil {
			return fmt.Errorf("failed to create DARS: %w", err)
		}
	}
	return err
}

func (service *ProjectService) ReviewDARS(secretary,
	projectId string,
	wto_notification_notified bool,
	unresolvedIssues,
	alternativeDeliverable,
	status string, userID *string, ipAddress, userAgent, sessionID, requestID string) error {
	if models.DARSStatus(status) == models.DARSRejected &&
		alternativeDeliverable == "" {
		return fmt.Errorf("alternative deliverables cannot be empty when status is rejected")

	}
	err := service.repo.ReviewDARS(secretary, projectId, wto_notification_notified, unresolvedIssues, alternativeDeliverable, status)
	if err == nil && status != "" && status == string(models.DARSApproved) {
		projectUUID, err := uuid.Parse(projectId)
		if err != nil {
			return err
		}
		project, err := service.GetProjectByID(projectUUID)
		if err != nil {
			return err
		}
		fileName := fmt.Sprintf("%s.docx", strings.ReplaceAll(project.Reference, "/", "-"))
		doc, errr := service.docService.CopyOneDriveFile(context.Background(), *project.SharepointDocID, fileName, project.Number)
		if errr != nil {
			return fmt.Errorf("failed to copy OneDrive file: %w", errr)
		}

		project.SharepointDocID = &doc.ID
		err = service.UpdateProject(project, userID, ipAddress, userAgent, sessionID, requestID)
		if err != nil {
			return fmt.Errorf("failed to update project after DARS review: %w", err)
		}

		err = service.docService.UpdateProjectDoc(projectId, "FDARS", fileName, project.MemberID)
		if err != nil {
			return fmt.Errorf("failed to create FDARS: %w", err)
		}
	}
	return err
}

func (service *ProjectService) ApproveFDARS(secretary, projectId string, approve bool, action string) error {
	if action == "" && !approve {
		return fmt.Errorf("action cannot be empty when not approving")
	}
	return service.repo.ApproveFDARS(secretary, projectId, approve, action)
}

func (service *ProjectService) ApproveFDRSForPublication(secretary, projectId string, approve bool, comment string, userID *string, ipAddress, userAgent, sessionID, requestID string) error {
	err := service.repo.ApproveFDRSForPublication(secretary, projectId, approve, comment)
	if approve {
		if err == nil && approve {
			projectUUID, err := uuid.Parse(projectId)
			if err != nil {
				return err
			}
			project, err := service.GetProjectByID(projectUUID)
			if err != nil {
				return err
			}
			fileName := fmt.Sprintf("%s.docx", strings.ReplaceAll(project.Reference, "/", "-"))
			doc, errr := service.docService.CopyOneDriveFile(context.Background(), *project.SharepointDocID, fileName, project.Number)
			if errr != nil {
				return fmt.Errorf("failed to copy OneDrive file: %w", errr)
			}

			project.SharepointDocID = &doc.ID
			err = service.UpdateProject(project, userID, ipAddress, userAgent, sessionID, requestID)
			if err != nil {
				return fmt.Errorf("failed to update project after FDARS review: %w", err)
			}

			err = service.docService.UpdateProjectDoc(projectId, "ARS", fileName, project.MemberID)
			if err != nil {
				return fmt.Errorf("failed to create ARS: %w", err)
			}
		}
	}
	return err
}

func (service *ProjectService) GetDashboardStats() (map[string]any, error) {
	return service.repo.GetDashboardStats()
}

func (service *ProjectService) GetAllDistributions() (map[string]map[string]float64, error) {
	return service.repo.GetAllDistributions()
}
