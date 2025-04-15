package repository

import (
	"errors"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LibraryRepository handles database operations for library-related entities (Committees and Projects)
type LibraryRepository struct {
	db *gorm.DB
}

// NewLibraryRepository creates a new LibraryRepository instance
func NewLibraryRepository(db *gorm.DB) *LibraryRepository {
	return &LibraryRepository{db: db}
}

// GetCommitteeByID retrieves a committee by its ID
func (r *LibraryRepository) GetCommitteeByID(id uuid.UUID) (*models.Committee, error) {
	var committee models.Committee
	result := r.db.Preload("Chairperson").Preload("Secretary").First(&committee, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("committee not found")
		}
		return nil, result.Error
	}
	return &committee, nil
}

// GetCommitteeByCode retrieves a committee by its code
func (r *LibraryRepository) GetCommitteeByCode(code string) (*models.Committee, error) {
	var committee models.Committee
	result := r.db.Preload("Chairperson").Preload("Secretary").First(&committee, "code = ?", code)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("committee not found")
		}
		return nil, result.Error
	}
	return &committee, nil
}

// ListCommittees retrieves all committees with pagination
func (r *LibraryRepository) ListCommittees(limit, offset int) ([]models.Committee, int64, error) {
	var committees []models.Committee
	var total int64

	if err := r.db.Model(&models.Committee{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.Preload("Chairperson").Preload("Secretary").
		Limit(limit).Offset(offset).Order("created_at DESC").
		Find(&committees)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return committees, total, nil
}

// SearchCommittees searches committees by name or code
func (r *LibraryRepository) SearchCommittees(query string, limit, offset int) ([]models.Committee, int64, error) {
	var committees []models.Committee
	var total int64
	searchQuery := "%" + query + "%"

	countQuery := r.db.Model(&models.Committee{}).
		Where("name LIKE ? OR code LIKE ?", searchQuery, searchQuery)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.Preload("Chairperson").Preload("Secretary").
		Where("name LIKE ? OR code LIKE ?", searchQuery, searchQuery).
		Limit(limit).Offset(offset).Order("created_at DESC").
		Find(&committees)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return committees, total, nil
}

// CountCommittees returns the total number of committees
func (r *LibraryRepository) CountCommittees() (int64, error) {
	var count int64
	err := r.db.Model(&models.Committee{}).Count(&count).Error
	return count, err
}

// GetProjectByID retrieves a project by its ID
func (r *LibraryRepository) GetProjectByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	result := r.db.Preload("Member").Preload("TechnicalCommittee").Preload("WorkingGroup").
		Preload("Stage").Preload("WorkingDraft").Preload("CommitteeDraft").
		First(&project, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, result.Error
	}
	return &project, nil
}

// GetProjectByReference retrieves a project by its reference
func (r *LibraryRepository) GetProjectByReference(reference string) (*models.Project, error) {
	var project models.Project
	result := r.db.Preload("Member").Preload("TechnicalCommittee").Preload("WorkingGroup").
		Preload("Stage").Preload("WorkingDraft").Preload("CommitteeDraft").
		First(&project, "reference = ?", reference)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, result.Error
	}
	return &project, nil
}

// ListProjects retrieves all projects with pagination
func (r *LibraryRepository) ListProjects(limit, offset int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	if err := r.db.Model(&models.Project{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.Preload("Member").Preload("TechnicalCommittee").Preload("WorkingGroup").
		Preload("Stage").Preload("WorkingDraft").Preload("CommitteeDraft").
		Limit(limit).Offset(offset).Order("created_at DESC").
		Find(&projects)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return projects, total, nil
}

// SearchProjects searches projects by title, description, or reference
func (r *LibraryRepository) SearchProjects(query string, limit, offset int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64
	searchQuery := "%" + query + "%"

	countQuery := r.db.Model(&models.Project{}).
		Where("title LIKE ? OR description LIKE ? OR reference LIKE ?", searchQuery, searchQuery, searchQuery)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.Preload("Member").Preload("TechnicalCommittee").Preload("WorkingGroup").
		Preload("Stage").Preload("WorkingDraft").Preload("CommitteeDraft").
		Where("title LIKE ? OR description LIKE ? OR reference LIKE ?", searchQuery, searchQuery, searchQuery).
		Limit(limit).Offset(offset).Order("created_at DESC").
		Find(&projects)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return projects, total, nil
}

// GetProjectsByCommitteeID retrieves all projects for a given committee
func (r *LibraryRepository) GetProjectsByCommitteeID(committeeID string) ([]models.Project, error) {
	var projects []models.Project
	result := r.db.Preload("Member").Preload("TechnicalCommittee").Preload("WorkingGroup").
		Preload("Stage").Preload("WorkingDraft").Preload("CommitteeDraft").
		Where("technical_committee_id = ?", committeeID).
		Order("created_at DESC").
		Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}

// GetProjectsCreatedBetween retrieves projects created within a time range
func (r *LibraryRepository) GetProjectsCreatedBetween(startDate, endDate time.Time) ([]models.Project, error) {
	var projects []models.Project
	result := r.db.Preload("Member").Preload("TechnicalCommittee").Preload("WorkingGroup").
		Preload("Stage").Preload("WorkingDraft").Preload("CommitteeDraft").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("created_at DESC").
		Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}

// CountProjects returns the total number of projects
func (r *LibraryRepository) CountProjects() (int64, error) {
	var count int64
	err := r.db.Model(&models.Project{}).Count(&count).Error
	return count, err
}
