package models

import (
	"time"

	"github.com/google/uuid"
)

// FDARSAction represents the action taken on a FDARS in it is not approved
type FDARSAction string

const (
	RESUBMIT_CD            FDARSAction = "RESUBMIT_CD"
	RESUBMIT_ENQUIRY_DRAFT FDARSAction = "RESUBMIT_ENQUIRY_DRAFT"
	RESUBMIT_FDARS         FDARSAction = "RESUBMIT_FDARS"
	TS                     FDARSAction = "Publish Technical Specification (TS)"
	PAS                    FDARSAction = "Publish Publicly Available Specification (PAS)"
	TR                     FDARSAction = "Publish Technical Report (TR)"
	Guide                  FDARSAction = "Publish Guide"
	CANCELLED              FDARSAction = "CANCELLED"
)

// AcceptanceCriteriaResult represents the result of checking project acceptance criteria
type AcceptanceCriteriaResult struct {
	CriteriaMet    bool    `json:"criteria_met"`
	AcceptanceRate float64 `json:"acceptance_rate"`
	RequiredRate   float64 `json:"required_rate"`
	TotalVotes     int64   `json:"total_votes"`
	AcceptedVotes  int64   `json:"accepted_votes"`
	Message        string  `json:"message"`
}

type Balloting struct {
	ID                 uuid.UUID   `json:"id"`
	ProjectID          string      `json:"project_id" binding:"required"`
	Project            *Project    `json:"project"`
	Votes              *[]Vote     `json:"nsb_submissions,omitempty"`
	StartDate          time.Time   `json:"start_date" binding:"required"`
	EndDate            time.Time   `json:"end_date" binding:"required"`
	Active             bool        `json:"active" gorm:"default:true"`
	Recommended        bool        `json:"recommended" gorm:"default:false"`
	RecommendedAt      *time.Time  `json:"recommended_at"`
	RecommendedByID    *string     `json:"recommended_by_id"`
	RecommendedBy      Member      `json:"recommended_by" gorm:"constraint:OnDelete:SET NULL"`
	VerifiedByID       *string     `json:"verified_by_id"`
	VerifiedBy         Member      `json:"verified_by" gorm:"constraint:OnDelete:SET NULL"`
	Approved           bool        `json:"approved" gorm:"default:false"`
	ApprovedByID       *string     `json:"approved_by_id"`
	ApprovedBy         Member      `json:"approved_by" gorm:"constraint:OnDelete:SET NULL"`
	ApprovedAt         *time.Time  `json:"approved_at"`
	NextCourseOfAction FDARSAction `json:"next_course_of_action"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}
