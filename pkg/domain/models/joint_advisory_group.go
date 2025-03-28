package models

// JointAdvisoryGroup represents the coordination group
type JointAdvisoryGroup struct {
	Committee
	RegionalEconomicCommunities []*Member `gorm:"many2many:jag_members;"`
	ObserverMembers             []*Member `gorm:"many2many:jag_observers;"`
	ChairRotationPeriod         int       // years
}
