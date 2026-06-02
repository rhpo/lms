package entity

import (
	"time"
)

// SupervisorEvaluation représente l'évaluation de l'encadrant (critère 5, /4).
type SupervisorEvaluation struct {
	ID              int64       `json:"id"`
	PfeAssignmentID int64       `json:"pfe_assignment_id"`
	EvaluatorID     int64       `json:"evaluator_id"`
	Criterion5      NullFloat64 `json:"criterion5"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`


	Assignment *PfeAssignment `json:"assignment,omitempty"`
	Evaluator  *Teacher       `json:"evaluator,omitempty"`
}
