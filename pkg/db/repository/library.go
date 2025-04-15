package repository

import (
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LibraryRepository struct {
	db *gorm.DB
}

func NewLibraryRepository(db *gorm.DB) *LibraryRepository {
	return &LibraryRepository{db: db}
}

func (r *LibraryRepository) FindStandards(params map[string]any, limit, offset int) ([]models.Project, error) {
	var standards []models.Project

	query := r.db.Where("published = ?", true)

	if sector, ok := params["sector"].(string); ok && sector != "" {
		query = query.Where("sector = ?", models.ProjectSector(sector))
	}

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

	// Apply pagination
	result := query.Limit(limit).Offset(offset).Preload("Standard").Preload("TechnicalCommittee").Find(&standards)

	if result.Error != nil {
		return nil, result.Error
	}

	return standards, nil
}
