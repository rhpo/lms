package entity

import (
	"database/sql"
	"time"
)

// PfeProgressReport représente un compte-rendu de réunion pour un PFE.
type PfeProgressReport struct {
	ID           string         `json:"id"`
	AssignmentID string         `json:"assignment_id"`
	MeetingDate  time.Time      `json:"meeting_date"`
	Duration     int            `json:"duration"`     // minutes
	MeetingType  string         `json:"meeting_type"` // presentiel/visio
	Topics       string         `json:"topics"`
	Status       string         `json:"status"` // en_cours/termine
	Observation  sql.NullString `json:"observation"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`

	// Relations
	Assignment *PfeAssignment `json:"assignment,omitempty"`
}
