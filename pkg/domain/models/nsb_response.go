package models

import (
	"time"

	"github.com/google/uuid"
)

type Response string

const (
	// Agreement response types
	// @Description Agreement to advance the proposal
	ResponseAgreeAdvance Response = "AGREE_ADVANCE"

	// @Description Agreement to accept as a working draft
	ResponseAgreeAcceptWorkingDraft Response = "AGREE_ACCEPT_WORKING_DRAFT"

	// @Description Agreement to circulate as a Committee Draft (CD)
	ResponseAgreeCirculateCD Response = "AGREE_CIRCULATE_CD"

	// @Description Agreement to circulate as a Draft ARSO Standard (DARS)
	ResponseAgreeCirculateDARF Response = "AGREE_CIRCULATE_DARS"

	// @Description No agreement on advancing the proposal
	ResponseNoAgreement Response = "NO_AGREEMENT"

	// @Description Abstention from providing a response
	ResponseAbstention Response = "ABSTENTION"
)

type NSBResponse struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id" binding:"required"`
	Project   *Project  `json:"-"`

	AcceptanceID uuid.UUID   `json:"-"`
	Acceptance   *Acceptance `json:"-"`
	// @Description Response provided by the NSB
	Response Response `json:"response"`

	// Relevant documents section
	// @Description Indicates if relevant standards exist
	HasRelevantStandards bool `json:"has_relevant_standards"`

	// @Description References to relevant standards
	RelevantStandardsRefs *[]Document `json:"relevant_standards_refs,omitempty" gorm:"many2many:relevant_standards_refs;"`

	// @Description Indicates if relevant regulations exist
	HasRelevantRegulations bool `json:"has_relevant_regulations"`

	// @Description References to relevant regulations
	RelevantRegulationsRefs string `json:"relevant_regulations_refs,omitempty"`

	// Comments section
	// @Description Additional comments from the NSB
	Comments string `json:"comments,omitempty"`

	// Participation section
	// @Description Indicates commitment to participate in the work
	IsCommittedToParticipate bool `json:"is_committed_to_participate"`

	// @Description Contact details for the NSB
	NationalTCSecretaryID uuid.UUID `json:"national_tc_secretary_id"`
	NationalTCSecretary   *Member   `json:"national_tc_secretary"`

	// Response details
	// @Description Name of the responding NSB
	RespondingNSBID uuid.UUID             `json:"responding_nsb_id"`
	RespondingNSB   *NationalStandardBody `json:"responding_nsb"`

	// @Description Name of the person providing the response
	ResponderID uuid.UUID `json:"responder_id"`
	Responder   *Member   `json:"responder"`

	// @Description Date of response submission
	ResponseDate time.Time `json:"response_date"`
	// @Description Timestamp when the form was created
	CreatedAt time.Time `json:"created_at"`
}
