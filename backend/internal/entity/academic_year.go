package entity

import (
	"database/sql"
	"time"
)

// AcademicYear représente une année universitaire.
type AcademicYear struct {
	ID                string       `json:"id"`
	Label             string       `json:"label"`  // ex: 2024-2025
	Status            string       `json:"status"` // active/cloturee
	SubmissionOpenAt  sql.NullTime `json:"submission_open_at"`
	SubmissionCloseAt sql.NullTime `json:"submission_close_at"`
	MaxWishes         int          `json:"max_wishes"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}
