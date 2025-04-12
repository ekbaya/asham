package models

import "github.com/lib/pq"

// SpecializedCommittee represents specific purpose committees
type SpecializedCommittee struct {
	Committee
	Type       string         // e.g., "Conformity Assessment", "Consumer"
	Objectives pq.StringArray `json:"objectives" gorm:"type:text[]"`
	Members    []*Member      `json:"specialized_committee_members" gorm:"many2many:specialized_committee_members;"`
}
