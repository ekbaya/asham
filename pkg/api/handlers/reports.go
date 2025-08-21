package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
)

type ReportsHandler struct {
	reportsService *services.ReportsService
}

func NewReportsHandler(reportsService *services.ReportsService) *ReportsHandler {
	return &ReportsHandler{
		reportsService: reportsService,
	}
}

// GenerateReportRequest represents the request payload for generating reports
type GenerateReportRequest struct {
	ReportType models.ReportType    `json:"report_type" binding:"required"`
	Filters    models.ReportFilters `json:"filters"`
	Format     models.ReportFormat  `json:"format,omitempty"`
}

// ReportResponse represents the response for report operations
type ReportResponse struct {
	ID          string                `json:"id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Type        models.ReportType     `json:"type"`
	Format      models.ReportFormat   `json:"format"`
	Status      models.ReportStatus   `json:"status"`
	Data        interface{}           `json:"data,omitempty"`
	Filters     models.ReportFilters  `json:"filters"`
	CreatedAt   time.Time             `json:"created_at"`
	CreatedBy   string                `json:"created_by"`
}

// DashboardResponse represents the dashboard metrics response
type DashboardResponse struct {
	Metrics   []models.DashboardMetric `json:"metrics"`
	Widgets   []models.DashboardWidget `json:"widgets"`
	UpdatedAt time.Time                `json:"updated_at"`
}

// @Summary Generate a new report
// @Description Generate a report based on specified type and filters
// @Tags reports
// @Accept json
// @Produce json
// @Param request body GenerateReportRequest true "Report generation request"
// @Success 200 {object} ReportResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/reports/generate [post]
func (h *ReportsHandler) GenerateReport(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utilities.ShowError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		utilities.ShowError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	var req GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utilities.ShowError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Set default format if not specified
	if req.Format == "" {
		req.Format = models.ReportFormatJSON
	}

	report, err := h.reportsService.GenerateReport(userID, req.ReportType, req.Filters)
	if err != nil {
		utilities.ShowError(c, http.StatusInternalServerError, "Failed to generate report: "+err.Error())
		return
	}

	response := &ReportResponse{
		ID:          report.ID.String(),
		Title:       report.Title,
		Description: report.Description,
		Type:        report.Type,
		Format:      report.Format,
		Status:      report.Status,
		Filters:     report.Filters,
		CreatedAt:   report.CreatedAt,
		CreatedBy:   report.CreatedByID,
	}

	// Include data if format is JSON
	if report.Format == models.ReportFormatJSON {
		var data interface{}
		if dataBytes, ok := report.Data.([]byte); ok {
			if err := json.Unmarshal(dataBytes, &data); err == nil {
				response.Data = data
			}
		}
	}

	utilities.Show(c, http.StatusOK, "Report generated successfully", response)
}

// @Summary Generate real-time report data
// @Description Generate real-time report data without saving to database
// @Tags reports
// @Accept json
// @Produce json
// @Param request body GenerateReportRequest true "Report generation request"
// @Success 200 {object} interface{}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/reports/realtime [post]
func (h *ReportsHandler) GenerateRealTimeReport(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utilities.ShowError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		utilities.ShowError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	var req GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utilities.ShowError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	data, err := h.reportsService.GenerateRealTimeReport(userID, req.ReportType, req.Filters)
	if err != nil {
		utilities.ShowError(c, http.StatusInternalServerError, "Failed to generate real-time report: "+err.Error())
		return
	}

	responseData := gin.H{
		"data":       data,
		"type":       req.ReportType,
		"filters":    req.Filters,
		"generated_at": time.Now(),
	}

	utilities.Show(c, http.StatusOK, "Real-time report generated successfully", responseData)
}

// @Summary Get dashboard metrics
// @Description Retrieve dashboard metrics and widgets for the authenticated user
// @Tags reports
// @Produce json
// @Success 200 {object} DashboardResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/reports/dashboard [get]
func (h *ReportsHandler) GetDashboard(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utilities.ShowError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		utilities.ShowError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	metrics, err := h.reportsService.GetDashboardMetrics(userID)
	if err != nil {
		utilities.ShowError(c, http.StatusInternalServerError, "Failed to get dashboard metrics: "+err.Error())
		return
	}

	// For now, return empty widgets array - this would be expanded later
	response := &DashboardResponse{
		Metrics:   metrics,
		Widgets:   []models.DashboardWidget{},
		UpdatedAt: time.Now(),
	}

	utilities.Show(c, http.StatusOK, "Dashboard metrics retrieved successfully", response)
}

// @Summary Get report by ID
// @Description Retrieve a specific report by its ID
// @Tags reports
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} ReportResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/reports/{id} [get]
func (h *ReportsHandler) GetReport(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utilities.ShowError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		utilities.ShowError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	reportID := c.Param("id")
	if reportID == "" {
		utilities.ShowError(c, http.StatusBadRequest, "Report ID is required")
		return
	}

	report, err := h.reportsService.GetReport(reportID, userID)
	if err != nil {
		utilities.ShowError(c, http.StatusInternalServerError, "Failed to get report: "+err.Error())
		return
	}

	if report == nil {
		utilities.ShowError(c, http.StatusNotFound, "Report not found")
		return
	}

	response := &ReportResponse{
		ID:          report.ID.String(),
		Title:       report.Title,
		Description: report.Description,
		Type:        report.Type,
		Format:      report.Format,
		Status:      report.Status,
		Filters:     report.Filters,
		CreatedAt:   report.CreatedAt,
		CreatedBy:   report.CreatedByID,
	}

	// Include data if format is JSON
	if report.Format == models.ReportFormatJSON {
		var data interface{}
		if dataBytes, ok := report.Data.([]byte); ok {
			if err := json.Unmarshal(dataBytes, &data); err == nil {
				response.Data = data
			}
		}
	}

	utilities.Show(c, http.StatusOK, "Report retrieved successfully", response)
}

// @Summary List reports
// @Description Get a paginated list of reports for the authenticated user
// @Tags reports
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param type query string false "Filter by report type"
// @Param status query string false "Filter by report status"
// @Success 200 {object} utils.PaginatedResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/reports [get]
func (h *ReportsHandler) ListReports(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utilities.ShowError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		utilities.ShowError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	reports, total, err := h.reportsService.ListReports(userID, limit, offset)
	if err != nil {
		utilities.ShowError(c, http.StatusInternalServerError, "Failed to list reports: "+err.Error())
		return
	}

	// Convert to response format
	responseData := make([]ReportResponse, len(reports))
	for i, report := range reports {
		responseData[i] = ReportResponse{
			ID:          report.ID.String(),
			Title:       report.Title,
			Description: report.Description,
			Type:        report.Type,
			Format:      report.Format,
			Status:      report.Status,
			Filters:     report.Filters,
			CreatedAt:   report.CreatedAt,
			CreatedBy:   report.CreatedByID,
		}
	}

	response := gin.H{
		"data": responseData,
		"pagination": utilities.GeneratePaginationData(limit, page, int(total)),
	}

	utilities.Show(c, http.StatusOK, "Reports retrieved successfully", response)
}

// @Summary Delete report
// @Description Delete a specific report by its ID
// @Tags reports
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/reports/{id} [delete]
func (h *ReportsHandler) DeleteReport(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utilities.ShowError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		utilities.ShowError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	reportID := c.Param("id")
	if reportID == "" {
		utilities.ShowError(c, http.StatusBadRequest, "Report ID is required")
		return
	}

	err := h.reportsService.DeleteReport(reportID, userID)
	if err != nil {
		utilities.ShowError(c, http.StatusInternalServerError, "Failed to delete report: "+err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "Report deleted successfully", nil)
}

// Report Template Handlers

// CreateReportTemplateRequest represents the request for creating report templates
type CreateReportTemplateRequest struct {
	Name        string               `json:"name" binding:"required"`
	Description string               `json:"description"`
	ReportType  models.ReportType    `json:"report_type" binding:"required"`
	Filters     models.ReportFilters `json:"filters"`
	IsPublic    bool                 `json:"is_public"`
}

// @Summary Create report template
// @Description Create a new report template for reuse
// @Tags report-templates
// @Accept json
// @Produce json
// @Param request body CreateReportTemplateRequest true "Template creation request"
// @Success 201 {object} models.ReportTemplate
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/reports/templates [post]
func (h *ReportsHandler) CreateReportTemplate(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utilities.ShowError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		utilities.ShowError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	var req CreateReportTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utilities.ShowError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	template := &models.ReportTemplate{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.ReportType,
		Filters:     req.Filters,
		IsPublic:    req.IsPublic,
		CreatedByID: userID,
	}

	err := h.reportsService.CreateReportTemplate(userID, template)
	if err != nil {
		utilities.ShowError(c, http.StatusInternalServerError, "Failed to create template: "+err.Error())
		return
	}

	utilities.Show(c, http.StatusCreated, "Template created successfully", template)
}

// @Summary List report templates
// @Description Get a list of available report templates
// @Tags report-templates
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/reports/templates [get]
func (h *ReportsHandler) ListReportTemplates(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utilities.ShowError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		utilities.ShowError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	templates, err := h.reportsService.ListReportTemplates(userID, true)
	if err != nil {
		utilities.ShowError(c, http.StatusInternalServerError, "Failed to list templates: "+err.Error())
		return
	}

	response := gin.H{
		"data": templates,
	}

	utilities.Show(c, http.StatusOK, "Templates retrieved successfully", response)
}

// @Summary Delete report template
// @Description Delete a specific report template by its ID
// @Tags report-templates
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/reports/templates/{id} [delete]
func (h *ReportsHandler) DeleteReportTemplate(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utilities.ShowError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		utilities.ShowError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	templateID := c.Param("id")
	if templateID == "" {
		utilities.ShowError(c, http.StatusBadRequest, "Template ID is required")
		return
	}

	err := h.reportsService.DeleteReportTemplate(templateID, userID)
	if err != nil {
		utilities.ShowError(c, http.StatusInternalServerError, "Failed to delete template: "+err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "Template deleted successfully", nil)
}

// Report Schedule Handlers

// CreateReportScheduleRequest represents the request for creating report schedules
type CreateReportScheduleRequest struct {
	Name        string                  `json:"name" binding:"required"`
	Description string                  `json:"description"`
	TemplateID  string                  `json:"template_id" binding:"required"`
	Frequency   models.ReportFrequency  `json:"frequency" binding:"required"`
	NextRun     time.Time               `json:"next_run" binding:"required"`
	Recipients  []string                `json:"recipients"`
	IsActive    bool                    `json:"is_active"`
}

// @Summary Create report schedule
// @Description Create a new scheduled report
// @Tags report-schedules
// @Accept json
// @Produce json
// @Param request body CreateReportScheduleRequest true "Schedule creation request"
// @Success 201 {object} models.ReportSchedule
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/reports/schedules [post]
func (h *ReportsHandler) CreateReportSchedule(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utilities.ShowError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		utilities.ShowError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	var req CreateReportScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utilities.ShowError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	templateID, err := uuid.Parse(req.TemplateID)
	if err != nil {
		utilities.ShowError(c, http.StatusBadRequest, "Invalid template ID: "+err.Error())
		return
	}

	schedule := &models.ReportSchedule{
		Name:        req.Name,
		Description: req.Description,
		TemplateID:  templateID,
		Frequency:   req.Frequency,
		Recipients:  req.Recipients,
		NextRunAt:   req.NextRun,
		IsActive:    req.IsActive,
		CreatedByID: userID,
	}

	err = h.reportsService.CreateReportSchedule(userID, schedule)
	if err != nil {
		utilities.ShowError(c, http.StatusInternalServerError, "Failed to create schedule: "+err.Error())
		return
	}

	utilities.Show(c, http.StatusCreated, "Schedule created successfully", schedule)
}

// @Summary List report schedules
// @Description Get a list of report schedules
// @Tags report-schedules
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/reports/schedules [get]
func (h *ReportsHandler) ListReportSchedules(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utilities.ShowError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		utilities.ShowError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	schedules, err := h.reportsService.ListReportSchedules(userID)
	if err != nil {
		utilities.ShowError(c, http.StatusInternalServerError, "Failed to list schedules: "+err.Error())
		return
	}

	response := gin.H{
		"data": schedules,
	}

	utilities.Show(c, http.StatusOK, "Schedules retrieved successfully", response)
}

// @Summary Delete report schedule
// @Description Delete a specific report schedule by its ID
// @Tags report-schedules
// @Produce json
// @Param id path string true "Schedule ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/reports/schedules/{id} [delete]
func (h *ReportsHandler) DeleteReportSchedule(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utilities.ShowError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	userID, ok := userIDInterface.(string)
	if !ok {
		utilities.ShowError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	scheduleID := c.Param("id")
	if scheduleID == "" {
		utilities.ShowError(c, http.StatusBadRequest, "Schedule ID is required")
		return
	}

	err := h.reportsService.DeleteReportSchedule(scheduleID, userID)
	if err != nil {
		utilities.ShowError(c, http.StatusInternalServerError, "Failed to delete schedule: "+err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "Schedule deleted successfully", nil)
}