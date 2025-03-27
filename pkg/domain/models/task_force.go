package models

import (
	"time"

	"github.com/google/uuid"
)

// TaskForce represents a task-specific group
type TaskForce struct {
	ID                  uuid.UUID
	Name                string
	Convenor            *Member
	NationalDeligations []*Member
	ParentTC            *TechnicalCommittee
	Task                string
	CreatedAt           time.Time
	CompletedAt         *time.Time
}
