package services

type ServiceContainer struct {
	OrganizationService *OrganizationService
	MemberService       *MemberService
	ProjectService      *ProjectService
	DocumentService     *DocumentService
	ProposalService     *ProposalService
	AcceptanceService   *AcceptanceService
}

func NewServiceContainer(
	organizationService *OrganizationService,
	memberService *MemberService,
	projectService *ProjectService,
	documentService *DocumentService,
	proposalService *ProposalService,
	acceptanceService *AcceptanceService,
) *ServiceContainer {
	return &ServiceContainer{
		OrganizationService: organizationService,
		MemberService:       memberService,
		ProjectService:      projectService,
		DocumentService:     documentService,
		ProposalService:     proposalService,
		AcceptanceService:   acceptanceService,
	}
}
