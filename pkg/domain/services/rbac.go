package services

import (
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type RbacService struct {
	repo *repository.RbacRepository
}

func NewRbacService(repo *repository.RbacRepository) *RbacService {
	return &RbacService{repo: repo}
}

func (s *RbacService) CreateRole(role *models.Role) error {
	role.ID = uuid.New()
	role.CreatedAt = time.Now()
	return s.repo.CreateRole(role)
}

func (s *RbacService) GetRoleByID(id string) (*models.Role, error) {
	return s.repo.GetRoleByID(id)
}

func (s *RbacService) ListRoles() ([]models.Role, error) {
	return s.repo.ListRoles()
}

func (s *RbacService) UpdateRole(role *models.Role) error {
	role.UpdatedAt = time.Now()
	return s.repo.UpdateRole(role)
}

func (s *RbacService) DeleteRole(id string) error {
	return s.repo.DeleteRole(id)
}

func (s *RbacService) DeletePermission(id string) error {
	return s.repo.DeletePermission(id)
}

func (s *RbacService) CreatePermission(permission *models.Permission) error {
	permission.ID = uuid.New()
	permission.CreatedAt = time.Now()
	return s.repo.CreatePermission(permission)
}

func (s *RbacService) ListPermissions() ([]models.Permission, error) {
	return s.repo.ListPermissions()
}

func (s *RbacService) AddPermissionToRole(roleID, permissionID string) error {
	return s.repo.AddPermissionToRole(roleID, permissionID)
}

func (s *RbacService) RemovePermissionFromRole(roleID, permissionID string) error {
	return s.repo.RemovePermissionFromRole(roleID, permissionID)
}

func (s *RbacService) AssignRoleToMember(memberID, roleID string) error {
	return s.repo.AssignRoleToMember(memberID, roleID)
}

func (s *RbacService) RemoveRoleFromMember(memberID, roleID string) error {
	return s.repo.RemoveRoleFromMember(memberID, roleID)
}

func (s *RbacService) ListRolesForMember(memberID string) ([]models.Role, error) {
	return s.repo.ListRolesForMember(memberID)
}

func (s *RbacService) CreateIfNotExists(permission models.Permission) (*models.Permission, error) {
	return s.repo.CreateIfNotExists(&permission)
}
