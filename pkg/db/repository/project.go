package repository

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) CreateProject(project *models.Project) error {
	// Create initial stage history entry
	now := time.Now()
	stageHistory := models.ProjectStageHistory{
		ID:        uuid.New(),
		ProjectID: project.ID.String(),
		StageID:   project.StageID,
		StartedAt: now,
		CreatedAt: now,
		UpdatedAt: now,
	}
	project.StageHistory = append(project.StageHistory, stageHistory)

	return r.db.Create(project).Error
}

func (r *ProjectRepository) UpdateProjectStage(projectID uuid.UUID, newStageID uuid.UUID, notes string) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Get current project
	var project models.Project
	if err := tx.Preload("StageHistory").First(&project, "id = ?", projectID).Error; err != nil {
		tx.Rollback()
		return err
	}

	now := time.Now()

	// Find the current active stage history entry and mark it as ended
	if err := tx.Model(&models.ProjectStageHistory{}).
		Where("project_id = ? AND stage_id = ? AND ended_at IS NULL", projectID, project.StageID).
		Update("ended_at", now).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add new stage to history
	stageHistory := models.ProjectStageHistory{
		ID:        uuid.New(),
		ProjectID: projectID.String(),
		StageID:   newStageID.String(),
		StartedAt: now,
		Notes:     notes,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := tx.Create(&stageHistory).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update current stage of the project
	if err := tx.Model(&models.Project{}).
		Where("id = ?", projectID).
		Updates(map[string]any{
			"stage_id":   newStageID,
			"updated_at": now,
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *ProjectRepository) GetProjectWithStageHistory(projectID uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("StageHistory.Stage").
		Preload("Stage").
		Preload("TechnicalCommittee").
		Preload("WorkingGroup").
		First(&project, "id = ?", projectID).Error

	return &project, err
}

func (r *ProjectRepository) GetProjectStageHistory(projectID uuid.UUID) ([]models.ProjectStageHistory, error) {
	var stageHistory []models.ProjectStageHistory
	err := r.db.Preload("Stage").
		Where("project_id = ?", projectID).
		Order("started_at ASC").
		Find(&stageHistory).Error

	return stageHistory, err
}

func (r *ProjectRepository) FindProjectsByStage(stageID uuid.UUID, limit, offset int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	query := r.db.Model(&models.Project{}).Where("stage_id = ?", stageID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get projects with pagination
	err := query.Preload("Stage").
		Preload("TechnicalCommittee").
		Preload("WorkingGroup").
		Limit(limit).
		Offset(offset).
		Find(&projects).Error

	return projects, total, err
}

func (r *ProjectRepository) FindProjectsByStageTimeline(stageID uuid.UUID, startDate, endDate time.Time) ([]models.Project, error) {
	var projectIDs []uuid.UUID

	// Find projects that were in the specified stage during the given time period
	if err := r.db.Model(&models.ProjectStageHistory{}).
		Where("stage_id = ? AND started_at <= ? AND (ended_at >= ? OR ended_at IS NULL)",
			stageID, endDate, startDate).
		Pluck("project_id", &projectIDs).Error; err != nil {
		return nil, err
	}

	// Get the projects with their associations
	var projects []models.Project
	if len(projectIDs) > 0 {
		err := r.db.Preload("Stage").
			Preload("TechnicalCommittee").
			Preload("WorkingGroup").
			Where("id IN ?", projectIDs).
			Find(&projects).Error
		return projects, err
	}

	return projects, nil
}

func (r *ProjectRepository) Exists(projectID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Project{}).Where("id = ?", projectID).Count(&count).Error
	return count > 0, err
}

func (r *ProjectRepository) GetNextAvailableNumber() (int64, error) {
	var count int64
	err := r.db.Model(&models.Project{}).Count(&count).Error
	return count + 1, err
}

// GetProjectByID retrieves a project by its ID
func (r *ProjectRepository) GetProjectByID(projectID uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.Where("id = ?", projectID).Preload(clause.Associations).
		First(&project).Error
	return &project, err
}

func (r *ProjectRepository) GetTCByID(id string) (*models.TechnicalCommittee, error) {
	var tc models.TechnicalCommittee
	err := r.db.Where("id = ?", id).First(&tc).Error
	return &tc, err
}

func (r *ProjectRepository) UpdateProject(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *ProjectRepository) ReviewWD(secretary, projectID, comment string, status models.WorkingDraftStatus) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// First, fetch the entire project
		var project models.Project
		if err := tx.Where("id = ?", projectID).Preload("TechnicalCommittee").First(&project).Error; err != nil {
			return fmt.Errorf("project with ID %s not found: %w", projectID, err)
		}

		if project.TechnicalCommittee.SecretaryId == nil || *project.TechnicalCommittee.SecretaryId != secretary {
			tx.Rollback()
			return fmt.Errorf("User is not allowed to perform this action")
		}

		// Update fields directly on the project object
		project.WorkingDraftStatus = status
		project.WorkingDraftComments = comment
		project.WDTCSecretaryID = &secretary

		// If status is ACCEPTED, prepare additional changes
		if status == models.ACCEPTED {
			var stage models.Stage
			if err := tx.Where("number = ?", 3).First(&stage).Error; err != nil {
				return err
			}

			// First save initial changes to the project
			if err := tx.Save(&project).Error; err != nil {
				return err
			}

			// Update the stage (which will update the reference)
			if err := UpdateProjectStageWithTx(tx, projectID, stage.ID.String(), "WD Elevated to a CD", "WD", stage.Abbreviation); err != nil {
				return err
			}

			// Fetch document and create new one
			var document models.Document
			if err := tx.Where("id = ?", project.WorkingDraftID).First(&document).Error; err != nil {
				return err
			}

			doc := models.Document{
				ID:          uuid.New(),
				CreatedByID: document.CreatedByID,
				Title:       document.Title,
				Description: document.Description,
				Reference:   strings.ReplaceAll(document.Reference, "WD", stage.Abbreviation),
				FileURL:     document.FileURL,
				CreatedAt:   time.Now(),
			}

			if err := tx.Create(&doc).Error; err != nil {
				return err
			}

			// Fetch the updated project (after stage update)
			if err := tx.Where("id = ?", projectID).First(&project).Error; err != nil {
				return err
			}

			// Update document ID
			docId := doc.ID.String()
			project.CommitteeDraftID = &docId

			// Save final changes
			if err := tx.Save(&project).Error; err != nil {
				return err
			}
		} else {
			// For non-ACCEPTED status, just save the initial changes
			if err := tx.Save(&project).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *ProjectRepository) ApproveProject(projectID string, approved bool, comment, approvedBy string) error {
	tx := r.db.Begin() // Start a transaction
	if tx.Error != nil {
		return tx.Error
	}

	var project models.Project

	// Retrieve the project by ID
	if err := tx.Preload("TechnicalCommittee").First(&project, "id = ?", projectID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if project.TechnicalCommittee.SecretaryId == nil || *project.TechnicalCommittee.SecretaryId != approvedBy {
		tx.Rollback()
		return fmt.Errorf("User is not allowed to perform this action")
	}

	// Update approval status and comment
	project.PWIApproved = approved
	project.PWIApprovalComment = comment
	project.ApprovedByID = &approvedBy

	// Save the updated project
	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update the current stage as "ended"
	now := time.Now()

	// Find the active stage and mark it as ended
	if err := tx.Model(&models.ProjectStageHistory{}).
		Where("project_id = ? AND stage_id = ? AND ended_at IS NULL", projectID, project.StageID).
		Update("ended_at", now).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update the project's current stage status
	if err := tx.Model(&models.Project{}).
		Where("id = ?", projectID).
		Update("updated_at", now).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *ProjectRepository) DeleteProject(projectID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete related stage history first
		if err := tx.Where("project_id = ?", projectID).Delete(&models.ProjectStageHistory{}).Error; err != nil {
			return err
		}
		// Then delete the project
		return tx.Delete(&models.Project{}, "id = ?", projectID).Error
	})
}

// FindProjects searches for projects with optional filters
func (r *ProjectRepository) FindProjects(params map[string]any, limit, offset int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	query := r.db.Model(&models.Project{})

	// Apply filters
	if title, ok := params["title"].(string); ok && title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	if projectType, ok := params["type"].(models.ProjectType); ok && projectType != "" {
		query = query.Where("type = ?", projectType)
	}

	if committeeID, ok := params["committee_id"].(uuid.UUID); ok && committeeID != uuid.Nil {
		query = query.Where("technical_committee_id = ?", committeeID)
	}

	if workingGroupID, ok := params["working_group_id"].(uuid.UUID); ok && workingGroupID != uuid.Nil {
		query = query.Where("working_group_id = ?", workingGroupID)
	}

	if visible, ok := params["visible_on_library"].(bool); ok {
		query = query.Where("visible_on_library = ?", visible)
	}

	if emergency, ok := params["is_emergency"].(bool); ok {
		query = query.Where("is_emergency = ?", emergency)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get projects with pagination
	err := query.
		Preload(clause.Associations).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&projects).Error

	return projects, total, err
}

// GetProjectsByTimeframe retrieves projects created within a specific timeframe
func (r *ProjectRepository) GetProjectsByTimeframe(startDate, endDate time.Time) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Preload("Stage").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Find(&projects).Error
	return projects, err
}

// GetProjectCountByType returns the number of projects grouped by project type
func (r *ProjectRepository) GetProjectCountByType() (map[models.ProjectType]int64, error) {
	var results []struct {
		Type  models.ProjectType `json:"type"`
		Count int64              `json:"count"`
	}

	err := r.db.Model(&models.Project{}).
		Select("type, count(*) as count").
		Group("type").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	countByType := make(map[models.ProjectType]int64)
	for _, result := range results {
		countByType[result.Type] = result.Count
	}

	return countByType, nil
}

// GetProjectsWithStageTransitions returns projects that have transitions between given stages
func (r *ProjectRepository) GetProjectsWithStageTransitions(fromStageID, toStageID uuid.UUID) ([]models.Project, error) {
	var projectIDs []uuid.UUID

	// Find projects that have transitioned from one stage to another
	subquery := r.db.Table("project_stage_histories as psh1").
		Joins("JOIN project_stage_histories as psh2 ON psh1.project_id = psh2.project_id").
		Where("psh1.stage_id = ? AND psh2.stage_id = ? AND psh1.ended_at IS NOT NULL AND psh2.started_at = psh1.ended_at",
			fromStageID, toStageID).
		Select("DISTINCT psh1.project_id")

	if err := r.db.Table("(?) as subquery", subquery).Pluck("project_id", &projectIDs).Error; err != nil {
		return nil, err
	}

	var projects []models.Project
	if len(projectIDs) > 0 {
		err := r.db.Preload("Stage").
			Preload("TechnicalCommittee").
			Preload("WorkingGroup").
			Where("id IN ?", projectIDs).
			Find(&projects).Error
		return projects, err
	}

	return projects, nil
}

// GetProjectsByReferenceBase finds all projects with the same reference base
func (r *ProjectRepository) GetProjectsByReferenceBase(referenceBase string) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Where("reference LIKE ?", referenceBase+"%").
		Order("edition_no DESC").
		Find(&projects).Error
	return projects, err
}

// CreateProjectRevision creates a new revision of an existing project
func (r *ProjectRepository) CreateProjectRevision(baseProjectID uuid.UUID) (*models.Project, error) {
	// Get the base project
	baseProject, err := r.GetProjectByID(baseProjectID)
	if err != nil {
		return nil, err
	}

	// Create a new project as a revision
	newProject := models.Project{
		ID:                   uuid.New(),
		Number:               baseProject.Number,
		PartNo:               baseProject.PartNo,
		EditionNo:            baseProject.EditionNo + 1, // Increment edition number
		Reference:            baseProject.Reference,
		ReferenceSuffix:      baseProject.ReferenceSuffix,
		Title:                baseProject.Title,
		Description:          baseProject.Description,
		TechnicalCommitteeID: baseProject.TechnicalCommitteeID,
		WorkingGroupID:       baseProject.WorkingGroupID,
		Timeframe:            baseProject.Timeframe,
		Type:                 models.REVISION, // Mark as revision
		VisibleOnLibrary:     false,           // Not visible until published
		PricePerPage:         baseProject.PricePerPage,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	err = r.db.Create(&newProject).Error
	if err != nil {
		return nil, err
	}

	return &newProject, nil
}

// GetProjectsApproachingDeadline finds projects nearing their expected completion date
func (r *ProjectRepository) GetProjectsApproachingDeadline(daysThreshold int) ([]models.Project, error) {
	var projects []models.Project

	// Calculate the deadline based on created_at + timeframe
	now := time.Now()
	query := r.db.Preload("Stage").
		Where("visible_on_library = ? AND DATE_PART('day', (created_at + (timeframe * interval '1 month')) - ?) <= ?",
			false, now, daysThreshold).
		Find(&projects)

	return projects, query.Error
}

// GetProjectsInStageForTooLong finds projects that have been in a particular stage longer than expected
func (r *ProjectRepository) GetProjectsInStageForTooLong(stageID uuid.UUID, dayThreshold int) ([]models.Project, error) {
	var projects []models.Project

	// Find projects that have been in the current stage for longer than the threshold
	subquery := r.db.Table("project_stage_histories").
		Select("project_id").
		Where("stage_id = ? AND ended_at IS NULL AND DATE_PART('day', NOW() - started_at) > ?",
			stageID, dayThreshold)

	err := r.db.Preload("Stage").
		Preload("TechnicalCommittee").
		Preload("WorkingGroup").
		Where("stage_id = ? AND id IN (?)", stageID, subquery).
		Find(&projects).Error

	return projects, err
}

// GetRelatedProjects finds projects related to the given project
// For example, different parts of the same standard
func (r *ProjectRepository) GetRelatedProjects(projectID uuid.UUID) ([]models.Project, error) {
	var project models.Project
	if err := r.db.First(&project, "id = ?", projectID).Error; err != nil {
		return nil, err
	}

	var relatedProjects []models.Project
	err := r.db.Where("reference = ? AND id != ?", project.Reference, projectID).
		Or("reference_suffix = ? AND reference_suffix != '' AND id != ?", project.ReferenceSuffix, projectID).
		Find(&relatedProjects).Error

	return relatedProjects, err
}

func (r *ProjectRepository) FetchStages() (*[]models.Stage, error) {
	var stages []models.Stage
	err := r.db.Preload(clause.Associations).
		Preload("Timeframe.Standard").
		Preload("Timeframe.IS").
		Preload("Timeframe.Emergency").
		Find(&stages).Error
	return &stages, err
}

func (r *ProjectRepository) GetStageByNumber(number int16) (*models.Stage, error) {
	var stage models.Stage
	err := r.db.Where("number = ?", number).Preload(clause.Associations).First(&stage).Error
	return &stage, err
}

func (repo *ProjectRepository) FindByDocumentID(documentID uuid.UUID) ([]models.Project, error) {
	var projects []models.Project

	// Query for projects where either WorkingDraftID or CommitteeDraftID matches the given documentID
	err := repo.db.Where("working_draft_id = ? OR committee_draft_id = ?", documentID, documentID).Find(&projects).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find projects with document ID %s: %w", documentID, err)
	}

	return projects, nil
}

func (r *ProjectRepository) ReviewCD(secretary, projectId string, isConsensusReached bool, action models.ProposalAction, meetingRequired bool) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var project models.Project
		if err := tx.Where("id = ?", projectId).Preload("TechnicalCommittee").First(&project).Error; err != nil {
			return err
		}

		// if project.TechnicalCommittee.SecretaryId == nil || *project.TechnicalCommittee.SecretaryId != secretary {
		// 	tx.Rollback()
		// 	return fmt.Errorf("User is not allowed to perform this action")
		// }

		project.IsConsensusReached = isConsensusReached
		project.ProposalAction = action
		project.MeetingRequired = meetingRequired
		project.CDTCSecretaryID = &secretary

		if isConsensusReached {
			now := time.Now()
			project.SubmissionDate = &now

			var stage models.Stage
			if err := tx.Where("number = ?", 4).First(&stage).Error; err != nil {
				return err
			}

			// First save the current changes
			if err := tx.Save(&project).Error; err != nil {
				return err
			}

			dars := models.DARS{
				ID:                    uuid.New(),
				ProjectID:             projectId,
				CreatedAt:             time.Now(),
				PublicReviewStartDate: now,
				PublicReviewEndDate:   now.AddDate(0, 2, 7),
			}

			if err := tx.Create(&dars).Error; err != nil {
				return err
			}

			// Then update the stage (which will fetch and update the project again)
			if err := UpdateProjectStageWithTx(tx, projectId, stage.ID.String(), "CD Consensus reached", "CD", stage.Abbreviation); err != nil {
				return err
			}

			// Don't save again after this point
			return nil
		}

		// Only save if we didn't reach consensus
		if err := tx.Save(&project).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *ProjectRepository) ReviewDARS(secretary,
	projectId string,
	wto_notification_notified bool,
	unresolvedIssues,
	alternativeDeliverable,
	status string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var project models.Project
		if err := tx.Where("id = ?", projectId).Preload("TechnicalCommittee").First(&project).Error; err != nil {
			return err
		}

		if project.TechnicalCommittee.SecretaryId == nil || *project.TechnicalCommittee.SecretaryId != secretary {
			tx.Rollback()
			return fmt.Errorf("User is not allowed to perform this action")
		}

		var dars models.DARS
		if err := tx.Where("project_id = ?", projectId).First(&dars).Error; err != nil {
			return err
		}

		dars.DARSTCSecretaryID = &secretary
		if unresolvedIssues != "" {
			dars.UnresolvedIssues = unresolvedIssues

		}

		if unresolvedIssues != "" {
			dars.AlternativeDeliverable = alternativeDeliverable
		}

		if status != "" {
			dars.Status = models.DARSStatus(status)

		}
		if wto_notification_notified && dars.WTONotificationDate == nil {
			now := time.Now()
			dars.WTONotificationDate = &now
		}

		if status != "" && status == string(models.DARSApproved) {
			dars.MoveToBalloting = true
			var stage models.Stage
			if err := tx.Where("number = ?", 5).First(&stage).Error; err != nil {
				return err
			}

			// First save the current changes
			if err := tx.Save(&dars).Error; err != nil {
				return err
			}

			balloting := models.Balloting{
				ID:        uuid.New(),
				ProjectID: projectId,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(0, 0, 30), // Set end date to 30 days from now
			}

			if err := tx.Create(&balloting).Error; err != nil {
				return err
			}

			// Then update the stage (which will fetch and update the project again)
			if err := UpdateProjectStageWithTx(tx, projectId, stage.ID.String(), "DARS is accepted to advance to the balloting stage as an FDARS", "DARS", stage.Abbreviation); err != nil {
				return err
			}

			// Don't save again after this point
			return nil
		}

		// Only save if we didn't reach consensus
		if err := tx.Save(&dars).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *ProjectRepository) ApproveFDARS(secretary,
	projectId string, approve bool, action string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var balloting models.Balloting
		if err := tx.Where("project_id = ?", projectId).First(&balloting).Error; err != nil {
			return err
		}

		now := time.Now()

		balloting.ApprovedByID = &secretary
		balloting.Approved = approve
		balloting.ApprovedAt = &now

		if action != "" {
			balloting.NextCourseOfAction = models.FDARSAction(action)

			if models.FDARSAction(action) == models.CANCELLED {
				// get project and update it
				var project models.Project
				if err := tx.Where("id = ?", projectId).First(&project).Error; err != nil {
					return err
				}
				project.Cancelled = true
				project.CancelledDate = &now

				if err := tx.Save(&project).Error; err != nil {
					return err
				}
			}

		}

		if approve {
			var stage models.Stage
			if err := tx.Where("number = ?", 6).First(&stage).Error; err != nil {
				return err
			}

			if err := UpdateProjectStageWithTx(tx, projectId, stage.ID.String(), "FDARS is approved in accordance with the conditions in 7.7.3", "FDARS", stage.Abbreviation); err != nil {
				return err
			}
		}

		if err := tx.Save(&balloting).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *ProjectRepository) ApproveFDRSForPublication(secretary, projectId string, approve bool, comment string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var project models.Project
		if err := tx.Where("id = ?", projectId).Preload("TechnicalCommittee").First(&project).Error; err != nil {
			return err
		}

		if project.TechnicalCommittee.SecretaryId == nil || *project.TechnicalCommittee.SecretaryId != secretary {
			tx.Rollback()
			return fmt.Errorf("User is not allowed to perform this action")
		}

		now := time.Now()

		project.ApprovedForPublicationByID = &secretary
		project.ApprovedForPublication = approve
		project.ApprovedForPublicationDate = &now
		project.ApprovedForPublicationComment = comment

		if approve {
			project.Reference = fmt.Sprintf("ARS %03d:%d", project.Number, time.Now().Year())
		}

		if err := tx.Save(&project).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *ProjectRepository) GetDashboardStats() (map[string]any, error) {
	stats := make(map[string]any)

	// Get total projects count
	var totalCount int64
	if err := r.db.Model(&models.Project{}).Where("cancelled = ?", false).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["total"] = totalCount

	// Get counts by stage
	stageStats := make(map[string]int)
	stageRows, err := r.db.Model(&models.Project{}).
		Select("s.name as stage_name, COUNT(projects.id) as count").
		Joins("JOIN stages s ON projects.stage_id = s.id").
		Where("projects.cancelled = ?", false).
		Group("s.name").
		Rows()

	if err != nil {
		return nil, err
	}
	defer stageRows.Close()

	for stageRows.Next() {
		var stageName string
		var count int
		if err := stageRows.Scan(&stageName, &count); err != nil {
			return nil, err
		}
		stageStats[stageName] = count
	}
	stats["by_stage"] = stageStats

	// Get counts by project type
	typeStats := make(map[string]int)
	typeRows, err := r.db.Model(&models.Project{}).
		Select("type, COUNT(id) as count").
		Where("cancelled = ?", false).
		Group("type").
		Rows()

	if err != nil {
		return nil, err
	}
	defer typeRows.Close()

	for typeRows.Next() {
		var typeName string
		var count int
		if err := typeRows.Scan(&typeName, &count); err != nil {
			return nil, err
		}
		typeStats[typeName] = count
	}
	stats["by_type"] = typeStats

	// Get counts by working draft status
	wdStats := make(map[string]int)
	wdRows, err := r.db.Model(&models.Project{}).
		Select("working_draft_status as status, COUNT(id) as count").
		Where("cancelled = ?", false).
		Where("working_draft_status IS NOT NULL").
		Group("working_draft_status").
		Rows()

	if err != nil {
		return nil, err
	}
	defer wdRows.Close()

	for wdRows.Next() {
		var status string
		var count int
		if err := wdRows.Scan(&status, &count); err != nil {
			return nil, err
		}
		wdStats[status] = count
	}
	stats["by_working_draft_status"] = wdStats

	// Get counts by technical committee
	tcStats := make(map[string]int)
	tcRows, err := r.db.Model(&models.Project{}).
		Select("tc.name as committee_name, COUNT(projects.id) as count").
		Joins("JOIN technical_committees tc ON projects.technical_committee_id = tc.id").
		Where("projects.cancelled = ?", false).
		Group("tc.name").
		Rows()

	if err != nil {
		return nil, err
	}
	defer tcRows.Close()

	for tcRows.Next() {
		var committeeName string
		var count int
		if err := tcRows.Scan(&committeeName, &count); err != nil {
			return nil, err
		}
		tcStats[committeeName] = count
	}
	stats["by_committee"] = tcStats

	// Get counts for emergency projects
	var emergencyCount int64
	if err := r.db.Model(&models.Project{}).
		Where("cancelled = ?", false).
		Where("is_emergency = ?", true).
		Count(&emergencyCount).Error; err != nil {
		return nil, err
	}
	stats["emergency_projects"] = emergencyCount

	// Get counts for projects approved for publication
	var publishedCount int64
	if err := r.db.Model(&models.Project{}).
		Where("cancelled = ?", false).
		Where("approved_for_publication = ?", true).
		Count(&publishedCount).Error; err != nil {
		return nil, err
	}
	stats["published_projects"] = publishedCount

	return stats, nil
}

func (r *ProjectRepository) GetAllDistributions() (map[string]map[string]float64, error) {
	allDistributions := make(map[string]map[string]float64)

	// Get total count of non-cancelled projects
	var totalCount int64
	if err := r.db.Model(&models.Project{}).
		Where("cancelled = ?", false).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// If there are no projects, return empty distributions
	if totalCount == 0 {
		allDistributions["by_stage"] = make(map[string]float64)
		allDistributions["by_type"] = make(map[string]float64)
		allDistributions["by_working_draft_status"] = make(map[string]float64)
		allDistributions["by_committee"] = make(map[string]float64)
		return allDistributions, nil
	}

	// Get stage distribution
	stageDistribution := make(map[string]float64)
	stageRows, err := r.db.Model(&models.Project{}).
		Select("s.name as stage_name, COUNT(projects.id) as count").
		Joins("JOIN stages s ON projects.stage_id = s.id").
		Where("projects.cancelled = ?", false).
		Group("s.name").
		Rows()

	if err != nil {
		return nil, err
	}
	defer stageRows.Close()

	for stageRows.Next() {
		var stageName string
		var count int64
		if err := stageRows.Scan(&stageName, &count); err != nil {
			return nil, err
		}

		percentage := float64(count) / float64(totalCount) * 100
		stageDistribution[stageName] = math.Round(percentage*100) / 100
	}
	allDistributions["by_stage"] = stageDistribution

	// Get type distribution
	typeDistribution := make(map[string]float64)
	typeRows, err := r.db.Model(&models.Project{}).
		Select("type, COUNT(id) as count").
		Where("cancelled = ?", false).
		Group("type").
		Rows()

	if err != nil {
		return nil, err
	}
	defer typeRows.Close()

	for typeRows.Next() {
		var typeName string
		var count int64
		if err := typeRows.Scan(&typeName, &count); err != nil {
			return nil, err
		}

		percentage := float64(count) / float64(totalCount) * 100
		typeDistribution[typeName] = math.Round(percentage*100) / 100
	}
	allDistributions["by_type"] = typeDistribution

	// Get working draft status distribution
	wdDistribution := make(map[string]float64)
	wdRows, err := r.db.Model(&models.Project{}).
		Select("working_draft_status as status, COUNT(id) as count").
		Where("cancelled = ?", false).
		Where("working_draft_status IS NOT NULL").
		Group("working_draft_status").
		Rows()

	if err != nil {
		return nil, err
	}
	defer wdRows.Close()

	for wdRows.Next() {
		var status string
		var count int64
		if err := wdRows.Scan(&status, &count); err != nil {
			return nil, err
		}

		percentage := float64(count) / float64(totalCount) * 100
		wdDistribution[status] = math.Round(percentage*100) / 100
	}
	allDistributions["by_working_draft_status"] = wdDistribution

	// Get technical committee distribution
	tcDistribution := make(map[string]float64)
	tcRows, err := r.db.Model(&models.Project{}).
		Select("tc.name as committee_name, COUNT(projects.id) as count").
		Joins("JOIN technical_committees tc ON projects.technical_committee_id = tc.id").
		Where("projects.cancelled = ?", false).
		Group("tc.name").
		Rows()

	if err != nil {
		return nil, err
	}
	defer tcRows.Close()

	for tcRows.Next() {
		var committeeName string
		var count int64
		if err := tcRows.Scan(&committeeName, &count); err != nil {
			return nil, err
		}

		percentage := float64(count) / float64(totalCount) * 100
		tcDistribution[committeeName] = math.Round(percentage*100) / 100
	}
	allDistributions["by_committee"] = tcDistribution

	return allDistributions, nil
}
