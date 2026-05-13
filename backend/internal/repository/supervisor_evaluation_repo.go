package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// SupervisorEvaluationRepository gère les opérations base de données pour les évaluations d'encadrants.
type SupervisorEvaluationRepository struct {
	db *sql.DB
}

// NewSupervisorEvaluationRepository crée un nouveau SupervisorEvaluationRepository.
func NewSupervisorEvaluationRepository(db *sql.DB) *SupervisorEvaluationRepository {
	return &SupervisorEvaluationRepository{db: db}
}

// FindByID cherche une évaluation par son ID.
func (r *SupervisorEvaluationRepository) FindByID(id string) (*entity.SupervisorEvaluation, error) {
	row := r.db.QueryRow(`SELECT id, pfe_assignment_id, evaluator_id, criterion5, created_at, updated_at FROM supervisor_evaluations WHERE id = ?`, id)
	e := &entity.SupervisorEvaluation{}
	err := row.Scan(&e.ID, &e.PfeAssignmentID, &e.EvaluatorID, &e.Criterion5, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return e, nil
}

// FindByAssignment retourne l'évaluation pour une assignation.
func (r *SupervisorEvaluationRepository) FindByAssignment(assignmentID string) (*entity.SupervisorEvaluation, error) {
	row := r.db.QueryRow(`SELECT id, pfe_assignment_id, evaluator_id, criterion5, created_at, updated_at
		FROM supervisor_evaluations WHERE pfe_assignment_id = ?`, assignmentID)
	e := &entity.SupervisorEvaluation{}
	err := row.Scan(&e.ID, &e.PfeAssignmentID, &e.EvaluatorID, &e.Criterion5, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return e, nil
}

// Insert crée ou met à jour une évaluation d'encadrant (upsert sur pfe_assignment_id).
func (r *SupervisorEvaluationRepository) Insert(e *entity.SupervisorEvaluation) error {
	_, err := r.db.Exec(`INSERT OR REPLACE INTO supervisor_evaluations (id, pfe_assignment_id, evaluator_id, criterion5) VALUES (?, ?, ?, ?)`,
		e.ID, e.PfeAssignmentID, e.EvaluatorID, e.Criterion5)
	return err
}

// Update met à jour une évaluation d'encadrant.
func (r *SupervisorEvaluationRepository) Update(e *entity.SupervisorEvaluation) error {
	_, err := r.db.Exec(`UPDATE supervisor_evaluations SET criterion5 = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, e.Criterion5, e.ID)
	return err
}

// FindAll retourne toutes les évaluations d'encadrants.
func (r *SupervisorEvaluationRepository) FindAll() ([]*entity.SupervisorEvaluation, error) {
	rows, err := r.db.Query(`SELECT id, pfe_assignment_id, evaluator_id, criterion5, created_at, updated_at FROM supervisor_evaluations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evaluations []*entity.SupervisorEvaluation
	for rows.Next() {
		e := &entity.SupervisorEvaluation{}
		if err := rows.Scan(&e.ID, &e.PfeAssignmentID, &e.EvaluatorID, &e.Criterion5, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		evaluations = append(evaluations, e)
	}
	return evaluations, rows.Err()
}
