package models

import (
	"time"

	"github.com/google/uuid"
)

// WorkingGroup represents a task-specific group
type WorkingGroup struct {
	ID          uuid.UUID           `json:"id" gorm:"primaryKey"`
	Name        string              `json:"name" gorm:"unique;not null"`
	ConvenorId  string              `json:"convenor_id" gorm:"index"`
	Convenor    *Member             `json:"convenor"`
	Experts     []*Member           `json:"experts" gorm:"many2many:working_group_experts;"`
	ParentTCID  string              `json:"parent_tc_id" gorm:"index"`
	ParentTC    *TechnicalCommittee `json:"parent_tc"`
	Task        string              `json:"task"`
	CreatedAt   time.Time           `json:"created_at" gorm:"autoCreateTime"`
	CompletedAt *time.Time          `json:"completed_at" gorm:"autoUpdateTime"`
}
