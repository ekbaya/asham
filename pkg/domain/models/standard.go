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
	UpdatedByID string         `json:"updated_by_id" gorm:"index"`
	UpdatedBy   *Member        `json:"updated_by" gorm:"foreignKey:UpdatedByID"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
	ProjectID   string         `json:"project_id" gorm:"index"`
	Project     *Project       `json:"project" gorm:"foreignKey:ProjectID"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type StandardVersion struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	StandardID uuid.UUID `gorm:"type:uuid;index"`
	Standard   *Standard `gorm:"foreignKey:StandardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Content    string    `json:"content" gorm:"type:jsonb"`
	Version    int       `json:"version"`
	SavedByID  string    `json:"saved_by_id" gorm:"index"`
	SavedBy    *Member   `gorm:"foreignKey:SavedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	SavedAt    time.Time `json:"saved_at" gorm:"autoCreateTime"`
}

type StandardAuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	StandardID uuid.UUID
	Version    int
	ChangedBy  string
	ChangeDiff string `gorm:"type:text"` // Can be a long string of ASCII diff
	CreatedAt  time.Time
}
