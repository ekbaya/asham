package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/utilities"
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

func (service *StandardService) SaveStandard(standard *models.Standard, memberId string) error {
	standard.UpdatedAt = time.Now()

	// Get current saved version
	current, err := service.repo.GetStandardByID(standard.ID.String())
	if err != nil {
		return err
	}

	// Diff the content
	diff, err := utilities.DiffJSON(current.Content, standard.Content)
	if err != nil {
		return err
	}

	if diff == "No changes" {
		return nil // Or you may choose to return an explicit "no-op" response
	}

	// Save the new version in versions table
	version := models.StandardVersion{
		ID:         uuid.New(),
		StandardID: standard.ID,
		Version:    current.Version + 1,
		Content:    standard.Content,
		SavedByID:  memberId,
		SavedAt:    time.Now(),
	}
	if err := service.repo.SaveVersion(&version); err != nil {
		return err
	}

	// Optional: Save audit log
	audit := models.StandardAuditLog{
		ID:         uuid.New(),
		StandardID: standard.ID,
		Version:    version.Version,
		ChangedBy:  memberId,
		ChangeDiff: diff,
		CreatedAt:  time.Now(),
	}
	if err := service.repo.SaveAuditLog(&audit); err != nil {
		return err
	}

	// Update the base standard
	return service.repo.SaveStandard(standard, memberId)
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

func (service *StandardService) GetVersion(standardID, versionStr string) (*models.StandardVersion, error) {
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return nil, errors.New("invalid version number")
	}
	return service.repo.GetStandardVersion(standardID, version)
}

func (service *StandardService) GetAuditLogs(standardID string) ([]models.StandardAuditLog, error) {
	return service.repo.GetAuditLogsByStandardID(standardID)
}
