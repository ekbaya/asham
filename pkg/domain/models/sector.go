package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Sector struct {
	ID        uuid.UUID      `json:"id" gorm:"primary_key"`
	Title     string         `json:"title" binding:"required" gorm:"unique"`
	Slug      string         `json:"slug" gorm:"unique"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// BeforeCreate hook for GORM to set UUID and slug before creating a record
func (status *Sector) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate UUID if not set
	if status.ID == uuid.Nil {
		status.ID, err = uuid.NewRandom()
		if err != nil {
			return err
		}
	}

	// Generate slug if not set
	if status.Slug == "" {
		status.Slug = slug.Make(status.Title)
	}

	return nil
}
