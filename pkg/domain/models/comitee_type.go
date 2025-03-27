package models

type CommitteeType string

const (
	ARSO_Council                   CommitteeType = "ARSOCouncil"
	Joint_Advisory_Group           CommitteeType = "JointAdvisoryGroup"
	Standards_Management_Committee CommitteeType = "StandardsManagementCommittee"
	Technical_Committee            CommitteeType = "TechnicalCommittee"
	Specialized_Committee          CommitteeType = "SpecializedCommittee"
	Joint_Technical_Committee      CommitteeType = "JointTechnicalCommittee"
)

// ValidateCommitteeType checks if a given type is valid
func ValidateCommitteeType(t string) bool {
	switch CommitteeType(t) {
	case ARSO_Council, Joint_Advisory_Group, Standards_Management_Committee,
		Technical_Committee, Specialized_Committee, Joint_Technical_Committee:
		return true
	default:
		return false
	}
}
