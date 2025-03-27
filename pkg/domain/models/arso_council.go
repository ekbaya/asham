package models

// ARSOCouncil represents the executive organ of the organization
type ARSOCouncil struct {
	Committee
	Members          []*Member
	MeetingFrequency int // meetings per year
}
