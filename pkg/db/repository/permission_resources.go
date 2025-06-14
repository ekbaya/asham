package repository

import (
	"errors"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionResourceRepository struct {
	db *gorm.DB
}

func NewPermissionResourceRepository(db *gorm.DB) *PermissionResourceRepository {
	return &PermissionResourceRepository{db: db}
}

// GetPermissionSlug returns the permission slug for the given HTTP method and path.
func (r *PermissionResourceRepository) GetPermissionSlug(method, path string) (string, error) {
	var resource models.ResourcePermission

	err := r.db.Preload("Permission").
		Where("method = ? AND path = ?", method, path).
		First(&resource).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil // No mapping found
	}

	if err != nil {
		return "", err
	}

	return resource.Permission.Slug, nil
}

// CreateMapping creates a new ResourcePermission
func (r *PermissionResourceRepository) CreateMapping(method, path string, permissionID uuid.UUID) error {
	rp := models.ResourcePermission{
		ID:           uuid.New(), // Generate a new UUID
		Method:       method,
		Path:         path,
		PermissionID: permissionID,
	}
	return r.db.Create(&rp).Error
}

// UpdateMapping updates a ResourcePermission's permission
func (r *PermissionResourceRepository) UpdateMapping(id uuid.UUID, newPermissionID uuid.UUID) error {
	return r.db.Model(&models.ResourcePermission{}).
		Where("id = ?", id).
		Update("permission_id", newPermissionID).Error
}

// ListMappings lists all ResourcePermissions with their permissions
func (r *PermissionResourceRepository) ListMappings() ([]models.ResourcePermission, error) {
	var resources []models.ResourcePermission
	err := r.db.Preload("Permission").Find(&resources).Error
	return resources, err
}

// DeleteMapping deletes a ResourcePermission
func (r *PermissionResourceRepository) DeleteMapping(id uuid.UUID) error {
	return r.db.Delete(&models.ResourcePermission{}, "id = ?", id).Error
}
