package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrProposalNotFound      = errors.New("proposal not found")
	ErrProposalAlreadyExists = errors.New("a proposal for this project already exists")
)

// ProposalRepository handles database operations for Proposal entities
type ProposalRepository struct {
	db *gorm.DB
}

// NewProposalRepository creates a new proposal repository instance
func NewProposalRepository(db *gorm.DB) *ProposalRepository {
	return &ProposalRepository{db: db}
}

// Create adds a new proposal to the database, ensuring only one proposal exists per project
func (r *ProposalRepository) Create(proposal *models.Proposal) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Defer a rollback in case anything fails
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if a proposal already exists for this project
	var count int64
	if err := tx.Model(&models.Proposal{}).Where("project_id = ?", proposal.ProjectID).Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}

	if count > 0 {
		tx.Rollback()
		return ErrProposalAlreadyExists
	}

	// Prepare proposal data
	if proposal.ID == uuid.Nil {
		proposal.ID = uuid.New()
	}
	if proposal.CreatedAt.IsZero() {
		proposal.CreatedAt = time.Now()
	}

	// Find the initial stage (number = 1)
	var stage models.Stage
	if err := tx.Where("number = ?", 1).First(&stage).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create the proposal first
	if err := tx.Create(proposal).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update the project stage
	if err := r.updateProjectStageWithTx(tx, proposal.ProjectID, stage.ID.String(), "Proposal submitted"); err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

// Helper method to update project stage within a transaction
func (r *ProposalRepository) updateProjectStageWithTx(tx *gorm.DB, projectID string, newStageID string, notes string) error {
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

	// Change Ref to NWIP
	project.Reference = strings.Replace(project.Reference, "PWI", "NWIP", -1)

	// Update current stage of the project
	if err := tx.Model(&models.Project{}).
		Where("id = ?", projectID).
		Updates(map[string]interface{}{
			"stage_id":   newStageID,
			"updated_at": now,
		}).Error; err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a proposal by its ID with all associated entities
func (r *ProposalRepository) GetByID(id uuid.UUID) (*models.Proposal, error) {
	var proposal models.Proposal
	err := r.db.Preload("Project").
		Preload("CreatedBy").
		Preload("ProposingNSB").
		Preload("ReferencedStandards").
		First(&proposal, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProposalNotFound
		}
		return nil, err
	}
	return &proposal, nil
}

// GetByProjectID retrieves the proposal for a specific project
// Since there's only one proposal per project, this returns a single proposal or an error
func (r *ProposalRepository) GetByProjectID(projectID uuid.UUID) (*models.Proposal, error) {
	var proposal models.Proposal
	err := r.db.
		Preload(clause.Associations).
		Preload("CreatedBy.NationalStandardBody").
		Where("project_id = ?", projectID).
		First(&proposal).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProposalNotFound
		}
		return nil, err
	}

	return &proposal, nil
}

// GetByProposingNSB retrieves proposals from a specific National Standards Body
func (r *ProposalRepository) GetByProposingNSB(nsbID uuid.UUID, limit, offset int) ([]models.Proposal, int64, error) {
	var proposals []models.Proposal
	var total int64

	// Get total count
	if err := r.db.Model(&models.Proposal{}).Where("proposing_nsb_id = ?", nsbID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get proposals with pagination
	err := r.db.Preload("Project").
		Preload("CreatedBy").
		Preload("ProposingNSB").
		Where("proposing_nsb_id = ?", nsbID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&proposals).Error

	return proposals, total, err
}

// GetByCreator retrieves proposals created by a specific member
func (r *ProposalRepository) GetByCreator(memberID uuid.UUID, limit, offset int) ([]models.Proposal, int64, error) {
	var proposals []models.Proposal
	var total int64

	// Get total count
	if err := r.db.Model(&models.Proposal{}).Where("created_by_id = ?", memberID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get proposals with pagination
	err := r.db.Preload("Project").
		Preload("CreatedBy").
		Preload("ProposingNSB").
		Where("created_by_id = ?", memberID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&proposals).Error

	return proposals, total, err
}

// Update updates an existing proposal
func (r *ProposalRepository) Update(proposal *models.Proposal) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(proposal).Error
}

// UpdatePartial updates specific fields of a proposal
func (r *ProposalRepository) UpdatePartial(id uuid.UUID, updates map[string]interface{}) error {
	// Block updating project_id to prevent violating the one-proposal-per-project constraint
	if _, exists := updates["project_id"]; exists {
		return errors.New("changing project_id is not allowed")
	}

	return r.db.Model(&models.Proposal{}).Where("id = ?", id).Updates(updates).Error
}

// Delete removes a proposal from the database
func (r *ProposalRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Proposal{}, "id = ?", id).Error
}

// List retrieves all proposals with pagination
func (r *ProposalRepository) List(limit, offset int) ([]models.Proposal, int64, error) {
	var proposals []models.Proposal
	var total int64

	// Get total count
	if err := r.db.Model(&models.Proposal{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get proposals with pagination
	err := r.db.Preload("Project").
		Preload("CreatedBy").
		Preload("ProposingNSB").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&proposals).Error

	if err != nil {
		return nil, 0, err
	}

	return proposals, total, nil
}

// Search searches for proposals by title or scope
func (r *ProposalRepository) Search(query string, limit, offset int) ([]models.Proposal, int64, error) {
	var proposals []models.Proposal
	var total int64

	searchQuery := "%" + query + "%"

	countQuery := r.db.Model(&models.Proposal{}).
		Where("full_title LIKE ? OR scope LIKE ? OR justification LIKE ?",
			searchQuery, searchQuery, searchQuery)

	// Get total count
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get proposals with pagination
	err := r.db.Preload("Project").
		Preload("CreatedBy").
		Preload("ProposingNSB").
		Where("full_title LIKE ? OR scope LIKE ? OR justification LIKE ?",
			searchQuery, searchQuery, searchQuery).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&proposals).Error

	if err != nil {
		return nil, 0, err
	}

	return proposals, total, nil
}

// Exists checks if a proposal already exists for a given project
func (r *ProposalRepository) Exists(projectID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Proposal{}).Where("project_id = ?", projectID).Count(&count).Error
	return count > 0, err
}

// AddReferencedStandard adds a document to the proposal's referenced standards
func (r *ProposalRepository) AddReferencedStandard(proposalID, documentID uuid.UUID) error {
	return r.db.Exec(
		"INSERT INTO referenced_standards (proposal_id, document_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
		proposalID, documentID,
	).Error
}

// RemoveReferencedStandard removes a document from the proposal's referenced standards
func (r *ProposalRepository) RemoveReferencedStandard(proposalID, documentID uuid.UUID) error {
	return r.db.Exec(
		"DELETE FROM referenced_standards WHERE proposal_id = ? AND document_id = ?",
		proposalID, documentID,
	).Error
}

// GetProposalCountByNSB gets the count of proposals submitted by each NSB
func (r *ProposalRepository) GetProposalCountByNSB() (map[uuid.UUID]int64, error) {
	type Result struct {
		NSBID uuid.UUID `gorm:"column:proposing_nsb_id"`
		Count int64     `gorm:"column:count"`
	}

	var results []Result
	err := r.db.Model(&models.Proposal{}).
		Select("proposing_nsb_id, COUNT(*) as count").
		Group("proposing_nsb_id").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	countMap := make(map[uuid.UUID]int64)
	for _, result := range results {
		countMap[result.NSBID] = result.Count
	}

	return countMap, nil
}

// Transfer transfers a proposal from one project to another
// This is a special operation that maintains the one-proposal-per-project constraint
func (r *ProposalRepository) Transfer(proposalID uuid.UUID, newProjectID string) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Check if the proposal exists
	var proposal models.Proposal
	if err := tx.First(&proposal, "id = ?", proposalID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProposalNotFound
		}
		return err
	}

	// Check if the target project already has a proposal
	var count int64
	if err := tx.Model(&models.Proposal{}).Where("project_id = ?", newProjectID).Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}

	if count > 0 {
		tx.Rollback()
		return ErrProposalAlreadyExists
	}

	// Update the proposal's project ID
	if err := tx.Model(&models.Proposal{}).Where("id = ?", proposalID).Update("project_id", newProjectID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
