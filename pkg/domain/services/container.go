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
	}
}
