package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Standard struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title       string         `json:"title"`
	Content     string         `json:"content" gorm:"type:jsonb"`
	Version     int            `json:"version"`
	UpdatedByID string         `json:"updated_by_id"`
	UpdatedBy   *Member        `json:"updated_by"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type StandardVersion struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	StandardID  uuid.UUID `gorm:"type:uuid;index"`
	Content     string    `json:"content" gorm:"type:jsonb"`
	Version     int       `json:"version"`
	SavedByByID string    `json:"updated_by_id"`
	SavedBy     *Member   `json:"saved_by"`
	SavedAt     time.Time `json:"saved_at" gorm:"autoCreateTime"`
}
