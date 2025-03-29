package services

type ServiceContainer struct {
	OrganizationService *OrganizationService
	MemberService       *MemberService
	ProjectService      *ProjectService
	DocumentService     *DocumentService
	ProposalService     *ProposalService
}

func NewServiceContainer(
	organizationService *OrganizationService,
	memberService *MemberService,
	projectService *ProjectService,
	documentService *DocumentService,
	proposalService *ProposalService,
) *ServiceContainer {
	return &ServiceContainer{
		OrganizationService: organizationService,
		MemberService:       memberService,
		ProjectService:      projectService,
		DocumentService:     documentService,
		ProposalService:     proposalService,
	}
}
