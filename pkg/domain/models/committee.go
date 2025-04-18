package models

import (
	"github.com/google/uuid"
)

// Committee is a base struct for different committees
type Committee struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name" binding:"required"`
	Code          string    `json:"code"`
	Description   string    `json:"description"`
	ChairpersonId *string   `json:"chairperson_id"`
	Chairperson   *Member   `json:"chairperson"`
	SecretaryId   *string   `json:"secretary_id"`
	Secretary     *Member   `json:"Secretary"`
}

// CommitteeDTO includes committee details, members, and counts
type CommitteeDTO struct {
	ID                 uuid.UUID       `json:"id"`
	Name               string          `json:"name"`
	Code               string          `json:"code"`
	Description        string          `json:"description"`
	Chairperson        *MemberMinified `json:"chairperson"`
	ChairpersonId      *string         `json:"chairperson_id"`
	WorkingGroupCount  int64           `json:"working_group_count"`
	MemberCount        int64           `json:"member_count"`
	WorkingMemberCount int64           `json:"working_member_count"`
	ActiveProjectCount int64           `json:"active_project_count"`
}
