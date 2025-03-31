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

// ProjectStageHistory tracks the history of stages a project has gone through
type ProjectStageHistory struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	ProjectID uuid.UUID  `json:"project_id" gorm:"type:uuid"`
	StageID   uuid.UUID  `json:"stage_id" gorm:"type:uuid"`
	Stage     *Stage     `json:"stage"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"` // Null until the project moves to a new stage
	Notes     string     `json:"notes"`    // Optional notes about this stage transition
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Project struct {
	ID                   uuid.UUID             `json:"id"`
	MemberID             *uuid.UUID            `json:"-"`
	Member               *Member               `json:"creator"`
	Number               int64                 `json:"number" binding:"required"`
	PartNo               int64                 `json:"part_number"`
	EditionNo            int64                 `json:"edition_number"`
	Reference            string                `json:"reference"`
	ReferenceSuffix      string                `json:"reference_suffix"`
	Title                string                `json:"title" binding:"required"`
	Description          string                `json:"description" binding:"required"`
	TechnicalCommitteeID uuid.UUID             `json:"technical_committee_id"`
	TechnicalCommittee   *TechnicalCommittee   `json:"committee"`
	WorkingGroupID       uuid.UUID             `json:"working_group_id"`
	WorkingGroup         *WorkingGroup         `json:"working_group"`
	StageID              *uuid.UUID            `json:"stage_id"`                // Current stage ID
	Stage                *Stage                `json:"stage"`                   // Current stage
	StageHistory         []ProjectStageHistory `json:"stage_history,omitempty"` // History of all stages
	Timeframe            int                   `json:"time_frame"`              // Timeframe In Months
	Type                 ProjectType           `json:"type" binding:"required" gorm:"default:NEW"`
	VisibleOnLibrary     bool                  `json:"visible_on_library" gorm:"default:true"`
	PricePerPage         float64               `json:"price_per_page"`
	IsEmergency          bool                  `json:"is_emergency" gorm:"default:false"`
	Proposal             *Proposal             `json:"proposal"`
	Acceptance           *Acceptance           `json:"acceptance"`
	Comments             []CommentObservation  `json:"comments,omitempty"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
}
