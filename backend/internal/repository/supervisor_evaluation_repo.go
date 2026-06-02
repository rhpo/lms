package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type SupervisorEvaluationRepository struct {
	db *sql.DB
}

func NewSupervisorEvaluationRepository(db *sql.DB) *SupervisorEvaluationRepository {
	return &SupervisorEvaluationRepository{db: db}
}

func (r *SupervisorEvaluationRepository) FindByID(id int64) (*entity.SupervisorEvaluation, error) {
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

func (r *SupervisorEvaluationRepository) FindByAssignment(assignmentID int64) (*entity.SupervisorEvaluation, error) {
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

func (r *SupervisorEvaluationRepository) Insert(e *entity.SupervisorEvaluation) error {
	result, err := r.db.Exec(`INSERT OR REPLACE INTO supervisor_evaluations (pfe_assignment_id, evaluator_id, criterion5) VALUES (?, ?, ?)`,
		e.PfeAssignmentID, e.EvaluatorID, e.Criterion5)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = id
	return nil
}

func (r *SupervisorEvaluationRepository) Update(e *entity.SupervisorEvaluation) error {
	_, err := r.db.Exec(`UPDATE supervisor_evaluations SET criterion5 = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, e.Criterion5, e.ID)
	return err
}

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
