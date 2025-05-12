package models

import "html/template"

// TechnicalCommittee represents a standards development group
type TechnicalCommittee struct {
	Committee
	Scope                  template.HTML         `json:"scope" binding:"required"`
	WorkProgram            string                `json:"work_program"`
	WorkingGroups          []*WorkingGroup       `gorm:"foreignKey:ParentTCID"`
	SubCommittees          []*SubCommittee       `gorm:"foreignKey:ParentTCID"`
	MinimumMembers         int                   `json:"minimum_members" gorm:"default:5"`
	CurrentMembers         []*Member             `json:"current_members" gorm:"many2many:current_members;"`
	ParticipatingCountries []*MemberState        `json:"participating_countries" gorm:"many2many:participating_countries;"`
	EquivalentCommittees   []*TechnicalCommittee `json:"equivalent_committees" gorm:"many2many:equivalent_committees;"`
	ObserverCountries      []*MemberState        `json:"tc_observers" gorm:"many2many:tc_observers;"`
	Projects               []*Project            `json:"projects"`
	EditingCommittee       *EditingCommittee     `gorm:"foreignKey:ParentTCID"`
}

// TechnicalCommitteeDTO represents a DTO for TechnicalCommittee
type TechnicalCommitteeDTO struct {
	CommitteeDTO
	Scope       template.HTML `json:"scope"`
	WorkProgram string        `json:"work_program"`
}
type TechnicalCommitteeDetailDTO struct {
	CommitteeDTO
	Scope          string          `json:"scope"`
	WorkProgram    string          `json:"work_program"`
	WorkingGroups  []*WorkingGroup `gorm:"foreignKey:ParentTCID"`
	SubCommittees  []*SubCommittee `gorm:"foreignKey:ParentTCID"`
	MinimumMembers int             `json:"minimum_members" gorm:"default:5"`
	CurrentMembers []*Member
	Projects       []*Project
}
