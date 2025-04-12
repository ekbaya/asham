package models

// JointAdvisoryGroup represents the coordination group
type JointAdvisoryGroup struct {
	Committee
	RegionalEconomicCommunities []*Member `json:"jag_members" gorm:"many2many:jag_members;"`
	ObserverMembers             []*Member `json:"jag_observers" gorm:"many2many:jag_observers;"`
	ChairRotationPeriod         int       // years
}
