package models

import (
	"github.com/google/uuid"
)

// Committee is a base struct for different committees
type Committee struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name" binding:"required"`
	Code          string    `json:"code" binding:"required"`
	ChairpersonId *string   `json:"chairperson_id"`
	Chairperson   *Member   `json:"chairperson"`
	SecretaryId   *string   `json:"secretary_id"`
	Secretary     *Member   `json:"Secretary"`
}
