package models

// StandardsManagementCommittee manages standards development
type StandardsManagementCommittee struct {
	Committee
	RegionalRepresentatives []*Member `json:"regional_representatives" gorm:"many2many:regional_representatives;"`
	ElectedMembers          []*Member `json:"elected_members" gorm:"many2many:elected_members;"`
	Observers               []*Member `json:"observers" gorm:"many2many:observers;"`
}
