package repository

import (
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

// Create adds a new NationalConsultation to the database
func (r *ConsultationRepository) Create(Consultation *models.NationalConsultation) error {
	Consultation.ID = uuid.New()
	Consultation.CreatedAt = time.Now()
	return r.db.Create(&Consultation).Error
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
		Joins("JOIN members ON members.id=Consultation_observations.national_secretary_id").
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
