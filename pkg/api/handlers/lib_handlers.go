package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LibraryApiHandler handles HTTP requests for library-related operations
type LibraryApiHandler struct {
	service *services.LibraryService
}

// NewLibraryApiHandler creates a new LibraryApiHandler
func NewLibraryApiHandler(service *services.LibraryService) *LibraryApiHandler {
	return &LibraryApiHandler{service: service}
}

// ListStandards lists all standards with pagination
func (h *LibraryApiHandler) ListStandards(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	projects, total, err := h.service.ListProjects(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standards": projects,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

// GetStandardByID retrieves a standard by its ID
func (h *LibraryApiHandler) GetStandardByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	project, err := h.service.GetProjectByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// GetStandardByReference retrieves a standard by its reference
func (h *LibraryApiHandler) GetStandardByReference(c *gin.Context) {
	reference := c.Param("reference")
	if reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reference is required"})
		return
	}

	project, err := h.service.GetProjectByReference(reference)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// SearchStandards searches standards by query
func (h *LibraryApiHandler) SearchStandards(c *gin.Context) {
	query := c.Query("query")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	projects, total, err := h.service.SearchProjects(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standards": projects,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

// GetStandardsByDateRange retrieves standards created within a date range
func (h *LibraryApiHandler) GetStandardsByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use YYYY-MM-DD"})
		return
	}

	projects, err := h.service.GetProjectsCreatedBetween(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standards": projects,
		"total":     len(projects),
	})
}

// CountStandards returns the total number of standards
func (h *LibraryApiHandler) CountStandards(c *gin.Context) {
	count, err := h.service.CountProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// ListCommittees lists all committees with pagination
func (h *LibraryApiHandler) ListCommittees(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	committees, total, err := h.service.ListCommittees(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"committees": committees,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
	})
}

// GetCommitteeByID retrieves a committee by its ID
func (h *LibraryApiHandler) GetCommitteeByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	committee, err := h.service.GetCommitteeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, committee)
}

// GetCommitteeByCode retrieves a committee by its code
func (h *LibraryApiHandler) GetCommitteeByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	committee, err := h.service.GetCommitteeByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, committee)
}

// SearchCommittees searches committees by query
func (h *LibraryApiHandler) SearchCommittees(c *gin.Context) {
	query := c.Query("query")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	committees, total, err := h.service.SearchCommittees(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"committees": committees,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
	})
}

// CountCommittees returns the total number of committees
func (h *LibraryApiHandler) CountCommittees(c *gin.Context) {
	count, err := h.service.CountCommittees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetStandardsByCommittee retrieves standards associated with a committee
func (h *LibraryApiHandler) GetStandardsByCommittee(c *gin.Context) {
	committeeID := c.Param("id")
	if _, err := uuid.Parse(committeeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid committee ID format"})
		return
	}

	projects, err := h.service.GetProjectsByCommittee(committeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standards": projects,
		"total":     len(projects),
	})
}
