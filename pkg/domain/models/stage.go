package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectDuration struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Min       int64     `json:"min,omitempty"`
	Max       int64     `json:"max,omitempty"`
}

type Timeframe struct {
	ID          uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	StandardID  *uuid.UUID       `gorm:"column:standard_id" json:"-"`
	Standard    *ProjectDuration `json:"standard,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	ISID        *uuid.UUID       `gorm:"column:is_id" json:"-"`
	IS          *ProjectDuration `json:"is,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	EmergencyID *uuid.UUID       `gorm:"column:emergency_id" json:"-"`
	Emergency   *ProjectDuration `json:"emergency,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Description string           `json:"description,omitempty"`
}

type Stage struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Number       int        `json:"number"`
	Name         string     `json:"name"`
	DocumentName string     `json:"document_name"`
	Abbreviation string     `json:"abbreviation"`
	CreatedAt    time.Time  `json:"created_at"`
	TimeframeID  *uuid.UUID `gorm:"column:timeframe_id" json:"-"`
	Timeframe    *Timeframe `json:"timeframe,omitempty" gorm:"foreignKey:TimeframeID;references:ID"`
}
