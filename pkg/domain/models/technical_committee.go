package models

// TechnicalCommittee represents a standards development group
type TechnicalCommittee struct {
	Committee
	Scope          string          `json:"scope" binding:"required"`
	WorkProgram    string          `json:"work_program"`
	WorkingGroups  []*WorkingGroup `gorm:"foreignKey:ParentTCID"`
	SubCommittees  []*SubCommittee `gorm:"foreignKey:ParentTCID"`
	MinimumMembers int             `json:"minimum_members" gorm:"default:5"`
	CurrentMembers []*Member       `gorm:"many2many:current_members;"`
}

// TechnicalCommitteeDTO represents a DTO for TechnicalCommittee
type TechnicalCommitteeDTO struct {
	CommitteeDTO
	Scope       string `json:"scope"`
	WorkProgram string `json:"work_program"`
}
type TechnicalCommitteeDetailDTO struct {
	CommitteeDTO
	Scope          string          `json:"scope"`
	WorkProgram    string          `json:"work_program"`
	WorkingGroups  []*WorkingGroup `gorm:"foreignKey:ParentTCID"`
	SubCommittees  []*SubCommittee `gorm:"foreignKey:ParentTCID"`
	MinimumMembers int             `json:"minimum_members" gorm:"default:5"`
	CurrentMembers []*Member       `gorm:"many2many:current_members;"`
}
