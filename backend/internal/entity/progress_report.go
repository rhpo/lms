package entity

import (
	"time"
)

// PfeProgressReport représente un compte-rendu de réunion pour un PFE.
type PfeProgressReport struct {
	ID           int64      `json:"id"`
	AssignmentID int64      `json:"assignment_id"`
	MeetingDate  time.Time  `json:"meeting_date"`
	Duration     int        `json:"duration"`     // minutes
	MeetingType  string     `json:"meeting_type"` // presentiel/visio
	Topics       string     `json:"topics"`
	Status       string     `json:"status"` // en_cours/termine
	Observation  NullString `json:"observation"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relations
	Assignment *PfeAssignment `json:"assignment,omitempty"`
}
