package repository

import (
	"errors"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

// NSBs
func (r *OrganizationRepository) CreateNSB(nsb models.NationalStandardBody) error {
	return r.db.Create(nsb).Error
}

func (r *OrganizationRepository) FetchNSBs() ([]models.NationalStandardBody, error) {
	var nsbs []models.NationalStandardBody
	err := r.db.Find(nsbs).Error
	return nsbs, err
}

// Committee Methods
func (r *OrganizationRepository) CreateCommittee(committee any) error {
	return r.db.Create(committee).Error
}

func (r *OrganizationRepository) GetCommitteeByID(id string, committee any) error {
	return r.db.First(committee, "id = ?", id).Error
}

func (r *OrganizationRepository) UpdateCommittee(committee any) error {
	return r.db.Save(committee).Error
}

func (r *OrganizationRepository) DeleteCommittee(committeeType string, id string) error {
	switch committeeType {
	case "ARSOCouncil":
		return r.db.Delete(&models.ARSOCouncil{}, "id = ?", id).Error
	case "JointAdvisoryGroup":
		return r.db.Delete(&models.JointAdvisoryGroup{}, "id = ?", id).Error
	case "StandardsManagementCommittee":
		return r.db.Delete(&models.StandardsManagementCommittee{}, "id = ?", id).Error
	case "TechnicalCommittee":
		return r.db.Delete(&models.TechnicalCommittee{}, "id = ?", id).Error
	case "SpecializedCommittee":
		return r.db.Delete(&models.SpecializedCommittee{}, "id = ?", id).Error
	case "JointTechnicalCommittee":
		return r.db.Delete(&models.JointTechnicalCommittee{}, "id = ?", id).Error
	default:
		return errors.New("unknown committee type")
	}
}

// Technical Committee Specific Methods
func (r *OrganizationRepository) AddWorkingGroupToTechnicalCommittee(tc *models.TechnicalCommittee, wg *models.WorkingGroup) error {
	wg.ParentTC = tc
	return r.db.Save(wg).Error
}

func (r *OrganizationRepository) CompleteWorkingGroup(wg *models.WorkingGroup) error {
	now := time.Now()
	wg.CompletedAt = &now
	return r.db.Save(wg).Error
}

// Working Group Methods
func (r *OrganizationRepository) CreateWorkingGroup(wg *models.WorkingGroup) error {
	return r.db.Create(wg).Error
}

func (r *OrganizationRepository) GetWorkingGroupByID(id string) (*models.WorkingGroup, error) {
	var wg models.WorkingGroup
	result := r.db.First(&wg, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wg, nil
}

// Task Force Methods
func (r *OrganizationRepository) CreateTaskForce(wg *models.TaskForce) error {
	return r.db.Create(wg).Error
}

func (r *OrganizationRepository) GetTaskForceByID(id string) (*models.TaskForce, error) {
	var wg models.TaskForce
	result := r.db.First(&wg, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wg, nil
}

// SubCommittee Methods
func (r *OrganizationRepository) CreateSubCommittee(sc *models.SubCommittee) error {
	return r.db.Create(sc).Error
}

func (r *OrganizationRepository) AddMemberToSubCommittee(sc *models.SubCommittee, member *models.Member) error {
	sc.Members = append(sc.Members, member)
	return r.db.Save(sc).Error
}

// Specialized Committee Methods
func (r *OrganizationRepository) CreateSpecializedCommittee(sc *models.SpecializedCommittee) error {
	return r.db.Create(sc).Error
}

func (r *OrganizationRepository) GetSpecializedCommitteeByType(committeeType string) ([]*models.SpecializedCommittee, error) {
	var committees []*models.SpecializedCommittee
	result := r.db.Where("type = ?", committeeType).Find(&committees)
	return committees, result.Error
}

// Utility Methods
func (r *OrganizationRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *OrganizationRepository) RollbackTransaction(tx *gorm.DB) {
	tx.Rollback()
}

func (r *OrganizationRepository) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}
