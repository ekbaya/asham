package repository

import (
	"strings"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Helper method to update project stage within a transaction
func UpdateProjectStageWithTx(tx *gorm.DB, projectID string, newStageID string, notes, currentDoc, newDoc string) error {
	// Get current project
	var project models.Project
	if err := tx.Preload("StageHistory").First(&project, "id = ?", projectID).Error; err != nil {
		return err
	}

	now := time.Now()

	// Find the current active stage history entry and mark it as ended
	if err := tx.Model(&models.ProjectStageHistory{}).
		Where("project_id = ? AND stage_id = ? AND ended_at IS NULL", projectID, project.StageID).
		Update("ended_at", now).Error; err != nil {
		return err
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
		return err
	}

	// Change Ref to $newDoc
	reference := strings.ReplaceAll(project.Reference, currentDoc, newDoc)

	// Update current stage of the project
	if err := tx.Model(&models.Project{}).
		Where("id = ?", projectID).
		Updates(map[string]any{
			"stage_id":   newStageID,
			"updated_at": now,
			"reference":  reference,
		}).Error; err != nil {
		return err
	}

	return nil
}
