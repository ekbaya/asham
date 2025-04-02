package wire

import (
	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(
	repository.NewOrganizationRepository,
	services.NewOrganizationService,
	repository.NewMemberRepository,
	GetEmailConfigurations,
	services.NewEmailService,
	services.NewMemberService,
	repository.NewProjectRepository,
	services.NewProjectService,
	repository.NewDocumentRepository,
	services.NewDocumentService,
	repository.NewProposalRepository,
	services.NewProposalService,
	repository.NewAcceptanceRepository,
	services.NewAcceptanceService,
	repository.NewCommentRepository,
	services.NewCommentService,
)

func GetEmailConfigurations() *services.EmailConfig {
	emailConfig := services.EmailConfig{
		Host:     "live.smtp.mailtrap.io",
		Port:     587,
		Username: "smtp@mailtrap.io",
		Password: "61dc207f67686fdb2aadbe5bc179fa71",
		From:     "no-reply@collectwave.com",
	}
	return &emailConfig
}
