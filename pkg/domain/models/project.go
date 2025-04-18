package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectType string

const (
	NEW           ProjectType = "NEW"
	REVISION      ProjectType = "REVISION"
	INTERNATIONAL ProjectType = "INTERNATIONAL"
)

type WorkingDraftStatus string

const (
	UNDER_REVIEW WorkingDraftStatus = "UNDER_REVIEW"
	ACCEPTED     WorkingDraftStatus = "ACCEPTED"
	DROPPED      WorkingDraftStatus = "DROPPED"
)

type ProposalAction string

const (
	DISCUSS_CD_COMMENTS  ProposalAction = "DISCUSS_CD_COMMENTS"
	REVISED_CD           ProposalAction = "REVISED_CD"
	REGISTER_FOR_ENQUIRY ProposalAction = "REGISTER_FOR_ENQUIRY"
)

// ProjectStageHistory tracks the history of stages a project has gone through
type ProjectStageHistory struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	ProjectID string     `json:"project_id" gorm:"type:uuid"`
	StageID   string     `json:"stage_id" gorm:"type:uuid"`
	Stage     *Stage     `json:"stage"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"` // Null until the project moves to a new stage
	Notes     string     `json:"notes"`    // Optional notes about this stage transition
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Project struct {
	ID                   uuid.UUID             `json:"id"`
	MemberID             *string               `json:"-"`
	Member               *Member               `json:"creator"`
	Sector               ProjectSector         `json:"sector" binding:"required"`
	Number               int64                 `json:"number"`
	PartNo               int64                 `json:"part_number"`
	EditionNo            int64                 `json:"edition_number"`
	Reference            string                `json:"reference"`
	ReferenceSuffix      string                `json:"reference_suffix"`
	Title                string                `json:"title" binding:"required"`
	Description          string                `json:"description" binding:"required"`
	TechnicalCommitteeID string                `json:"technical_committee_id" binding:"required"`
	TechnicalCommittee   *TechnicalCommittee   `json:"committee"`
	WorkingGroupID       *string               `json:"working_group_id"`
	WorkingGroup         *WorkingGroup         `json:"working_group"`
	StageID              string                `json:"stage_id"`                // Current stage ID
	Stage                *Stage                `json:"stage"`                   // Current stage
	StageHistory         []ProjectStageHistory `json:"stage_history,omitempty"` // History of all stages
	Timeframe            int                   `json:"time_frame"`              // Timeframe In Months
	Type                 ProjectType           `json:"type" binding:"required" gorm:"default:NEW"`
	VisibleOnLibrary     bool                  `json:"visible_on_library" gorm:"default:true"`
	PricePerPage         float64               `json:"price_per_page"`
	IsEmergency          bool                  `json:"is_emergency" gorm:"default:false"`
	PWIApproved          bool                  `json:"pwi_approved" gorm:"default:false"`
	PWIApprovalComment   string                `json:"pwi_approval_comment"`
	ApprovedByID         *string               `json:"-"`
	ApprovedBy           *Member               `json:"approved_by"`
	Proposal             *Proposal             `json:"proposal"`
	Acceptance           *Acceptance           `json:"acceptance"`
	WorkingDraftID       *string               `json:"working_draft_id"`
	WorkingDraft         *Document             `json:"working_draft"`
	WorkingDraftStatus   WorkingDraftStatus    `json:"wd_status" gorm:"default:UNDER_REVIEW"`
	WDTCSecretaryID      *string               `json:"wd_tc_secretary_id" gorm:"column:wd_tc_secretary_id"`
	WDTCSecretary        *Member               `json:"wd_tc_secretary"`
	WorkingDraftComments string                `json:"wd_comments"`
	CommitteeDraftID     *string               `json:"committee_draft_id"`
	CommitteeDraft       *Document             `json:"committee_draft"`
	Comments             []CommentObservation  `json:"comments,omitempty"`
	CDTCSecretaryID      *string               `json:"cd_tc_secretary_id" gorm:"column:cd_tc_secretary_id"`
	CDTCSecretary        *Member               `json:"cd_tc_secretary"`
	IsConsensusReached   bool                  `json:"is_consensus_reached"`
	ProposalAction       ProposalAction        `json:"proposal_action"`
	MeetingRequired      bool                  `json:"meeting_required"`
	SubmissionDate       *time.Time            `json:"submission_date,omitempty"`
	DARS                 *DARS                 `json:"dars"`
	Balloting            *Balloting            `json:"ballot"`
	RelatedDocuments     *[]Document           `json:"project_related_documents" gorm:"many2many:project_related_documents;"`
	// Keep track of cancellation at ballot level
	Cancelled                     bool       `json:"cancelled" gorm:"default:false"`
	CancelledDate                 *time.Time `json:"cancelled_date"`
	ApprovedForPublication        bool       `json:"approved_for_publication" gorm:"default:false"`
	ApprovedForPublicationDate    *time.Time `json:"approved_for_publication_date"`
	ApprovedForPublicationByID    *string    `json:"approved_for_publication_by_id"`
	ApprovedForPublicationBy      *Member    `json:"approved_for_publication_by"`
	ApprovedForPublicationComment string     `json:"approved_for_publication_comment"`
	DARSDocID                     *string    `json:"dars_doc_id"`
	DARSDoc                       *Document  `json:"dars_doc"`
	FDARSDocID                    *string    `json:"fdars_doc_id"`
	FDARSDoc                      *Document  `json:"fdars_doc"`
	StandardID                    *string    `json:"standard_id"`
	Standard                      *Document  `json:"standard"`
	Published                     bool       `json:"published" gorm:"default:false"`
	PublishedDate                 *time.Time `json:"published_date"`
	CreatedAt                     time.Time  `json:"created_at"`
	UpdatedAt                     time.Time  `json:"updated_at"`
}

// ProjectDTO represents a subset of Project fields for repository queries
type ProjectDTO struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Title       string    `json:"title"`
	Reference   string    `json:"reference"`
	Sector      string    `json:"sector"`
	Description string    `json:"description"`
	Published   bool      `json:"published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
