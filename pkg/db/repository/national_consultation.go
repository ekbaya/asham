package repository

import (
	"errors"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConsultationRepository handles database operations for NationalConsultation
type ConsultationRepository struct {
	db *gorm.DB
}

// NewConsultationRepository initializes a new ConsultationRepository
func NewConsultationRepository(db *gorm.DB) *ConsultationRepository {
	return &ConsultationRepository{db: db}
}

func (r *ConsultationRepository) Create(consultation *models.NationalConsultation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var dars models.DARS

		// Check if DARS exists for the given project
		if err := tx.Where("project_id = ?", consultation.ProjectID).First(&dars).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Get Project
				var project models.Project
				if err := tx.Where("id = ?", consultation.ProjectID).First(&project).Error; err != nil {
					tx.Rollback()
					return err
				}

				now := time.Now()

				// Create a new DARS if it does not exist
				dars = models.DARS{
					ID:                    uuid.New(),
					ProjectID:             consultation.ProjectID,
					CreatedAt:             time.Now(),
					PublicReviewStartDate: now,
					PublicReviewEndDate:   now.AddDate(0, 2, 7),
				}

				if err := tx.Create(&dars).Error; err != nil {
					tx.Rollback()
					return err
				}
			} else {
				tx.Rollback()
				return err
			}
		}

		consultation.DARSID = dars.ID
		consultation.ID = uuid.New()
		consultation.CreatedAt = time.Now()

		if err := tx.Create(&consultation).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}

// GetByID retrieves a NationalConsultation by its ID
func (r *ConsultationRepository) GetByID(id uuid.UUID) (*models.NationalConsultation, error) {
	var Consultation models.NationalConsultation
	err := r.db.Preload("Project").Preload("Member").First(&Consultation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &Consultation, nil
}

// GetAll retrieves all NationalConsultations
func (r *ConsultationRepository) GetAll() ([]models.NationalConsultation, error) {
	var Consultations []models.NationalConsultation
	err := r.db.Preload("Project").Preload("Member").Find(&Consultations).Error
	if err != nil {
		return nil, err
	}
	return Consultations, nil
}

// Update modifies an existing NationalConsultation
func (r *ConsultationRepository) Update(Consultation *models.NationalConsultation) error {
	return r.db.Save(Consultation).Error
}

// Delete removes a NationalConsultation by ID
func (r *ConsultationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.NationalConsultation{}, "id = ?", id).Error
}

// GetByProjectID retrieves all Consultations associated with a specific ProjectID
func (r *ConsultationRepository) GetByProjectID(projectID uuid.UUID) ([]models.NationalConsultation, error) {
	var Consultations []models.NationalConsultation
	err := r.db.
		Preload("NationalSecretary").
		Preload("NationalSecretary.NationalStandardBody").
		Preload("NationalSecretary.NationalStandardBody.MemberState").
		Where("project_id = ?", projectID).
		Find(&Consultations).Error

	if err != nil {
		return nil, err
	}
	return Consultations, nil
}

func (r *ConsultationRepository) GetByProjectIDAndMemberState(projectID, memberState string) ([]models.NationalConsultation, error) {
	var Consultations []models.NationalConsultation
	err := r.db.
		Joins("JOIN members ON members.id=consultation_observations.national_secretary_id").
		Joins("JOIN national_standard_bodys ON national_standard_bodys.id=members.national_standard_body_id").
		Joins("JOIN member_states ON member_states.id=national_standard_bodys.member_state_id").
		Where("Consultation_observations.project_id = ? AND member_states.id = ?", projectID, memberState).
		Preload("Member").
		Preload("Member.NationalStandardBody").
		Preload("Member.NationalStandardBody.MemberState").
		Find(&Consultations).Error

	if err != nil {
		return nil, err
	}
	return Consultations, nil
}
