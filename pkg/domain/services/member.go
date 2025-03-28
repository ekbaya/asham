package services

import (
	"errors"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/utilities"
)

type MemberService struct {
	repo *repository.MemberRepository
}

func NewMemberService(repo *repository.MemberRepository) *MemberService {
	return &MemberService{repo: repo}
}

func (service *MemberService) CreateMember(member *models.Member) error {
	// Hash password
	hashedPassword, err := utilities.HashPassword(member.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Create user
	member.HashedPassword = hashedPassword
	return service.repo.CreateMember(member)
}

func (service *MemberService) Login(email, password string) (string, string, error) {
	var user *models.Member
	var err error

	user, err = service.repo.GetMemberByEmail(email)

	// Handle error if user is not found
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Verify password
	if !utilities.CheckPasswordHash(password, user.HashedPassword) {
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
