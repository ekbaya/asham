package repository

import (
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CommentRepository handles database operations for CommentObservation
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository initializes a new CommentRepository
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create adds a new CommentObservation to the database
func (r *CommentRepository) Create(comment *models.CommentObservation) error {
	comment.ID = uuid.New()
	comment.CreatedAt = time.Now()
	return r.db.Create(&comment).Error
}

// GetByID retrieves a CommentObservation by its ID
func (r *CommentRepository) GetByID(id uuid.UUID) (*models.CommentObservation, error) {
	var comment models.CommentObservation
	err := r.db.Preload("Project").Preload("Member").First(&comment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetAll retrieves all CommentObservations
func (r *CommentRepository) GetAll() ([]models.CommentObservation, error) {
	var comments []models.CommentObservation
	err := r.db.Preload("Project").Preload("Member").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// Update modifies an existing CommentObservation
func (r *CommentRepository) Update(comment *models.CommentObservation) error {
	return r.db.Save(comment).Error
}

// Delete removes a CommentObservation by ID
func (r *CommentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.CommentObservation{}, "id = ?", id).Error
}

// GetByProjectID retrieves all comments associated with a specific ProjectID
func (r *CommentRepository) GetByProjectID(projectID uuid.UUID) ([]models.CommentObservation, error) {
	var comments []models.CommentObservation
	err := r.db.
		Preload("Member").
		Preload("Member.NationalStandardBody").
		Preload("Member.NationalStandardBody.MemberState").
		Where("project_id = ?", projectID).
		Find(&comments).Error

	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) GetByProjectIDAndMemberState(projectID, memberState string) ([]models.CommentObservation, error) {
	var comments []models.CommentObservation
	err := r.db.
		Joins("JOIN members ON members.id=comment_observations.national_secretary_id").
		Joins("JOIN national_standard_bodys ON national_standard_bodys.id=members.national_standard_body_id").
		Joins("JOIN member_states ON member_states.id=national_standard_bodys.member_state_id").
		Where("comment_observations.project_id = ? AND member_states.id = ?", projectID, memberState).
		Preload("Member").
		Preload("Member.NationalStandardBody").
		Preload("Member.NationalStandardBody.MemberState").
		Find(&comments).Error

	if err != nil {
		return nil, err
	}
	return comments, nil
}
