package models

import (
	"time"

	"github.com/google/uuid"
)

// TaskForce represents a task-specific group
type TaskForce struct {
	ID                  uuid.UUID           `json:"id" gorm:"primaryKey"`
	Name                string              `json:"name" gorm:"unique;not null"`
	ConvenorId          string              `json:"convenor_id" gorm:"index"`
	Convenor            *Member             `json:"convenor"`
	NationalDeligations []*Member           `gorm:"many2many:national_deligations;"`
	ParentTCID          string              `json:"parent_tc_id" gorm:"index"`
	ParentTC            *TechnicalCommittee `json:"parent_tc"`
	Task                string              `json:"task"`
	CreatedAt           time.Time           `json:"created_at" gorm:"autoCreateTime"`
	CompletedAt         *time.Time          `json:"completed_at" gorm:"autoUpdateTime"`
}
