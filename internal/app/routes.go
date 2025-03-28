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
	projectHandler := handlers.NewProjectHandler(*services.ProjectService)

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

	// Project Route
	projects := api.Group("projects")
	projects.Use(middleware.AuthMiddleware())
	{
		projects.POST("/create", projectHandler.CreateProject)
	}

	return router, nil
}
