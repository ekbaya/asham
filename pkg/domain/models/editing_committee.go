package models

// EditingCommittee is charged with the updating and editorial work of the standards
type EditingCommittee struct {
	Committee
	Editors    []*Member           `json:"editors" gorm:"many2many:editors;"`
	ParentTCID string              `json:"parent_tc_id" gorm:"index"`
	ParentTC   *TechnicalCommittee `json:"parent_tc"`
}
