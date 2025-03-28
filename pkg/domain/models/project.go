package models

import (
	"time"

	"github.com/google/uuid"
)

/*
	For new WDs, shall indicate WD/TC NN/XXX/YYYY, where NN is the TC code, XXX is the serial
	number allocated to the Working Draft by the TC Secretariat and YYYY is the year of circulation.
	For WDs on revision of the standard, shall indicate WD/XXX: YYYY where XXX is the ARS
	number of the current standard and YYYY is the year of circulation. For example, when revising
	ARS 461:2021 in 2024, the corresponding drafts shall be numbered as WD/461:2024. This kind
	of numbering WDs applies also for various stages (CD, DARS and FDARS).
*/

type Project struct {
	ID                   uuid.UUID           `json:"id"`
	Number               int64               `json:"number" binding:"required"`
	PartNo               int64               `json:"part_number"`
	Reference            string              `json:"reference"`
	Title                string              `json:"title" binding:"required"`
	Description          string              `json:"description" binding:"required"`
	TechnicalCommitteeID uuid.UUID           `json:"technical_committee_id"`
	TechnicalCommittee   *TechnicalCommittee `json:"committee"`
	WorkingGroupID       uuid.UUID           `json:"working_group_id"`
	WorkingGroup         *WorkingGroup       `json:"working_group"`
	StageID              *uuid.UUID          `json:"stage_id"`
	Stage                *Stage              `json:"stage"`
	Timeframe            int                 `json:"time_frame" binding:"required"` // Timeframe In Months
	CreatedAt            time.Time           `json:"created_at"`
}
