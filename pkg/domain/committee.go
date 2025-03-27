package domain

import (
	"time"

	"github.com/google/uuid"
)

// Committee is a base struct for different committees
type Committee struct {
	ID          uuid.UUID
	Name        string
	Chairperson *Member
	Secretary   *Member
	CreatedAt   time.Time
}
