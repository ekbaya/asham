package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LibraryHandler struct {
	libraryService services.LibraryService
}

func NewLibraryHandler(libraryService services.LibraryService) *LibraryHandler {
	return &LibraryHandler{
		libraryService: libraryService,
	}
}

func (h *LibraryHandler) FindStandards(c *gin.Context) {
	// Parse query parameters
	params := make(map[string]any)

	if sector := c.Query("sector"); sector != "" {
		params["sector"] = sector
	}

	if title := c.Query("title"); title != "" {
		params["title"] = title
	}

	if typeStr := c.Query("type"); typeStr != "" {
		params["type"] = models.ProjectType(typeStr)
	}

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

	if limitStr := c.Query("pageSize"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	if offsetStr := c.Query("page"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val > 0 {
			offset = (val - 1) * limit // Convert page to offset
		}
	}

	projects, total, err := h.libraryService.FindStandards(params, limit, offset)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standards": projects,
		"total":     total,
		"limit":     limit,
		"page":      offset/limit + 1,
	})
}

func (h *LibraryHandler) GetStandardByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "invalid ID format")
		return
	}

	project, err := h.libraryService.GetProjectByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *LibraryHandler) GetStandardByReference(c *gin.Context) {
	reference := c.Param("reference")
	if reference == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "reference is required")
		return
	}

	project, err := h.libraryService.GetProjectByReference(reference)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *LibraryHandler) SearchStandards(c *gin.Context) {
	query := c.Query("query")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	projects, total, err := h.libraryService.SearchProjects(query, limit, offset)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standards": projects,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *LibraryHandler) GetStandardsByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "invalid start_date format, use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "invalid end_date format, use YYYY-MM-DD")
		return
	}

	projects, err := h.libraryService.GetProjectsCreatedBetween(startDate, endDate)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standards": projects,
		"total":     len(projects),
	})
}

func (h *LibraryHandler) CountStandards(c *gin.Context) {
	count, err := h.libraryService.CountProjects()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *LibraryHandler) ListCommittees(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	committees, total, err := h.libraryService.ListCommittees(limit, offset)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"committees": committees,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
	})
}

func (h *LibraryHandler) GetCommitteeByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "invalid ID format")
		return
	}

	committee, err := h.libraryService.GetCommitteeByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, committee)
}

func (h *LibraryHandler) GetCommitteeByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "code is required")
		return
	}

	committee, err := h.libraryService.GetCommitteeByCode(code)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, committee)
}

func (h *LibraryHandler) SearchCommittees(c *gin.Context) {
	query := c.Query("query")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	committees, total, err := h.libraryService.SearchCommittees(query, limit, offset)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"committees": committees,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
	})
}

func (h *LibraryHandler) CountCommittees(c *gin.Context) {
	count, err := h.libraryService.CountCommittees()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *LibraryHandler) GetStandardsByCommittee(c *gin.Context) {
	committeeID := c.Param("id")
	if _, err := uuid.Parse(committeeID); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "invalid committee ID format")
		return
	}

	projects, err := h.libraryService.GetProjectsByCommittee(committeeID)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standards": projects,
		"total":     len(projects),
	})
}
