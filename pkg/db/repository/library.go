package repository

import (
	"errors"
	"time"

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

func (r *LibraryRepository) FindStandards(params map[string]any, limit, offset int) ([]models.ProjectDTO, int64, error) {
	var standards []models.ProjectDTO
	var total int64

	query := r.db.Model(&models.Project{}).Where("published = ?", true)

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

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and preload
	result := query.Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&standards)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return standards, total, nil
}

func (r *LibraryRepository) GetProjectByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	result := r.db.Preload("Standard").Preload("TechnicalCommittee").Preload("WorkingGroup").
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

func (r *LibraryRepository) GetProjectByReference(reference string) (*models.Project, error) {
	var project models.Project
	result := r.db.Preload("Standard").Preload("TechnicalCommittee").Preload("WorkingGroup").
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

func (r *LibraryRepository) SearchProjects(query string, limit, offset int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64
	searchQuery := "%" + query + "%"

	countQuery := r.db.Model(&models.Project{}).
		Where("published = ? AND (title ILIKE ? OR description ILIKE ? OR reference ILIKE ?)", true, searchQuery, searchQuery, searchQuery)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.Preload("Standard").Preload("TechnicalCommittee").Preload("WorkingGroup").
		Preload("Stage").Preload("WorkingDraft").Preload("CommitteeDraft").
		Where("published = ? AND (title ILIKE ? OR description ILIKE ? OR reference ILIKE ?)", true, searchQuery, searchQuery, searchQuery).
		Limit(limit).Offset(offset).Order("created_at DESC").
		Find(&projects)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return projects, total, nil
}

func (r *LibraryRepository) GetProjectsCreatedBetween(startDate, endDate time.Time) ([]models.Project, error) {
	var projects []models.Project
	result := r.db.Preload("Standard").Preload("TechnicalCommittee").Preload("WorkingGroup").
		Preload("Stage").Preload("WorkingDraft").Preload("CommitteeDraft").
		Where("published = ? AND created_at BETWEEN ? AND ?", true, startDate, endDate).
		Order("created_at DESC").
		Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}

func (r *LibraryRepository) CountProjects() (int64, error) {
	var count int64
	err := r.db.Model(&models.Project{}).Where("published = ?", true).Count(&count).Error
	return count, err
}

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

func (r *LibraryRepository) SearchCommittees(query string, limit, offset int) ([]models.Committee, int64, error) {
	var committees []models.Committee
	var total int64
	searchQuery := "%" + query + "%"

	countQuery := r.db.Model(&models.Committee{}).
		Where("name ILIKE ? OR code ILIKE ?", searchQuery, searchQuery)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.Preload("Chairperson").Preload("Secretary").
		Where("name ILIKE ? OR code ILIKE ?", searchQuery, searchQuery).
		Limit(limit).Offset(offset).Order("created_at DESC").
		Find(&committees)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return committees, total, nil
}

func (r *LibraryRepository) CountCommittees() (int64, error) {
	var count int64
	err := r.db.Model(&models.Committee{}).Count(&count).Error
	return count, err
}

func (r *LibraryRepository) GetProjectsByCommitteeID(committeeID string) ([]models.Project, error) {
	var projects []models.Project
	result := r.db.Preload("Standard").Preload("TechnicalCommittee").Preload("WorkingGroup").
		Preload("Stage").Preload("WorkingDraft").Preload("CommitteeDraft").
		Where("published = ? AND technical_committee_id = ?", true, committeeID).
		Order("created_at DESC").
		Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}
