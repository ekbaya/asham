package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
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

func (r *OrganizationRepository) FetchMemberStates(limit, offset int) (*[]models.MemberState, error) {
	var states []models.MemberState
	err := r.db.Limit(limit).Offset(offset).Preload(clause.Associations).Find(&states).Error
	return &states, err
}

// NSBs
func (r *OrganizationRepository) CreateNSB(nsb *models.NationalStandardBody) error {
	return r.db.Create(nsb).Error
}

func (r *OrganizationRepository) UpdateNationalTCSecretary(nsbID, newSecretaryID string) error {
	return r.db.Model(&models.NationalStandardBody{}).
		Where("id = ?", nsbID).
		Update("national_tc_secretary_id", newSecretaryID).Error
}

func (r *OrganizationRepository) FetchNSBs(limit, offset int) (*[]models.NationalStandardBody, error) {
	var nsbs []models.NationalStandardBody
	err := r.db.
		Limit(limit).
		Offset(offset).
		Preload("MemberState").
		Preload("NationalTCSecretary").
		Find(&nsbs).Error
	return &nsbs, err
}

// Committee Methods
func (r *OrganizationRepository) CreateCommittee(committee any) error {
	return r.db.Create(committee).Error
}

func (r *OrganizationRepository) GetCommitteeByID(id string, committee any) (any, error) {
	err := r.db.Model(&committee).
		Where("id = ?", id).
		Preload(clause.Associations).First(&committee).Error
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

func (r *OrganizationRepository) UpdateCommitteeSecretary(committeeType string, committeeID string, newSecretaryID string) error {
	switch committeeType {
	case string(models.ARSO_Council):
		return r.db.Model(&models.ARSOCouncil{}).
			Where("id = ?", committeeID).
			Update("secretary_id", newSecretaryID).Error

	case string(models.Joint_Advisory_Group):
		return r.db.Model(&models.JointAdvisoryGroup{}).
			Where("id = ?", committeeID).
			Update("secretary_id", newSecretaryID).Error

	case string(models.Standards_Management_Committee):
		return r.db.Model(&models.StandardsManagementCommittee{}).
			Where("id = ?", committeeID).
			Update("secretary_id", newSecretaryID).Error

	case string(models.Technical_Committee):
		return r.db.Model(&models.TechnicalCommittee{}).
			Where("id = ?", committeeID).
			Update("secretary_id", newSecretaryID).Error

	case string(models.Specialized_Committee):
		return r.db.Model(&models.SpecializedCommittee{}).
			Where("id = ?", committeeID).
			Update("secretary_id", newSecretaryID).Error

	case string(models.Joint_Technical_Committee):
		return r.db.Model(&models.JointTechnicalCommittee{}).
			Where("id = ?", committeeID).
			Update("secretary_id", newSecretaryID).Error

	default:
		return errors.New("unknown committee type")
	}
}

func (r *OrganizationRepository) UpdateCommitteeChairperson(committeeType string, committeeID string, newChairpersonID string) error {
	switch committeeType {
	case string(models.ARSO_Council):
		return r.db.Model(&models.ARSOCouncil{}).
			Where("id = ?", committeeID).
			Update("chairperson_id", newChairpersonID).Error

	case string(models.Joint_Advisory_Group):
		return r.db.Model(&models.JointAdvisoryGroup{}).
			Where("id = ?", committeeID).
			Update("chairperson_id", newChairpersonID).Error

	case string(models.Standards_Management_Committee):
		return r.db.Model(&models.StandardsManagementCommittee{}).
			Where("id = ?", committeeID).
			Update("chairperson_id", newChairpersonID).Error

	case string(models.Technical_Committee):
		return r.db.Model(&models.TechnicalCommittee{}).
			Where("id = ?", committeeID).
			Update("chairperson_id", newChairpersonID).Error

	case string(models.Specialized_Committee):
		return r.db.Model(&models.SpecializedCommittee{}).
			Where("id = ?", committeeID).
			Update("chairperson_id", newChairpersonID).Error

	case string(models.Joint_Technical_Committee):
		return r.db.Model(&models.JointTechnicalCommittee{}).
			Where("id = ?", committeeID).
			Update("chairperson_id", newChairpersonID).Error

	default:
		return fmt.Errorf("unknown committee type: %s", committeeType)
	}
}

func (r *OrganizationRepository) AddMemberToARSOCouncil(id string, memberID string) error {
	var council models.ARSOCouncil
	if err := r.db.First(&council, "id = ?", id).Error; err != nil {
		return err
	}

	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&council).Association("Members").Append(&member)
}

func (r *OrganizationRepository) AddRegionalEconomicCommunityToJointAdvisoryGroup(id string, memberID string) error {
	var group models.JointAdvisoryGroup
	if err := r.db.First(&group, "id = ?", id).Error; err != nil {
		return err
	}

	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&group).Association("RegionalEconomicCommunities").Append(&member)
}

func (r *OrganizationRepository) AddObserverMemberToJointAdvisoryGroup(id string, memberID string) error {
	var group models.JointAdvisoryGroup
	if err := r.db.First(&group, "id = ?", id).Error; err != nil {
		return err
	}

	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&group).Association("ObserverMembers").Append(&member)
}

func (r *OrganizationRepository) AddRegionalRepresentativeToStandardsManagementCommittee(id string, memberID string) error {
	var committee models.StandardsManagementCommittee
	if err := r.db.First(&committee, "id = ?", id).Error; err != nil {
		return err
	}

	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&committee).Association("RegionalRepresentatives").Append(&member)
}

func (r *OrganizationRepository) AddElectedMemberToStandardsManagementCommittee(id string, memberID string) error {
	var committee models.StandardsManagementCommittee
	if err := r.db.First(&committee, "id = ?", id).Error; err != nil {
		return err
	}

	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&committee).Association("ElectedMembers").Append(&member)
}

func (r *OrganizationRepository) AddObserverToStandardsManagementCommittee(id string, memberID string) error {
	var committee models.StandardsManagementCommittee
	if err := r.db.First(&committee, "id = ?", id).Error; err != nil {
		return err
	}

	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&committee).Association("Observers").Append(&member)
}

func (r *OrganizationRepository) AddMemberToTechnicalCommittee(id string, memberID string) error {
	var committee models.TechnicalCommittee
	if err := r.db.First(&committee, "id = ?", id).Error; err != nil {
		return err
	}

	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&committee).Association("CurrentMembers").Append(&member)
}

func (r *OrganizationRepository) AddMemberToJointTechnicalCommittee(id string, memberID string) error {
	var committee models.JointTechnicalCommittee
	if err := r.db.First(&committee, "id = ?", id).Error; err != nil {
		return err
	}

	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&committee).Association("JointMembers").Append(&member)
}

func (r *OrganizationRepository) AddMemberToSpecializedCommittee(id string, memberID string) error {
	var committee models.SpecializedCommittee
	if err := r.db.First(&committee, "id = ?", id).Error; err != nil {
		return err
	}

	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&committee).Association("Members").Append(&member)
}

func (r *OrganizationRepository) AddMemberToTaskForce(id string, memberID string) error {
	var committee models.TaskForce
	if err := r.db.First(&committee, "id = ?", id).Error; err != nil {
		return err
	}

	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&committee).Association("NationalDeligations").Append(&member)
}

func (r *OrganizationRepository) AddMemberToWorkingGroup(id string, memberID string) error {
	var committee models.WorkingGroup
	if err := r.db.First(&committee, "id = ?", id).Error; err != nil {
		return err
	}

	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&committee).Association("Experts").Append(&member)
}
func (r *OrganizationRepository) AddMemberToEditingCommittee(id string, memberID string) error {
	var committee models.EditingCommittee
	if err := r.db.First(&committee, "id = ?", id).Error; err != nil {
		return err
	}

	var member models.Member
	if err := r.db.First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&committee).Association("Editors").Append(&member)
}

func (r *OrganizationRepository) RemoveMemberFromARSOCouncil(committeeID string, memberID string) error {
	memberUUID, err := uuid.Parse(memberID)
	if err != nil {
		return err
	}
	var committee models.ARSOCouncil
	if err := r.db.First(&committee, "id = ?", committeeID).Error; err != nil {
		return err
	}
	return r.db.Model(&committee).Association("Members").Delete(&models.Member{ID: memberUUID})
}

func (r *OrganizationRepository) RemoveRECFromJointAdvisoryGroup(committeeID string, memberID string) error {
	memberUUID, err := uuid.Parse(memberID)
	if err != nil {
		return err
	}
	var committee models.JointAdvisoryGroup
	if err := r.db.First(&committee, "id = ?", committeeID).Error; err != nil {
		return err
	}
	return r.db.Model(&committee).Association("RegionalEconomicCommunities").Delete(&models.Member{ID: memberUUID})
}

func (r *OrganizationRepository) RemoveObserverFromJointAdvisoryGroup(committeeID string, memberID string) error {
	memberUUID, err := uuid.Parse(memberID)
	if err != nil {
		return err
	}
	var committee models.JointAdvisoryGroup
	if err := r.db.First(&committee, "id = ?", committeeID).Error; err != nil {
		return err
	}
	return r.db.Model(&committee).Association("ObserverMembers").Delete(&models.Member{ID: memberUUID})
}

func (r *OrganizationRepository) RemoveRegionalRepresentativeFromStandardsManagementCommittee(committeeID string, memberID string) error {
	memberUUID, err := uuid.Parse(memberID)
	if err != nil {
		return err
	}
	var committee models.StandardsManagementCommittee
	if err := r.db.First(&committee, "id = ?", committeeID).Error; err != nil {
		return err
	}
	return r.db.Model(&committee).Association("RegionalRepresentatives").Delete(&models.Member{ID: memberUUID})
}

func (r *OrganizationRepository) RemoveRegionalElectedMemberFromStandardsManagementCommittee(committeeID string, memberID string) error {
	memberUUID, err := uuid.Parse(memberID)
	if err != nil {
		return err
	}
	var committee models.StandardsManagementCommittee
	if err := r.db.First(&committee, "id = ?", committeeID).Error; err != nil {
		return err
	}
	return r.db.Model(&committee).Association("ElectedMembers").Delete(&models.Member{ID: memberUUID})
}

func (r *OrganizationRepository) RemoveMemberFromTechnicalCommittee(committeeID string, memberID string) error {
	memberUUID, err := uuid.Parse(memberID)
	if err != nil {
		return err
	}
	var committee models.TechnicalCommittee
	if err := r.db.First(&committee, "id = ?", committeeID).Error; err != nil {
		return err
	}
	return r.db.Model(&committee).Association("CurrentMembers").Delete(&models.Member{ID: memberUUID})
}

func (r *OrganizationRepository) RemoveMemberFromSpecializedCommittee(committeeID string, memberID string) error {
	memberUUID, err := uuid.Parse(memberID)
	if err != nil {
		return err
	}
	var committee models.SpecializedCommittee
	if err := r.db.First(&committee, "id = ?", committeeID).Error; err != nil {
		return err
	}
	return r.db.Model(&committee).Association("Members").Delete(&models.Member{ID: memberUUID})
}

func (r *OrganizationRepository) RemoveMemberFromJointTechnicalCommittee(committeeID string, memberID string) error {
	memberUUID, err := uuid.Parse(memberID)
	if err != nil {
		return err
	}
	var committee models.JointTechnicalCommittee
	if err := r.db.First(&committee, "id = ?", committeeID).Error; err != nil {
		return err
	}
	return r.db.Model(&committee).Association("JointMembers").Delete(&models.Member{ID: memberUUID})
}

func (r *OrganizationRepository) RemoveMemberFromEditingCommittee(committeeID string, memberID string) error {
	memberUUID, err := uuid.Parse(memberID)
	if err != nil {
		return err
	}
	var committee models.EditingCommittee
	if err := r.db.First(&committee, "id = ?", committeeID).Error; err != nil {
		return err
	}
	return r.db.Model(&committee).Association("Editors").Delete(&models.Member{ID: memberUUID})
}

func (r *OrganizationRepository) GetArsoCouncilMembers(committeeID string) ([]*models.Member, error) {
	var committee models.ARSOCouncil
	if err := r.db.Preload("Members").First(&committee, "id = ?", committeeID).Error; err != nil {
		return nil, err
	}
	return committee.Members, nil
}

func (r *OrganizationRepository) GetJointAdvisoryGroupMembers(committeeID string) ([]*models.Member, error) {
	var committee models.JointAdvisoryGroup
	if err := r.db.Preload("ObserverMembers").First(&committee, "id = ?", committeeID).Error; err != nil {
		return nil, err
	}
	return committee.ObserverMembers, nil
}

func (r *OrganizationRepository) GetStandardsManagementCommitteeMembers(committeeID string) ([]*models.Member, error) {
	var committee models.StandardsManagementCommittee
	if err := r.db.Preload("ElectedMembers").First(&committee, "id = ?", committeeID).Error; err != nil {
		return nil, err
	}
	return committee.ElectedMembers, nil
}

func (r *OrganizationRepository) GetTechnicalCommitteeMembers(committeeID string) ([]*models.Member, error) {
	var committee models.TechnicalCommittee
	if err := r.db.Preload("CurrentMembers").First(&committee, "id = ?", committeeID).Error; err != nil {
		return nil, err
	}
	return committee.CurrentMembers, nil
}

func (r *OrganizationRepository) GetSpecializedCommitteeMembers(committeeID string) ([]*models.Member, error) {
	var committee models.SpecializedCommittee
	if err := r.db.Preload("Members").First(&committee, "id = ?", committeeID).Error; err != nil {
		return nil, err
	}
	return committee.Members, nil
}

func (r *OrganizationRepository) GetJointTechnicalCommitteeMembers(committeeID string) ([]*models.Member, error) {
	var committee models.JointTechnicalCommittee
	if err := r.db.Preload("JointMembers").First(&committee, "id = ?", committeeID).Error; err != nil {
		return nil, err
	}
	return committee.JointMembers, nil
}

func (r *OrganizationRepository) GetArsoCouncil() ([]models.ARSOCouncil, error) {
	var committees []models.ARSOCouncil
	if err := r.db.Preload(clause.Associations).Find(&committees).Error; err != nil {
		return nil, err
	}
	return committees, nil
}

func (r *OrganizationRepository) GetJointAdvisoryGroups() ([]models.JointAdvisoryGroup, error) {
	var committees []models.JointAdvisoryGroup
	if err := r.db.Preload(clause.Associations).Find(&committees).Error; err != nil {
		return nil, err
	}
	return committees, nil
}

func (r *OrganizationRepository) GetStandardsManagementCommittees() ([]models.StandardsManagementCommittee, error) {
	var committees []models.StandardsManagementCommittee
	if err := r.db.Preload(clause.Associations).Find(&committees).Error; err != nil {
		return nil, err
	}
	return committees, nil
}

func (r *OrganizationRepository) GetTechnicalCommittees() ([]models.TechnicalCommittee, error) {
	var committees []models.TechnicalCommittee
	if err := r.db.Preload(clause.Associations).Find(&committees).Error; err != nil {
		return nil, err
	}
	return committees, nil
}

func (r *OrganizationRepository) GetSpecializedCommittees() ([]models.SpecializedCommittee, error) {
	var committees []models.SpecializedCommittee
	if err := r.db.Preload(clause.Associations).Find(&committees).Error; err != nil {
		return nil, err
	}
	return committees, nil
}

func (r *OrganizationRepository) GetJointTechnicalCommittees() ([]models.JointTechnicalCommittee, error) {
	var committees []models.JointTechnicalCommittee
	if err := r.db.Preload(clause.Associations).Find(&committees).Error; err != nil {
		return nil, err
	}
	return committees, nil
}

// Technical Committee Specific Methods
func (r *OrganizationRepository) AddWorkingGroupToTechnicalCommittee(tc *models.TechnicalCommittee, wg *models.WorkingGroup) error {
	wg.ParentTC = tc
	return r.db.Save(wg).Error
}

func (r *OrganizationRepository) AddEditingCommitteeToTechnicalCommittee(tc *models.TechnicalCommittee, ed *models.EditingCommittee) error {
	ed.ParentTC = tc
	return r.db.Save(ed).Error
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

func (r *OrganizationRepository) GetEditingCommitteeByID(id string) (*models.EditingCommittee, error) {
	var ec models.EditingCommittee
	result := r.db.First(&ec, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ec, nil
}

func (r *OrganizationRepository) GetCommitteeWorkingGroups(committeeID string) (*[]models.WorkingGroup, error) {
	var wgs []models.WorkingGroup
	err := r.db.Where("parent_tc_id = ?", committeeID).Find(&wgs).Error
	return &wgs, err
}

func (r *OrganizationRepository) GetCommitteeEditingCommittee(committeeID string) (*models.EditingCommittee, error) {
	var ec models.EditingCommittee
	err := r.db.Where("parent_tc_id = ?", committeeID).First(&ec).Error
	return &ec, err
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

func (r *OrganizationRepository) AddMemberStateToTCParticipatingCountries(id string, stateId string) error {
	var tc models.TechnicalCommittee
	if err := r.db.First(&tc, "id = ?", id).Error; err != nil {
		return err
	}

	var state models.MemberState
	if err := r.db.First(&state, "id = ?", stateId).Error; err != nil {
		return err
	}

	return r.db.Model(&tc).Association("ParticipatingCountries").Append(&state)
}

func (r *OrganizationRepository) AddMemberStateToTCObserverCountries(id string, stateId string) error {
	var tc models.TechnicalCommittee
	if err := r.db.First(&tc, "id = ?", id).Error; err != nil {
		return err
	}

	var state models.MemberState
	if err := r.db.First(&state, "id = ?", stateId).Error; err != nil {
		return err
	}

	return r.db.Model(&tc).Association("ObserverCountries").Append(&state)
}

func (r *OrganizationRepository) AddTCToTCEquivalentCommittees(id string, equivalentTCId string) error {
	var tc models.TechnicalCommittee
	if err := r.db.First(&tc, "id = ?", id).Error; err != nil {
		return err
	}

	var equivalent_committee models.TechnicalCommittee
	if err := r.db.First(&equivalent_committee, "id = ?", equivalentTCId).Error; err != nil {
		return err
	}

	return r.db.Model(&tc).Association("EquivalentCommittees").Append(&equivalent_committee)
}

func (r *OrganizationRepository) GetTCParticipatingCountries(id string) ([]*models.MemberState, error) {
	var tc models.TechnicalCommittee
	if err := r.db.Where("id = ?", id).Preload("ParticipatingCountries").First(&tc).Error; err != nil {
		return nil, err
	}
	return tc.ParticipatingCountries, nil
}

func (r *OrganizationRepository) GetTCObserverCountries(id string) ([]*models.MemberState, error) {
	var tc models.TechnicalCommittee
	if err := r.db.Preload("ObserverCountries").Where("id = ?", id).First(&tc).Error; err != nil {
		return nil, err
	}
	return tc.ObserverCountries, nil
}

func (r *OrganizationRepository) GetTCEquivalentCommittees(id string) ([]*models.TechnicalCommittee, error) {
	var tc models.TechnicalCommittee
	if err := r.db.Where("id = ?", id).Preload("EquivalentCommittees").First(&tc).Error; err != nil {
		return nil, err
	}
	return tc.EquivalentCommittees, nil
}

func (r *OrganizationRepository) RemoveMemberStateFromTCParticipatingCountries(id string, stateId string) error {
	var tc models.TechnicalCommittee
	if err := r.db.First(&tc, "id = ?", id).Error; err != nil {
		return err
	}

	var state models.MemberState
	if err := r.db.First(&state, "id = ?", stateId).Error; err != nil {
		return err
	}

	return r.db.Model(&tc).Association("ParticipatingCountries").Delete(&state)
}

func (r *OrganizationRepository) RemoveMemberStateFromTCObserverCountries(id string, stateId string) error {
	var tc models.TechnicalCommittee
	if err := r.db.First(&tc, "id = ?", id).Error; err != nil {
		return err
	}

	var state models.MemberState
	if err := r.db.First(&state, "id = ?", stateId).Error; err != nil {
		return err
	}

	return r.db.Model(&tc).Association("ObserverCountries").Delete(&state)
}

func (r *OrganizationRepository) RemoveTCFromTCEquivalentCommittees(id string, equivalentTCId string) error {
	var tc models.TechnicalCommittee
	if err := r.db.First(&tc, "id = ?", id).Error; err != nil {
		return err
	}

	var equivalent_committee models.TechnicalCommittee
	if err := r.db.First(&equivalent_committee, "id = ?", equivalentTCId).Error; err != nil {
		return err
	}

	return r.db.Model(&tc).Association("EquivalentCommittees").Delete(&equivalent_committee)
}

func (r *OrganizationRepository) GetTCProjects(id string) ([]*models.Project, error) {
	var tc models.TechnicalCommittee
	if err := r.db.Where("id = ?", id).Preload("Projects").First(&tc).Error; err != nil {
		return nil, err
	}
	return tc.Projects, nil
}

func (r *OrganizationRepository) GetTCWorkingGroups(id string) ([]*models.WorkingGroup, error) {
	var tc models.TechnicalCommittee
	if err := r.db.Where("id = ?", id).Preload("WorkingGroups").First(&tc).Error; err != nil {
		return nil, err
	}
	return tc.WorkingGroups, nil
}

func (r *OrganizationRepository) GetTCEditingCommittee(id string) (*models.EditingCommittee, error) {
	var tc models.TechnicalCommittee
	if err := r.db.Where("id = ?", id).Preload("EditingCommittee").First(&tc).Error; err != nil {
		return nil, err
	}
	return tc.EditingCommittee, nil
}

func (r *OrganizationRepository) GetCommitteeMeetings(id string) ([]models.Meeting, error) {
	var meetings []models.Meeting
	if err := r.db.Where("committee_id = ?", id).Preload(clause.Associations).Find(&meetings).Error; err != nil {
		return nil, err
	}
	return meetings, nil
}

func (r *OrganizationRepository) CreateEditingCommittee(ec *models.EditingCommittee) error {
	return r.db.Create(ec).Error
}
