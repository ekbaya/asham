package services

type ServiceContainer struct {
	OrganizationService *OrganizationService
	MemberService       *MemberService
}

func NewServiceContainer(
	organizationService *OrganizationService,
	memberService *MemberService,
) *ServiceContainer {
	return &ServiceContainer{
		OrganizationService: organizationService,
		MemberService:       memberService,
	}
}
