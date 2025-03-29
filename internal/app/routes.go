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
		organization.GET("/committee/:type/:id", organizationHandler.GetCommitteeByID)
		organization.PUT("/committee", organizationHandler.UpdateCommittee)
		organization.DELETE("/committee/:type/:id", organizationHandler.DeleteCommittee)
	}

	// documents Route
	document := api.Group("documents")
	document.Use(middleware.AuthMiddleware())
	{
		document.POST("/", documentHandler.CreateDocument)
		document.GET("/:id", documentHandler.GetDocumentByID)
		document.POST("/reference/:reference", documentHandler.GetDocumentByReference)
		document.GET("/title/:title", documentHandler.GetDocumentByTitle)
		document.PUT("/", documentHandler.UpdateDocument)
		document.PUT("/partial", documentHandler.UpdateDocumentPartial)
		document.PUT("/fileUrl/:id", documentHandler.UpdateFileURL)
		document.DELETE("/:id", documentHandler.DeleteDocument)
		document.GET("/list", documentHandler.ListDocuments)
		document.GET("/search", documentHandler.SearchDocuments)
		document.GET("/date", documentHandler.GetDocumentsByDateRange)
		document.GET("/count", documentHandler.CountDocuments)
	}

	// Project Route
	projects := api.Group("projects")
	projects.Use(middleware.AuthMiddleware())
	{
		projects.POST("/create", projectHandler.CreateProject)
	}

	// proposal Route
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

	// acceptance Route
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

	return router, nil
}
