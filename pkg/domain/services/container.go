package services

type ServiceContainer struct {
	OrganizationService *OrganizationService
	MemberService       *MemberService
	ProjectService      *ProjectService
	DocumentService     *DocumentService
}

func NewServiceContainer(
	organizationService *OrganizationService,
	memberService *MemberService,
	projectService *ProjectService,
	documentService *DocumentService,
) *ServiceContainer {
	return &ServiceContainer{
		OrganizationService: organizationService,
		MemberService:       memberService,
		ProjectService:      projectService,
		DocumentService:     documentService,
	}
}
