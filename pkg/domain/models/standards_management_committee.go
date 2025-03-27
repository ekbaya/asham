package models

// StandardsManagementCommittee manages standards development
type StandardsManagementCommittee struct {
	Committee
	RegionalRepresentatives []*Member
	ElectedMembers          []*Member
	Observers               []*Member
}
