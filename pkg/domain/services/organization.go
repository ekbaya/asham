package services

import (
	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
)

type OrganizationService struct {
	repo *repository.OrganizationRepository
}

func NewOrganizationService(repo *repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repo: repo}
}

func (service *OrganizationService) CreateNSB(nsb *models.NationalStandardBody) error {
	return service.repo.CreateNSB(nsb)
}

func (service *OrganizationService) FetchNSBs() (*[]models.NationalStandardBody, error) {
	return service.repo.FetchNSBs()
}

func (service *OrganizationService) CreateCommittee(committee any) error {
	return service.repo.CreateCommittee(committee)
}

func (service *OrganizationService) GetCommitteeByID(id string, committee any) error {
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

func (service *OrganizationService) CompleteWorkingGroup(wg *models.WorkingGroup) error {
	return service.repo.CompleteWorkingGroup(wg)
}

func (service *OrganizationService) CreateWorkingGroup(wg *models.WorkingGroup) error {
	return service.repo.CreateWorkingGroup(wg)
}

func (service *OrganizationService) GetWorkingGroupByID(id string) (*models.WorkingGroup, error) {
	return service.repo.GetWorkingGroupByID(id)
}

func (service *OrganizationService) CreateTaskForce(tf *models.TaskForce) error {
	return service.repo.CreateTaskForce(tf)
}

func (service *OrganizationService) GetTaskForceByID(id string) (*models.TaskForce, error) {
	return service.repo.GetTaskForceByID(id)
}

func (service *OrganizationService) CreateSubCommittee(sc *models.SubCommittee) error {
	return service.repo.CreateSubCommittee(sc)
}

func (service *OrganizationService) AddMemberToSubCommittee(sc *models.SubCommittee, member *models.Member) error {
	return service.repo.AddMemberToSubCommittee(sc, member)
}

func (service *OrganizationService) CreateSpecializedCommittee(sc *models.SpecializedCommittee) error {
	return service.repo.CreateSpecializedCommittee(sc)
}

func (service *OrganizationService) GetSpecializedCommitteeByType(committeeType string) ([]*models.SpecializedCommittee, error) {
	return service.repo.GetSpecializedCommitteeByType(committeeType)
}

func (service *OrganizationService) CreateMemberState(state *models.MemberState) error {
	return service.repo.CreateMemberState(state)
}

func (service *OrganizationService) FetchMemberStates() (*[]models.MemberState, error) {
	return service.repo.FetchMemberStates()
}
