package models

import (
	"time"

	"github.com/google/uuid"
)

// CommentType represents the type of comment
type CommentType string

const (
	General   CommentType = "ge" // General comment
	Technical CommentType = "te" // Technical comment
	Editorial CommentType = "ed" // Editorial comment
)

// CommentObservation represents a single comment and observation entry
type CommentObservation struct {
	ID                 uuid.UUID   `json:"id"`
	ProjectID          uuid.UUID   `json:"-"`
	Project            *Project    `json:"project"`
	MemberID           uuid.UUID   `json:"-"`
	Member             *Member     `json:"member"`
	ClauseNo           string      `json:"clause_no" binding:"required"`
	ParagraphRef       string      `json:"paragraph_ref" binding:"required"`
	CommentType        CommentType `json:"comment_type" binding:"required"`
	Comment            string      `json:"comment" binding:"required"`
	ProposedChange     string      `json:"proposed_change"`
	SecretariatRemarks string      `json:"secretariat_remarks"`
	CreatedAt          time.Time   `json:"created_at"`
}
