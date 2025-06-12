package services

import (
	"log"
	"time"

	"github.com/ekbaya/asham/pkg/utilities"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LibraryService struct {
	repo          *repository.LibraryRepository
	memberService *MemberService
}

func NewLibraryService(repo *repository.LibraryRepository, memberService *MemberService) *LibraryService {
	return &LibraryService{
		repo:          repo,
		memberService: memberService,
	}
}

func (s *LibraryService) RegisterMember(user *models.Member) error {
	return s.memberService.CreateMember(user)
}

func (s *LibraryService) Login(email, password string) (string, string, error) {
	return s.memberService.Login(email, password, models.External)
}

func (s *LibraryService) GetTopStandards(limit, offset int) ([]models.ProjectDTO, int64, error) {
	return s.repo.GetTopStandards(limit, offset)
}

func (s *LibraryService) GetLatestStandards(limit, offset int) ([]models.ProjectDTO, int64, error) {
	return s.repo.GetLatestStandards(limit, offset)
}

func (s *LibraryService) GetTopCommittees(limit, offset int) ([]models.CommitteeDTO, int64, error) {
	return s.repo.GetTopCommittees(limit, offset)
}

func (s *LibraryService) FindStandards(params map[string]any, limit, offset int) ([]map[string]any, int64, error) {
	var standards []map[string]any

	projects, total, err := s.repo.FindStandards(params, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	if len(projects) > 0 {
		for _, project := range projects {
			pageCount := 20
			if project.Standard != nil && project.Standard.FileURL != "" {
				calculatedPages, err := utilities.GetPDFPageCount(project.Standard.FileURL)
				if err == nil {
					pageCount = calculatedPages
				} else {
					log.Printf("Error calculating PDF pages for standard ID %v: %v", project.ID, err)
				}
			}

			code := ""
			if project.TechnicalCommittee != nil {
				code = project.TechnicalCommittee.Code
			}

			standard := map[string]any{
				"id":             project.ID,
				"title":          project.Title,
				"reference":      project.Reference,
				"description":    project.Description,
				"sector":         "Pending",
				"committee":      code,
				"language":       project.Language,
				"published":      project.Published,
				"published_date": project.PublishedDate,
				"pages":          pageCount,
				"created_at":     project.CreatedAt,
				"updated_at":     project.UpdatedAt,
			}
			standards = append(standards, standard)
		}
		return standards, total, nil
	} else {
		return nil, total, err
	}
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

func (s *LibraryService) GetCommitteeByID(id uuid.UUID) (*models.TechnicalCommitteeDetailDTO, error) {
	return s.repo.GetCommitteeByID(id)
}

func (s *LibraryService) GetCommitteeByCode(code string) (*models.TechnicalCommitteeDTO, error) {
	return s.repo.GetCommitteeByCode(code)
}

func (s *LibraryService) ListCommittees(limit, offset int, query string) ([]models.TechnicalCommitteeDTO, int64, error) {
	return s.repo.ListCommittees(limit, offset, query)
}

func (s *LibraryService) SearchCommittees(query string, limit, offset int) ([]models.TechnicalCommitteeDTO, int64, error) {
	return s.repo.SearchCommittees(query, limit, offset)
}

func (s *LibraryService) CountCommittees() (int64, error) {
	return s.repo.CountCommittees()
}

func (s *LibraryService) GetProjectsByCommittee(committeeID string) ([]models.Project, error) {
	return s.repo.GetProjectsByCommitteeID(committeeID)
}

func (s *LibraryService) GetSectors() ([]models.Sector, error) {
	return s.repo.GetSectors()
}

func (s *LibraryService) GetBaseQuery() *gorm.DB {
	return s.repo.GetBaseQuery()
}
