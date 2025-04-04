package app

import (
	"github.com/ekbaya/asham/pkg/api/handlers"
	middleware "github.com/ekbaya/asham/pkg/api/middlewares"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/gin-gonic/gin"
)

func InitRoutes(services *services.ServiceContainer) (*gin.Engine, error) {
	router := gin.Default()
	// Cors Middleware
	router.Use(middleware.CORSMiddleware())

	authHandler := handlers.NewUsersHandler(*services.MemberService)
	organizationHandler := handlers.NewOrganizationHandler(*services.OrganizationService)
	documentHandler := handlers.NewDocumentHandler(*services.DocumentService)
	projectHandler := handlers.NewProjectHandler(*services.ProjectService)
	proposalHandler := handlers.NewProposalHandler(services.ProposalService, services.DocumentService, services.ProjectService)
	acceptanceHandler := handlers.NewAcceptanceHandler(*services.AcceptanceService)
	commentHandler := handlers.NewCommentHandler(*&services.CommentService)
	publicCommentHandler := handlers.NewNationalConsultationHandler(*&services.NationalConsultationService)

	api := router.Group("/api")

	// Health check route
	api.GET("/health", handlers.HealthCheckHandler)

	// Auth routes
	auth := api.Group("auth")
	{
		auth.POST("/register", authHandler.RegisterMember)
		auth.POST("/login", authHandler.Login)

		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/account", authHandler.Account)
			protected.GET("/account/:id", authHandler.GetUserDetails)
			protected.GET("/user", authHandler.GetUserDetails)
			protected.PUT("/user", authHandler.UpdateUser)
			protected.GET("/users", authHandler.GetAllUsers)
			protected.DELETE("/users/:id", authHandler.DeleteMember)
		}

		refreshTokenGroup := auth.Group("/")
		refreshTokenGroup.Use(middleware.TokenMiddleware())
		{
			refreshTokenGroup.POST("/refresh-token", authHandler.GenerateRefreshToken)
		}
	}

	// Organization Route
	organization := api.Group("organization")
	organization.Use(middleware.AuthMiddleware())
	{
		organization.POST("/member_states", organizationHandler.CreateMemberState)
		organization.GET("/member_states", organizationHandler.FetchMemberStates)
		organization.POST("/nsbs", organizationHandler.CreateNSB)
		organization.GET("/nsbs", organizationHandler.FetchNSBs)
		organization.POST("/committee", organizationHandler.CreateCommittee)
		organization.GET("/technical_committees", organizationHandler.FetchTechnicalCommittees)
		organization.GET("/technical_committees/search", organizationHandler.SearchTechnicalCommittees)
		organization.GET("/committee/:type/:id", organizationHandler.GetCommitteeByID)
		organization.PUT("/committee", organizationHandler.UpdateCommittee)
		organization.DELETE("/committee/:type/:id", organizationHandler.DeleteCommittee)
		organization.POST("/add_working_group_to_committee", organizationHandler.CreateWorkingGroup)
		organization.POST("/working_groups", organizationHandler.CreateWorkingGroup)
		organization.POST("/complete_working_group", organizationHandler.CompleteWorkingGroup)
		organization.GET("/working_groups/tc/:id", organizationHandler.GetCommitteeWorkingGroups)
		organization.GET("/working_groups/:id", organizationHandler.GetWorkingGroupByID)
		organization.POST("/task_force", organizationHandler.CreateTaskForce)
		organization.POST("/task_force/:id", organizationHandler.GetTaskForceByID)
		organization.POST("/sub_commitee", organizationHandler.CreateSubCommittee)
		organization.POST("/add_member_to_sub_commitee", organizationHandler.AddMemberToSubCommittee)
		organization.POST("/specialized_committee", organizationHandler.CreateSpecializedCommittee)
		organization.GET("/specialized_committee/type/:id", organizationHandler.GetSpecializedCommitteeByType)
	}

	// documents Route
	document := api.Group("documents")
	document.Use(middleware.AuthMiddleware())
	{
		document.POST("/", documentHandler.CreateDocument)
		document.POST("/upload", documentHandler.UploadDocument)
		document.GET("/:id", documentHandler.GetDocumentByID)
		document.POST("/reference/:reference", documentHandler.GetDocumentByReference)
		document.GET("/title/:title", documentHandler.GetDocumentByTitle)
		document.PUT("/", documentHandler.UpdateDocument)
		document.PUT("/partial", documentHandler.UpdateDocumentPartial)
		document.PUT("/fileUrl/:id", documentHandler.UpdateFileURL)
		document.DELETE("/:id", documentHandler.DeleteDocument)
		document.GET("/list", documentHandler.ListDocuments)
		document.GET("/list/:projectId", documentHandler.ProjectDocuments)
		document.GET("/search", documentHandler.SearchDocuments)
		document.GET("/date", documentHandler.GetDocumentsByDateRange)
		document.GET("/count", documentHandler.CountDocuments)
	}

	// Project Route
	projects := api.Group("projects")
	projects.Use(middleware.AuthMiddleware())
	{
		// Basic CRUD operations
		projects.POST("/", projectHandler.CreateProject)
		projects.POST("/approve", projectHandler.ApproveProject)
		projects.POST("/wd/review", projectHandler.ReviewWD)
		projects.POST("/cd/review", projectHandler.ReviewCD)
		projects.GET("/:id", projectHandler.GetProjectByID)
		projects.PUT("/:id", projectHandler.UpdateProject)
		projects.DELETE("/:id", projectHandler.DeleteProject)

		// Project listings and searches
		projects.GET("/", projectHandler.FindProjects)
		projects.GET("/stages", projectHandler.FetchStages)
		projects.GET("/next-number", projectHandler.GetNextAvailableNumber)

		// Project stage management
		projects.GET("/:id/with-stage-history", projectHandler.GetProjectWithStageHistory)
		projects.GET("/:id/stage-history", projectHandler.GetProjectStageHistory)
		projects.PUT("/:id/stage", projectHandler.UpdateProjectStage)

		// Project analytics
		projects.GET("/by-timeframe", projectHandler.GetProjectsByTimeframe)
		projects.GET("/count-by-type", projectHandler.GetProjectCountByType)
		projects.GET("/stage-transitions", projectHandler.GetProjectsWithStageTransitions)
		projects.GET("/approaching-deadline", projectHandler.GetProjectsApproachingDeadline)
		projects.GET("/stage-delayed", projectHandler.GetProjectsInStageForTooLong)

		// Project relationships
		projects.GET("/by-reference-base", projectHandler.GetProjectsByReferenceBase)
		projects.GET("/:id/related", projectHandler.GetRelatedProjects)

		// Project versioning
		projects.POST("/:id/revision", projectHandler.CreateProjectRevision)
	}

	// Proposal Route
	proposal := api.Group("proposals")
	proposal.Use(middleware.AuthMiddleware())
	{
		proposal.POST("/", proposalHandler.CreateProposal)
		proposal.GET("/:id", proposalHandler.GetProposal)
		proposal.GET("/project/:projectId", proposalHandler.GetProposalByProject)
		proposal.GET("/list", proposalHandler.ListProposals)
		proposal.PUT("/", proposalHandler.UpdateProposal)
		proposal.PUT("/partial", proposalHandler.UpdatePartialProposal)
		proposal.DELETE("/:id", proposalHandler.DeleteProposal)
		proposal.GET("/search", proposalHandler.SearchProposals)
		proposal.GET("/creator/:memberId", proposalHandler.GetProposalsByCreator)
		proposal.POST("/reference/:id", proposalHandler.AddReferencedStandard)
		proposal.DELETE("/reference/:id/:documentId", proposalHandler.RemoveReferencedStandard)
		proposal.POST("/transfer/:id", proposalHandler.TransferProposal)
	}

	// Acceptance Route
	acceptance := api.Group("acceptance")
	acceptance.Use(middleware.AuthMiddleware())
	{
		acceptance.POST("/submission", acceptanceHandler.CreateNSBResponse)
		acceptance.GET("/submission/:id", acceptanceHandler.GetNSBResponse)
		acceptance.GET("/submission/project/:id", acceptanceHandler.GetNSBResponsesByProjectID)
		acceptance.PUT("/submission", acceptanceHandler.UpdateNSBResponse)
		acceptance.DELETE("/submission/:id", acceptanceHandler.DeleteNSBResponse)

		// Compilation of Submissions from various NSBs
		acceptance.GET("/submission/compilation/list", acceptanceHandler.GetAcceptances)
		acceptance.GET("/submission/compilation/:id", acceptanceHandler.GetAcceptance)
		acceptance.GET("/submission/compilation/project/:id", acceptanceHandler.GetAcceptanceByProjectID)
		acceptance.PUT("/submission/compilation", acceptanceHandler.UpdateAcceptance)
		acceptance.GET("/submission/compilation/:id/responses", acceptanceHandler.GetAcceptanceWithResponses)
		acceptance.GET("/submission/compilation/:id/count-by-type", acceptanceHandler.CountNSBResponsesByType)
		acceptance.GET("/submission/compilation/:id/calculate-stats", acceptanceHandler.CalculateNSBResponseStats)
		acceptance.POST("/submission/compilation/approve", acceptanceHandler.SetNSBResponseacceptanceApproval)
		acceptance.GET("/submission/compilation/:id/results", acceptanceHandler.GetAcceptanceResults)
	}

	// Comments and  Observations Route
	comment := api.Group("comments")
	comment.Use(middleware.AuthMiddleware())
	{
		comment.POST("/", commentHandler.CreateComment)
		comment.GET("/:id", commentHandler.GetCommentByID)
		comment.GET("/list", commentHandler.GetAllComments)
		comment.PUT("/:comment_id", commentHandler.UpdateComment)
		comment.DELETE("/:id", commentHandler.DeleteComment)
		comment.GET("/project/:id", commentHandler.GetCommentsByProjectID)

		// Public comments
		comment.POST("/public/", publicCommentHandler.CreateNationalConsultation)
		comment.GET("/public/:id", publicCommentHandler.GetNationalConsultationByID)
		comment.GET("/public/list", publicCommentHandler.GetAllNationalConsultations)
		comment.PUT("/public/:comment_id", publicCommentHandler.UpdateNationalConsultation)
		comment.DELETE("/public/:id", publicCommentHandler.DeleteNationalConsultation)
		comment.GET("/public/project/:id", publicCommentHandler.GetNationalConsultationsByProjectID)
	}

	return router, nil
}
