package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"

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
func (h *LibraryHandler) GetTopStandards(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "3"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 {
		log.Printf("Invalid limit parameter: %d", limit)
		utilities.ShowMessage(c, http.StatusBadRequest, "limit must be a positive integer")
		return
	}

	standards, total, err := h.libraryService.GetTopStandards(limit, offset)
	if err != nil {
		log.Printf("Error fetching top standards: %v", err)
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standards": standards,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *LibraryHandler) GetLatestStandards(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "3"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 {
		log.Printf("Invalid limit parameter: %d", limit)
		utilities.ShowMessage(c, http.StatusBadRequest, "limit must be a positive integer")
		return
	}

	standards, total, err := h.libraryService.GetLatestStandards(limit, offset)
	if err != nil {
		log.Printf("Error fetching latest standards: %v", err)
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standards": standards,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *LibraryHandler) GetTopCommittees(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "3"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 {
		log.Printf("Invalid limit parameter: %d", limit)
		utilities.ShowMessage(c, http.StatusBadRequest, "limit must be a positive integer")
		return
	}

	committees, total, err := h.libraryService.GetTopCommittees(limit, offset)
	if err != nil {
		log.Printf("Error fetching top committees: %v", err)
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
func (h *LibraryHandler) RegisterMember(c *gin.Context) {
	// Read and log request body for debugging
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		utilities.ShowMessage(c, http.StatusBadRequest, "Failed to read request body")
		return
	}
	log.Printf("Request body: %s", string(body))

	// Restore request body for ShouldBindJSON
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var payload models.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("Error binding JSON: %v", err)
		if err.Error() == "EOF" {
			utilities.ShowMessage(c, http.StatusBadRequest, "Empty or invalid request body")
			return
		}
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.ShowError(c, http.StatusBadRequest, formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.libraryService.RegisterMember(&payload)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusCreated, "User registered successfully")
}

func (h *LibraryHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON for login: %v", err)
		if err.Error() == "EOF" {
			utilities.ShowMessage(c, http.StatusBadRequest, "Empty or invalid request body")
			return
		}
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			formattedErrors := utilities.FormatValidationErrors(validationErrors)
			utilities.ShowError(c, http.StatusBadRequest, formattedErrors)
			return
		}
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	token, refreshToken, err := h.libraryService.Login(req.Username, req.Password)
	if err != nil {
		log.Printf("Error logging in: %v", err)
		utilities.ShowMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  token,
		"refresh_token": refreshToken,
		"expires_in":    86400,
	})
}

func (h *LibraryHandler) FindStandards(c *gin.Context) {
	params := make(map[string]any)

	if sector := c.Query("sector"); sector != "" {
		params["sector"] = sector
	}

	if query := c.Query("query"); query != "" {
		params["query"] = query
	}

	if language := c.Query("language"); language != "" {
		params["language"] = language
	}

	if year := c.Query("year"); year != "" {
		params["year"] = year
	}

	if sortBy := c.Query("sortBy"); sortBy != "" {
		params["sortBy"] = sortBy
	}

	limit := 10
	offset := 0

	if limitStr := c.Query("pageSize"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	if offsetStr := c.Query("page"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val > 0 {
			offset = (val - 1) * limit
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

	committee := map[string]any{}
	if project.TechnicalCommittee != nil {
		committee = map[string]any{
			"id":   project.TechnicalCommitteeID,
			"code": project.TechnicalCommittee.Code,
			"name": project.TechnicalCommittee.Name,
		}
	}

	pageCount := 20
	if project.Standard != nil && project.Standard.FileURL != "" {
		calculatedPages, err := utilities.GetPDFPageCount(project.Standard.FileURL)
		if err == nil {
			pageCount = calculatedPages
		} else {
			log.Printf("Error calculating PDF pages for standard ID %v: %v", project.ID, err)
		}
	}

	utilities.Show(c, http.StatusOK, "standard", map[string]any{
		"id":             project.ID,
		"reference":      project.Reference,
		"title":          project.Title,
		"description":    project.Description,
		"language":       project.Language,
		"version":        "1.0",
		"committee":      committee,
		"published_date": project.PublishedDate,
		"sector":         "Manufacturing",
		"file_url":       project.Standard.FileURL,
		"pages":          pageCount,
	})
}

func (h *LibraryHandler) GetPreviewStandard(c *gin.Context) {
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

	// Check if Standard is nil
	if project.Standard == nil {
		utilities.ShowMessage(c, http.StatusNotFound, "standard not found")
		return
	}

	// Check if FileURL is empty
	if project.Standard.FileURL == "" {
		utilities.ShowMessage(c, http.StatusNotFound, "no file available for download")
		return
	}

	path := project.Standard.FileURL
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "failed to open file")
		return
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "failed to get file info")
		return
	}

	// Set headers and send file
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name()))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.File(path)
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

func (h *LibraryHandler) ListCommittees(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	query := c.DefaultQuery("query", "")

	committees, total, err := h.libraryService.ListCommittees(limit, offset, query)
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
func (h *LibraryHandler) GetSectors(c *gin.Context) {
	sectors, err := h.libraryService.GetSectors()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sectors": sectors,
		"total":   len(sectors),
	})
}
