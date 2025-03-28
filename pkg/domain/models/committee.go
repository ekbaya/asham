package models

import (
	"time"

	"github.com/google/uuid"
)

// Committee is a base struct for different committees
type Committee struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Chairperson *Member   `json:"chairperson"`
	Secretary   *Member   `json:"Secretary"`
	CreatedAt   time.Time `json:"created_at"`
}
