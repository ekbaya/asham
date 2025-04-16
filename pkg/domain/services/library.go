package services

import (
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type LibraryService struct {
	repo *repository.LibraryRepository
}

func NewLibraryService(repo *repository.LibraryRepository) *LibraryService {
	return &LibraryService{repo: repo}
}

func (s *LibraryService) FindStandards(params map[string]any, limit, offset int) ([]models.ProjectDTO, int64, error) {
	return s.repo.FindStandards(params, limit, offset)
}

func (s *LibraryService) GetProjectByID(id uuid.UUID) (*models.Project, error) {
	return s.repo.GetProjectByID(id)
}

func (s *LibraryService) GetProjectByReference(reference string) (*models.Project, error) {
	return s.repo.GetProjectByReference(reference)
}

func (s *LibraryService) SearchProjects(query string, limit, offset int) ([]models.Project, int64, error) {
	return s.repo.SearchProjects(query, limit, offset)
}

func (s *LibraryService) GetProjectsCreatedBetween(startDate, endDate time.Time) ([]models.Project, error) {
	return s.repo.GetProjectsCreatedBetween(startDate, endDate)
}

func (s *LibraryService) CountProjects() (int64, error) {
	return s.repo.CountProjects()
}

func (s *LibraryService) GetCommitteeByID(id uuid.UUID) (*models.Committee, error) {
	return s.repo.GetCommitteeByID(id)
}

func (s *LibraryService) GetCommitteeByCode(code string) (*models.Committee, error) {
	return s.repo.GetCommitteeByCode(code)
}

func (s *LibraryService) ListCommittees(limit, offset int) ([]models.Committee, int64, error) {
	return s.repo.ListCommittees(limit, offset)
}

func (s *LibraryService) SearchCommittees(query string, limit, offset int) ([]models.Committee, int64, error) {
	return s.repo.SearchCommittees(query, limit, offset)
}

func (s *LibraryService) CountCommittees() (int64, error) {
	return s.repo.CountCommittees()
}

func (s *LibraryService) GetProjectsByCommittee(committeeID string) ([]models.Project, error) {
	return s.repo.GetProjectsByCommitteeID(committeeID)
}
