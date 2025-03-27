package models

import (
	"time"

	"github.com/google/uuid"
)

type MemberState struct {
	ID        uuid.UUID               `json:"id"`
	Name      string                  `json:"name" binding:"required" gorm:"unique"`
	Code      string                  `json:"code" binding:"required" gorm:"index;unique"`
	NSBs      *[]NationalStandardBody `json:"nsbs"`
	CreatedAt time.Time
}
