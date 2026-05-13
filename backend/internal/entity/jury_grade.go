package entity

import (
	"database/sql"
	"time"
)

// JuryGrade représente une note soumise par un membre du jury.
type JuryGrade struct {
	ID           string          `json:"id"`
	DefenseID    string          `json:"defense_id"`
	JuryMemberID string          `json:"jury_member_id"`
	Criterion1   sql.NullFloat64 `json:"criterion1"` // /4
	Criterion2   sql.NullFloat64 `json:"criterion2"` // /4
	Criterion3   sql.NullFloat64 `json:"criterion3"` // /4
	Criterion4   sql.NullFloat64 `json:"criterion4"` // /4
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`

	// Relations
	Defense    *Defense `json:"defense,omitempty"`
	JuryMember *Teacher `json:"jury_member,omitempty"`
}
