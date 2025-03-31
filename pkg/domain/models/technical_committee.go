package models

// TechnicalCommittee represents a standards development group
type TechnicalCommittee struct {
	Committee
	Scope          string          `json:"scope"`
	WorkProgram    string          `json:"work_program"`
	WorkingGroups  []*WorkingGroup `gorm:"foreignKey:ParentTCID"`
	SubCommittees  []*SubCommittee `gorm:"foreignKey:ParentTCID"`
	MinimumMembers int
	CurrentMembers []*Member `gorm:"many2many:current_members;"`
}
