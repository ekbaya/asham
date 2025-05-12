package repository

import (
	"github.com/ekbaya/asham/pkg/domain/models"
	"gorm.io/gorm"
)

type StandardRepository struct {
	db *gorm.DB
}

func NewStandardRepository(db *gorm.DB) *StandardRepository {
	return &StandardRepository{db: db}
}

// Create a new Standard
func (r *StandardRepository) CreateStandard(standard *models.Standard) error {
	return r.db.Create(standard).Error
}

// Update Standard and store version history
func (r *StandardRepository) SaveStandard(standard *models.Standard) error {
	// Save current content to history before updating
	version := models.StandardVersion{
		StandardID: standard.ID,
		Content:    standard.Content,
		Version:    standard.Version,
		SavedBy:    standard.UpdatedBy,
	}

	if err := r.db.Create(&version).Error; err != nil {
		return err
	}

	// Increment version and update the main standard
	standard.Version += 1
	return r.db.Save(standard).Error
}

// Fetch a Standard by ID
func (r *StandardRepository) GetStandardByID(id string) (*models.Standard, error) {
	var standard models.Standard
	if err := r.db.First(&standard, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &standard, nil
}

// Get all versions for a given Standard ID
func (r *StandardRepository) GetStandardVersions(standardID string) ([]models.StandardVersion, error) {
	var versions []models.StandardVersion
	if err := r.db.Where("standard_id = ?", standardID).Order("version desc").Find(&versions).Error; err != nil {
		return nil, err
	}
	return versions, nil
}

// Restore a specific version
func (r *StandardRepository) RestoreStandardVersion(standardID string, version int) error {
	var oldVersion models.StandardVersion
	if err := r.db.Where("standard_id = ? AND version = ?", standardID, version).First(&oldVersion).Error; err != nil {
		return err
	}

	return r.db.Model(&models.Standard{}).
		Where("id = ?", standardID).
		Updates(map[string]interface{}{
			"content":    oldVersion.Content,
			"version":    oldVersion.Version,
			"updated_by": oldVersion.SavedBy,
		}).Error
}
