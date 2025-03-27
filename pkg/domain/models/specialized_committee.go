package models

// SpecializedCommittee represents specific purpose committees
type SpecializedCommittee struct {
	Committee
	Type       string // e.g., "Conformity Assessment", "Consumer"
	Objectives []string
	Members    []*Member
}
