package models

import (
	"time"

	"github.com/google/uuid"
)

// Member represents a member state in the organization
type Member struct {
	ID                     uuid.UUID             `json:"id"`
	Email                  string                `json:"email" gorm:"index;unique" binding:"required"`
	FirstName              string                `json:"first_name" binding:"required"`
	LastName               string                `json:"last_name" binding:"required"`
	PhotoUrl               string                `json:"photo_url"`
	NationalStandardBodyID string                `json:"nsb_id" binding:"required"`
	NationalStandardBody   *NationalStandardBody `json:"nsb"`
	Password               string                `json:"password,omitempty" gorm:"-" binding:"required"` // Exclude from output; custom handling for input
	HashedPassword         string                `json:"-" gorm:"column:password"`
	CreatedAt              time.Time
}
