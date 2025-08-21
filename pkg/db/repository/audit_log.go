package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLogRepository handles database operations for audit logs
type AuditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository creates a new audit log repository instance
func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

// Create adds a new audit log entry to the database
func (r *AuditLogRepository) Create(auditLog *models.AuditLog) error {
	if auditLog.ID == uuid.Nil {
		auditLog.ID = uuid.New()
	}
	if auditLog.CreatedAt.IsZero() {
		auditLog.CreatedAt = time.Now()
	}
	return r.db.Create(auditLog).Error
}

// CreateBatch creates multiple audit log entries in a single transaction
func (r *AuditLogRepository) CreateBatch(auditLogs []models.AuditLog) error {
	if len(auditLogs) == 0 {
		return nil
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i := range auditLogs {
		if auditLogs[i].ID == uuid.Nil {
			auditLogs[i].ID = uuid.New()
		}
		if auditLogs[i].CreatedAt.IsZero() {
			auditLogs[i].CreatedAt = time.Now()
		}
	}

	if err := tx.CreateInBatches(auditLogs, 100).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetByID retrieves an audit log by its ID
func (r *AuditLogRepository) GetByID(id uuid.UUID) (*models.AuditLog, error) {
	var auditLog models.AuditLog
	err := r.db.Preload("User").First(&auditLog, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &auditLog, nil
}

// List retrieves audit logs with filtering and pagination
func (r *AuditLogRepository) List(filter models.AuditLogFilter) ([]models.AuditLog, int64, error) {
	var auditLogs []models.AuditLog
	var total int64

	query := r.db.Model(&models.AuditLog{})

	// Apply filters
	query = r.applyFilters(query, filter)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	orderBy := "created_at DESC"
	if filter.OrderBy != "" {
		orderBy = filter.OrderBy
	}

	limit := 50 // Default limit
	if filter.Limit > 0 && filter.Limit <= 1000 {
		limit = filter.Limit
	}

	offset := 0
	if filter.Offset > 0 {
		offset = filter.Offset
	}

	err := query.Preload("User").
		Order(orderBy).
		Limit(limit).
		Offset(offset).
		Find(&auditLogs).Error

	return auditLogs, total, err
}

// GetByResourceID retrieves audit logs for a specific resource
func (r *AuditLogRepository) GetByResourceID(resourceType, resourceID string, limit int) ([]models.AuditLog, error) {
	var auditLogs []models.AuditLog

	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	err := r.db.Preload("User").
		Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).
		Order("created_at DESC").
		Limit(limit).
		Find(&auditLogs).Error

	return auditLogs, err
}

// GetByUserID retrieves audit logs for a specific user
func (r *AuditLogRepository) GetByUserID(userID string, limit int) ([]models.AuditLog, error) {
	var auditLogs []models.AuditLog

	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	err := r.db.Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&auditLogs).Error

	return auditLogs, err
}

// GetSummary generates audit log summary statistics
func (r *AuditLogRepository) GetSummary(filter models.AuditLogFilter) ([]models.AuditLogSummary, error) {
	var summaries []models.AuditLogSummary

	query := r.db.Model(&models.AuditLog{}).
		Select("module, action, COUNT(*) as count, AVG(CASE WHEN success THEN 1.0 ELSE 0.0 END) as success_rate, MAX(created_at) as last_occurred").
		Group("module, action")

	// Apply date filters if provided
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}
	if filter.Module != nil {
		query = query.Where("module = ?", *filter.Module)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	err := query.Order("count DESC").Find(&summaries).Error
	return summaries, err
}

// GetActivityTimeline retrieves audit logs grouped by time periods
func (r *AuditLogRepository) GetActivityTimeline(filter models.AuditLogFilter, interval string) (map[string]int64, error) {
	result := make(map[string]int64)

	// Validate interval
	if interval != "hour" && interval != "day" && interval != "week" && interval != "month" {
		interval = "day"
	}

	var timeFormat string
	switch interval {
	case "hour":
		timeFormat = "DATE_TRUNC('hour', created_at)"
	case "day":
		timeFormat = "DATE_TRUNC('day', created_at)"
	case "week":
		timeFormat = "DATE_TRUNC('week', created_at)"
	case "month":
		timeFormat = "DATE_TRUNC('month', created_at)"
	}

	query := r.db.Model(&models.AuditLog{}).
		Select(fmt.Sprintf("%s as time_period, COUNT(*) as count", timeFormat)).
		Group("time_period")

	// Apply filters
	query = r.applyFilters(query, filter)

	rows, err := query.Order("time_period").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var timePeriod time.Time
		var count int64
		if err := rows.Scan(&timePeriod, &count); err != nil {
			return nil, err
		}
		result[timePeriod.Format("2006-01-02 15:04:05")] = count
	}

	return result, nil
}

// GetFailedActions retrieves audit logs for failed actions
func (r *AuditLogRepository) GetFailedActions(filter models.AuditLogFilter) ([]models.AuditLog, error) {
	var auditLogs []models.AuditLog

	query := r.db.Model(&models.AuditLog{}).Where("success = ?", false)
	query = r.applyFilters(query, filter)

	limit := 100
	if filter.Limit > 0 && filter.Limit <= 1000 {
		limit = filter.Limit
	}

	err := query.Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Find(&auditLogs).Error

	return auditLogs, err
}

// GetUserActivity retrieves user activity statistics
func (r *AuditLogRepository) GetUserActivity(filter models.AuditLogFilter) (map[string]int64, error) {
	result := make(map[string]int64)

	query := r.db.Model(&models.AuditLog{}).
		Joins("LEFT JOIN members ON audit_logs.user_id = members.id::text").
		Select("CONCAT(members.first_name, ' ', members.last_name) as user_name, COUNT(*) as count").
		Where("audit_logs.user_id IS NOT NULL").
		Group("members.id, members.first_name, members.last_name")

	// Apply filters
	query = r.applyFilters(query, filter)

	rows, err := query.Order("count DESC").Limit(50).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userName string
		var count int64
		if err := rows.Scan(&userName, &count); err != nil {
			return nil, err
		}
		result[userName] = count
	}

	return result, nil
}

// DeleteOldLogs removes audit logs older than the specified duration
func (r *AuditLogRepository) DeleteOldLogs(olderThan time.Duration) (int64, error) {
	cutoffDate := time.Now().Add(-olderThan)
	result := r.db.Unscoped().Where("created_at < ?", cutoffDate).Delete(&models.AuditLog{})
	return result.RowsAffected, result.Error
}

// GetModuleActivity retrieves activity statistics by module
func (r *AuditLogRepository) GetModuleActivity(filter models.AuditLogFilter) (map[string]int64, error) {
	result := make(map[string]int64)

	query := r.db.Model(&models.AuditLog{}).
		Select("module, COUNT(*) as count").
		Group("module")

	// Apply filters
	query = r.applyFilters(query, filter)

	rows, err := query.Order("count DESC").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var module string
		var count int64
		if err := rows.Scan(&module, &count); err != nil {
			return nil, err
		}
		result[module] = count
	}

	return result, nil
}

// applyFilters applies the provided filters to the query
func (r *AuditLogRepository) applyFilters(query *gorm.DB, filter models.AuditLogFilter) *gorm.DB {
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Action != nil {
		query = query.Where("action = ?", *filter.Action)
	}
	if filter.Module != nil {
		query = query.Where("module = ?", *filter.Module)
	}
	if filter.ResourceType != nil {
		query = query.Where("resource_type = ?", *filter.ResourceType)
	}
	if filter.ResourceID != nil {
		query = query.Where("resource_id = ?", *filter.ResourceID)
	}
	if filter.Success != nil {
		query = query.Where("success = ?", *filter.Success)
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}
	if filter.IPAddress != nil {
		query = query.Where("ip_address = ?", *filter.IPAddress)
	}

	return query
}

// SearchLogs performs a text search across audit log descriptions and metadata
func (r *AuditLogRepository) SearchLogs(searchTerm string, filter models.AuditLogFilter) ([]models.AuditLog, int64, error) {
	var auditLogs []models.AuditLog
	var total int64

	if searchTerm == "" {
		return r.List(filter)
	}

	// Clean and prepare search term
	searchTerm = strings.TrimSpace(searchTerm)
	searchPattern := fmt.Sprintf("%%%s%%", searchTerm)

	query := r.db.Model(&models.AuditLog{}).Where(
		"description ILIKE ? OR resource_title ILIKE ? OR metadata::text ILIKE ?",
		searchPattern, searchPattern, searchPattern,
	)

	// Apply additional filters
	query = r.applyFilters(query, filter)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	orderBy := "created_at DESC"
	if filter.OrderBy != "" {
		orderBy = filter.OrderBy
	}

	limit := 50
	if filter.Limit > 0 && filter.Limit <= 1000 {
		limit = filter.Limit
	}

	offset := 0
	if filter.Offset > 0 {
		offset = filter.Offset
	}

	err := query.Preload("User").
		Order(orderBy).
		Limit(limit).
		Offset(offset).
		Find(&auditLogs).Error

	return auditLogs, total, err
}

// GetComplianceReport generates a compliance report for audit logs
func (r *AuditLogRepository) GetComplianceReport(filter models.AuditLogFilter) (map[string]interface{}, error) {
	report := make(map[string]interface{})

	// Total logs count
	var totalLogs int64
	query := r.db.Model(&models.AuditLog{})
	query = r.applyFilters(query, filter)
	if err := query.Count(&totalLogs).Error; err != nil {
		return nil, err
	}
	report["total_logs"] = totalLogs

	// Success rate
	var successCount int64
	successQuery := r.db.Model(&models.AuditLog{}).Where("success = ?", true)
	successQuery = r.applyFilters(successQuery, filter)
	if err := successQuery.Count(&successCount).Error; err != nil {
		return nil, err
	}

	successRate := float64(0)
	if totalLogs > 0 {
		successRate = float64(successCount) / float64(totalLogs) * 100
	}
	report["success_rate"] = successRate

	// Module activity
	moduleActivity, err := r.GetModuleActivity(filter)
	if err != nil {
		return nil, err
	}
	report["module_activity"] = moduleActivity

	// User activity
	userActivity, err := r.GetUserActivity(filter)
	if err != nil {
		return nil, err
	}
	report["user_activity"] = userActivity

	// Failed actions
	failedActions, err := r.GetFailedActions(filter)
	if err != nil {
		return nil, err
	}
	report["failed_actions_count"] = len(failedActions)
	report["recent_failed_actions"] = failedActions

	return report, nil
}