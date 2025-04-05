package models

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	ID          uuid.UUID  `json:"id"`
	ProjectID   string     `json:"project_id" binding:"required"`
	Project     *Project   `json:"project"`
	MemberID    string     `json:"member_id"`
	Member      *Member    `json:"national_secretary"`
	BallotingID uuid.UUID  `json:"-"`
	Balloting   *Balloting `json:"-"`
	Acceptance  bool       `json:"acceptance" gorm:"default:false"`
	Comment     string     `json:"comment"`
	CreatedAt   time.Time  `json:"created_at"`
}
