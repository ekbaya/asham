package handlers

import (
	"net/http"
	"strconv"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LibraryHandler struct {
	libraryService services.LibraryService
}

func NewLibraryHandler(LibraryService services.LibraryService) *LibraryHandler {
	return &LibraryHandler{
		libraryService: LibraryService,
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

	if limitStr := c.Query("pageSize"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	if offsetStr := c.Query("page"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
			offset = val
		}
	}

	projects, err := h.libraryService.FindStandards(params, limit, offset)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standards": projects,
		"total":     len(projects),
		"limit":     limit,
		"page":      offset,
	})
}
