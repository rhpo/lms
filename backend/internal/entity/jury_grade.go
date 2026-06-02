package entity

import (
	"time"
)

// JuryGrade représente une note soumise par un membre du jury.
type JuryGrade struct {
	ID           int64       `json:"id"`
	DefenseID    int64       `json:"defense_id"`
	JuryMemberID int64       `json:"jury_member_id"`
	Criterion1      NullFloat64 `json:"criterion1"` // /4
	Criterion2      NullFloat64 `json:"criterion2"` // /4
	Criterion3      NullFloat64 `json:"criterion3"` // /4
	Criterion4      NullFloat64 `json:"criterion4"` // /4
	ArchiveDecision NullString  `json:"archive_decision"` // archivable | minor_corrections | major_corrections
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`

	// Relations
	Defense    *Defense `json:"defense,omitempty"`
	JuryMember *Teacher `json:"jury_member,omitempty"`
}
