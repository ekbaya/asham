package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

// AuditLogHandler handles HTTP requests for audit logs
type AuditLogHandler struct {
	auditLogService *services.AuditLogService
}

// NewAuditLogHandler creates a new audit log handler instance
func NewAuditLogHandler(auditLogService *services.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{
		auditLogService: auditLogService,
	}
}

// GetAuditLogs retrieves audit logs with filtering and pagination
// @Summary Get audit logs
// @Description Retrieve audit logs with optional filtering and pagination
// @Tags audit-logs
// @Accept json
// @Produce json
// @Param user_id query string false "Filter by user ID"
// @Param action query string false "Filter by action type"
// @Param module query string false "Filter by module"
// @Param resource_type query string false "Filter by resource type"
// @Param resource_id query string false "Filter by resource ID"
// @Param success query bool false "Filter by success status"
// @Param date_from query string false "Filter from date (RFC3339 format)"
// @Param date_to query string false "Filter to date (RFC3339 format)"
// @Param ip_address query string false "Filter by IP address"
// @Param limit query int false "Limit number of results (max 1000)" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Param order_by query string false "Order by field" default("created_at DESC")
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /audit-logs [get]
func (h *AuditLogHandler) GetAuditLogs(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse query parameters
	filter, err := h.parseAuditLogFilter(c)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, fmt.Sprintf("Invalid filter parameters: %v", err))
		return
	}

	// Validate user access
	if err := h.auditLogService.ValidateUserAccess(userID, filter); err != nil {
		utilities.ShowMessage(c, http.StatusForbidden, err.Error())
		return
	}

	// Get audit logs
	auditLogs, total, err := h.auditLogService.GetAuditLogs(filter)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to retrieve audit logs")
		return
	}

	response := map[string]interface{}{
		"audit_logs": auditLogs,
		"total":      total,
		"limit":      filter.Limit,
		"offset":     filter.Offset,
	}

	utilities.Show(c, http.StatusOK, "audit_logs", response)
}

// GetAuditLogByID retrieves a specific audit log by ID
// @Summary Get audit log by ID
// @Description Retrieve a specific audit log entry by its ID
// @Tags audit-logs
// @Accept json
// @Produce json
// @Param id path string true "Audit log ID"
// @Success 200 {object} models.AuditLog
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /audit-logs/{id} [get]
func (h *AuditLogHandler) GetAuditLogByID(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid audit log ID")
		return
	}

	// Validate user access (basic filter for access check)
	filter := models.AuditLogFilter{}
	if err := h.auditLogService.ValidateUserAccess(userID, filter); err != nil {
		utilities.ShowMessage(c, http.StatusForbidden, err.Error())
		return
	}

	auditLog, err := h.auditLogService.GetAuditLogByID(id)
	if err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Audit log not found")
		return
	}

	utilities.Show(c, http.StatusOK, "audit_log", auditLog)
}

// GetResourceAuditTrail retrieves the audit trail for a specific resource
// @Summary Get resource audit trail
// @Description Retrieve the audit trail for a specific resource
// @Tags audit-logs
// @Accept json
// @Produce json
// @Param resource_type path string true "Resource type (e.g., Project, Document)"
// @Param resource_id path string true "Resource ID"
// @Param limit query int false "Limit number of results" default(100)
// @Success 200 {object} []models.AuditLog
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /audit-logs/resource/{resource_type}/{resource_id} [get]
func (h *AuditLogHandler) GetResourceAuditTrail(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	resourceType := c.Param("resource_type")
	resourceID := c.Param("resource_id")

	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 1000 {
		limit = 100
	}

	// Validate user access
	filter := models.AuditLogFilter{
		ResourceType: &resourceType,
		ResourceID:   &resourceID,
	}
	if err := h.auditLogService.ValidateUserAccess(userID, filter); err != nil {
		utilities.ShowMessage(c, http.StatusForbidden, err.Error())
		return
	}

	auditLogs, err := h.auditLogService.GetResourceAuditTrail(resourceType, resourceID, limit)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to retrieve audit trail")
		return
	}

	utilities.Show(c, http.StatusOK, "audit_trail", auditLogs)
}

// GetUserActivity retrieves audit logs for a specific user
// @Summary Get user activity
// @Description Retrieve audit logs for a specific user
// @Tags audit-logs
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param limit query int false "Limit number of results" default(100)
// @Success 200 {object} []models.AuditLog
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /audit-logs/user/{user_id} [get]
func (h *AuditLogHandler) GetUserActivity(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	if currentUserID == "" {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	targetUserID := c.Param("user_id")

	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 1000 {
		limit = 100
	}

	// Validate user access
	filter := models.AuditLogFilter{
		UserID: &targetUserID,
	}
	if err := h.auditLogService.ValidateUserAccess(currentUserID, filter); err != nil {
		utilities.ShowMessage(c, http.StatusForbidden, err.Error())
		return
	}

	auditLogs, err := h.auditLogService.GetUserActivity(targetUserID, limit)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to retrieve user activity")
		return
	}

	utilities.Show(c, http.StatusOK, "user_activity", auditLogs)
}

// SearchAuditLogs performs a text search across audit logs
// @Summary Search audit logs
// @Description Perform a text search across audit log descriptions and metadata
// @Tags audit-logs
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param user_id query string false "Filter by user ID"
// @Param action query string false "Filter by action type"
// @Param module query string false "Filter by module"
// @Param date_from query string false "Filter from date (RFC3339 format)"
// @Param date_to query string false "Filter to date (RFC3339 format)"
// @Param limit query int false "Limit number of results" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /audit-logs/search [get]
func (h *AuditLogHandler) SearchAuditLogs(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	searchQuery := c.Query("q")
	if searchQuery == "" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Search query is required")
		return
	}

	// Parse filter parameters
	filter, err := h.parseAuditLogFilter(c)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, fmt.Sprintf("Invalid filter parameters: %v", err))
		return
	}

	// Validate user access
	if err := h.auditLogService.ValidateUserAccess(userID, filter); err != nil {
		utilities.ShowMessage(c, http.StatusForbidden, err.Error())
		return
	}

	auditLogs, total, err := h.auditLogService.SearchAuditLogs(searchQuery, filter)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to search audit logs")
		return
	}

	response := map[string]interface{}{
		"audit_logs": auditLogs,
		"total":      total,
		"query":      searchQuery,
		"limit":      filter.Limit,
		"offset":     filter.Offset,
	}

	utilities.Show(c, http.StatusOK, "search_results", response)
}

// GetAuditSummary generates audit log summary statistics
// @Summary Get audit summary
// @Description Generate audit log summary statistics
// @Tags audit-logs
// @Accept json
// @Produce json
// @Param user_id query string false "Filter by user ID"
// @Param module query string false "Filter by module"
// @Param date_from query string false "Filter from date (RFC3339 format)"
// @Param date_to query string false "Filter to date (RFC3339 format)"
// @Success 200 {object} []models.AuditLogSummary
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /audit-logs/summary [get]
func (h *AuditLogHandler) GetAuditSummary(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse filter parameters
	filter, err := h.parseAuditLogFilter(c)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, fmt.Sprintf("Invalid filter parameters: %v", err))
		return
	}

	// Validate user access
	if err := h.auditLogService.ValidateUserAccess(userID, filter); err != nil {
		utilities.ShowMessage(c, http.StatusForbidden, err.Error())
		return
	}

	summary, err := h.auditLogService.GetAuditSummary(filter)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to generate audit summary")
		return
	}

	utilities.Show(c, http.StatusOK, "summary", summary)
}

// GetActivityTimeline retrieves audit logs grouped by time periods
// @Summary Get activity timeline
// @Description Retrieve audit logs grouped by time periods
// @Tags audit-logs
// @Accept json
// @Produce json
// @Param interval query string false "Time interval (hour, day, week, month)" default("day")
// @Param date_from query string false "Filter from date (RFC3339 format)"
// @Param date_to query string false "Filter to date (RFC3339 format)"
// @Param module query string false "Filter by module"
// @Success 200 {object} map[string]int64
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /audit-logs/timeline [get]
func (h *AuditLogHandler) GetActivityTimeline(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	interval := c.DefaultQuery("interval", "day")

	// Parse filter parameters
	filter, err := h.parseAuditLogFilter(c)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, fmt.Sprintf("Invalid filter parameters: %v", err))
		return
	}

	// Validate user access
	if err := h.auditLogService.ValidateUserAccess(userID, filter); err != nil {
		utilities.ShowMessage(c, http.StatusForbidden, err.Error())
		return
	}

	timeline, err := h.auditLogService.GetActivityTimeline(filter, interval)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to generate activity timeline")
		return
	}

	utilities.Show(c, http.StatusOK, "timeline", timeline)
}

// GetComplianceReport generates a comprehensive compliance report
// @Summary Get compliance report
// @Description Generate a comprehensive compliance report
// @Tags audit-logs
// @Accept json
// @Produce json
// @Param date_from query string false "Filter from date (RFC3339 format)"
// @Param date_to query string false "Filter to date (RFC3339 format)"
// @Param module query string false "Filter by module"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /audit-logs/compliance-report [get]
func (h *AuditLogHandler) GetComplianceReport(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse filter parameters
	filter, err := h.parseAuditLogFilter(c)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, fmt.Sprintf("Invalid filter parameters: %v", err))
		return
	}

	// Validate user access
	if err := h.auditLogService.ValidateUserAccess(userID, filter); err != nil {
		utilities.ShowMessage(c, http.StatusForbidden, err.Error())
		return
	}

	report, err := h.auditLogService.GetComplianceReport(filter)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to generate compliance report")
		return
	}

	utilities.Show(c, http.StatusOK, "compliance_report", report)
}

// ExportAuditLogs exports audit logs in various formats
// @Summary Export audit logs
// @Description Export audit logs in CSV, Excel, or PDF format
// @Tags audit-logs
// @Accept json
// @Produce application/octet-stream
// @Param format query string true "Export format (csv, excel, pdf)"
// @Param user_id query string false "Filter by user ID"
// @Param action query string false "Filter by action type"
// @Param module query string false "Filter by module"
// @Param date_from query string false "Filter from date (RFC3339 format)"
// @Param date_to query string false "Filter to date (RFC3339 format)"
// @Success 200 {file} file
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /audit-logs/export [get]
func (h *AuditLogHandler) ExportAuditLogs(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utilities.ShowMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	format := c.Query("format")
	if format == "" {
		format = "csv"
	}

	if format != "csv" && format != "excel" && format != "pdf" {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid export format. Supported formats: csv, excel, pdf")
		return
	}

	// Parse filter parameters
	filter, err := h.parseAuditLogFilter(c)
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, fmt.Sprintf("Invalid filter parameters: %v", err))
		return
	}

	// Set a reasonable limit for exports
	if filter.Limit == 0 || filter.Limit > 10000 {
		filter.Limit = 10000
	}

	// Validate user access
	if err := h.auditLogService.ValidateUserAccess(userID, filter); err != nil {
		utilities.ShowMessage(c, http.StatusForbidden, err.Error())
		return
	}

	// Get audit logs
	auditLogs, _, err := h.auditLogService.GetAuditLogs(filter)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to retrieve audit logs for export")
		return
	}

	switch format {
	case "csv":
		h.exportCSV(c, auditLogs)
	case "excel":
		h.exportExcel(c, auditLogs)
	case "pdf":
		h.exportPDF(c, auditLogs)
	}
}

// Helper methods

// parseAuditLogFilter parses query parameters into an AuditLogFilter
func (h *AuditLogHandler) parseAuditLogFilter(c *gin.Context) (models.AuditLogFilter, error) {
	filter := models.AuditLogFilter{}

	// Parse user_id
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		filter.UserID = &userIDStr
	}

	// Parse action
	if actionStr := c.Query("action"); actionStr != "" {
		action := models.ActionType(actionStr)
		filter.Action = &action
	}

	// Parse module
	if moduleStr := c.Query("module"); moduleStr != "" {
		module := models.ModuleType(moduleStr)
		filter.Module = &module
	}

	// Parse resource_type
	if resourceTypeStr := c.Query("resource_type"); resourceTypeStr != "" {
		filter.ResourceType = &resourceTypeStr
	}

	// Parse resource_id
	if resourceIDStr := c.Query("resource_id"); resourceIDStr != "" {
		filter.ResourceID = &resourceIDStr
	}

	// Parse success
	if successStr := c.Query("success"); successStr != "" {
		success, err := strconv.ParseBool(successStr)
		if err != nil {
			return filter, fmt.Errorf("invalid success parameter: %v", err)
		}
		filter.Success = &success
	}

	// Parse date_from
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		dateFrom, err := time.Parse(time.RFC3339, dateFromStr)
		if err != nil {
			return filter, fmt.Errorf("invalid date_from parameter: %v", err)
		}
		filter.DateFrom = &dateFrom
	}

	// Parse date_to
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		dateTo, err := time.Parse(time.RFC3339, dateToStr)
		if err != nil {
			return filter, fmt.Errorf("invalid date_to parameter: %v", err)
		}
		filter.DateTo = &dateTo
	}

	// Parse ip_address
	if ipAddressStr := c.Query("ip_address"); ipAddressStr != "" {
		filter.IPAddress = &ipAddressStr
	}

	// Parse limit
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	if limit > 1000 {
		limit = 1000
	}
	filter.Limit = limit

	// Parse offset
	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	filter.Offset = offset

	// Parse order_by
	orderBy := c.DefaultQuery("order_by", "created_at DESC")
	filter.OrderBy = orderBy

	return filter, nil
}

// exportCSV exports audit logs as CSV
func (h *AuditLogHandler) exportCSV(c *gin.Context, auditLogs []models.AuditLog) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=audit_logs_%s.csv", time.Now().Format("2006-01-02")))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	header := []string{
		"ID", "User", "Action", "Module", "Resource Type", "Resource ID", "Resource Title",
		"Description", "Success", "Error Message", "IP Address", "Duration (ms)", "Created At",
	}
	writer.Write(header)

	// Write data
	for _, log := range auditLogs {
		userName := ""
		if log.User != nil {
			userName = fmt.Sprintf("%s %s", log.User.FirstName, log.User.LastName)
		}

		resourceID := ""
		if log.ResourceID != nil {
			resourceID = *log.ResourceID
		}

		row := []string{
			log.ID.String(),
			userName,
			string(log.Action),
			string(log.Module),
			log.ResourceType,
			resourceID,
			log.ResourceTitle,
			log.Description,
			strconv.FormatBool(log.Success),
			log.ErrorMessage,
			log.IPAddress,
			strconv.FormatInt(log.Duration, 10),
			log.CreatedAt.Format(time.RFC3339),
		}
		writer.Write(row)
	}
}

// exportExcel exports audit logs as Excel
func (h *AuditLogHandler) exportExcel(c *gin.Context, auditLogs []models.AuditLog) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create a new worksheet
	sheetName := "Audit Logs"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to create Excel worksheet")
		return
	}

	// Set headers
	headers := []string{
		"ID", "User", "Action", "Module", "Resource Type", "Resource ID", "Resource Title",
		"Description", "Success", "Error Message", "IP Address", "Duration (ms)", "Created At",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Set header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})
	f.SetRowStyle(sheetName, 1, 1, headerStyle)

	// Add data
	for i, log := range auditLogs {
		row := i + 2
		userName := ""
		if log.User != nil {
			userName = fmt.Sprintf("%s %s", log.User.FirstName, log.User.LastName)
		}

		resourceID := ""
		if log.ResourceID != nil {
			resourceID = *log.ResourceID
		}

		data := []interface{}{
			log.ID.String(),
			userName,
			string(log.Action),
			string(log.Module),
			log.ResourceType,
			resourceID,
			log.ResourceTitle,
			log.Description,
			log.Success,
			log.ErrorMessage,
			log.IPAddress,
			log.Duration,
			log.CreatedAt.Format(time.RFC3339),
		}

		for j, value := range data {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// Auto-fit columns
	for i := range headers {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheetName, colName, colName, 15)
	}

	// Set active sheet
	f.SetActiveSheet(index)

	// Generate buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to generate Excel file")
		return
	}

	// Set headers and send file
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=audit_logs_%s.xlsx", time.Now().Format("2006-01-02")))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// exportPDF exports audit logs as PDF (simplified text-based approach)
func (h *AuditLogHandler) exportPDF(c *gin.Context, auditLogs []models.AuditLog) {
	// Create a detailed text report
	var content bytes.Buffer
	content.WriteString("AUDIT LOGS REPORT\n")
	content.WriteString(fmt.Sprintf("Generated on: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("Total Records: %d\n", len(auditLogs)))
	content.WriteString("=================================================\n\n")

	for i, log := range auditLogs {
		userName := "Unknown User"
		if log.User != nil {
			userName = fmt.Sprintf("%s %s", log.User.FirstName, log.User.LastName)
		}

		resourceID := "N/A"
		if log.ResourceID != nil {
			resourceID = *log.ResourceID
		}

		content.WriteString(fmt.Sprintf("%d. %s\n", i+1, log.ResourceTitle))
		content.WriteString(fmt.Sprintf("   ID: %s\n", log.ID.String()))
		content.WriteString(fmt.Sprintf("   User: %s\n", userName))
		content.WriteString(fmt.Sprintf("   Action: %s\n", log.Action))
		content.WriteString(fmt.Sprintf("   Module: %s\n", log.Module))
		content.WriteString(fmt.Sprintf("   Resource: %s (%s)\n", log.ResourceType, resourceID))
		content.WriteString(fmt.Sprintf("   Success: %t\n", log.Success))
		if log.ErrorMessage != "" {
			content.WriteString(fmt.Sprintf("   Error: %s\n", log.ErrorMessage))
		}
		content.WriteString(fmt.Sprintf("   IP Address: %s\n", log.IPAddress))
		content.WriteString(fmt.Sprintf("   Duration: %dms\n", log.Duration))
		content.WriteString(fmt.Sprintf("   Timestamp: %s\n", log.CreatedAt.Format("2006-01-02 15:04:05")))
		content.WriteString(fmt.Sprintf("   Description: %s\n", log.Description))
		content.WriteString("\n-------------------------------------------------\n\n")
	}

	content.WriteString("\n\nEnd of Report")

	// For now, return as plain text with PDF content type
	// In a production environment, you would use a proper PDF library like gofpdf
	c.Header("Content-Type", "text/plain")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=audit_logs_%s.txt", time.Now().Format("2006-01-02")))
	c.String(http.StatusOK, content.String())
}