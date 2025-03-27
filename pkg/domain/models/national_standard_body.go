package models

import (
	"time"

	"github.com/google/uuid"
)

// NSB
type NationalStandardBody struct {
	ID          uuid.UUID
	Name        string
	MemberState string
	Members     *[]Member
	CreatedAt   time.Time
}
