package models

import (
	"time"

	"github.com/google/uuid"
)

// DARSStatus represents the status of a Draft African Standard
type DARSStatus string

const (
	DARSUnderReview DARSStatus = "UNDER_REVIEW" // DARS is under public review
	DARSApproved    DARSStatus = "APPROVED"     // DARS is approved for Balloting Stage
	DARSRejected    DARSStatus = "REJECTED"     // DARS is rejected
)

type DARS struct {
	ID        uuid.UUID `json:"id"`
	ProjectID string    `json:"project_id" binding:"required"`
	Project   *Project  `json:"-"`
	// Public review information
	PublicReviewStartDate time.Time `json:"public_review_start_date"` // Start date of public review
	PublicReviewEndDate   time.Time `json:"public_review_end_date"`   // End date of public review (60 days after start)

	// WTO notification (optional)
	WTONotificationDate *time.Time              `json:"wto_notification_date,omitempty"` // Date of WTO notification
	Submissions         *[]NationalConsultation `json:"nsb_submissions,omitempty"`

	// Status and unresolved issues
	Status           DARSStatus `json:"status" gorm:"default:UNDER_REVIEW"` // Current status of the DARS
	UnresolvedIssues string     `json:"unresolved_issues,omitempty"`        // Description of unresolved issues

	// Decision for the next stage
	MoveToBalloting        bool   `json:"move_to_balloting" gorm:"default:false"` // Whether to move to Balloting Stage
	AlternativeDeliverable string `json:"alternative_deliverable,omitempty"`      // If another deliverable is needed

	DARSTCSecretaryID *string `json:"dars_tc_secretary_id" gorm:"column:dars_tc_secretary_id"`
	DARSTCSecretary   *Member `json:"dars_tc_secretary"`

	// Standard timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
