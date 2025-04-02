package models

import (
	"time"

	"github.com/google/uuid"
)

// DevelopmentTrack represents the possible development tracks
// @Description Defines the timeline for the development of a standard
type DevelopmentTrack string

const (
	// Default track (21 months)
	// @Description Standard development track with a duration of 21 months
	TrackDefault DevelopmentTrack = "DEFAULT"

	// International track (7 months)
	// @Description Accelerated development track aligning with international standards (7 months)
	TrackInternational DevelopmentTrack = "INTERNATIONAL"

	// Fast track (3 months)
	// @Description Expedited development track for urgent standardization (3 months)
	TrackFastTrack DevelopmentTrack = "FAST_TRACK"
)

// DraftStatus represents the status of the associated draft
// @Description Indicates the current stage of the draft in the standardization process
type DraftStatus string

const (
	// No draft associated
	// @Description No draft is currently linked to this standard
	DraftNone DraftStatus = "NONE"

	// Working Draft
	// @Description The standard is in the Working Draft (WD) stage
	DraftWD DraftStatus = "WD"

	// Committee Draft
	// @Description The standard is in the Committee Draft (CD) stage
	DraftCD DraftStatus = "CD"

	// Draft African Standard
	// @Description The standard has progressed to the Draft African Standard (DARS) stage
	DraftDARS DraftStatus = "DARS"
)

type Acceptance struct {
	ID              uuid.UUID      `json:"id"`
	ProjectID       string         `json:"project_id" binding:"required"`
	Project         *Project       `json:"-"`
	CirculationDate time.Time      `json:"circulation_date"`
	ClosingDate     time.Time      `json:"closing_date"`
	Submissions     *[]NSBResponse `json:"nsb_submissions,omitempty"`

	// Aggregate statistics
	// @Description Summary of response counts
	TotalResponses    int `json:"total_responses"`    // Total number of responses received
	AgreementCount    int `json:"agreement_count"`    // Number of NSBs that agreed with the proposal
	DisagreementCount int `json:"disagreement_count"` // Number of NSBs that disagreed with the proposal
	AbstentionCount   int `json:"abstention_count"`   // Number of NSBs that abstained

	// Approval status
	// @Description Indicates whether the proposal met the required approval criteria and was accepted
	ApprovalCriteriaMet bool `json:"approval_criteria_met"`
	IsApproved          bool `json:"is_approved"`

	// Draft status
	// @Description Indicates the current drafting stage of the proposal
	DraftStatus DraftStatus `json:"draft_status" gorm:"default:NONE"`

	// Expected date for first draft
	// @Description The anticipated date for the submission of the first draft, if applicable
	DraftExpectedDate *time.Time `json:"draft_expected_date,omitempty"`

	// Project registration details
	// @Description Specifies how the project is categorized within the standardization workflow
	IsPreliminaryWork bool `json:"is_preliminary_work"`
	IsActiveWork      bool `json:"is_active_work"`

	// Documents to consider
	// @Description Additional documents relevant to the proposal
	DocumentsToConsider *[]Document `json:"documents_to_consider" gorm:"many2many:documents_to_consider;"`

	// Development track selection
	// @Description The selected development track determining the timeline for standardization
	DevelopmentTrack DevelopmentTrack `json:"development_track"`

	// Target dates for various milestones
	// @Description Expected target dates for Committee Draft (CD), Draft African Standard (DARS), and Final Draft African Standard (FDARS)
	TargetDateCD    *time.Time `json:"target_date_cd,omitempty"`
	TargetDateDARS  *time.Time `json:"target_date_dars,omitempty"`
	TargetDateFDARS *time.Time `json:"target_date_fdars,omitempty"`

	// Secretariat details
	// @Description Information about the secretariat responsible for managing the project
	TCSecretaryID *string `json:"tc_secretary_id"`
	TCSecretary   *Member `json:"tc_secretary"`

	// Additional information
	// @Description Any other relevant information about the project
	OtherInformation string `json:"other_information,omitempty"`

	// Approval timestamps
	// @Description Records the date when the proposal was approved by the Standards Management Committee (SMC)
	SMCApprovalDate *time.Time `json:"smc_approval_date,omitempty"`

	// Standard fields
	// @Description Metadata fields tracking creation, updates, and deletions
	CreatedAt time.Time `json:"created_at"`
}
