package repository

import (
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Helper method to update project stage within a transaction
func UpdateProjectStageWithTx(tx *gorm.DB, projectID string, newStageID string, notes, currentDoc, newDoc string) error {
	// Direct SQL update to avoid fetching the project again
	if err := tx.Exec("UPDATE projects SET reference = REPLACE(reference, ?, ?), stage_id = ?, updated_at = ? WHERE id = ?",
		currentDoc, newDoc, newStageID, time.Now(), projectID).Error; err != nil {
		return err
	}

	// Handle stage history separately
	now := time.Now()

	// Close previous stage
	if err := tx.Exec("UPDATE project_stage_histories SET ended_at = ? WHERE project_id = ? AND ended_at IS NULL",
		now, projectID).Error; err != nil {
		return err
	}

	// Add new stage history
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

	return nil
}
