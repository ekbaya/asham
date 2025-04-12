package models

// SubCommittee represents a specialized group within a Technical Committee
type SubCommittee struct {
	Committee
	ParentTCID string              `json:"parent_tc_id" gorm:"index"`
	ParentTC   *TechnicalCommittee `json:"parent_tc"`
	Scope      string              `json:"scope"`
	Members    []*Member           `json:"sc_members" gorm:"many2many:sc_members;"`
}
