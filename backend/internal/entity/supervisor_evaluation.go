package entity

import (
	"database/sql"
	"time"
)

// SupervisorEvaluation représente l'évaluation de l'encadrant (critère 5, /4).
type SupervisorEvaluation struct {
	ID              string          `json:"id"`
	PfeAssignmentID string          `json:"pfe_assignment_id"`
	EvaluatorID     string          `json:"evaluator_id"`
	Criterion5      sql.NullFloat64 `json:"criterion5"` // /4
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	// Relations
	Assignment *PfeAssignment `json:"assignment,omitempty"`
	Evaluator  *Teacher       `json:"evaluator,omitempty"`
}
