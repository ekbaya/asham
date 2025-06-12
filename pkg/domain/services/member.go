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

	go func() {
		err := service.emailService.SendWelcomeEmail(member.Email, member.FirstName, clearPassword)
		if err != nil {
			fmt.Printf("Failed to send welcome email to %s: %v\n", member.Email, err)
		}
	}()

	return nil
}

func (service *MemberService) Login(email, password string, channel models.UserType) (string, string, error) {
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

	if channel == models.Internal && user.Type == models.External {
		return "", "", errors.New("user is not authorized to login to admin panel")
	}

	// Generate JWT token
	token, err := models.GenerateJWT(*user)
	if err != nil {
		return "", "", errors.New("failed to generate token")
	}

	// Generate JWT refresh token
	refreshToken, err := models.GenerateRefreshToken(*user)
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

func (service *MemberService) GetAllMembers(limit, offset int) (*[]models.Member, int64, error) {
	return service.repo.GetAllMembers(limit, offset)
}

func (service *MemberService) Logout(accessToken, refreshToken string) error {
	// Validate the access token to make sure it's not already invalid
	claims, err := models.ValidateJWT(accessToken)
	if err != nil {
		return nil
	}

	// Use Redis to blacklist the tokens
	err = models.Logout(accessToken, refreshToken)
	if err != nil {
		return errors.New("failed to invalidate tokens: " + err.Error())
	}

	// log this activity
	userID := claims.UserID
	service.logUserActivity(userID, "logout")

	return nil
}

func (service *MemberService) LogoutAll(userID string) error {
	err := models.LogoutUser(userID)
	if err != nil {
		return errors.New("failed to invalidate all tokens: " + err.Error())
	}

	service.logUserActivity(userID, "logout-all-devices")
	return nil
}

func (service *MemberService) logUserActivity(userID, activity string) {
	fmt.Printf("User %s: %s at %s\n", userID, activity, time.Now().Format(time.RFC3339))
}
