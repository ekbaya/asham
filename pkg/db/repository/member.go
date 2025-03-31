package repository

import (
	"errors"

	"github.com/ekbaya/asham/pkg/domain/models"
	"gorm.io/gorm"
)

type MemberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) CreateMember(member *models.Member) error {
	return r.db.Create(member).Error
}

func (r *MemberRepository) GetMemberByID(id string) (*models.Member, error) {
	var member models.Member
	result := r.db.Preload("NationalStandardBody").First(&member, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &member, nil
}

func (r *MemberRepository) GetAllMembers() (*[]models.Member, error) {
	var members []models.Member
	result := r.db.Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return &members, nil
}

func (r *MemberRepository) GetMemberByEmail(email string) (*models.Member, error) {
	var member models.Member
	result := r.db.Where("email = ?", email).First(&member)

	if result.Error != nil {
		// Check if the error is "record not found"
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil instead of an error when member not found
		}
		return nil, result.Error // Return other errors as is
	}

	return &member, nil
}

func (r *MemberRepository) EmailExists(email string) (bool, error) {
	var count int64
	result := r.db.Model(&models.Member{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *MemberRepository) GetMembersByCountryCode(countryCode string) (*[]models.Member, error) {
	var members []models.Member
	result := r.db.Where("country_code = ?", countryCode).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return &members, nil
}

func (r *MemberRepository) UpdateMember(member *models.Member) error {
	return r.db.Save(member).Error
}

func (r *MemberRepository) DeleteMember(memberID string) error {
	return r.db.Delete(&models.Member{}, "id = ?", memberID).Error
}
