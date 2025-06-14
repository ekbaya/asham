package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourcePermission struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Method       string         `json:"method"`        // e.g. GET, POST
	Path         string         `json:"path"`          // e.g. /projects/approve
	PermissionID uuid.UUID      `json:"permission_id"` // FK to permissions table
	Permission   Permission     `json:"permission" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (permission *ResourcePermission) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate UUID if not set
	if permission.ID == uuid.Nil {
		permission.ID, err = uuid.NewRandom()
		if err != nil {
			return err
		}
	}
	return nil
}
