package repository

import (
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) CreateProject(project *models.Project) error {
	// If project has a stage, create initial stage history entry
	if project.StageID != nil {
		now := time.Now()
		stageHistory := models.ProjectStageHistory{
			ID:        uuid.New(),
			ProjectID: project.ID,
			StageID:   *project.StageID,
			StartedAt: now,
			CreatedAt: now,
			UpdatedAt: now,
		}
		project.StageHistory = append(project.StageHistory, stageHistory)
	}

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

	// If there's a current stage, mark it as ended in the history
	if project.StageID != nil {
		// Find the current active stage history entry and mark it as ended
		if err := tx.Model(&models.ProjectStageHistory{}).
			Where("project_id = ? AND stage_id = ? AND ended_at IS NULL", projectID, *project.StageID).
			Update("ended_at", now).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Add new stage to history
	stageHistory := models.ProjectStageHistory{
		ID:        uuid.New(),
		ProjectID: projectID,
		StageID:   newStageID,
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
		Updates(map[string]interface{}{
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
