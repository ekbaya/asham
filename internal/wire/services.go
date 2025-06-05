package wire

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
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
	emailConfig := services.EmailConfig{
		Host:     "live.smtp.mailtrap.io",
		Port:     587,
		Username: "smtp@mailtrap.io",
		Password: "61dc207f67686fdb2aadbe5bc179fa71",
		From:     "no-reply@collectwave.com",
	}
	return &emailConfig
}

func GetMSAzureConfig() *services.MSAzureConfig {
	AZURE_TENANT_ID := "1f7b2203-1b88-4587-a5e9-3d5d0bf4136f"
	AZURE_CLIENT_ID := "179a2979-0ea8-44c1-b194-c0a3e053911c"
	AZURE_CLIENT_SECRET := "6Mr8Q~M1LMKf_36GFhBmhppkyRuBHpsDzBFTYbD-"
	config := services.MSAzureConfig{
		TenantID:     AZURE_TENANT_ID,
		ClientID:     AZURE_CLIENT_ID,
		ClientSecret: AZURE_CLIENT_SECRET,
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
	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		return nil
	}
	return client
}
