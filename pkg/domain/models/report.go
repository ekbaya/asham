package models

import (
	"time"

	"github.com/google/uuid"
)

// ReportType defines the type of report
type ReportType string

const (
	ReportTypeProject           ReportType = "PROJECT"
	ReportTypeBallot            ReportType = "BALLOT"
	ReportTypeCommittee         ReportType = "COMMITTEE"
	ReportTypeDocument          ReportType = "DOCUMENT"
	ReportTypeUserActivity      ReportType = "USER_ACTIVITY"
	ReportTypePerformance       ReportType = "PERFORMANCE"
	ReportTypeDashboard         ReportType = "DASHBOARD"
	ReportTypeCustom            ReportType = "CUSTOM"
)

// ReportFormat defines the export format for reports
type ReportFormat string

const (
	ReportFormatPDF   ReportFormat = "PDF"
	ReportFormatExcel ReportFormat = "EXCEL"
	ReportFormatCSV   ReportFormat = "CSV"
	ReportFormatJSON  ReportFormat = "JSON"
)

// ReportStatus defines the status of report generation
type ReportStatus string

const (
	ReportStatusPending    ReportStatus = "PENDING"
	ReportStatusProcessing ReportStatus = "PROCESSING"
	ReportStatusCompleted  ReportStatus = "COMPLETED"
	ReportStatusFailed     ReportStatus = "FAILED"
)

// ReportFrequency defines how often scheduled reports are generated
type ReportFrequency string

const (
	ReportFrequencyDaily   ReportFrequency = "DAILY"
	ReportFrequencyWeekly  ReportFrequency = "WEEKLY"
	ReportFrequencyMonthly ReportFrequency = "MONTHLY"
	ReportFrequencyYearly  ReportFrequency = "YEARLY"
)

// ReportFilters represents the filtering criteria for reports
type ReportFilters struct {
	DateFrom         *time.Time `json:"date_from,omitempty"`
	DateTo           *time.Time `json:"date_to,omitempty"`
	ProjectType      *string    `json:"project_type,omitempty"`
	CommitteeID      *string    `json:"committee_id,omitempty"`
	WorkingGroupID   *string    `json:"working_group_id,omitempty"`
	NSBID            *string    `json:"nsb_id,omitempty"`
	Status           *string    `json:"status,omitempty"`
	Sector           *string    `json:"sector,omitempty"`
	Language         *string    `json:"language,omitempty"`
	IsEmergency      *bool      `json:"is_emergency,omitempty"`
	Published        *bool      `json:"published,omitempty"`
	StageID          *string    `json:"stage_id,omitempty"`
	CreatedByID      *string    `json:"created_by_id,omitempty"`
}

// Report represents a generated report
type Report struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description"`
	Type        ReportType     `json:"type" binding:"required"`
	Format      ReportFormat   `json:"format" binding:"required"`
	Status      ReportStatus   `json:"status" gorm:"default:PENDING"`
	Filters     ReportFilters  `json:"filters" gorm:"type:jsonb"`
	Data        interface{}    `json:"data,omitempty" gorm:"type:jsonb"`
	FileURL     string         `json:"file_url,omitempty"`
	FileSize    int64          `json:"file_size,omitempty"`
	GeneratedAt *time.Time     `json:"generated_at,omitempty"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	CreatedByID string         `json:"created_by_id" binding:"required"`
	CreatedBy   *Member        `json:"created_by"`
	TemplateID  *uuid.UUID     `json:"template_id,omitempty"`
	Template    *ReportTemplate `json:"template,omitempty"`
	ScheduleID  *uuid.UUID     `json:"schedule_id,omitempty"`
	Schedule    *ReportSchedule `json:"schedule,omitempty"`
	ErrorMessage string        `json:"error_message,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ReportTemplate represents a saved report configuration
type ReportTemplate struct {
	ID          uuid.UUID     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string        `json:"name" binding:"required"`
	Description string        `json:"description"`
	Type        ReportType    `json:"type" binding:"required"`
	Filters     ReportFilters `json:"filters" gorm:"type:jsonb"`
	Columns     []string      `json:"columns" gorm:"type:jsonb"`
	SortBy      string        `json:"sort_by,omitempty"`
	SortOrder   string        `json:"sort_order,omitempty" gorm:"default:DESC"`
	IsPublic    bool          `json:"is_public" gorm:"default:false"`
	CreatedByID string        `json:"created_by_id" binding:"required"`
	CreatedBy   *Member       `json:"created_by"`
	UsageCount  int64         `json:"usage_count" gorm:"default:0"`
	LastUsedAt  *time.Time    `json:"last_used_at,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// ReportSchedule represents a scheduled report configuration
type ReportSchedule struct {
	ID          uuid.UUID       `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string          `json:"name" binding:"required"`
	Description string          `json:"description"`
	TemplateID  uuid.UUID       `json:"template_id" binding:"required"`
	Template    *ReportTemplate `json:"template"`
	Frequency   ReportFrequency `json:"frequency" binding:"required"`
	Format      ReportFormat    `json:"format" binding:"required"`
	Recipients  []string        `json:"recipients" gorm:"type:jsonb" binding:"required"`
	NextRunAt   time.Time       `json:"next_run_at"`
	LastRunAt   *time.Time      `json:"last_run_at,omitempty"`
	IsActive    bool            `json:"is_active" gorm:"default:true"`
	CreatedByID string          `json:"created_by_id" binding:"required"`
	CreatedBy   *Member         `json:"created_by"`
	RunCount    int64           `json:"run_count" gorm:"default:0"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// DashboardMetric represents a key performance indicator
type DashboardMetric struct {
	ID          uuid.UUID   `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description"`
	Value       interface{} `json:"value" gorm:"type:jsonb"`
	PreviousValue interface{} `json:"previous_value,omitempty" gorm:"type:jsonb"`
	ChangePercent float64   `json:"change_percent,omitempty"`
	Trend       string      `json:"trend,omitempty"` // UP, DOWN, STABLE
	Category    string      `json:"category"`        // PROJECT, BALLOT, COMMITTEE, etc.
	UpdatedAt   time.Time   `json:"updated_at"`
}

// DashboardWidget represents a dashboard visualization component
type DashboardWidget struct {
	ID          uuid.UUID   `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Title       string      `json:"title" binding:"required"`
	Type        string      `json:"type" binding:"required"` // CHART, METRIC, TABLE, etc.
	ChartType   string      `json:"chart_type,omitempty"`    // BAR, PIE, LINE, etc.
	Data        interface{} `json:"data" gorm:"type:jsonb"`
	Config      interface{} `json:"config,omitempty" gorm:"type:jsonb"`
	Position    int         `json:"position" gorm:"default:0"`
	Size        string      `json:"size" gorm:"default:MEDIUM"` // SMALL, MEDIUM, LARGE
	IsVisible   bool        `json:"is_visible" gorm:"default:true"`
	CreatedByID string      `json:"created_by_id" binding:"required"`
	CreatedBy   *Member     `json:"created_by"`
	UpdatedAt   time.Time   `json:"updated_at"`
	CreatedAt   time.Time   `json:"created_at"`
}

// ReportData represents structured data for different report types
type ProjectReportData struct {
	TotalProjects       int64                    `json:"total_projects"`
	ActiveProjects      int64                    `json:"active_projects"`
	CompletedProjects   int64                    `json:"completed_projects"`
	ProjectsByType      map[string]int64         `json:"projects_by_type"`
	ProjectsByStage     map[string]int64         `json:"projects_by_stage"`
	ProjectsByCommittee map[string]int64         `json:"projects_by_committee"`
	ProjectsByTimeframe map[string]int64         `json:"projects_by_timeframe"`
	Projects            []Project                `json:"projects,omitempty"`
	Summary             ProjectReportSummary     `json:"summary"`
}

type BallotReportData struct {
	TotalBallots        int64                `json:"total_ballots"`
	ActiveBallots       int64                `json:"active_ballots"`
	CompletedBallots    int64                `json:"completed_ballots"`
	AverageSuccessRate  float64              `json:"average_success_rate"`
	BallotsByCommittee  map[string]int64     `json:"ballots_by_committee"`
	BallotsByTimeframe  map[string]int64     `json:"ballots_by_timeframe"`
	VotingParticipation map[string]float64   `json:"voting_participation"`
	Ballots             []Balloting          `json:"ballots,omitempty"`
	Summary             BallotReportSummary  `json:"summary"`
}

type CommitteeReportData struct {
	TotalCommittees       int64                     `json:"total_committees"`
	ActiveCommittees      int64                     `json:"active_committees"`
	CommitteePerformance  map[string]interface{}    `json:"committee_performance"`
	MemberParticipation   map[string]float64        `json:"member_participation"`
	ProjectDistribution   map[string]int64          `json:"project_distribution"`
	Committees            []TechnicalCommittee      `json:"committees,omitempty"`
	Summary               CommitteeReportSummary    `json:"summary"`
}

type DocumentReportData struct {
	TotalDocuments      int64                   `json:"total_documents"`
	PublishedDocuments  int64                   `json:"published_documents"`
	DocumentsByType     map[string]int64        `json:"documents_by_type"`
	DocumentsByLanguage map[string]int64        `json:"documents_by_language"`
	DocumentsByTimeframe map[string]int64       `json:"documents_by_timeframe"`
	Documents           []Document              `json:"documents,omitempty"`
	Summary             DocumentReportSummary   `json:"summary"`
}

type UserActivityReportData struct {
	TotalUsers          int64                      `json:"total_users"`
	ActiveUsers         int64                      `json:"active_users"`
	UsersByNSB          map[string]int64           `json:"users_by_nsb"`
	LoginActivity       map[string]int64           `json:"login_activity"`
	UserEngagement      map[string]interface{}     `json:"user_engagement"`
	Users               []Member                   `json:"users,omitempty"`
	Summary             UserActivityReportSummary  `json:"summary"`
}

// Report summary structures
type ProjectReportSummary struct {
	AverageTimeframe    float64 `json:"average_timeframe"`
	OnTimeCompletion    float64 `json:"on_time_completion"`
	EmergencyProjects   int64   `json:"emergency_projects"`
	MostActiveCommittee string  `json:"most_active_committee"`
}

type BallotReportSummary struct {
	AverageVotingTime     float64 `json:"average_voting_time"`
	HighestParticipation  float64 `json:"highest_participation"`
	LowestParticipation   float64 `json:"lowest_participation"`
	MostActiveBallotType  string  `json:"most_active_ballot_type"`
}

type CommitteeReportSummary struct {
	AverageMemberCount    float64 `json:"average_member_count"`
	MostProductiveTC      string  `json:"most_productive_tc"`
	HighestEngagement     float64 `json:"highest_engagement"`
	TotalWorkingGroups    int64   `json:"total_working_groups"`
}

type DocumentReportSummary struct {
	AverageProcessingTime float64 `json:"average_processing_time"`
	MostCommonLanguage    string  `json:"most_common_language"`
	PublicationRate       float64 `json:"publication_rate"`
	TotalDownloads        int64   `json:"total_downloads"`
}

type UserActivityReportSummary struct {
	AverageSessionTime    float64 `json:"average_session_time"`
	MostActiveNSB         string  `json:"most_active_nsb"`
	UserGrowthRate        float64 `json:"user_growth_rate"`
	EngagementScore       float64 `json:"engagement_score"`
}