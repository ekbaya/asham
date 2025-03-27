package domain

// SubCommittee represents a specialized group within a Technical Committee
type SubCommittee struct {
	Committee
	ParentTC *TechnicalCommittee
	Scope    string
	Members  []*Member
}
