package services

import (
	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type PermissionResourceService struct {
	repo *repository.PermissionResourceRepository
}

func NewPermissionResourceService(repo *repository.PermissionResourceRepository) *PermissionResourceService {
	return &PermissionResourceService{repo: repo}
}

func (service *PermissionResourceService) GetPermissionSlug(method, path string) (string, error) {
	return service.repo.GetPermissionSlug(method, path)
}

func (service *PermissionResourceService) CreateMapping(method, path string, permissionID uuid.UUID) error {
	return service.repo.CreateMapping(method, path, permissionID)
}

func (service *PermissionResourceService) UpdateMapping(id uuid.UUID, newPermissionID uuid.UUID) error {
	return service.repo.UpdateMapping(id, newPermissionID)
}

func (service *PermissionResourceService) ListMappings() ([]models.ResourcePermission, error) {
	return service.repo.ListMappings()
}

func (service *PermissionResourceService) DeleteMapping(id uuid.UUID) error {
	return service.repo.DeleteMapping(id)
}
