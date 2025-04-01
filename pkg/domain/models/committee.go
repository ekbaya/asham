package models

import (
	"github.com/google/uuid"
)

// Committee is a base struct for different committees
type Committee struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Code          int64      `json:"code"`
	ChairpersonId *uuid.UUID `json:"chairperson_id"`
	Chairperson   *Member    `json:"chairperson"`
	SecretaryId   *uuid.UUID `json:"secretary_id"`
	Secretary     *Member    `json:"Secretary"`
}
