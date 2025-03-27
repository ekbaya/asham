package domain

// TechnicalCommittee represents a standards development group
type TechnicalCommittee struct {
	Committee
	Scope          string
	WorkProgram    string
	WorkingGroups  []*WorkingGroup
	SubCommittees  []*SubCommittee
	MinimumMembers int
	CurrentMembers []*Member
}
