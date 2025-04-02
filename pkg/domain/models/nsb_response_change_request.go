package models

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	PENDING  Status = "PENDING"
	APPROVED Status = "APPROVED"
	REJECTED Status = "REJECTED"
)

type NSBResponseStatusChange struct {
	ID                       uuid.UUID `json:"id"`
	ResponderID              string    `json:"responder_id"`
	Responder                *Member   `json:"responder"`
	InitialResponseID        string    `json:"nsb_response_id"`
	InitialResponse          *NSBResponse
	Response                 Response  `json:"response" binding:"required"`
	IsCommittedToParticipate bool      `json:"is_committed_to_participate"`
	Status                   Status    `json:"status" gorm:"default:PENDING;"`
	TCSecretariatID          *string   `json:"tc_secretariat_id"`
	TCSecretariat            *Member   `json:"tc_secretariat"`
	TCSecretariatComment     string    `json:"tc_secretariat_comment"`
	CreatedAt                time.Time `json:"created_at"`
}
