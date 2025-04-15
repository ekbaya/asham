package repository

import (
	"github.com/ekbaya/asham/pkg/domain/models"
	"gorm.io/gorm"
)

type LibraryRepository struct {
	db *gorm.DB
}

func NewLibraryRepository(db *gorm.DB) *LibraryRepository {
	return &LibraryRepository{db: db}
}

func (r *LibraryRepository) FindStandards(sector string, limit, offset int) ([]models.Project, error) {
	var standards []models.Project

	query := r.db.Where("published = ?", true)

	// If sector is provided and not empty, filter by sector
	if sector != "" {
		query = query.Where("sector = ?", models.ProjectSector(sector))
	}

	// Apply pagination
	result := query.Limit(limit).Offset(offset).Preload("Standard").Preload("TechnicalCommittee").Find(&standards)

	if result.Error != nil {
		return nil, result.Error
	}

	return standards, nil
}
