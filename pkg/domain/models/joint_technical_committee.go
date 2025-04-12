package models

import "github.com/lib/pq"

// JointTechnicalCommittee for cross-organizational work
type JointTechnicalCommittee struct {
	Committee
	CollaboratingOrganizations pq.StringArray `json:"collaborating_organizations" gorm:"type:text[]"`
	JointMembers               []*Member      `json:"joint_members" gorm:"many2many:joint_members;"`
	Scope                      string         `json:"scope"`
}
