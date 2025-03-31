package services

type ServiceContainer struct {
	OrganizationService *OrganizationService
	MemberService       *MemberService
	ProjectService      *ProjectService
	DocumentService     *DocumentService
	ProposalService     *ProposalService
	AcceptanceService   *AcceptanceService
	CommentService      *CommentService
}

func NewServiceContainer(
	organizationService *OrganizationService,
	memberService *MemberService,
	projectService *ProjectService,
	documentService *DocumentService,
	proposalService *ProposalService,
	acceptanceService *AcceptanceService,
	commentService *CommentService,
) *ServiceContainer {
	return &ServiceContainer{
		OrganizationService: organizationService,
		MemberService:       memberService,
		ProjectService:      projectService,
		DocumentService:     documentService,
		ProposalService:     proposalService,
		AcceptanceService:   acceptanceService,
		CommentService:      commentService,
	}
}
