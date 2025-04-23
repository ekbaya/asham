package models

import (
	"time"

	"github.com/google/uuid"
)

// NSB
type NationalStandardBody struct {
	ID                    uuid.UUID    `json:"id"`
	Name                  string       `json:"name" binding:"required" gorm:"unique"`
	NationalTCSecretaryID *string      `json:"national_tc_secretary_id" binding:"required" gorm:"column:national_tc_secretary_id"`
	NationalTCSecretary   *Member      `json:"national_tc_secretary" gorm:"foreignKey:NationalTCSecretaryID;references:ID"`
	MemberStateID         string       `json:"member_state_id" binding:"required"`
	MemberState           *MemberState `json:"member_state" gorm:"foreignKey:MemberStateID;references:ID"`
	Members               *[]Member    `json:"members"`
	TelephoneNumber       string       `json:"telephone_number" example:"+254700000000"`
	CreatedAt             time.Time    `json:"created_at"`
}
