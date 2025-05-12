package services

import (
	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type OrganizationService struct {
	repo *repository.OrganizationRepository
}

func NewOrganizationService(repo *repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repo: repo}
}

func (service *OrganizationService) CreateNSB(nsb *models.NationalStandardBody) error {
	nsb.ID = uuid.New()
	return service.repo.CreateNSB(nsb)
}

func (service *OrganizationService) UpdateNationalTCSecretary(nsbID, newSecretaryID string) error {
	return service.repo.UpdateNationalTCSecretary(nsbID, newSecretaryID)
}

func (service *OrganizationService) FetchNSBs() (*[]models.NationalStandardBody, error) {
	return service.repo.FetchNSBs()
}

func (service *OrganizationService) SearchTechnicalCommittees(params map[string]interface{}) ([]models.TechnicalCommittee, error) {
	return service.repo.SearchTechnicalCommittees(params)
}

func (service *OrganizationService) CreateCommittee(committee any) error {
	return service.repo.CreateCommittee(committee)
}

func (service *OrganizationService) GetCommitteeByID(id string, committee any) (any, error) {
	return service.repo.GetCommitteeByID(id, committee)
}

func (service *OrganizationService) UpdateCommittee(committee any) error {
	return service.repo.UpdateCommittee(committee)
}

func (service *OrganizationService) DeleteCommittee(committeeType string, id string) error {
	return service.repo.DeleteCommittee(committeeType, id)
}

func (service *OrganizationService) AddWorkingGroupToTechnicalCommittee(tc *models.TechnicalCommittee, wg *models.WorkingGroup) error {
	return service.repo.AddWorkingGroupToTechnicalCommittee(tc, wg)
}

func (service *OrganizationService) AddEditingCommitteeToTechnicalCommittee(tc *models.TechnicalCommittee, ec *models.EditingCommittee) error {
	return service.repo.AddEditingCommitteeToTechnicalCommittee(tc, ec)
}

func (service *OrganizationService) GetCommitteeWorkingGroups(committeeID string) (*[]models.WorkingGroup, error) {
	return service.repo.GetCommitteeWorkingGroups(committeeID)
}

func (service *OrganizationService) GetCommitteeEditingCommittee(committeeID string) (*models.EditingCommittee, error) {
	return service.repo.GetCommitteeEditingCommittee(committeeID)
}

func (service *OrganizationService) FetchTechnicalCommittees() (*[]models.TechnicalCommittee, error) {
	return service.repo.FetchTechnicalCommittees()
}

func (service *OrganizationService) CompleteWorkingGroup(wg *models.WorkingGroup) error {
	return service.repo.CompleteWorkingGroup(wg)
}

func (service *OrganizationService) CreateWorkingGroup(wg *models.WorkingGroup) error {
	wg.ID = uuid.New()
	return service.repo.CreateWorkingGroup(wg)
}

func (service *OrganizationService) CreateEditingCommittee(ec *models.EditingCommittee) error {
	ec.ID = uuid.New()
	return service.repo.CreateEditingCommittee(ec)
}

func (service *OrganizationService) GetWorkingGroupByID(id string) (*models.WorkingGroup, error) {
	return service.repo.GetWorkingGroupByID(id)
}

func (service *OrganizationService) GetEditingCommitteeByID(id string) (*models.EditingCommittee, error) {
	return service.repo.GetEditingCommitteeByID(id)
}

func (service *OrganizationService) CreateTaskForce(tf *models.TaskForce) error {
	tf.ID = uuid.New()
	return service.repo.CreateTaskForce(tf)
}

func (service *OrganizationService) GetTaskForceByID(id string) (*models.TaskForce, error) {
	return service.repo.GetTaskForceByID(id)
}

func (service *OrganizationService) CreateSubCommittee(sc *models.SubCommittee) error {
	sc.ID = uuid.New()
	return service.repo.CreateSubCommittee(sc)
}

func (service *OrganizationService) AddMemberToSubCommittee(sc *models.SubCommittee, member *models.Member) error {
	return service.repo.AddMemberToSubCommittee(sc, member)
}

func (service *OrganizationService) CreateSpecializedCommittee(sc *models.SpecializedCommittee) error {
	sc.ID = uuid.New()
	return service.repo.CreateSpecializedCommittee(sc)
}

func (service *OrganizationService) GetSpecializedCommitteeByType(committeeType string) ([]*models.SpecializedCommittee, error) {
	return service.repo.GetSpecializedCommitteeByType(committeeType)
}

func (service *OrganizationService) CreateMemberState(state *models.MemberState) error {
	state.ID = uuid.New()
	return service.repo.CreateMemberState(state)
}

func (service *OrganizationService) FetchMemberStates() (*[]models.MemberState, error) {
	return service.repo.FetchMemberStates()
}

func (service *OrganizationService) UpdateCommitteeSecretary(committeeType string, committeeID string, newSecretaryID string) error {
	return service.repo.UpdateCommitteeSecretary(committeeType, committeeID, newSecretaryID)
}

func (service *OrganizationService) UpdateCommitteeChairperson(committeeType string, committeeID string, newChairpersonID string) error {
	return service.repo.UpdateCommitteeChairperson(committeeType, committeeID, newChairpersonID)
}

func (service *OrganizationService) AddMemberToARSOCouncil(id string, memberID string) error {
	return service.repo.AddMemberToARSOCouncil(id, memberID)
}

func (service *OrganizationService) AddRegionalEconomicCommunityToJointAdvisoryGroup(id string, memberID string) error {
	return service.repo.AddRegionalEconomicCommunityToJointAdvisoryGroup(id, memberID)
}

func (service *OrganizationService) AddObserverMemberToJointAdvisoryGroup(id string, memberID string) error {
	return service.repo.AddObserverMemberToJointAdvisoryGroup(id, memberID)
}

func (service *OrganizationService) AddRegionalRepresentativeToStandardsManagementCommittee(id string, memberID string) error {
	return service.repo.AddRegionalRepresentativeToStandardsManagementCommittee(id, memberID)
}

func (service *OrganizationService) AddElectedMemberToStandardsManagementCommittee(id string, memberID string) error {
	return service.repo.AddElectedMemberToStandardsManagementCommittee(id, memberID)
}

func (service *OrganizationService) AddObserverToStandardsManagementCommittee(id string, memberID string) error {
	return service.repo.AddObserverToStandardsManagementCommittee(id, memberID)
}

func (service *OrganizationService) AddMemberToTechnicalCommittee(id string, memberID string) error {
	return service.repo.AddMemberToTechnicalCommittee(id, memberID)
}

func (service *OrganizationService) AddMemberToJointTechnicalCommittee(id string, memberID string) error {
	return service.repo.AddMemberToJointTechnicalCommittee(id, memberID)
}

func (service *OrganizationService) AddMemberToSpecializedCommittee(id string, memberID string) error {
	return service.repo.AddMemberToSpecializedCommittee(id, memberID)
}

func (service *OrganizationService) AddMemberToTaskForce(id string, memberID string) error {
	return service.repo.AddMemberToTaskForce(id, memberID)
}

func (service *OrganizationService) AddMemberToWorkingGroup(id string, memberID string) error {
	return service.repo.AddMemberToWorkingGroup(id, memberID)
}

func (service *OrganizationService) RemoveMemberFromARSOCouncil(id string, memberID string) error {
	return service.repo.RemoveMemberFromARSOCouncil(id, memberID)
}

func (service *OrganizationService) RemoveRECFromJointAdvisoryGroup(id string, memberID string) error {
	return service.repo.RemoveRECFromJointAdvisoryGroup(id, memberID)
}

func (service *OrganizationService) RemoveObserverFromJointAdvisoryGroup(id string, memberID string) error {
	return service.repo.RemoveObserverFromJointAdvisoryGroup(id, memberID)
}

func (service *OrganizationService) RemoveRegionalRepresentativeFromStandardsManagementCommittee(id string, memberID string) error {
	return service.repo.RemoveRegionalRepresentativeFromStandardsManagementCommittee(id, memberID)
}

func (service *OrganizationService) RemoveRegionalElectedMemberFromStandardsManagementCommittee(id string, memberID string) error {
	return service.repo.RemoveRegionalElectedMemberFromStandardsManagementCommittee(id, memberID)
}

func (service *OrganizationService) RemoveMemberFromTechnicalCommittee(id string, memberID string) error {
	return service.repo.RemoveMemberFromTechnicalCommittee(id, memberID)
}

func (service *OrganizationService) RemoveMemberFromSpecializedCommittee(id string, memberID string) error {
	return service.repo.RemoveMemberFromSpecializedCommittee(id, memberID)
}

func (service *OrganizationService) RemoveMemberFromJointTechnicalCommittee(id string, memberID string) error {
	return service.repo.RemoveMemberFromJointTechnicalCommittee(id, memberID)
}

func (service *OrganizationService) GetArsoCouncilMembers(committeeID string) ([]*models.Member, error) {
	return service.repo.GetArsoCouncilMembers(committeeID)
}

func (service *OrganizationService) GetJointAdvisoryGroupMembers(committeeID string) ([]*models.Member, error) {
	return service.repo.GetJointAdvisoryGroupMembers(committeeID)
}

func (service *OrganizationService) GetStandardsManagementCommitteeMembers(committeeID string) ([]*models.Member, error) {
	return service.repo.GetStandardsManagementCommitteeMembers(committeeID)
}

func (service *OrganizationService) GetTechnicalCommitteeMembers(committeeID string) ([]*models.Member, error) {
	return service.repo.GetTechnicalCommitteeMembers(committeeID)
}

func (service *OrganizationService) GetSpecializedCommitteeMembers(committeeID string) ([]*models.Member, error) {
	return service.repo.GetSpecializedCommitteeMembers(committeeID)
}

func (service *OrganizationService) GetJointTechnicalCommitteeMembers(committeeID string) ([]*models.Member, error) {
	return service.repo.GetJointTechnicalCommitteeMembers(committeeID)
}

func (service *OrganizationService) GetArsoCouncil() ([]models.ARSOCouncil, error) {
	return service.repo.GetArsoCouncil()
}

func (service *OrganizationService) GetJointAdvisoryGroups() ([]models.JointAdvisoryGroup, error) {
	return service.repo.GetJointAdvisoryGroups()
}

func (service *OrganizationService) GetStandardsManagementCommittees() ([]models.StandardsManagementCommittee, error) {
	return service.repo.GetStandardsManagementCommittees()
}

func (service *OrganizationService) GetTechnicalCommittees() ([]models.TechnicalCommittee, error) {
	return service.repo.GetTechnicalCommittees()
}

func (service *OrganizationService) GetSpecializedCommittees() ([]models.SpecializedCommittee, error) {
	return service.repo.GetSpecializedCommittees()
}

func (service *OrganizationService) GetJointTechnicalCommittees() ([]models.JointTechnicalCommittee, error) {
	return service.repo.GetJointTechnicalCommittees()
}

func (service *OrganizationService) AddMemberStateToTCParticipatingCountries(id string, stateId string) error {
	return service.repo.AddMemberStateToTCParticipatingCountries(id, stateId)
}

func (service *OrganizationService) AddMemberStateToTCObserverCountries(id string, stateId string) error {
	return service.repo.AddMemberStateToTCObserverCountries(id, stateId)
}

func (service *OrganizationService) AddTCToTCEquivalentCommittees(id string, equivalentTCId string) error {
	return service.repo.AddTCToTCEquivalentCommittees(id, equivalentTCId)
}

func (service *OrganizationService) GetTCParticipatingCountries(id string) ([]*models.MemberState, error) {
	return service.repo.GetTCParticipatingCountries(id)
}

func (service *OrganizationService) GetTCObserverCountries(id string) ([]*models.MemberState, error) {
	return service.repo.GetTCObserverCountries(id)
}

func (service *OrganizationService) GetTCEquivalentCommittees(id string) ([]*models.TechnicalCommittee, error) {
	return service.repo.GetTCEquivalentCommittees(id)
}

func (service *OrganizationService) RemoveMemberStateFromTCParticipatingCountries(id string, stateId string) error {
	return service.repo.RemoveMemberStateFromTCParticipatingCountries(id, stateId)
}

func (service *OrganizationService) RemoveMemberStateFromTCObserverCountries(id string, stateId string) error {
	return service.repo.RemoveMemberStateFromTCObserverCountries(id, stateId)
}

func (service *OrganizationService) RemoveTCFromTCEquivalentCommittees(id string, equivalentTCId string) error {
	return service.repo.RemoveTCFromTCEquivalentCommittees(id, equivalentTCId)
}

func (service *OrganizationService) GetTCProjects(id string) ([]*models.Project, error) {
	return service.repo.GetTCProjects(id)
}

func (service *OrganizationService) GetTCWorkingGroups(id string) ([]*models.WorkingGroup, error) {
	return service.repo.GetTCWorkingGroups(id)
}

func (service *OrganizationService) GetTCEditingCommittee(id string) (*models.EditingCommittee, error) {
	return service.repo.GetTCEditingCommittee(id)
}

func (service *OrganizationService) GetCommitteeMeetings(id string) ([]models.Meeting, error) {
	return service.repo.GetCommitteeMeetings(id)
}
