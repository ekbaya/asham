package models

import (
	"time"

	"github.com/google/uuid"
)

// TaskForce represents a task-specific group
type TaskForce struct {
	ID                  uuid.UUID
	Name                string
	ConvenorId          uuid.UUID
	Convenor            *Member
	NationalDeligations []*Member `gorm:"many2many:national_deligations;"`
	ParentTCID          uuid.UUID
	ParentTC            *TechnicalCommittee
	Task                string
	CreatedAt           time.Time
	CompletedAt         *time.Time
}
