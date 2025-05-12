package services

import (
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type StandardService struct {
	repo *repository.StandardRepository
}

func NewStandardService(repo *repository.StandardRepository) *StandardService {
	return &StandardService{repo: repo}
}

func (service *StandardService) CreateStandard(standard *models.Standard) error {
	standard.ID = uuid.New()
	standard.CreatedAt = time.Now()
	return service.repo.CreateStandard(standard)
}

func (service *StandardService) SaveStandard(standard *models.Standard) error {
	standard.UpdatedAt = time.Now()
	return service.repo.SaveStandard(standard)
}

func (service *StandardService) GetStandardByID(id string) (*models.Standard, error) {
	return service.repo.GetStandardByID(id)
}

func (service *StandardService) GetStandardVersions(standardID string) ([]models.StandardVersion, error) {
	return service.repo.GetStandardVersions(standardID)
}

func (service *StandardService) RestoreStandardVersion(standardID string, version int) error {
	return service.repo.RestoreStandardVersion(standardID, version)
}
