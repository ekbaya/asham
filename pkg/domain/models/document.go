package models

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID              uuid.UUID `json:"id"`
	CreatedByID     string    `json:"-"`
	CreatedBy       *Member   `json:"created_by"`
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	Reference       string    `json:"reference" binding:"required" gorm:"unique;index"`
	FileURL         string    `json:"file_url"`
	SharepointDocID *string   `json:"sharepoint_doc_id"`
	CreatedAt       time.Time `json:"created_at"`
}
