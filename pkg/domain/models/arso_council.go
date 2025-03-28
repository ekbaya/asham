package models

// ARSOCouncil represents the executive organ of the organization
type ARSOCouncil struct {
	Committee
	Members          []*Member `gorm:"many2many:arsocouncil_members;"`
	MeetingFrequency int
}
