package services

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LibraryService struct {
	repo          *repository.LibraryRepository
	memberService *MemberService
}

func NewLibraryService(repo *repository.LibraryRepository) *LibraryService {
	return &LibraryService{
		repo: repo,
	}
}

func (s *LibraryService) RegisterMember(user *models.User) error {
	// Validate password
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Hash password
	hashedPassword, err := utilities.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Create user
	user.HashedPassword = hashedPassword
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	err = s.repo.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *LibraryService) Login(email, password string) (string, string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		fmt.Print("User Not Found: ", err)
		return "", "", errors.New("invalid credentials")
	}

	// Verify password
	if !utilities.CheckPasswordHash(password, user.HashedPassword) {
		fmt.Print("Wrong Username Or Password")
		return "", "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := models.GenerateJWT(user.ID.String())
	if err != nil {
		return "", "", errors.New("failed to generate token")
	}

	// Generate JWT refresh token
	refreshToken, err := models.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	return token, refreshToken, nil
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
				calculatedPages, err := s.getPDFPageCount(project.Standard.FileURL)
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
				"sector":         project.Sector,
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

// getPDFPageCount downloads the PDF from the URL and calculates the number of pages
func (service *LibraryService) getPDFPageCount(pdfURL string) (int, error) {
	// Parse the URL
	parsedURL, err := url.Parse(pdfURL)
	if err != nil {
		return 0, err
	}

	// Handle local file path (for assets directory)
	if parsedURL.Scheme == "" || parsedURL.Scheme == "file" {
		// Assuming this is a local file path
		return api.PageCountFile(pdfURL)
	}

	// For remote URLs, download the file temporarily
	resp, err := http.Get(pdfURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Create a temporary file to store the PDF
	tmpFile, err := os.CreateTemp("", "standard-*.pdf")
	if err != nil {
		return 0, err
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	// Copy the PDF data to the temporary file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		tmpFile.Close()
		return 0, err
	}
	tmpFile.Close()

	// Get the page count from the downloaded file
	pageCount, err := api.PageCountFile(tmpFile.Name())
	if err != nil {
		return 0, err
	}

	return pageCount, nil
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
