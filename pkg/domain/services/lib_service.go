package services

import (
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

// LibraryService handles operations related to searching committees and standards
type LibraryService struct {
	repo *repository.LibraryRepository
}

// NewLibraryService creates a new LibraryService instance
func NewLibraryService(repo *repository.LibraryRepository) *LibraryService {
	return &LibraryService{repo: repo}
}

// GetCommitteeByID retrieves a committee by its ID
func (s *LibraryService) GetCommitteeByID(id uuid.UUID) (*models.Committee, error) {
	return s.repo.GetCommitteeByID(id)
}

// GetCommitteeByCode retrieves a committee by its code
func (s *LibraryService) GetCommitteeByCode(code string) (*models.Committee, error) {
	return s.repo.GetCommitteeByCode(code)
}

// SearchCommittees searches committees by name or code
func (s *LibraryService) SearchCommittees(query string, limit, offset int) ([]models.Committee, int64, error) {
	return s.repo.SearchCommittees(query, limit, offset)
}

// ListCommittees retrieves all committees with pagination
func (s *LibraryService) ListCommittees(limit, offset int) ([]models.Committee, int64, error) {
	return s.repo.ListCommittees(limit, offset)
}

// GetProjectByID retrieves a project (standard) by its ID
func (s *LibraryService) GetProjectByID(id uuid.UUID) (*models.Project, error) {
	return s.repo.GetProjectByID(id)
}

// GetProjectByReference retrieves a project by its reference
func (s *LibraryService) GetProjectByReference(reference string) (*models.Project, error) {
	return s.repo.GetProjectByReference(reference)
}

// SearchProjects searches projects by title, description, or reference
func (s *LibraryService) SearchProjects(query string, limit, offset int) ([]models.Project, int64, error) {
	return s.repo.SearchProjects(query, limit, offset)
}

// ListProjects retrieves all projects with pagination
func (s *LibraryService) ListProjects(limit, offset int) ([]models.Project, int64, error) {
	return s.repo.ListProjects(limit, offset)
}

// GetProjectsByCommittee retrieves all projects associated with a committee
func (s *LibraryService) GetProjectsByCommittee(committeeID string) ([]models.Project, error) {
	return s.repo.GetProjectsByCommitteeID(committeeID)
}

// GetProjectsCreatedBetween retrieves projects created within a time range
func (s *LibraryService) GetProjectsCreatedBetween(startDate, endDate time.Time) ([]models.Project, error) {
	return s.repo.GetProjectsCreatedBetween(startDate, endDate)
}

// CountCommittees returns the total number of committees
func (s *LibraryService) CountCommittees() (int64, error) {
	return s.repo.CountCommittees()
}

// CountProjects returns the total number of projects
func (s *LibraryService) CountProjects() (int64, error) {
	return s.repo.CountProjects()
}
