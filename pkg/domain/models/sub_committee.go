package models

import "github.com/google/uuid"

// SubCommittee represents a specialized group within a Technical Committee
type SubCommittee struct {
	Committee
	ParentTCID uuid.UUID
	ParentTC   *TechnicalCommittee
	Scope      string
	Members    []*Member `gorm:"many2many:sc_members;"`
}
