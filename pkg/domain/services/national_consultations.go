package services

import (
	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type NationalConsultationService struct {
	repo *repository.ConsultationRepository
}

func NewNationalConsultationService(repo *repository.ConsultationRepository) *NationalConsultationService {
	return &NationalConsultationService{repo: repo}
}

func (service *NationalConsultationService) Create(NationalConsultation *models.NationalConsultation) error {
	return service.repo.Create(NationalConsultation)
}

func (service *NationalConsultationService) GetByID(id uuid.UUID) (*models.NationalConsultation, error) {
	return service.repo.GetByID(id)
}

func (service *NationalConsultationService) GetAll() ([]models.NationalConsultation, error) {
	return service.repo.GetAll()
}

func (service *NationalConsultationService) Update(NationalConsultation *models.NationalConsultation) error {
	return service.repo.Update(NationalConsultation)
}

func (service *NationalConsultationService) Delete(id uuid.UUID) error {
	return service.repo.Delete(id)
}

func (service *NationalConsultationService) GetByProjectID(projectID uuid.UUID) ([]models.NationalConsultation, error) {
	return service.repo.GetByProjectID(projectID)
}

func (service *NationalConsultationService) GetByProjectIDAndMemberState(projectID, memberState string) ([]models.NationalConsultation, error) {
	return service.repo.GetByProjectIDAndMemberState(projectID, memberState)
}
