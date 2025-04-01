package models

import (
	"time"

	"github.com/google/uuid"
)

type Proposal struct {
	ID             uuid.UUID             `json:"id"`
	ProjectID      string                `json:"project_id" binding:"required"`
	Project        *Project              `json:"-"`
	CreatedByID    string                `json:"-"`
	CreatedBy      *Member               `json:"created_by"`
	ProposingNSBID *string               `json:"proposing_nsb_id" binding:"required"`
	ProposingNSB   *NationalStandardBody `json:"proposing_nsb"`
	FullTitle      string                `json:"full_title" example:"Milk Quality and Safety Standard"`
	Scope          string                `json:"scope" example:"Defines requirements for milk quality"`
	Justification  string                `json:"justification" example:"To enhance milk safety and market access"`
	// @Description Estimated time for technical project (section 4a)
	EstimatedTime string `json:"estimated_time" example:"12 months"`
	// @Description Proposed deadline for FDARS submission (section 4b)
	ProposedDeadline string `json:"proposed_deadline" example:"2026-01-01"`
	// @Description List of standards to base on (section 5a)
	ReferencedStandards *[]Document `json:"referenced_standards" gorm:"many2many:referenced_standards;"`
	// International Standard information
	// @Description Does an existing International Standard exist? (section 5b)
	ExistingIntlStandard bool `json:"existing_intl_standard" example:"true"`

	// @Description Details if yes
	ExistingIntlStandardDetails string `json:"existing_intl_standard_details" example:"ISO 22000:2018"`

	// @Description Is it suitable for endorsement?
	SuitableForEndorsement bool `json:"suitable_for_endorsement" example:"true"`

	// @Description Is it suitable for reference?
	SuitableForReference bool `json:"suitable_for_reference" example:"true"`

	// @Description Reasons if not suitable
	ReasonIfNotSuitable string `json:"reason_if_not_suitable" example:"Not applicable"`

	// Draft text
	// @Description Is draft text or outline attached? (section 5c)
	IsDraftTextAttached    bool   `json:"is_draft_text_attached" example:"true"`
	DraftTextAttachmentURL string `json:"draft_url"`

	// Legislation
	// @Description Existing national legislation (section 6)
	ExistingLegislation string `json:"existing_legislation" example:"Food Safety Act 2024"`

	// @Description YES, NO, or NOT KNOWN
	LegislationStatus string `json:"legislation_status" example:"YES"`

	// @Description Details of legislation
	LegislationDetails string `json:"legislation_details" example:"Mandatory milk quality standards enforced"`

	// Patents
	// @Description Would any aspect conflict with patented items? (section 7)
	ConflictWithPatents string `json:"conflict_with_patents" example:"NO"`

	// @Description Details if YES
	PatentDetails string `json:"patent_details" example:"Not applicable"`

	// Participation
	// @Description Is member prepared to participate? (section 8a)
	WillParticipateInWork bool `json:"will_participate_in_work" example:"true"`

	// @Description Is member prepared to undertake secretariat? (section 8b)
	WillUndertakeSecretariat bool `json:"will_undertake_secretariat" example:"false"`

	// @Description Is member prepared to undertake preparatory work? (section 8c)
	WillUndertakePrepWork bool `json:"will_undertake_prep_work" example:"true"`
	CreatedAt             time.Time
}
