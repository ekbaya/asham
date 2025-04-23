package models

import (
	"github.com/google/uuid"
)

type Committee struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name" binding:"required"`
	Code          string    `json:"code"`
	Description   string    `json:"description" gorm:"default:Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem."`
	ChairpersonId *string   `json:"chairperson_id" gorm:"column:chairperson_id;type:uuid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Chairperson   *Member   `json:"chairperson" gorm:"foreignKey:ChairpersonId;references:ID"`
	SecretaryId   *string   `json:"secretary_id" gorm:"column:secretary_id;type:uuid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Secretary     *Member   `json:"secretary" gorm:"foreignKey:SecretaryId;references:ID"`
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
