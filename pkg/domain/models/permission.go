package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Permission struct {
	ID          uuid.UUID      `json:"id" gorm:"primary_key"`
	Title       string         `json:"title" binding:"required" gorm:"unique"`
	Slug        string         `json:"slug" gorm:"unique"`
	Description string         `json:"description"`
	Resource    string         `json:"resource"`
	Action      string         `json:"action"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func (permission *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate UUID if not set
	if permission.ID == uuid.Nil {
		permission.ID, err = uuid.NewRandom()
		if err != nil {
			return err
		}
	}

	// Generate slug if not set
	if permission.Slug == "" {
		permission.Slug = slug.Make(permission.Title)
	}

	return nil
}
