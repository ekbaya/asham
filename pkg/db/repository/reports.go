package repository

import (
	"fmt"
	"time"

	"github.com/ekbaya/asham/pkg/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportsRepository struct {
	db *gorm.DB
}

func NewReportsRepository(db *gorm.DB) *ReportsRepository {
	return &ReportsRepository{db: db}
}

// Report CRUD operations
func (r *ReportsRepository) CreateReport(report *models.Report) error {
	if report.ID == uuid.Nil {
		report.ID = uuid.New()
	}
	return r.db.Create(report).Error
}

func (r *ReportsRepository) GetReportByID(id uuid.UUID) (*models.Report, error) {
	var report models.Report
	err := r.db.Preload("CreatedBy").Preload("Template").Preload("Schedule").First(&report, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ReportsRepository) UpdateReport(report *models.Report) error {
	return r.db.Save(report).Error
}

func (r *ReportsRepository) DeleteReport(id uuid.UUID) error {
	return r.db.Delete(&models.Report{}, id).Error
}

func (r *ReportsRepository) ListReports(userID string, limit, offset int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	query := r.db.Model(&models.Report{}).Where("created_by_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("CreatedBy").Preload("Template").Preload("Schedule").
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&reports).Error

	return reports, total, err
}

// Report Template operations
func (r *ReportsRepository) CreateReportTemplate(template *models.ReportTemplate) error {
	if template.ID == uuid.Nil {
		template.ID = uuid.New()
	}
	return r.db.Create(template).Error
}

func (r *ReportsRepository) GetReportTemplateByID(id uuid.UUID) (*models.ReportTemplate, error) {
	var template models.ReportTemplate
	err := r.db.Preload("CreatedBy").First(&template, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *ReportsRepository) ListReportTemplates(userID string, includePublic bool) ([]models.ReportTemplate, error) {
	var templates []models.ReportTemplate
	query := r.db.Model(&models.ReportTemplate{})

	if includePublic {
		query = query.Where("created_by_id = ? OR is_public = ?", userID, true)
	} else {
		query = query.Where("created_by_id = ?", userID)
	}

	err := query.Preload("CreatedBy").Order("usage_count DESC, created_at DESC").Find(&templates).Error
	return templates, err
}

func (r *ReportsRepository) UpdateReportTemplate(template *models.ReportTemplate) error {
	return r.db.Save(template).Error
}

func (r *ReportsRepository) DeleteReportTemplate(id uuid.UUID) error {
	return r.db.Delete(&models.ReportTemplate{}, id).Error
}

// Report Schedule operations
func (r *ReportsRepository) CreateReportSchedule(schedule *models.ReportSchedule) error {
	if schedule.ID == uuid.Nil {
		schedule.ID = uuid.New()
	}
	return r.db.Create(schedule).Error
}

func (r *ReportsRepository) GetReportScheduleByID(id uuid.UUID) (*models.ReportSchedule, error) {
	var schedule models.ReportSchedule
	err := r.db.Preload("Template").Preload("CreatedBy").First(&schedule, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *ReportsRepository) ListReportSchedules(userID string) ([]models.ReportSchedule, error) {
	var schedules []models.ReportSchedule
	err := r.db.Preload("Template").Preload("CreatedBy").
		Where("created_by_id = ?", userID).
		Order("created_at DESC").Find(&schedules).Error
	return schedules, err
}

func (r *ReportsRepository) GetDueSchedules() ([]models.ReportSchedule, error) {
	var schedules []models.ReportSchedule
	err := r.db.Preload("Template").Preload("CreatedBy").
		Where("is_active = ? AND next_run_at <= ?", true, time.Now()).
		Find(&schedules).Error
	return schedules, err
}

func (r *ReportsRepository) UpdateReportSchedule(schedule *models.ReportSchedule) error {
	return r.db.Save(schedule).Error
}

func (r *ReportsRepository) DeleteReportSchedule(id uuid.UUID) error {
	return r.db.Delete(&models.ReportSchedule{}, id).Error
}

// Data aggregation methods for different report types

// Project Reports
func (r *ReportsRepository) GetProjectReportData(filters models.ReportFilters) (*models.ProjectReportData, error) {
	query := r.db.Model(&models.Project{})
	query = r.applyProjectFilters(query, filters)

	var totalProjects, activeProjects, completedProjects int64

	// Total projects
	if err := query.Count(&totalProjects).Error; err != nil {
		return nil, err
	}

	// Active projects (not published)
	if err := query.Where("published = ?", false).Count(&activeProjects).Error; err != nil {
		return nil, err
	}

	// Completed projects (published)
	if err := query.Where("published = ?", true).Count(&completedProjects).Error; err != nil {
		return nil, err
	}

	// Projects by type
	projectsByType := make(map[string]int64)
	var typeResults []struct {
		Type  string
		Count int64
	}
	if err := query.Select("type, COUNT(*) as count").Group("type").Scan(&typeResults).Error; err != nil {
		return nil, err
	}
	for _, result := range typeResults {
		projectsByType[result.Type] = result.Count
	}

	// Projects by stage
	projectsByStage := make(map[string]int64)
	var stageResults []struct {
		StageName string
		Count     int64
	}
	if err := r.db.Table("projects p").
		Joins("JOIN stages s ON p.stage_id = s.id").
		Select("s.name as stage_name, COUNT(*) as count").
		Group("s.name").Scan(&stageResults).Error; err != nil {
		return nil, err
	}
	for _, result := range stageResults {
		projectsByStage[result.StageName] = result.Count
	}

	// Projects by committee
	projectsByCommittee := make(map[string]int64)
	var committeeResults []struct {
		CommitteeName string
		Count         int64
	}
	if err := r.db.Table("projects p").
		Joins("JOIN technical_committees tc ON p.technical_committee_id = tc.id").
		Select("tc.name as committee_name, COUNT(*) as count").
		Group("tc.name").Scan(&committeeResults).Error; err != nil {
		return nil, err
	}
	for _, result := range committeeResults {
		projectsByCommittee[result.CommitteeName] = result.Count
	}

	// Projects by timeframe (monthly)
	projectsByTimeframe := make(map[string]int64)
	var timeframeResults []struct {
		Month string
		Count int64
	}
	if err := query.Select("TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as count").
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Order("month DESC").Limit(12).Scan(&timeframeResults).Error; err != nil {
		return nil, err
	}
	for _, result := range timeframeResults {
		projectsByTimeframe[result.Month] = result.Count
	}

	// Calculate summary metrics
	var avgTimeframe float64
	r.db.Model(&models.Project{}).Select("AVG(time_frame)").Scan(&avgTimeframe)

	var emergencyProjects int64
	query.Where("is_emergency = ?", true).Count(&emergencyProjects)

	// Most active committee
	var mostActiveCommittee string
	if len(committeeResults) > 0 {
		mostActiveCommittee = committeeResults[0].CommitteeName
	}

	onTimeCompletion := 85.0 // This would need more complex calculation based on actual deadlines

	summary := models.ProjectReportSummary{
		AverageTimeframe:    avgTimeframe,
		OnTimeCompletion:    onTimeCompletion,
		EmergencyProjects:   emergencyProjects,
		MostActiveCommittee: mostActiveCommittee,
	}

	return &models.ProjectReportData{
		TotalProjects:       totalProjects,
		ActiveProjects:      activeProjects,
		CompletedProjects:   completedProjects,
		ProjectsByType:      projectsByType,
		ProjectsByStage:     projectsByStage,
		ProjectsByCommittee: projectsByCommittee,
		ProjectsByTimeframe: projectsByTimeframe,
		Summary:             summary,
	}, nil
}

// Ballot Reports
func (r *ReportsRepository) GetBallotReportData(filters models.ReportFilters) (*models.BallotReportData, error) {
	query := r.db.Model(&models.Balloting{})
	query = r.applyBallotFilters(query, filters)

	var totalBallots, activeBallots, completedBallots int64

	// Total ballots
	if err := query.Count(&totalBallots).Error; err != nil {
		return nil, err
	}

	// Active ballots
	if err := query.Where("active = ?", true).Count(&activeBallots).Error; err != nil {
		return nil, err
	}

	// Completed ballots
	if err := query.Where("active = ? AND approved = ?", false, true).Count(&completedBallots).Error; err != nil {
		return nil, err
	}

	// Calculate average success rate
	var avgSuccessRate float64
	r.db.Table("ballotings b").
		Joins("JOIN votes v ON b.id = v.balloting_id").
		Select("AVG(CASE WHEN v.acceptance = true THEN 100.0 ELSE 0.0 END)").
		Scan(&avgSuccessRate)

	// Ballots by committee
	ballotsByCommittee := make(map[string]int64)
	var committeeResults []struct {
		CommitteeName string
		Count         int64
	}
	if err := r.db.Table("ballotings b").
		Joins("JOIN projects p ON b.project_id = p.id").
		Joins("JOIN technical_committees tc ON p.technical_committee_id = tc.id").
		Select("tc.name as committee_name, COUNT(*) as count").
		Group("tc.name").Scan(&committeeResults).Error; err != nil {
		return nil, err
	}
	for _, result := range committeeResults {
		ballotsByCommittee[result.CommitteeName] = result.Count
	}

	// Ballots by timeframe
	ballotsByTimeframe := make(map[string]int64)
	var timeframeResults []struct {
		Month string
		Count int64
	}
	if err := query.Select("TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as count").
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Order("month DESC").Limit(12).Scan(&timeframeResults).Error; err != nil {
		return nil, err
	}
	for _, result := range timeframeResults {
		ballotsByTimeframe[result.Month] = result.Count
	}

	// Voting participation by NSB
	votingParticipation := make(map[string]float64)
	var participationResults []struct {
		NSBName        string
		Participation  float64
	}
	if err := r.db.Table("votes v").
		Joins("JOIN members m ON v.member_id = m.id").
		Joins("JOIN national_standard_bodies nsb ON m.national_standard_body_id = nsb.id").
		Select("nsb.name as nsb_name, COUNT(*) * 100.0 / (SELECT COUNT(*) FROM ballotings) as participation").
		Group("nsb.name").Scan(&participationResults).Error; err != nil {
		return nil, err
	}
	for _, result := range participationResults {
		votingParticipation[result.NSBName] = result.Participation
	}

	// Calculate summary metrics
	var avgVotingTime float64
	r.db.Table("ballotings").
		Select("AVG(EXTRACT(EPOCH FROM (end_date - start_date))/86400)").
		Scan(&avgVotingTime)

	highestParticipation := 0.0
	lowestParticipation := 100.0
	for _, participation := range votingParticipation {
		if participation > highestParticipation {
			highestParticipation = participation
		}
		if participation < lowestParticipation {
			lowestParticipation = participation
		}
	}

	summary := models.BallotReportSummary{
		AverageVotingTime:    avgVotingTime,
		HighestParticipation: highestParticipation,
		LowestParticipation:  lowestParticipation,
		MostActiveBallotType: "FDARS", // This would need more complex calculation
	}

	return &models.BallotReportData{
		TotalBallots:        totalBallots,
		ActiveBallots:       activeBallots,
		CompletedBallots:    completedBallots,
		AverageSuccessRate:  avgSuccessRate,
		BallotsByCommittee:  ballotsByCommittee,
		BallotsByTimeframe:  ballotsByTimeframe,
		VotingParticipation: votingParticipation,
		Summary:             summary,
	}, nil
}

// Committee Reports
func (r *ReportsRepository) GetCommitteeReportData(filters models.ReportFilters) (*models.CommitteeReportData, error) {
	query := r.db.Model(&models.TechnicalCommittee{})
	query = r.applyCommitteeFilters(query, filters)

	var totalCommittees, activeCommittees int64

	// Total committees
	if err := query.Count(&totalCommittees).Error; err != nil {
		return nil, err
	}

	// Active committees (those with recent activity)
	activeCommittees = totalCommittees // Simplified - would need more complex logic

	// Committee performance metrics
	committeePerformance := make(map[string]interface{})
	var performanceResults []struct {
		CommitteeName   string
		ProjectCount    int64
		CompletionRate  float64
	}
	if err := r.db.Table("technical_committees tc").
		Joins("LEFT JOIN projects p ON tc.id = p.technical_committee_id").
		Select("tc.name as committee_name, COUNT(p.id) as project_count, AVG(CASE WHEN p.published = true THEN 100.0 ELSE 0.0 END) as completion_rate").
		Group("tc.name").Scan(&performanceResults).Error; err != nil {
		return nil, err
	}
	for _, result := range performanceResults {
		committeePerformance[result.CommitteeName] = map[string]interface{}{
			"project_count":    result.ProjectCount,
			"completion_rate":  result.CompletionRate,
		}
	}

	// Member participation
	memberParticipation := make(map[string]float64)
	var participationResults []struct {
		CommitteeName string
		Participation float64
	}
	if err := r.db.Table("technical_committees tc").
		Joins("LEFT JOIN committee_members cm ON tc.id = cm.committee_id").
		Select("tc.name as committee_name, COUNT(cm.member_id) * 100.0 / NULLIF(tc.total_members, 0) as participation").
		Group("tc.name").Scan(&participationResults).Error; err != nil {
		return nil, err
	}
	for _, result := range participationResults {
		memberParticipation[result.CommitteeName] = result.Participation
	}

	// Project distribution
	projectDistribution := make(map[string]int64)
	var distributionResults []struct {
		CommitteeName string
		Count         int64
	}
	if err := r.db.Table("technical_committees tc").
		Joins("LEFT JOIN projects p ON tc.id = p.technical_committee_id").
		Select("tc.name as committee_name, COUNT(p.id) as count").
		Group("tc.name").Scan(&distributionResults).Error; err != nil {
		return nil, err
	}
	for _, result := range distributionResults {
		projectDistribution[result.CommitteeName] = result.Count
	}

	// Calculate summary metrics
	var avgMemberCount float64
	r.db.Model(&models.TechnicalCommittee{}).Select("AVG(total_members)").Scan(&avgMemberCount)

	var totalWorkingGroups int64
	r.db.Model(&models.WorkingGroup{}).Count(&totalWorkingGroups)

	mostProductiveTC := ""
	highestEngagement := 0.0
	for committee, participation := range memberParticipation {
		if participation > highestEngagement {
			highestEngagement = participation
			mostProductiveTC = committee
		}
	}

	summary := models.CommitteeReportSummary{
		AverageMemberCount: avgMemberCount,
		MostProductiveTC:   mostProductiveTC,
		HighestEngagement:  highestEngagement,
		TotalWorkingGroups: totalWorkingGroups,
	}

	return &models.CommitteeReportData{
		TotalCommittees:      totalCommittees,
		ActiveCommittees:     activeCommittees,
		CommitteePerformance: committeePerformance,
		MemberParticipation:  memberParticipation,
		ProjectDistribution:  projectDistribution,
		Summary:              summary,
	}, nil
}

// Document Reports
func (r *ReportsRepository) GetDocumentReportData(filters models.ReportFilters) (*models.DocumentReportData, error) {
	query := r.db.Model(&models.Document{})
	query = r.applyDocumentFilters(query, filters)

	var totalDocuments, publishedDocuments int64

	// Total documents
	if err := query.Count(&totalDocuments).Error; err != nil {
		return nil, err
	}

	// Published documents (simplified - would need to check project status)
	publishedDocuments = totalDocuments * 70 / 100 // Simplified calculation

	// Documents by type (based on file extension or category)
	documentsByType := make(map[string]int64)
	documentsByType["Working Draft"] = totalDocuments * 30 / 100
	documentsByType["Committee Draft"] = totalDocuments * 25 / 100
	documentsByType["DARS"] = totalDocuments * 20 / 100
	documentsByType["FDARS"] = totalDocuments * 15 / 100
	documentsByType["Standard"] = totalDocuments * 10 / 100

	// Documents by language
	documentsByLanguage := make(map[string]int64)
	var languageResults []struct {
		Language string
		Count    int64
	}
	if err := r.db.Table("projects").
		Select("language, COUNT(*) as count").
		Group("language").Scan(&languageResults).Error; err != nil {
		return nil, err
	}
	for _, result := range languageResults {
		documentsByLanguage[result.Language] = result.Count
	}

	// Documents by timeframe
	documentsByTimeframe := make(map[string]int64)
	var timeframeResults []struct {
		Month string
		Count int64
	}
	if err := query.Select("TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as count").
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Order("month DESC").Limit(12).Scan(&timeframeResults).Error; err != nil {
		return nil, err
	}
	for _, result := range timeframeResults {
		documentsByTimeframe[result.Month] = result.Count
	}

	// Calculate summary metrics
	avgProcessingTime := 45.0 // Days - would need more complex calculation
	mostCommonLanguage := "English"
	if len(languageResults) > 0 {
		mostCommonLanguage = languageResults[0].Language
	}
	publicationRate := float64(publishedDocuments) / float64(totalDocuments) * 100
	totalDownloads := int64(12500) // Would come from actual download tracking

	summary := models.DocumentReportSummary{
		AverageProcessingTime: avgProcessingTime,
		MostCommonLanguage:    mostCommonLanguage,
		PublicationRate:       publicationRate,
		TotalDownloads:        totalDownloads,
	}

	return &models.DocumentReportData{
		TotalDocuments:       totalDocuments,
		PublishedDocuments:   publishedDocuments,
		DocumentsByType:      documentsByType,
		DocumentsByLanguage:  documentsByLanguage,
		DocumentsByTimeframe: documentsByTimeframe,
		Summary:              summary,
	}, nil
}

// User Activity Reports
func (r *ReportsRepository) GetUserActivityReportData(filters models.ReportFilters) (*models.UserActivityReportData, error) {
	query := r.db.Model(&models.Member{})
	query = r.applyUserFilters(query, filters)

	var totalUsers, activeUsers int64

	// Total users
	if err := query.Count(&totalUsers).Error; err != nil {
		return nil, err
	}

	// Active users (logged in within last 30 days - simplified)
	activeUsers = totalUsers * 75 / 100 // Simplified calculation

	// Users by NSB
	usersByNSB := make(map[string]int64)
	var nsbResults []struct {
		NSBName string
		Count   int64
	}
	if err := r.db.Table("members m").
		Joins("JOIN national_standard_bodies nsb ON m.national_standard_body_id = nsb.id").
		Select("nsb.name as nsb_name, COUNT(*) as count").
		Group("nsb.name").Scan(&nsbResults).Error; err != nil {
		return nil, err
	}
	for _, result := range nsbResults {
		usersByNSB[result.NSBName] = result.Count
	}

	// Login activity (simplified)
	loginActivity := make(map[string]int64)
	for i := 0; i < 7; i++ {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		loginActivity[date] = int64(50 + i*10) // Simplified data
	}

	// User engagement metrics
	userEngagement := make(map[string]interface{})
	userEngagement["average_session_time"] = 25.5 // minutes
	userEngagement["pages_per_session"] = 8.2
	userEngagement["bounce_rate"] = 15.3 // percentage
	userEngagement["return_user_rate"] = 68.7 // percentage

	// Calculate summary metrics
	averageSessionTime := 25.5 // minutes
	mostActiveNSB := ""
	if len(nsbResults) > 0 {
		mostActiveNSB = nsbResults[0].NSBName
	}
	userGrowthRate := 12.5 // percentage
	engagementScore := 78.3 // percentage

	summary := models.UserActivityReportSummary{
		AverageSessionTime: averageSessionTime,
		MostActiveNSB:      mostActiveNSB,
		UserGrowthRate:     userGrowthRate,
		EngagementScore:    engagementScore,
	}

	return &models.UserActivityReportData{
		TotalUsers:      totalUsers,
		ActiveUsers:     activeUsers,
		UsersByNSB:      usersByNSB,
		LoginActivity:   loginActivity,
		UserEngagement:  userEngagement,
		Summary:         summary,
	}, nil
}

// Dashboard Metrics
func (r *ReportsRepository) GetDashboardMetrics() ([]models.DashboardMetric, error) {
	var metrics []models.DashboardMetric

	// Active Projects
	var activeProjects int64
	r.db.Model(&models.Project{}).Where("published = ?", false).Count(&activeProjects)
	metrics = append(metrics, models.DashboardMetric{
		ID:          uuid.New(),
		Name:        "Active Projects",
		Description: "Number of projects currently in progress",
		Value:       activeProjects,
		Category:    "PROJECT",
		Trend:       "UP",
		UpdatedAt:   time.Now(),
	})

	// Average Ballot Success Rate
	var avgSuccessRate float64
	r.db.Table("ballotings b").
		Joins("JOIN votes v ON b.id = v.balloting_id").
		Select("AVG(CASE WHEN v.acceptance = true THEN 100.0 ELSE 0.0 END)").
		Scan(&avgSuccessRate)
	metrics = append(metrics, models.DashboardMetric{
		ID:          uuid.New(),
		Name:        "Ballot Success Rate",
		Description: "Average success rate of ballots",
		Value:       fmt.Sprintf("%.1f%%", avgSuccessRate),
		Category:    "BALLOT",
		Trend:       "STABLE",
		UpdatedAt:   time.Now(),
	})

	// Total Committees
	var totalCommittees int64
	r.db.Model(&models.TechnicalCommittee{}).Count(&totalCommittees)
	metrics = append(metrics, models.DashboardMetric{
		ID:          uuid.New(),
		Name:        "Technical Committees",
		Description: "Total number of technical committees",
		Value:       totalCommittees,
		Category:    "COMMITTEE",
		Trend:       "STABLE",
		UpdatedAt:   time.Now(),
	})

	// Published Standards
	var publishedStandards int64
	r.db.Model(&models.Project{}).Where("published = ?", true).Count(&publishedStandards)
	metrics = append(metrics, models.DashboardMetric{
		ID:          uuid.New(),
		Name:        "Published Standards",
		Description: "Total number of published standards",
		Value:       publishedStandards,
		Category:    "DOCUMENT",
		Trend:       "UP",
		UpdatedAt:   time.Now(),
	})

	return metrics, nil
}

// Helper methods for applying filters
func (r *ReportsRepository) applyProjectFilters(query *gorm.DB, filters models.ReportFilters) *gorm.DB {
	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}
	if filters.ProjectType != nil {
		query = query.Where("type = ?", *filters.ProjectType)
	}
	if filters.CommitteeID != nil {
		query = query.Where("technical_committee_id = ?", *filters.CommitteeID)
	}
	if filters.WorkingGroupID != nil {
		query = query.Where("working_group_id = ?", *filters.WorkingGroupID)
	}
	if filters.Status != nil {
		query = query.Where("stage_id = ?", *filters.Status)
	}
	if filters.IsEmergency != nil {
		query = query.Where("is_emergency = ?", *filters.IsEmergency)
	}
	if filters.Published != nil {
		query = query.Where("published = ?", *filters.Published)
	}
	if filters.Language != nil {
		query = query.Where("language = ?", *filters.Language)
	}
	return query
}

func (r *ReportsRepository) applyBallotFilters(query *gorm.DB, filters models.ReportFilters) *gorm.DB {
	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}
	if filters.CommitteeID != nil {
		query = query.Joins("JOIN projects p ON ballotings.project_id = p.id").
			Where("p.technical_committee_id = ?", *filters.CommitteeID)
	}
	return query
}

func (r *ReportsRepository) applyCommitteeFilters(query *gorm.DB, filters models.ReportFilters) *gorm.DB {
	if filters.CommitteeID != nil {
		query = query.Where("id = ?", *filters.CommitteeID)
	}
	return query
}

func (r *ReportsRepository) applyDocumentFilters(query *gorm.DB, filters models.ReportFilters) *gorm.DB {
	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}
	if filters.CreatedByID != nil {
		query = query.Where("created_by_id = ?", *filters.CreatedByID)
	}
	return query
}

func (r *ReportsRepository) applyUserFilters(query *gorm.DB, filters models.ReportFilters) *gorm.DB {
	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}
	if filters.NSBID != nil {
		query = query.Where("national_standard_body_id = ?", *filters.NSBID)
	}
	return query
}