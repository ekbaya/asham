package services

import (
	"fmt"
	"time"

	"github.com/ekbaya/asham/pkg/db/repository"
	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
)

type ReportsService struct {
	reportsRepo *repository.ReportsRepository
	projectRepo *repository.ProjectRepository
	memberRepo  *repository.MemberRepository
}

func NewReportsService(
	reportsRepo *repository.ReportsRepository,
	projectRepo *repository.ProjectRepository,
	memberRepo *repository.MemberRepository,
) *ReportsService {
	return &ReportsService{
		reportsRepo: reportsRepo,
		projectRepo: projectRepo,
		memberRepo:  memberRepo,
	}
}

// Report Generation
func (s *ReportsService) GenerateReport(userID string, reportType models.ReportType, filters models.ReportFilters) (*models.Report, error) {
	// Validate user permissions
	user, err := s.memberRepo.GetMemberByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Apply role-based filtering
	filters = s.applyRoleBasedFilters(user, filters)

	// Generate report data based on type
	var data interface{}
	switch reportType {
	case models.ReportTypeProject:
		data, err = s.reportsRepo.GetProjectReportData(filters)
	case models.ReportTypeBallot:
		data, err = s.reportsRepo.GetBallotReportData(filters)
	case models.ReportTypeCommittee:
		data, err = s.reportsRepo.GetCommitteeReportData(filters)
	case models.ReportTypeDocument:
		data, err = s.reportsRepo.GetDocumentReportData(filters)
	case models.ReportTypeUserActivity:
		data, err = s.reportsRepo.GetUserActivityReportData(filters)
	default:
		return nil, fmt.Errorf("unsupported report type: %s", reportType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate report data: %w", err)
	}

	// Create report record
	report := &models.Report{
		ID:          uuid.New(),
		Type:        reportType,
		Title:       s.generateReportTitle(reportType, filters),
		Description: s.generateReportDescription(reportType, filters),
		Filters:     filters,
		Data:        data,
		Format:      models.ReportFormatJSON, // Default format
		Status:      models.ReportStatusCompleted,
		CreatedByID: userID,
		CreatedAt:   time.Now(),
	}

	// Save report
	if err := s.reportsRepo.CreateReport(report); err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}

	return report, nil
}

// Real-time Report Generation
func (s *ReportsService) GenerateRealTimeReport(userID string, reportType models.ReportType, filters models.ReportFilters) (interface{}, error) {
	// Validate user permissions
	user, err := s.memberRepo.GetMemberByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Apply role-based filtering
	filters = s.applyRoleBasedFilters(user, filters)

	// Generate real-time data without saving
	switch reportType {
	case models.ReportTypeProject:
		return s.reportsRepo.GetProjectReportData(filters)
	case models.ReportTypeBallot:
		return s.reportsRepo.GetBallotReportData(filters)
	case models.ReportTypeCommittee:
		return s.reportsRepo.GetCommitteeReportData(filters)
	case models.ReportTypeDocument:
		return s.reportsRepo.GetDocumentReportData(filters)
	case models.ReportTypeUserActivity:
		return s.reportsRepo.GetUserActivityReportData(filters)
	default:
		return nil, fmt.Errorf("unsupported report type: %s", reportType)
	}
}

// Dashboard Metrics
func (s *ReportsService) GetDashboardMetrics(userID string) ([]models.DashboardMetric, error) {
	// Validate user permissions
	user, err := s.memberRepo.GetMemberByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get base metrics
	metrics, err := s.reportsRepo.GetDashboardMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard metrics: %w", err)
	}

	// Filter metrics based on user role
	filteredMetrics := s.filterMetricsByRole(user, metrics)

	return filteredMetrics, nil
}

// Report Template Management
func (s *ReportsService) CreateReportTemplate(userID string, template *models.ReportTemplate) error {
	template.ID = uuid.New()
	template.CreatedByID = userID
	template.CreatedAt = time.Now()
	return s.reportsRepo.CreateReportTemplate(template)
}

func (s *ReportsService) GetReportTemplate(userID, templateID string) (*models.ReportTemplate, error) {
	templateUUID, err := uuid.Parse(templateID)
	if err != nil {
		return nil, fmt.Errorf("invalid template ID: %w", err)
	}
	template, err := s.reportsRepo.GetReportTemplateByID(templateUUID)
	if err != nil {
		return nil, err
	}

	// Check if user has access to this template
	user, err := s.memberRepo.GetMemberByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Allow access if template is public or user is the creator
	if !template.IsPublic && template.CreatedByID != userID {
		// Check if user has admin role
		hasAdminAccess := false
		for _, role := range user.Roles {
			if role.Title == "ARSO_SECRETARIAT" || role.Title == "SMC_MEMBER" {
				hasAdminAccess = true
				break
			}
		}
		if !hasAdminAccess {
			return nil, fmt.Errorf("access denied to template")
		}
	}

	return template, nil
}

func (s *ReportsService) ListReportTemplates(userID string, includePublic bool) ([]models.ReportTemplate, error) {
	return s.reportsRepo.ListReportTemplates(userID, includePublic)
}

func (s *ReportsService) UpdateReportTemplate(userID string, template *models.ReportTemplate) error {
	// Check if user owns the template or has admin access
	existing, err := s.reportsRepo.GetReportTemplateByID(template.ID)
	if err != nil {
		return err
	}

	if existing.CreatedByID != userID {
		user, err := s.memberRepo.GetMemberByID(userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		// Check admin access
		hasAdminAccess := false
		for _, role := range user.Roles {
			if role.Title == "ARSO_SECRETARIAT" || role.Title == "SMC_MEMBER" {
				hasAdminAccess = true
				break
			}
		}
		if !hasAdminAccess {
			return fmt.Errorf("access denied to update template")
		}
	}

	template.UpdatedAt = time.Now()
	return s.reportsRepo.UpdateReportTemplate(template)
}

func (s *ReportsService) DeleteReportTemplate(userID, templateID string) error {
	// Check if user owns the template or has admin access
	templateUUID, err := uuid.Parse(templateID)
	if err != nil {
		return fmt.Errorf("invalid template ID: %w", err)
	}
	template, err := s.reportsRepo.GetReportTemplateByID(templateUUID)
	if err != nil {
		return err
	}

	if template.CreatedByID != userID {
		user, err := s.memberRepo.GetMemberByID(userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		// Check admin access
		hasAdminAccess := false
		for _, role := range user.Roles {
			if role.Title == "ARSO_SECRETARIAT" || role.Title == "SMC_MEMBER" {
				hasAdminAccess = true
				break
			}
		}
		if !hasAdminAccess {
			return fmt.Errorf("access denied to delete template")
		}
	}

	return s.reportsRepo.DeleteReportTemplate(templateUUID)
}

// Report Schedule Management
func (s *ReportsService) CreateReportSchedule(userID string, schedule *models.ReportSchedule) error {
	schedule.ID = uuid.New()
	schedule.CreatedByID = userID
	schedule.CreatedAt = time.Now()
	schedule.NextRunAt = s.calculateNextRunTime(schedule.Frequency, time.Now())
	return s.reportsRepo.CreateReportSchedule(schedule)
}

func (s *ReportsService) GetReportSchedule(userID, scheduleID string) (*models.ReportSchedule, error) {
	scheduleUUID, err := uuid.Parse(scheduleID)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule ID: %w", err)
	}
	schedule, err := s.reportsRepo.GetReportScheduleByID(scheduleUUID)
	if err != nil {
		return nil, err
	}

	// Check if user owns the schedule or has admin access
	if schedule.CreatedByID != userID {
		user, err := s.memberRepo.GetMemberByID(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
		// Check admin access
		hasAdminAccess := false
		for _, role := range user.Roles {
			if role.Title == "ARSO_SECRETARIAT" || role.Title == "SMC_MEMBER" {
				hasAdminAccess = true
				break
			}
		}
		if !hasAdminAccess {
			return nil, fmt.Errorf("access denied to schedule")
		}
	}

	return schedule, nil
}

func (s *ReportsService) ListReportSchedules(userID string) ([]models.ReportSchedule, error) {
	return s.reportsRepo.ListReportSchedules(userID)
}

func (s *ReportsService) UpdateReportSchedule(userID string, schedule *models.ReportSchedule) error {
	// Check if user owns the schedule or has admin access
	existing, err := s.reportsRepo.GetReportScheduleByID(schedule.ID)
	if err != nil {
		return err
	}

	if existing.CreatedByID != userID {
		user, err := s.memberRepo.GetMemberByID(userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		// Check admin access
		hasAdminAccess := false
		for _, role := range user.Roles {
			if role.Title == "ARSO_SECRETARIAT" || role.Title == "SMC_MEMBER" {
				hasAdminAccess = true
				break
			}
		}
		if !hasAdminAccess {
			return fmt.Errorf("access denied to update schedule")
		}
	}

	schedule.UpdatedAt = time.Now()
	schedule.NextRunAt = s.calculateNextRunTime(schedule.Frequency, time.Now())
	return s.reportsRepo.UpdateReportSchedule(schedule)
}

func (s *ReportsService) DeleteReportSchedule(userID, scheduleID string) error {
	// Check if user owns the schedule or has admin access
	scheduleUUID, err := uuid.Parse(scheduleID)
	if err != nil {
		return fmt.Errorf("invalid schedule ID: %w", err)
	}
	schedule, err := s.reportsRepo.GetReportScheduleByID(scheduleUUID)
	if err != nil {
		return err
	}

	if schedule.CreatedByID != userID {
		user, err := s.memberRepo.GetMemberByID(userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		// Check admin access
		hasAdminAccess := false
		for _, role := range user.Roles {
			if role.Title == "ARSO_SECRETARIAT" || role.Title == "SMC_MEMBER" {
				hasAdminAccess = true
				break
			}
		}
		if !hasAdminAccess {
			return fmt.Errorf("access denied to delete schedule")
		}
	}

	return s.reportsRepo.DeleteReportSchedule(scheduleUUID)
}

// Scheduled Report Processing
func (s *ReportsService) ProcessScheduledReports() error {
	// Get all active schedules that are due
	schedules, err := s.reportsRepo.GetDueSchedules()
	if err != nil {
		return fmt.Errorf("failed to get due schedules: %w", err)
	}

	// Process each schedule
	for _, schedule := range schedules {
		if err := s.processScheduledReport(&schedule); err != nil {
			// Log error but continue processing other schedules
			fmt.Printf("Failed to process scheduled report %s: %v\n", schedule.ID, err)
			continue
		}

		// Update next run time
		schedule.NextRunAt = s.calculateNextRunTime(schedule.Frequency, time.Now())
		now := time.Now()
		schedule.LastRunAt = &now
		if err := s.reportsRepo.UpdateReportSchedule(&schedule); err != nil {
			fmt.Printf("Failed to update schedule %s: %v\n", schedule.ID, err)
		}
	}

	return nil
}

func (s *ReportsService) processScheduledReport(schedule *models.ReportSchedule) error {
	// Generate report using the schedule's template
	template, err := s.reportsRepo.GetReportTemplateByID(schedule.TemplateID)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	// Generate report
	report, err := s.GenerateReport(schedule.CreatedByID, template.Type, template.Filters)
	if err != nil {
		return fmt.Errorf("failed to generate scheduled report: %w", err)
	}

	// Update report with schedule information
	report.Title = fmt.Sprintf("[Scheduled] %s", report.Title)
	report.Description = fmt.Sprintf("Automatically generated report from schedule: %s", schedule.Name)

	// Save the updated report
	if err := s.reportsRepo.UpdateReport(report); err != nil {
		return fmt.Errorf("failed to update scheduled report: %w", err)
	}

	// TODO: Send notification or email to recipients

	return nil
}

// Report Access and Listing
func (s *ReportsService) GetReport(userID, reportID string) (*models.Report, error) {
	reportUUID, err := uuid.Parse(reportID)
	if err != nil {
		return nil, fmt.Errorf("invalid report ID: %w", err)
	}
	report, err := s.reportsRepo.GetReportByID(reportUUID)
	if err != nil {
		return nil, err
	}

	// Check if user has access to this report
	user, err := s.memberRepo.GetMemberByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !s.hasReportAccess(user, report) {
		return nil, fmt.Errorf("access denied to report")
	}

	return report, nil
}

func (s *ReportsService) ListReports(userID string, limit, offset int) ([]models.Report, int64, error) {
	// Apply role-based filtering in the repository query
	return s.reportsRepo.ListReports(userID, limit, offset)
}

func (s *ReportsService) DeleteReport(userID, reportID string) error {
	reportUUID, err := uuid.Parse(reportID)
	if err != nil {
		return fmt.Errorf("invalid report ID: %w", err)
	}
	report, err := s.reportsRepo.GetReportByID(reportUUID)
	if err != nil {
		return err
	}

	// Check if user has access to delete this report
	user, err := s.memberRepo.GetMemberByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Allow deletion if user created the report or has admin access
	if report.CreatedByID != userID {
		hasAdminAccess := false
		for _, role := range user.Roles {
			if role.Title == "ARSO_SECRETARIAT" || role.Title == "SMC_MEMBER" {
				hasAdminAccess = true
				break
			}
		}
		if !hasAdminAccess {
			return fmt.Errorf("access denied to delete report")
		}
	}

	return s.reportsRepo.DeleteReport(reportUUID)
}

// Helper methods
func (s *ReportsService) applyRoleBasedFilters(user *models.Member, filters models.ReportFilters) models.ReportFilters {
	// Apply role-based filtering based on user's role and NSB
	for _, role := range user.Roles {
		switch role.Title {
		case "ARSO_SECRETARIAT", "SMC_MEMBER", "COUNCIL_MEMBER":
			// These roles can see all data - no additional filtering
			return filters
		case "NSB_REPRESENTATIVE", "COMMITTEE_CHAIR", "COMMITTEE_MEMBER":
			// Filter to only show data related to user's NSB/committees
			if filters.NSBID == nil {
				filters.NSBID = user.NationalStandardBodyID
			}
			// Additional committee filtering would be added here
		case "OBSERVER":
			// Observers have limited access - filter to public data only
			published := true
			filters.Published = &published
		}
	}

	return filters
}

func (s *ReportsService) filterMetricsByRole(user *models.Member, metrics []models.DashboardMetric) []models.DashboardMetric {
	// Filter dashboard metrics based on user role
	var filteredMetrics []models.DashboardMetric

	for _, metric := range metrics {
		if s.hasMetricAccess(user, metric) {
			filteredMetrics = append(filteredMetrics, metric)
		}
	}

	return filteredMetrics
}

func (s *ReportsService) hasMetricAccess(user *models.Member, metric models.DashboardMetric) bool {
	// Check if user has access to specific metric based on role
	for _, role := range user.Roles {
		switch role.Title {
		case "ARSO_SECRETARIAT", "SMC_MEMBER", "COUNCIL_MEMBER":
			return true // Full access
		case "NSB_REPRESENTATIVE", "COMMITTEE_CHAIR", "COMMITTEE_MEMBER":
			// Access to most metrics except sensitive ones
			return metric.Category != "SENSITIVE"
		case "OBSERVER":
			// Limited access to basic metrics only
			return metric.Category == "PROJECT" || metric.Category == "DOCUMENT"
		}
	}
	return false
}

func (s *ReportsService) hasReportAccess(user *models.Member, report *models.Report) bool {
	// Check if user has access to specific report based on role
	for _, role := range user.Roles {
		switch role.Title {
		case "ARSO_SECRETARIAT", "SMC_MEMBER", "COUNCIL_MEMBER":
			return true // Full access
		case "NSB_REPRESENTATIVE", "COMMITTEE_CHAIR", "COMMITTEE_MEMBER":
			// Access to reports related to their NSB/committees
			return true // Simplified - would need more complex logic
		case "OBSERVER":
			// Limited access to public reports only
			return report.Type == models.ReportTypeProject || report.Type == models.ReportTypeDocument
		}
	}
	return false
}

func (s *ReportsService) generateReportTitle(reportType models.ReportType, filters models.ReportFilters) string {
	baseTitle := fmt.Sprintf("%s Report", reportType)
	if filters.DateFrom != nil && filters.DateTo != nil {
		return fmt.Sprintf("%s (%s to %s)", baseTitle, filters.DateFrom.Format("2006-01-02"), filters.DateTo.Format("2006-01-02"))
	}
	return baseTitle
}

func (s *ReportsService) generateReportDescription(reportType models.ReportType, filters models.ReportFilters) string {
	description := fmt.Sprintf("Generated %s report", reportType)
	
	var filterParts []string
	if filters.NSBID != nil {
		filterParts = append(filterParts, "filtered by NSB")
	}
	if filters.CommitteeID != nil {
		filterParts = append(filterParts, "filtered by Committee")
	}
	// ProjectID field doesn't exist in ReportFilters, skipping
	if filters.Status != nil {
		filterParts = append(filterParts, fmt.Sprintf("status: %s", *filters.Status))
	}
	if filters.DateFrom != nil && filters.DateTo != nil {
		filterParts = append(filterParts, fmt.Sprintf("date range: %s to %s", 
			filters.DateFrom.Format("2006-01-02"), filters.DateTo.Format("2006-01-02")))
	}
	
	if len(filterParts) > 0 {
		description += " with filters: " + fmt.Sprintf("[%s]", filterParts[0])
		for i := 1; i < len(filterParts); i++ {
			description += ", " + filterParts[i]
		}
	}
	
	return description
}

func (s *ReportsService) calculateNextRunTime(frequency models.ReportFrequency, from time.Time) time.Time {
	switch frequency {
	case models.ReportFrequencyDaily:
		return from.AddDate(0, 0, 1)
	case models.ReportFrequencyWeekly:
		return from.AddDate(0, 0, 7)
	case models.ReportFrequencyMonthly:
		return from.AddDate(0, 1, 0)
	// ReportFrequencyQuarterly doesn't exist, using monthly as fallback
	case models.ReportFrequencyYearly:
		return from.AddDate(1, 0, 0)
	default:
		return from.AddDate(0, 1, 0) // Default to monthly
	}
}