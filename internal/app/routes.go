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
	commentHandler := handlers.NewCommentHandler(services.CommentService)
	publicCommentHandler := handlers.NewNationalConsultationHandler(services.NationalConsultationService)
	voteHandler := handlers.NewVoteHandler(services.BallotingService)
	ballotingHandler := handlers.NewBallotingHandler(services.BallotingService)
	meetingHandler := handlers.NewMeetingHandler(*services.MeetingService)
	libraryHandler := handlers.NewLibraryHandler(*services.LibraryService, *services.MemberService)
	standardHandler := handlers.NewStandardHandler(services.StandardService)
	rbacHandler := handlers.NewRbacHandler(services.RbacService)
	notificationHandler := handlers.NewNotificationHandler(*services.NotificationService)
	reportsHandler := handlers.NewReportsHandler(services.ReportsService)

	api := router.Group("/api")

	// Serve static assets
	api.Static("/assets", "../assets")

	// Health check route
	api.GET("/health", handlers.HealthCheckHandler)

	// Auth routes
	auth := api.Group("auth")
	{
		auth.POST("/register", authHandler.RegisterMember)
		auth.POST("/login", authHandler.Login)
		auth.GET("/logout/:userId", authHandler.LogoutAll)

		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/account", authHandler.Account)
			protected.GET("/account/:id", authHandler.GetUserDetails)
			protected.GET("/user/:id", authHandler.GetUserDetails)
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

	// RBAC Routes
	rbac := api.Group("rbac")
	rbac.Use(middleware.AuthMiddleware(), middleware.DynamicAuthorize(services.PermissionResourceService))
	{
		rbac.POST("/roles", rbacHandler.CreateRole)
		rbac.GET("/roles", rbacHandler.ListRoles)
		rbac.GET("/roles/:id", rbacHandler.GetRoleByID)
		rbac.DELETE("/roles/:id", rbacHandler.DeleteRole)
		rbac.PUT("/roles", rbacHandler.UpdateRole)
		rbac.POST("/permissions", rbacHandler.CreatePermission)
		rbac.GET("/permissions", rbacHandler.ListPermissions)
		rbac.DELETE("/permissions/:id", rbacHandler.DeletePermission)
		rbac.POST("/assign/:member_id/:role_id", rbacHandler.AssignRoleToMember)
		rbac.GET("/members/:member_id/roles", rbacHandler.ListMemberRoles)
		rbac.POST("/roles/permission/:role_id/:permission_id", rbacHandler.AddPermissionToRole)
	}

	// Organization Route
	organization := api.Group("organization")
	organization.Use(middleware.AuthMiddleware(), middleware.DynamicAuthorize(services.PermissionResourceService))
	{
		organization.POST("/member_states", organizationHandler.CreateMemberState)
		organization.GET("/member_states", organizationHandler.FetchMemberStates)
		organization.PATCH("/member_states/:id", organizationHandler.UpdateMemberState)
		organization.DELETE("/member_states/:id", organizationHandler.DeleteMemberState)
		organization.POST("/nsbs", organizationHandler.CreateNSB)
		organization.PATCH("/nsbs", organizationHandler.UpdateNSB)
		organization.POST("/nsbs/secretary", organizationHandler.UpdateNationalTCSecretary)
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

		// editing committee
		organization.POST("/add_editing_committee_to_committee", organizationHandler.CreateEditingCommittee)
		organization.POST("/editing_committees", organizationHandler.CreateEditingCommittee)
		organization.GET("/editing_committees/tc/:id", organizationHandler.GetCommitteeEditingCommittee)
		organization.GET("/editing_committees/:id", organizationHandler.GetEditingCommitteeByID)

		organization.POST("/task_force", organizationHandler.CreateTaskForce)
		organization.POST("/task_force/:id", organizationHandler.GetTaskForceByID)
		organization.POST("/sub_commitee", organizationHandler.CreateSubCommittee)
		organization.POST("/add_member_to_sub_commitee", organizationHandler.AddMemberToSubCommittee)
		organization.POST("/specialized_committee", organizationHandler.CreateSpecializedCommittee)
		organization.GET("/specialized_committee/type/:id", organizationHandler.GetSpecializedCommitteeByType)

		// Update Secretary / Chairperson
		organization.PUT("/committees/:type/:id/secretary", organizationHandler.UpdateCommitteeSecretary)
		organization.PUT("/committees/:type/:id/chairperson", organizationHandler.UpdateCommitteeChairperson)

		// Add members
		organization.POST("/committees/arso_council/:id/members", organizationHandler.AddMemberToARSOCouncil)
		organization.POST("/committees/joint_advisory_group/:id/rec_members", organizationHandler.AddRegionalEconomicCommunityToJointAdvisoryGroup)
		organization.POST("/committees/joint_advisory_group/:id/observers", organizationHandler.AddObserverMemberToJointAdvisoryGroup)
		organization.POST("/committees/standards_management/:id/representatives", organizationHandler.AddRegionalRepresentativeToStandardsManagementCommittee)
		organization.POST("/committees/standards_management/:id/elected_members", organizationHandler.AddElectedMemberToStandardsManagementCommittee)
		organization.POST("/committees/standards_management/:id/observers", organizationHandler.AddObserverToStandardsManagementCommittee)
		organization.POST("/committees/technical/:id/members", organizationHandler.AddMemberToTechnicalCommittee)
		organization.POST("/committees/joint_technical/:id/members", organizationHandler.AddMemberToJointTechnicalCommittee)
		organization.POST("/committees/specialized/:id/members", organizationHandler.AddMemberToSpecializedCommittee)
		organization.POST("/committees/task_force/:id/members", organizationHandler.AddMemberToTaskForce)
		organization.POST("/committees/working_group/:id/members", organizationHandler.AddMemberToWorkingGroup)
		organization.POST("/committees/editing_committee/:id/members", organizationHandler.AddMemberToEditingCommittee)

		// Remove members
		organization.DELETE("/committees/arso_council/:id/members", organizationHandler.RemoveMemberFromARSOCouncil)
		organization.DELETE("/committees/joint_advisory_group/:id/rec_members", organizationHandler.RemoveRECFromJointAdvisoryGroup)
		organization.DELETE("/committees/joint_advisory_group/:id/observers", organizationHandler.RemoveObserverFromJointAdvisoryGroup)
		organization.DELETE("/committees/standards_management/:id/representatives", organizationHandler.RemoveRegionalRepresentativeFromStandardsManagementCommittee)
		organization.DELETE("/committees/standards_management/:id/elected_members", organizationHandler.RemoveRegionalElectedMemberFromStandardsManagementCommittee)
		organization.DELETE("/committees/technical/:id/members", organizationHandler.RemoveMemberFromTechnicalCommittee)
		organization.DELETE("/committees/specialized/:id/members", organizationHandler.RemoveMemberFromSpecializedCommittee)
		organization.DELETE("/committees/joint_technical/:id/members", organizationHandler.RemoveMemberFromJointTechnicalCommittee)
		organization.DELETE("/committees/editing_committee/:id/members", organizationHandler.RemoveMemberFromEditingCommittee)

		// Get members
		organization.GET("/committees/arso_council/:id/members", organizationHandler.GetArsoCouncilMembers)
		organization.GET("/committees/joint_advisory_group/:id/members", organizationHandler.GetJointAdvisoryGroupMembers)
		organization.GET("/committees/standards_management/:id/members", organizationHandler.GetStandardsManagementCommitteeMembers)
		organization.GET("/committees/technical/:id/members", organizationHandler.GetTechnicalCommitteeMembers)
		organization.GET("/committees/specialized/:id/members", organizationHandler.GetSpecializedCommitteeMembers)
		organization.GET("/committees/joint_technical/:id/members", organizationHandler.GetJointTechnicalCommitteeMembers)

		// Get all committees
		organization.GET("/committees/arso_council", organizationHandler.GetArsoCouncil)
		organization.GET("/committees/joint_advisory_group", organizationHandler.GetJointAdvisoryGroups)
		organization.GET("/committees/standards_management", organizationHandler.GetStandardsManagementCommittees)
		organization.GET("/committees/technical", organizationHandler.GetTechnicalCommittees)
		organization.GET("/committees/specialized", organizationHandler.GetSpecializedCommittees)
		organization.GET("/committees/joint_technical", organizationHandler.GetJointTechnicalCommittees)

		organization.POST("/technical_committees/participating_countries", organizationHandler.AddMemberStateToTCParticipatingCountries)
		organization.POST("/technical_committees/observer_countries", organizationHandler.AddMemberStateToTCObserverCountries)
		organization.POST("/technical_committees/equivalent_committees", organizationHandler.AddTCToTCEquivalentCommittees)
		organization.GET("/technical_committees/participating_countries/:id", organizationHandler.GetTCParticipatingCountries)
		organization.GET("/technical_committees/observer_countries/:id", organizationHandler.GetTCObserverCountries)
		organization.GET("/technical_committees/equivalent_committees/:id", organizationHandler.GetTCEquivalentCommittees)
		organization.PATCH("/technical_committees/participating_countries/:id", organizationHandler.RemoveMemberStateFromTCParticipatingCountries)
		organization.PATCH("/technical_committees/observer_countries/:id", organizationHandler.RemoveMemberStateFromTCObserverCountries)
		organization.PATCH("/technical_committees/equivalent_committees/:id", organizationHandler.RemoveTCFromTCEquivalentCommittees)
		organization.GET("/technical_committees/projects/:id", organizationHandler.GetTCProjects)
		organization.GET("/committees/meetings/:id", organizationHandler.GetCommitteeMeetings) // id can be TC,SC or WG
		organization.GET("/technical_committees/working_groups/:id", organizationHandler.GetTCWorkingGroups)

	}

	// documents Route
	document := api.Group("documents")
	document.Use(middleware.AuthMiddleware(), middleware.DynamicAuthorize(services.PermissionResourceService))
	{
		document.POST("/", documentHandler.CreateDocument)
		document.POST("/upload", documentHandler.UploadDocument)
		document.POST("/standards", documentHandler.UploadStandard)
		document.POST("/project", documentHandler.UploadRelatedDocument)
		document.POST("/minutes", documentHandler.UpdateMeetingMinutes)
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
		// Sharepoint documents
		document.GET("/sharepoint", documentHandler.ListSharepointDocuments)
		document.GET("/sharepoint/:id", documentHandler.GetSharepointDocument)
		document.POST("/sharepoint/copy", documentHandler.CopySharepointDocument)
		document.POST("/sharepoint/invite", documentHandler.InviteExternalUsersToDocument)
	}

	// Project Route
	projects := api.Group("projects")
	projects.Use(middleware.AuthMiddleware(), middleware.DynamicAuthorize(services.PermissionResourceService))
	{
		// Basic CRUD operations
		projects.POST("/", projectHandler.CreateProject)
		projects.POST("/approve", projectHandler.ApproveProject)
		projects.POST("/wd/review", projectHandler.ReviewWD)
		projects.POST("/cd/review", projectHandler.ReviewCD)
		projects.POST("/dars/review", projectHandler.ReviewDARS)
		projects.POST("/fdars/approve", projectHandler.ApproveFDRSForPublication)
		projects.GET("/:id", projectHandler.GetProjectByID)
		projects.PUT("/:id", projectHandler.UpdateProject)
		projects.DELETE("/:id", projectHandler.DeleteProject)

		// Project listings and searches
		projects.GET("/", projectHandler.FindProjects)
		projects.GET("/requests", projectHandler.FindProjectRequests)
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

		// Dashboard and statistics
		projects.GET("/statistics", projectHandler.GetDashboardStats)
		projects.GET("/distributions", projectHandler.GetAllDistributions)
	}

	// Proposal Route
	proposal := api.Group("proposals")
	proposal.Use(middleware.AuthMiddleware(), middleware.DynamicAuthorize(services.PermissionResourceService))
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
		proposal.POST("/approve", projectHandler.ApproveProjectProposal)
	}

	// Acceptance Route
	acceptance := api.Group("acceptance")
	acceptance.Use(middleware.AuthMiddleware(), middleware.DynamicAuthorize(services.PermissionResourceService))
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
	comment.Use(middleware.AuthMiddleware(), middleware.DynamicAuthorize(services.PermissionResourceService))
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

	// Balloting and  Observations Route
	balloting := api.Group("balloting")
	balloting.Use(middleware.AuthMiddleware(), middleware.DynamicAuthorize(services.PermissionResourceService))
	{
		balloting.POST("/votes", voteHandler.CreateVote)
		balloting.GET("/votes/:id", voteHandler.GetVoteByID)
		balloting.GET("/votes/ballot/:id", voteHandler.GetVotesByBallotingID)
		balloting.GET("/votes/project/:id", voteHandler.GetVotesByProjectID)
		balloting.PUT("/votes", voteHandler.UpdateVote)
		balloting.DELETE("/votes/:id", voteHandler.DeleteVote)
		balloting.GET("/votes/all", voteHandler.GetAllVotesWithAssociations)
		balloting.GET("/votes/count", voteHandler.CountVotesByBalloting)
		balloting.GET("/votes/criteria/:id", voteHandler.CheckProjectAcceptanceCriteria)

		// Decision of Balloting
		balloting.POST("/", ballotingHandler.CreateBalloting)
		balloting.GET("/list", ballotingHandler.GetAllBallotings)
		balloting.GET("/:id", ballotingHandler.GetBallotingByID)
		balloting.POST("/recommendation", ballotingHandler.RecommendFDARS)
		balloting.GET("/recommendation/verify/:project_id", ballotingHandler.VerifyFDARSRecommendation)
		balloting.POST("/approve", projectHandler.ApproveFDARS)
	}

	// Meeting Route
	meeting := api.Group("meetings")
	meeting.Use(middleware.AuthMiddleware(), middleware.DynamicAuthorize(services.PermissionResourceService))
	{
		meeting.POST("/create", meetingHandler.CreateMeeting)
		meeting.GET("/:id", meetingHandler.GetMeetingByID)
		meeting.GET("/list", meetingHandler.GetAllMeetings)
		meeting.GET("/committee/:committee_id", meetingHandler.GetMeetingsByCommittee)
		meeting.GET("/project/:project_id", meetingHandler.GetMeetingsByProject)
		meeting.GET("/upcoming", meetingHandler.GetUpcomingMeetings) // Can also be used in dashboard
		meeting.PUT("/:id", meetingHandler.UpdateMeeting)
		meeting.DELETE("/:id", meetingHandler.DeleteMeeting)
		meeting.PATCH("/status/:id", meetingHandler.UpdateMeetingStatus)
		meeting.POST("/attendees/:id", meetingHandler.AddAttendeeToMeeting)
		meeting.DELETE("/attendees/:meeting_id/:member_id", meetingHandler.RemoveAttendeeFromMeeting)
		meeting.POST("/documents/:meeting_id", meetingHandler.AddRelatedDocumentToMeeting)
		meeting.GET("/check-quorum/:meeting_id", meetingHandler.CheckQuorum)
	}

	// library Route
	library := api.Group("library")
	{
		library.GET("/top_standards", libraryHandler.GetTopStandards)
		library.GET("/latest_standards", libraryHandler.GetLatestStandards)
		library.GET("/top_committee", libraryHandler.GetTopCommittees)
		library.POST("/register", libraryHandler.RegisterMember)
		library.POST("/login", libraryHandler.Login)
		library.GET("/logout/:userId", authHandler.LogoutAll)
		library.GET("/account", middleware.AuthMiddleware(), authHandler.Account)
		library.GET("/standards", libraryHandler.FindStandards)
		library.GET("/standards/:id", middleware.AuthMiddleware(), libraryHandler.GetStandardByID)
		library.GET("/standards/preview/:id", middleware.AuthMiddleware(), libraryHandler.GetPreviewStandard)
		library.GET("/standards/reference/:reference", libraryHandler.GetStandardByReference)
		library.GET("/standards/search", libraryHandler.SearchStandards)
		library.GET("/standards/date-range", libraryHandler.GetStandardsByDateRange)
		library.GET("/standards/count", libraryHandler.CountStandards)
		library.GET("/committees", libraryHandler.ListCommittees)
		library.GET("/committees/:id", libraryHandler.GetCommitteeByID)
		library.GET("/committees/code/:code", libraryHandler.GetCommitteeByCode)
		library.GET("/committees/search", libraryHandler.SearchCommittees)
		library.GET("/committees/count", libraryHandler.CountCommittees)
		library.GET("/standards/committee/:id", libraryHandler.GetStandardsByCommittee)
		library.GET("/sectors", libraryHandler.GetSectors)
	}

	// Standard Development Route
	standard := api.Group("standards")
	{
		standard.POST("/", standardHandler.CreateStandard)
		standard.PUT("/:id/save", standardHandler.SaveStandard) // Auto-save / webhook-style
		standard.GET("/:id", standardHandler.GetStandard)
		standard.GET("/:id/versions", standardHandler.GetStandardVersions)
		standard.POST("/:id/restore", standardHandler.RestoreVersion)
		standard.GET("/:id/diff", standardHandler.DiffVersions)
		standard.GET("/:id/audit-log", standardHandler.GetAuditLogs)
	}

	// Notification Route
	notifications := api.Group("notifications")
	notifications.Use(middleware.AuthMiddleware(), middleware.DynamicAuthorize(services.PermissionResourceService))
	{
		// Dashboard and overview
		notifications.GET("/dashboard", notificationHandler.GetNotificationDashboard)
		notifications.GET("/", notificationHandler.GetNotifications)
		notifications.GET("/search", notificationHandler.SearchNotifications)
		notifications.GET("/history", notificationHandler.GetNotificationHistory)

		// Individual notification actions
		notifications.PUT("/:id/read", notificationHandler.MarkNotificationAsRead)
		notifications.PUT("/read-all", notificationHandler.MarkAllNotificationsAsRead)
		notifications.DELETE("/:id", notificationHandler.DeleteNotification)

		// Preferences
		notifications.GET("/preferences", notificationHandler.GetNotificationPreferences)
		notifications.PUT("/preferences", notificationHandler.UpdateNotificationPreferences)

		// Admin announcements (admin only)
		notifications.POST("/announcements", notificationHandler.SendAdminAnnouncement)
	}

	// Reports API
	reports := api.Group("/reports")
	reports.Use(middleware.AuthMiddleware())
	reports.Use(middleware.DynamicAuthorize(services.PermissionResourceService))
	{
		// Report generation
		reports.POST("/generate", reportsHandler.GenerateReport)
		reports.GET("/", reportsHandler.ListReports)
		reports.GET("/:id", reportsHandler.GetReport)
		reports.DELETE("/:id", reportsHandler.DeleteReport)

		// Report templates
		templates := reports.Group("/templates")
		{
			templates.POST("/", reportsHandler.CreateReportTemplate)
			templates.GET("/", reportsHandler.ListReportTemplates)
			templates.DELETE("/:id", reportsHandler.DeleteReportTemplate)
		}

		// Report schedules
		schedules := reports.Group("/schedules")
		{
			schedules.POST("/", reportsHandler.CreateReportSchedule)
			schedules.GET("/", reportsHandler.ListReportSchedules)
			schedules.DELETE("/:id", reportsHandler.DeleteReportSchedule)
		}
	}

	return router, nil
}
