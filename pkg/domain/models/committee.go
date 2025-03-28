package models

import (
	"time"

	"github.com/google/uuid"
)

// Committee is a base struct for different committees
type Committee struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Code          int64     `json:"code"`
	ChairpersonId uuid.UUID `json:"-"`
	Chairperson   *Member   `json:"chairperson"`
	SecretaryId   uuid.UUID `json:"-"`
	Secretary     *Member   `json:"Secretary"`
	CreatedAt     time.Time `json:"created_at"`
}
