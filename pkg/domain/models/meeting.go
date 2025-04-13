package models

import (
	"time"

	"github.com/google/uuid"
)

type MeetingType string

const (
	MeetingTypeTC MeetingType = "TC"
	MeetingTypeSC MeetingType = "SC"
	MeetingTypeWG MeetingType = "WG"
)

type MeetingFormat string

const (
	MeetingFormatElectronic MeetingFormat = "ELECTRONIC"
	MeetingFormatInPerson   MeetingFormat = "IN_PERSON"
	MeetingFormatHybrid     MeetingFormat = "HYBRID"
)

type MeetingStatus string

const (
	MeetingStatusPlanned   MeetingStatus = "PLANNED"
	MeetingStatusConfirmed MeetingStatus = "CONFIRMED"
	MeetingStatusCancelled MeetingStatus = "CANCELLED"
	MeetingStatusPostponed MeetingStatus = "POSTPONED"
	MeetingStatusCompleted MeetingStatus = "COMPLETED"
)

type Meeting struct {
	ID          uuid.UUID     `json:"id"`
	MeetingType MeetingType   `json:"meeting_type" binding:"required"` // TC, SC, or WG meeting
	Date        time.Time     `json:"date" binding:"required"`
	StartTime   string        `json:"start_time" binding:"required"` // Start time of the meeting
	EndTime     string        `json:"end_time" binding:"required"`   // End time of the meeting
	Venue       string        `json:"venue" binding:"required"`
	Title       string        `json:"title" binding:"required"`    // City, country, and/or virtual link
	Format      MeetingFormat `json:"format" binding:"required"`   // Electronic, in-person, hybrid
	Language    string        `json:"language" binding:"required"` // English and/or French
	Agenda      string        `json:"agenda" binding:"required"`
	Minutes     string        `json:"minutes"`

	// Committee/Working Group Information
	CommitteeID   string `json:"committee_id" binding:"required"`   // ID of TC, SC, or WG
	CommitteeName string `json:"committee_name" binding:"required"` // Name of TC, SC, or WG

	// Host information
	HostOrganizationID *string               `json:"host_organization_id"` // National body acting as host
	HostOrganization   *NationalStandardBody `json:"host_organization"`

	// Quorum tracking
	TotalPMembers   int  `json:"total_p_members"`   // Total number of P-members
	PresentPMembers int  `json:"present_p_members"` // Number of P-members present
	HasQuorum       bool `json:"has_quorum"`        // Whether quorum was achieved

	// Document distribution tracking
	AgendaDistributionDate time.Time `json:"agenda_distribution_date"` // When agenda was distributed
	DocsDistributionDate   time.Time `json:"docs_distribution_date"`   // When supporting docs were distributed

	// Approval information
	ChairApproval       bool `json:"chair_approval"`       // Chairperson approved date/place
	SecretariatApproval bool `json:"secretariat_approval"` // Secretariat approved date/place

	// Meeting status
	Status             MeetingStatus `json:"status" gorm:"default:PLANNED"` // Planned, confirmed, cancelled, etc.
	CancellationReason string        `json:"cancellation_reason,omitempty"` // If cancelled, why

	// Project information
	ProjectID *string  `json:"project_id"`
	Project   *Project `json:"project"`

	// Documents
	RelatedDocuments *[]Document `json:"meeting_related_documents" gorm:"many2many:meeting_related_documents;"` // CDs to be discussed and Other documents

	// Management
	CreatedByID             string    `json:"created_by"`
	CreatedBy               *Member   `json:"created_by_member"`
	Comments                string    `json:"comments"`
	MinutesApprovalDeadline time.Time `json:"minutes_approval_deadline"` // 30 days after minutes circulation
	Attendees               *[]Member `json:"attendees,omitempty" gorm:"many2many:meeting_attendees;"`
	CreatedAt               time.Time `json:"created_at"`
}
