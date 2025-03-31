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
