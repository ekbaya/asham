package models

import (
	"time"

	"github.com/google/uuid"
)

// WorkingGroup represents a task-specific group
type WorkingGroup struct {
	ID          uuid.UUID
	Name        string
	ConvenorId  uuid.UUID
	Convenor    *Member
	Experts     []*Member `json:"experts" gorm:"many2many:working_group_experts;"`
	ParentTCID  uuid.UUID
	ParentTC    *TechnicalCommittee
	Task        string
	CreatedAt   time.Time
	CompletedAt *time.Time
}
