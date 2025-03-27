package models

import (
	"time"

	"github.com/google/uuid"
)

// Member represents a member state in the organization
type Member struct {
	ID                     uuid.UUID
	Email                  string `gorm:"index;unique"`
	FirstName              string
	LastName               string
	PhotoUrl               string
	CountryCode            string
	NationalStandardBodyID string
	NationalStandardBody   *NationalStandardBody
	HashedPassword         string
	CreatedAt              time.Time
}
