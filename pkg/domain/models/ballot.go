package models

import (
	"time"

	"github.com/google/uuid"
)

type Balloting struct {
	ID        uuid.UUID `json:"id"`
	ProjectID string    `json:"project_id" binding:"required"`
	Project   *Project  `json:"-"`
	Votes     *[]Vote   `json:"nsb_submissions,omitempty"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
