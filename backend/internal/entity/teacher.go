package entity

import (
	"database/sql"
	"time"
)

// Teacher représente un enseignant.
type Teacher struct {
	ID                 string         `json:"id"`
	ProfileID          string         `json:"profile_id"`
	Grade              sql.NullString `json:"grade"`
	Department         sql.NullString `json:"department"`
	AvailabilityStatus string         `json:"availability_status"`
	UnavailableUntil   sql.NullTime   `json:"unavailable_until"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`

	// Relations
	Profile  *Profile  `json:"profile,omitempty"`
	Domaines []*Domain `json:"domaines,omitempty"`
}
