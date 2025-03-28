package models

import "github.com/lib/pq"

// SpecializedCommittee represents specific purpose committees
type SpecializedCommittee struct {
	Committee
	Type       string         // e.g., "Conformity Assessment", "Consumer"
	Objectives pq.StringArray `json:"objectives" gorm:"type:text[]"`
	Members    []*Member      `gorm:"many2many:specialized_committee_members;"`
}
