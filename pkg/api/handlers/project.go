package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	projectService services.ProjectService
}

func NewProjectHandler(projectService services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var payload models.Project
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			// Convert validation errors into human-readable messages
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.Show(c, http.StatusBadRequest, "errors", formattedErrors)
			return
		}

		// For non-validation errors
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, ok := userID.(*uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}
	payload.MemberID = userUUID

	err := h.projectService.CreateProject(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "Project added successfully")
}

func (h *ProjectHandler) GetNextAvailableNumber(c *gin.Context) {
	number, err := h.projectService.GetNextAvailableNumber()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"number":   number,
		"previous": number - 1,
	})
}

// GetProjectByID handles retrieving a project by its ID
func (h *ProjectHandler) GetProjectByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.projectService.GetProjectByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Project not found")
		return
	}

	utilities.Show(c, http.StatusOK, "project", project)
}

// UpdateProject handles updating an existing project
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var payload models.Project
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.Show(c, http.StatusBadRequest, "errors", formattedErrors)
			return
		}

		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Ensure the ID in the URL matches the payload
	payload.ID = id

	err = h.projectService.UpdateProject(&payload)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Project updated successfully")
}

// DeleteProject handles deleting a project
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	err = h.projectService.DeleteProject(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Project deleted successfully")
}

// FindProjects handles searching for projects with various filters
func (h *ProjectHandler) FindProjects(c *gin.Context) {
	// Parse query parameters
	params := make(map[string]interface{})

	// String parameters
	if title := c.Query("title"); title != "" {
		params["title"] = title
	}

	if typeStr := c.Query("type"); typeStr != "" {
		params["type"] = models.ProjectType(typeStr)
	}

	// UUID parameters
	if committeeID := c.Query("committee_id"); committeeID != "" {
		if id, err := uuid.Parse(committeeID); err == nil {
			params["committee_id"] = id
		}
	}

	if workingGroupID := c.Query("working_group_id"); workingGroupID != "" {
		if id, err := uuid.Parse(workingGroupID); err == nil {
			params["working_group_id"] = id
		}
	}

	// Boolean parameters
	if visibleStr := c.Query("visible_on_library"); visibleStr != "" {
		if visible, err := strconv.ParseBool(visibleStr); err == nil {
			params["visible_on_library"] = visible
		}
	}

	if emergencyStr := c.Query("is_emergency"); emergencyStr != "" {
		if emergency, err := strconv.ParseBool(emergencyStr); err == nil {
			params["is_emergency"] = emergency
		}
	}

	// Pagination
	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
			offset = val
		}
	}

	projects, total, err := h.projectService.FindProjects(params, limit, offset)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

// GetProjectsByTimeframe handles retrieving projects within a given timeframe
func (h *ProjectHandler) GetProjectsByTimeframe(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid start date format. Use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid end date format. Use YYYY-MM-DD")
		return
	}

	projects, err := h.projectService.GetProjectsByTimeframe(startDate, endDate)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "projects", projects)
}

// GetProjectCountByType handles retrieving project counts grouped by type
func (h *ProjectHandler) GetProjectCountByType(c *gin.Context) {
	counts, err := h.projectService.GetProjectCountByType()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "counts", counts)
}

// GetProjectsWithStageTransitions handles retrieving projects that have transitioned between stages
func (h *ProjectHandler) GetProjectsWithStageTransitions(c *gin.Context) {
	fromStageID, err := uuid.Parse(c.Query("from_stage_id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid from_stage_id")
		return
	}

	toStageID, err := uuid.Parse(c.Query("to_stage_id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid to_stage_id")
		return
	}

	projects, err := h.projectService.GetProjectsWithStageTransitions(fromStageID, toStageID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "projects", projects)
}

// GetProjectsByReferenceBase handles retrieving projects with the same reference base
func (h *ProjectHandler) GetProjectsByReferenceBase(c *gin.Context) {
	referenceBase := c.Query("reference_base")
	if referenceBase == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "reference_base is required")
		return
	}

	projects, err := h.projectService.GetProjectsByReferenceBase(referenceBase)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "projects", projects)
}

// CreateProjectRevision handles creating a new revision of an existing project
func (h *ProjectHandler) CreateProjectRevision(c *gin.Context) {
	baseProjectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	revision, err := h.projectService.CreateProjectRevision(baseProjectID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusCreated, "revision", revision)
}

// GetProjectsApproachingDeadline handles retrieving projects approaching their deadline
func (h *ProjectHandler) GetProjectsApproachingDeadline(c *gin.Context) {
	daysThresholdStr := c.DefaultQuery("days_threshold", "30")
	daysThreshold, err := strconv.Atoi(daysThresholdStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid days_threshold")
		return
	}

	projects, err := h.projectService.GetProjectsApproachingDeadline(daysThreshold)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "projects", projects)
}

// GetProjectsInStageForTooLong handles retrieving projects that have been in a stage for too long
func (h *ProjectHandler) GetProjectsInStageForTooLong(c *gin.Context) {
	stageID, err := uuid.Parse(c.Query("stage_id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid stage_id")
		return
	}

	dayThresholdStr := c.DefaultQuery("day_threshold", "30")
	dayThreshold, err := strconv.Atoi(dayThresholdStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid day_threshold")
		return
	}

	projects, err := h.projectService.GetProjectsInStageForTooLong(stageID, dayThreshold)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "projects", projects)
}

// GetRelatedProjects handles retrieving projects related to a given project
func (h *ProjectHandler) GetRelatedProjects(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	projects, err := h.projectService.GetRelatedProjects(projectID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "related_projects", projects)
}

// UpdateProjectStage handles updating a project's stage
func (h *ProjectHandler) UpdateProjectStage(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var payload struct {
		StageID uuid.UUID `json:"stage_id" binding:"required"`
		Notes   string    `json:"notes"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.projectService.UpdateProjectStage(projectID, payload.StageID, payload.Notes)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Project stage updated successfully")
}

// GetProjectWithStageHistory handles retrieving a project with its stage history
func (h *ProjectHandler) GetProjectWithStageHistory(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.projectService.GetProjectWithStageHistory(projectID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Project not found")
		return
	}

	utilities.Show(c, http.StatusOK, "project", project)
}

// GetProjectStageHistory handles retrieving a project's stage history
func (h *ProjectHandler) GetProjectStageHistory(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	history, err := h.projectService.GetProjectStageHistory(projectID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "stage_history", history)
}
