package models

import (
	"time"

	"github.com/google/uuid"
)

type Duration struct {
	Min time.Duration `json:"min,omitempty"` // Minimum duration
	Max time.Duration `json:"max,omitempty"` // Maximum duration (if applicable)
}

type Timeframe struct {
	Standard    *Duration `json:"standard,omitempty"`    // Default timeframe
	IS          *Duration `json:"is,omitempty"`          // Timeframe for International Standard (IS)
	Emergency   *Duration `json:"emergency,omitempty"`   // Timeframe for Emergency cases
	Description string    `json:"description,omitempty"` // Additional context
}

type Stage struct {
	ID           uuid.UUID  `json:"id"`
	Number       int        `json:"number"`
	Name         string     `json:"name"`
	DocumentName string     `json:"document_name"`
	Abbreviation string     `json:"abbreviation"`
	CreatedAt    time.Time  `json:"created_at"`
	Timeframe    *Timeframe `json:"timeframe,omitempty"` // Nullable timeframe field
}
