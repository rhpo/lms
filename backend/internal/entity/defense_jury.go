package entity

import (
	"time"
)

// DefenseJury représente le jury d'une soutenance (président + membre).
type DefenseJury struct {
	ID                    string    `json:"id"`
	AssignmentID          string    `json:"assignment_id"`
	PresidentID           string    `json:"president_id"`
	MemberID              string    `json:"member_id"`
	PresidentConfirmed    bool      `json:"president_confirmed"`
	MemberConfirmed       bool      `json:"member_confirmed"`
	PresidentWantsPrinted bool      `json:"president_wants_printed"`
	MemberWantsPrinted    bool      `json:"member_wants_printed"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`

	// Relations
	Assignment *PfeAssignment `json:"assignment,omitempty"`
	President  *Teacher       `json:"president,omitempty"`
	Member     *Teacher       `json:"member,omitempty"`
}
