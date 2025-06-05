package wire

import (
	"strconv"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/ekbaya/asham/pkg/config"
	"log"
	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/google/wire"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
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
	GetGraphServiceClient,
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
	repository.NewRbacRepository,
	services.NewRbacService,
	GetMSAzureConfig,
	services.NewTokenManager,
)

func GetEmailConfigurations() *services.EmailConfig {
	globalConfig := config.GetConfig()
	port, err := strconv.Atoi(globalConfig.EmailConfig.Port)
	if err != nil {
		port = 587
	}
	emailConfig := services.EmailConfig{
		Host:     globalConfig.EmailConfig.Host,
		Port:     port,
		Username: globalConfig.EmailConfig.Username,
		Password: globalConfig.EmailConfig.Password,
		From:     globalConfig.EmailConfig.From,
	}

	log.Printf("[EmailConfig] Host: %s, Port: %d, Username: %s, From: %s", host, port, username, from)

	return &emailConfig
}

func GetMSAzureConfig() *services.MSAzureConfig {
	globalConfig := config.GetConfig()
	config := services.MSAzureConfig{
		TenantID:     globalConfig.AZURE_TENANT_ID,
		ClientID:     globalConfig.AZURE_CLIENT_ID,
		ClientSecret: globalConfig.AZURE_CLIENT_SECRET,
	}
	return &config
}

func GetGraphServiceClient() *msgraphsdk.GraphServiceClient {
	config := GetMSAzureConfig()
	cred, err := azidentity.NewClientSecretCredential(
		config.TenantID,
		config.ClientID,
		config.ClientSecret,
		nil,
	)
	if err != nil {
		return nil
	}
	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		return nil
	}
	return client
}
