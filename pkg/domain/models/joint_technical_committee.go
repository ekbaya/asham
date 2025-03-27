package models

// JointTechnicalCommittee for cross-organizational work
type JointTechnicalCommittee struct {
	Committee
	CollaboratingOrganizations []string
	JointMembers               []*Member
	Scope                      string
}
