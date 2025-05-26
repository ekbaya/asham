package wire

import (
	"log"
	"os"
	"strconv"

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
	repository.NewConsultationRepository,
	services.NewNationalConsultationService,
	repository.NewBallotingRepository,
	services.NewBallotingService,
	repository.NewMeetingRepository,
	services.NewMeetingService,
	repository.NewLibraryRepository,
	services.NewLibraryService,
	repository.NewStandardRepository,
	services.NewStandardService,
)

func GetEmailConfigurations() *services.EmailConfig {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		host = "live.smtp.mailtrap.io"
	}
	portStr := os.Getenv("SMTP_PORT")
	port := 587
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}
	username := os.Getenv("SMTP_USERNAME")
	if username == "" {
		username = "smtp@mailtrap.io"
	}
	password := os.Getenv("SMTP_PASSWORD")
	if password == "" {
		password = "61dc207f67686fdb2aadbe5bc179fa71"
	}
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = "no-reply@collectwave.com"
	}

	emailConfig := services.EmailConfig{
		Host:              host,
		Port:              port,
		Username:          username,
		Password:          password,
		From:              from,
		EmailTemplatePath: "../templates/welcome_email.html",
	}

	log.Printf("[EmailConfig] Host: %s, Port: %d, Username: %s, From: %s", host, port, username, from)

	return &emailConfig
}
