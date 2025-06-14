package services

type ServiceContainer struct {
	OrganizationService         *OrganizationService
	MemberService               *MemberService
	ProjectService              *ProjectService
	DocumentService             *DocumentService
	ProposalService             *ProposalService
	AcceptanceService           *AcceptanceService
	CommentService              *CommentService
	EmailService                *EmailService
	NationalConsultationService *NationalConsultationService
	BallotingService            *BallotingService
	MeetingService              *MeetingService
	LibraryService              *LibraryService
	StandardService             *StandardService
	RbacService                 *RbacService
	TokenManager                *TokenManager
	PermissionResourceService   *PermissionResourceService
}

func NewServiceContainer(
	organizationService *OrganizationService,
	memberService *MemberService,
	projectService *ProjectService,
	documentService *DocumentService,
	proposalService *ProposalService,
	acceptanceService *AcceptanceService,
	commentService *CommentService,
	emailService *EmailService,
	nationalConsultationService *NationalConsultationService,
	ballotingService *BallotingService,
	meetingService *MeetingService,
	libraryService *LibraryService,
	standardService *StandardService,
	rbacService *RbacService,
	tokenManager *TokenManager,
	permissionResourceService *PermissionResourceService,
) *ServiceContainer {
	return &ServiceContainer{
		OrganizationService:         organizationService,
		MemberService:               memberService,
		ProjectService:              projectService,
		DocumentService:             documentService,
		ProposalService:             proposalService,
		AcceptanceService:           acceptanceService,
		CommentService:              commentService,
		EmailService:                emailService,
		NationalConsultationService: nationalConsultationService,
		BallotingService:            ballotingService,
		MeetingService:              meetingService,
		LibraryService:              libraryService,
		StandardService:             standardService,
		RbacService:                 rbacService,
		TokenManager:                tokenManager,
		PermissionResourceService:   permissionResourceService,
	}
}
