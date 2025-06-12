package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserType string

const (
	Internal UserType = "internal"
	External UserType = "external"
)

// Member represents a member state in the organization
type Member struct {
	ID                     uuid.UUID             `json:"id"`
	Phone                  string                `json:"phone" binding:"required" example:"+1234567890"`
	Email                  string                `json:"email" gorm:"index;unique" binding:"required"`
	FirstName              string                `json:"first_name" binding:"required"`
	LastName               string                `json:"last_name" binding:"required"`
	PhotoUrl               string                `json:"photo_url"`
	NationalStandardBodyID *string               `json:"nsb_id" binding:"required"`
	NationalStandardBody   *NationalStandardBody `json:"nsb"`
	HashedPassword         string                `json:"-" gorm:"column:password"`
	Organization           string                `json:"organization" `
	Country                string                `json:"country"`
	Type                   UserType              `json:"type" gorm:"default:internal"`
	CanPreviewStandard     bool                  `json:"can_preview_standard" gorm:"default:true"`
	CanDownloadStandard    bool                  `json:"can_download_standard" gorm:"default:true"`
	Roles                  []Role                `gorm:"many2many:user_roles;"`
	CreatedAt              time.Time
}
type MemberMinified struct {
	ID                     uuid.UUID             `json:"id"`
	FirstName              string                `json:"first_name" binding:"required"`
	LastName               string                `json:"last_name" binding:"required"`
	NationalStandardBodyID *string               `json:"nsb_id" binding:"required"`
	NationalStandardBody   *NationalStandardBody `json:"nsb"`
}

// User represents a library user
type User struct {
	ID           uuid.UUID `json:"id"`
	Phone        string    `json:"phone" binding:"required" example:"+1234567890"`
	Email        string    `json:"email" gorm:"index;unique" binding:"required"`
	FirstName    string    `json:"first_name" binding:"required"`
	LastName     string    `json:"last_name" binding:"required"`
	Organization string    `json:"organization" binding:"required"`
	Country      string    `json:"country" binding:"required"`
}

// MemberResponse is a struct used for securely sending Member data in API responses
type MemberResponse struct {
	ID                     uuid.UUID             `json:"id"`
	Phone                  string                `json:"phone"`
	Email                  string                `json:"email"`
	FirstName              string                `json:"first_name"`
	LastName               string                `json:"last_name"`
	PhotoUrl               string                `json:"photo_url"`
	NationalStandardBodyID *string               `json:"nsb_id"`
	NationalStandardBody   *NationalStandardBody `json:"nsb,omitempty"`
	CanPreviewStandard     bool                  `json:"can_preview_standard"`
	CanDownloadStandard    bool                  `json:"can_download_standard"`
	Type                   string                `json:"type"`
	CreatedAt              time.Time             `json:"created_at"`
}

// HashPhone masks a phone number showing only first 2 and last 4 digits
func HashPhone(phone string) string {
	if len(phone) <= 6 {
		return phone // Return as is if too short to mask properly
	}

	// Remove any non-digit characters for consistent handling
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	print(digits)

	prefix := phone[:2]            // Keep first 2 characters as they are
	suffix := phone[len(phone)-4:] // Keep last 4 characters
	maskedLength := len(phone) - 6 // Calculate the middle part to be masked
	masked := strings.Repeat("*", maskedLength)

	return prefix + masked + suffix
}

// HashEmail masks an email address showing only first character and domain
func HashEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // Return as is if not a valid email format
	}

	username := parts[0]
	domain := parts[1]

	if len(username) <= 1 {
		return email // Return as is if username too short
	}

	// Keep first character of username, mask the rest
	maskedUsername := username[:1] + strings.Repeat("*", len(username)-1)

	return maskedUsername + "@" + domain
}

// ToSecureResponse converts a Member to a MemberResponse with sensitive data masked
func (m *Member) ToSecureResponse() MemberResponse {
	return MemberResponse{
		ID:                     m.ID,
		Phone:                  HashPhone(m.Phone),
		Email:                  HashEmail(m.Email),
		FirstName:              m.FirstName,
		LastName:               m.LastName,
		PhotoUrl:               m.PhotoUrl,
		NationalStandardBodyID: m.NationalStandardBodyID,
		NationalStandardBody:   m.NationalStandardBody,
		Type:                   string(m.Type),
		CanPreviewStandard:     m.CanPreviewStandard,
		CanDownloadStandard:    m.CanDownloadStandard,
		CreatedAt:              m.CreatedAt,
	}
}

// MarshalJSON custom JSON marshaling for Member that uses the secure response
func (m Member) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.ToSecureResponse())
}

// UserResponse is a struct used for securely sending User data in API responses
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
}

// ToSecureResponse converts a User to a UserResponse with sensitive data masked
func (u *User) ToSecureResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Phone:     HashPhone(u.Phone),
		Email:     HashEmail(u.Email),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
	}
}

// MarshalJSON custom JSON marshaling for User that uses the secure response
func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.ToSecureResponse())
}
