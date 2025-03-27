package domain

import (
	"time"

	"github.com/google/uuid"
)

// WorkingGroup represents a task-specific group
type WorkingGroup struct {
	ID          uuid.UUID
	Name        string
	Convenor    *Member
	Experts     []*Member
	ParentTC    *TechnicalCommittee
	Task        string
	CreatedAt   time.Time
	CompletedAt *time.Time
}
