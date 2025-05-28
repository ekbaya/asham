package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Role struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key"`
	Title       string    `json:"title" binding:"required" gorm:"unique"`
	Slug        string    `json:"slug" gorm:"unique"`
	Description string
	Permissions []Permission   `gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func (role *Role) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate UUID if not set
	if role.ID == uuid.Nil {
		role.ID, err = uuid.NewRandom()
		if err != nil {
			return err
		}
	}

	// Generate slug if not set
	if role.Slug == "" {
		role.Slug = slug.Make(role.Title)
	}

	return nil
}
