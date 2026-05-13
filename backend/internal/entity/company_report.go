package entity

import (
	"database/sql"
	"time"
)

// CompanyReport représente une demande de correction soumise par une entreprise.
type CompanyReport struct {
	ID             string       `json:"id"`
	CompanyID      string       `json:"company_id"`
	SubmittedBy    string       `json:"submitted_by"`
	CorrectionType string       `json:"correction_type"`
	Description    string       `json:"description"`
	RequestedValue string       `json:"requested_value"`
	Status         string       `json:"status"` // en_attente/resolu/rejete
	ResolvedAt     sql.NullTime `json:"resolved_at"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`

	// Relations
	Company *Company `json:"company,omitempty"`
}
