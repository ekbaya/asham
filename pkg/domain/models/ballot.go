package models

import (
	"time"

	"github.com/google/uuid"
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
	ID        uuid.UUID `json:"id"`
	ProjectID string    `json:"project_id" binding:"required"`
	Project   *Project  `json:"-"`
	Votes     *[]Vote   `json:"nsb_submissions,omitempty"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
