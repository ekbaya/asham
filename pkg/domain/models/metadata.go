package models

import "time"

// Metadata for tracking and management
type Metadata struct {
	CreatedAt   time.Time
	LastUpdated time.Time
	Status      string
	Version     string
}
