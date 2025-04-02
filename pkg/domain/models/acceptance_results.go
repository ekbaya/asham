package models

// AcceptanceResults represents the aggregated results of a Acceptance
type AcceptanceResults struct {
	AcceptanceID           string                  `json:"accepatance_id"`
	ProjectID              string                  `json:"project_id"`
	IndividualNSBResponses []IndividualNSBResponse `json:"nsb_responses"`
	Totals                 ResponseTotals          `json:"totals"`
}

// IndividualNSBResponse represents an individual NSB's response for the table
type IndividualNSBResponse struct {
	NSB                    string `json:"nsb"`
	FeasibleYes            bool   `json:"feasible_yes"`
	FeasibleNo             bool   `json:"feasible_no"`
	Abstention             bool   `json:"abstention"`
	AcceptedAsNWIP         string `json:"accepted_as_nwip"`         // Y/N
	AcceptedForProgressing string `json:"accepted_for_progressing"` // Y/N
	AcceptedAsWD           string `json:"accepted_as_wd"`           // Y/N
	AcceptedAsCD           string `json:"accepted_as_cd"`           // Y/N
	AcceptedAsDARF         string `json:"accepted_as_dars"`         // Y/N
	CommentsEnclosed       bool   `json:"comments_enclosed"`
	Participation          bool   `json:"participation"`
}

// ResponseTotals aggregates all the counts for the "Total" row
type ResponseTotals struct {
	TotalResponses              int `json:"total_responses"`
	ValidResponses              int `json:"valid_responses"`
	FeasibleYesCount            int `json:"feasible_yes_count"`
	FeasibleNoCount             int `json:"feasible_no_count"`
	AbstentionCount             int `json:"abstention_count"`
	AcceptedAsNWIPCount         int `json:"accepted_as_nwip_count"`
	AcceptedForProgressingCount int `json:"accepted_for_progressing_count"`
	AcceptedAsWDCount           int `json:"accepted_as_wd_count"`
	AcceptedAsCDCount           int `json:"accepted_as_cd_count"`
	AcceptedAsDARFCount         int `json:"accepted_as_dars_count"`
	CommentsCount               int `json:"comments_count"`
	ParticipationCount          int `json:"participation_count"`
}
