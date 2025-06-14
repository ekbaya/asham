package repository

import (
	"errors"
	"fmt"

	"github.com/ekbaya/asham/pkg/domain/models"
	"gorm.io/gorm"
)

type RbacRepository struct {
	db *gorm.DB
}

func NewRbacRepository(db *gorm.DB) *RbacRepository {
	return &RbacRepository{db: db}
}

// Role CRUD

func (r *RbacRepository) CreateRole(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *RbacRepository) GetRoleByID(id string) (*models.Role, error) {
	var role models.Role
	if err := r.db.Preload("Permissions").First(&role, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RbacRepository) ListRoles() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *RbacRepository) UpdateRole(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *RbacRepository) DeleteRole(id string) error {
	return r.db.Delete(&models.Role{}, "id = ?", id).Error
}

// Permission CRUD

func (r *RbacRepository) CreatePermission(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

func (r *RbacRepository) GetPermissionByID(id string) (*models.Permission, error) {
	var perm models.Permission
	if err := r.db.First(&perm, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *RbacRepository) ListPermissions() ([]models.Permission, error) {
	var perms []models.Permission
	err := r.db.Find(&perms).Error
	return perms, err
}

func (r *RbacRepository) UpdatePermission(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

func (r *RbacRepository) DeletePermission(id string) error {
	return r.db.Delete(&models.Permission{}, "id = ?", id).Error
}

// Role-Permission Management

func (r *RbacRepository) AddPermissionToRole(roleID string, permissionID string) error {
	var role models.Role
	if err := r.db.Preload("Permissions").First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}

	var permission models.Permission
	if err := r.db.First(&permission, "id = ?", permissionID).Error; err != nil {
		return err
	}

	return r.db.Model(&role).Association("Permissions").Append(&permission)
}

func (r *RbacRepository) RemovePermissionFromRole(roleID string, permissionID string) error {
	var role models.Role
	if err := r.db.Preload("Permissions").First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}

	var permission models.Permission
	if err := r.db.First(&permission, "id = ?", permissionID).Error; err != nil {
		return err
	}

	return r.db.Model(&role).Association("Permissions").Delete(&permission)
}

// Member-Role Management

func (r *RbacRepository) AssignRoleToMember(memberID string, roleID string) error {
	var member models.Member
	if err := r.db.Preload("Roles").First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	var role models.Role
	if err := r.db.First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}

	return r.db.Model(&member).Association("Roles").Append(&role)
}

func (r *RbacRepository) RemoveRoleFromMember(memberID string, roleID string) error {
	var member models.Member
	if err := r.db.Preload("Roles").First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	var role models.Role
	if err := r.db.First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}

	return r.db.Model(&member).Association("Roles").Delete(&role)
}

func (r *RbacRepository) ListRolesForMember(memberID string) ([]models.Role, error) {
	var member models.Member
	if err := r.db.Preload("Roles").First(&member, "id = ?", memberID).Error; err != nil {
		return nil, err
	}
	return member.Roles, nil
}

func (r *RbacRepository) CreateIfNotExists(permission *models.Permission) (*models.Permission, error) {
	var existing models.Permission

	// Check if permission already exists by title
	err := r.db.
		Where("title = ?", permission.Title).
		First(&existing).Error

	if err == nil {
		// Already exists, optionally ensure it's assigned to ROLE_ADMIN
		r.ensurePermissionAssignedToAdmin(&existing)
		return &existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Not found, create the permission
	if err := r.db.Create(permission).Error; err != nil {
		return nil, err
	}

	// Assign to ROLE_ADMIN
	if err := r.ensurePermissionAssignedToAdmin(permission); err != nil {
		return nil, err
	}

	return permission, nil
}

func (r *RbacRepository) ensurePermissionAssignedToAdmin(permission *models.Permission) error {
	var adminRole models.Role

	// Fetch admin role
	err := r.db.Preload("Permissions").
		Where("title = ?", "ROLE_ADMIN").
		First(&adminRole).Error

	if err != nil {
		return fmt.Errorf("failed to find ROLE_ADMIN: %w", err)
	}

	// Check if permission already assigned
	for _, p := range adminRole.Permissions {
		if p.ID == permission.ID {
			return nil // already assigned
		}
	}

	// Append permission and save
	if err := r.db.Model(&adminRole).Association("Permissions").Append(permission); err != nil {
		return fmt.Errorf("failed to assign permission to ROLE_ADMIN: %w", err)
	}

	return nil
}
