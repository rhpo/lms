package entity

import (
	"database/sql"
	"time"
)

// Defense représente une soutenance.
type Defense struct {
	ID              string          `json:"id"`
	AssignmentID    string          `json:"assignment_id"`
	JuryID          string          `json:"jury_id"`
	ScheduledAt     sql.NullTime    `json:"scheduled_at"`
	Room            sql.NullString  `json:"room"`
	DefenseDeadline sql.NullTime    `json:"defense_deadline"`
	Status          string          `json:"status"` // scheduled/done/postponed
	Result          sql.NullString  `json:"result"` // admitted/corrections_required/not_admitted
	FinalGrade      sql.NullFloat64 `json:"final_grade"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	// Relations
	Assignment *PfeAssignment `json:"assignment,omitempty"`
	Jury       *DefenseJury   `json:"jury,omitempty"`
}
