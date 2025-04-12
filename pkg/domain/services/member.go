package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/google/uuid"
)

type MemberService struct {
	repo         *repository.MemberRepository
	emailService *EmailService
}

func NewMemberService(repo *repository.MemberRepository, emailService *EmailService) *MemberService {
	return &MemberService{
		repo:         repo,
		emailService: emailService,
	}
}

func (service *MemberService) CreateMember(member *models.Member) error {
	// Generate a random password
	clearPassword, err := utilities.GenerateRandomPassword(12)
	if err != nil {
		return errors.New("failed to generate password")
	}

	// Hash password
	hashedPassword, err := utilities.HashPassword(clearPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Create user
	member.HashedPassword = hashedPassword
	member.ID = uuid.New()
	member.CreatedAt = time.Now()
	err = service.repo.CreateMember(member)
	if err != nil {
		return err
	}

	go service.emailService.SendWelcomeEmail(member.Email, member.FirstName, clearPassword)

	return nil
}

func (service *MemberService) Login(email, password string) (string, string, error) {
	user, err := service.repo.GetMemberByEmail(email)

	// Handle error if user is not found
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

func (service *MemberService) Account(memberId string) (models.Member, error) {
	// Retrieve member
	member, err := service.repo.GetMemberByID(memberId)
	if err != nil {
		return models.Member{}, errors.New("member does not exist")
	}
	return *member, nil
}

func (service *MemberService) AccountWithResponsibilities(memberId string) (map[string]any, error) {
	// Retrieve member
	member, err := service.repo.GetMemberByID(memberId)
	if err != nil {
		return nil, errors.New("member does not exist")
	}

	// Retrieve responsibilities
	responsibilities, err := service.repo.GetMemberResponsibilities(memberId)
	if err != nil {
		return nil, errors.New("failed to retrieve responsibilities")
	}

	// Construct response map
	result := map[string]any{
		"member":           member,
		"responsibilities": responsibilities,
	}

	return result, nil
}

func (service *MemberService) GetMemberByEmail(email string) (*models.Member, error) {
	return service.repo.GetMemberByEmail(email)
}

func (service *MemberService) EmailExists(email string) (bool, error) {
	return service.repo.EmailExists(email)
}

func (service *MemberService) GetMembersByCountryCode(countryCode string) (*[]models.Member, error) {
	return service.repo.GetMembersByCountryCode(countryCode)
}

func (service *MemberService) UpdateMember(member *models.Member) error {
	return service.repo.UpdateMember(member)
}

func (service *MemberService) DeleteMember(memberID string) error {
	return service.repo.DeleteMember(memberID)
}

func (service *MemberService) GetAllMembers() (*[]models.Member, error) {
	return service.repo.GetAllMembers()
}
