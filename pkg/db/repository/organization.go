package repository

import (
	"errors"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

// Member states
func (r *OrganizationRepository) CreateMemberState(state *models.MemberState) error {
	return r.db.Create(state).Error
}

func (r *OrganizationRepository) FetchMemberStates() (*[]models.MemberState, error) {
	var states []models.MemberState
	err := r.db.Find(&states).Error
	return &states, err
}

// NSBs
func (r *OrganizationRepository) CreateNSB(nsb *models.NationalStandardBody) error {
	return r.db.Create(nsb).Error
}

func (r *OrganizationRepository) FetchNSBs() (*[]models.NationalStandardBody, error) {
	var nsbs []models.NationalStandardBody
	err := r.db.Find(&nsbs).Error
	return &nsbs, err
}

// Committee Methods
func (r *OrganizationRepository) CreateCommittee(committee any) error {
	return r.db.Create(committee).Error
}

func (r *OrganizationRepository) GetCommitteeByID(id string, committee any) (any, error) {
	err := r.db.Preload(clause.Associations).First(committee, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return committee, nil
}

func (r *OrganizationRepository) UpdateCommittee(committee any) error {
	return r.db.Save(&committee).Error
}

func (r *OrganizationRepository) DeleteCommittee(committeeType string, id string) error {
	switch committeeType {
	case string(models.ARSO_Council):
		return r.db.Delete(&models.ARSOCouncil{}, "id = ?", id).Error
	case string(models.Joint_Advisory_Group):
		return r.db.Delete(&models.JointAdvisoryGroup{}, "id = ?", id).Error
	case string(models.Standards_Management_Committee):
		return r.db.Delete(&models.StandardsManagementCommittee{}, "id = ?", id).Error
	case string(models.Technical_Committee):
		return r.db.Delete(&models.TechnicalCommittee{}, "id = ?", id).Error
	case string(models.Specialized_Committee):
		return r.db.Delete(&models.SpecializedCommittee{}, "id = ?", id).Error
	case string(models.Joint_Technical_Committee):
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

func (r *OrganizationRepository) FetchTechnicalCommittees() (*[]models.TechnicalCommittee, error) {
	var committees []models.TechnicalCommittee
	err := r.db.Find(&committees).Error
	return &committees, err
}

func (r *OrganizationRepository) SearchTechnicalCommittees(params map[string]interface{}) ([]models.TechnicalCommittee, error) {
	var committees []models.TechnicalCommittee
	var total int64

	query := r.db.Model(&models.TechnicalCommittee{})

	// Apply filters
	if name, ok := params["name"].(string); ok && name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	if code, ok := params["code"].(string); ok && code != "" {
		query = query.Where("code ILIKE ?", "%"+code+"%")
	}

	if scope, ok := params["scope"].(string); ok && scope != "" {
		query = query.Where("scope ILIKE ?", "%"+scope+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Get projects with pagination
	err := query.
		Preload(clause.Associations).
		Find(&committees).Error

	return committees, err
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

func (r *OrganizationRepository) GetCommitteeWorkingGroups(committeeID string) (*[]models.WorkingGroup, error) {
	var wgs []models.WorkingGroup
	err := r.db.Where("parent_tc_id = ?", committeeID).Find(&wgs).Error
	return &wgs, err
}

// Task Force Methods
func (r *OrganizationRepository) CreateTaskForce(tf *models.TaskForce) error {
	return r.db.Create(tf).Error
}

func (r *OrganizationRepository) GetTaskForceByID(id string) (*models.TaskForce, error) {
	var tf models.TaskForce
	result := r.db.First(&tf, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tf, nil
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
