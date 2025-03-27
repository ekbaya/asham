package domain

import (
	"time"

	"github.com/google/uuid"
)

// Member represents a member state in the organization
type Member struct {
	ID                   uuid.UUID
	Email                string
	FirstName            string
	LastName             string
	PhotoUrl             string
	CountryCode          string
	NationalStandardBody string
	CreatedAt            time.Time
}
