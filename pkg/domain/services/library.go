package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/ekbaya/asham/pkg/utilities"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
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

func (s *LibraryService) GetCommitteeByID(id uuid.UUID) (*models.TechnicalCommitteeDTO, error) {
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
func (s *LibraryService) GetSectors() ([]models.ProjectSector, error) {
	return s.repo.GetSectors()
}
