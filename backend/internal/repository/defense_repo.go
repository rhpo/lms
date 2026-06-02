package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type DefenseRepository struct {
	db *sql.DB
}

func NewDefenseRepository(db *sql.DB) *DefenseRepository {
	return &DefenseRepository{db: db}
}

func (r *DefenseRepository) FindByID(id int64) (*entity.Defense, error) {
	row := r.db.QueryRow(`SELECT id, assignment_id, jury_id, scheduled_at, room, defense_deadline, status, result, final_grade, created_at, updated_at
		FROM defenses WHERE id = ?`, id)
	d := &entity.Defense{}
	err := row.Scan(&d.ID, &d.AssignmentID, &d.JuryID, &d.ScheduledAt, &d.Room, &d.DefenseDeadline, &d.Status, &d.Result, &d.FinalGrade, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return d, nil
}

func (r *DefenseRepository) FindAll() ([]*entity.Defense, error) {
	rows, err := r.db.Query(`SELECT id, assignment_id, jury_id, scheduled_at, room, defense_deadline, status, result, final_grade, created_at, updated_at
		FROM defenses ORDER BY scheduled_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var defenses []*entity.Defense
	for rows.Next() {
		d := &entity.Defense{}
		if err := rows.Scan(&d.ID, &d.AssignmentID, &d.JuryID, &d.ScheduledAt, &d.Room, &d.DefenseDeadline, &d.Status, &d.Result, &d.FinalGrade, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		defenses = append(defenses, d)
	}
	return defenses, rows.Err()
}

func (r *DefenseRepository) Insert(d *entity.Defense) error {
	result, err := r.db.Exec(`INSERT INTO defenses (assignment_id, jury_id, scheduled_at, room, defense_deadline, status) VALUES (?, ?, ?, ?, ?, ?)`,
		d.AssignmentID, d.JuryID, d.ScheduledAt, d.Room, d.DefenseDeadline, d.Status)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	d.ID = id
	return nil
}

func (r *DefenseRepository) UpdateStatus(id int64, status string) error {
	_, err := r.db.Exec(`UPDATE defenses SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
	return err
}

func (r *DefenseRepository) UpdateResult(id int64, result string, finalGrade float64) error {
	_, err := r.db.Exec(`UPDATE defenses SET result = ?, final_grade = ?, status = 'done', updated_at = CURRENT_TIMESTAMP WHERE id = ?`, result, finalGrade, id)
	return err
}

func (r *DefenseRepository) FindByAssignment(assignmentID int64) (*entity.Defense, error) {
	row := r.db.QueryRow(`SELECT id, assignment_id, jury_id, scheduled_at, room, defense_deadline, status, result, final_grade, created_at, updated_at
		FROM defenses WHERE assignment_id = ? LIMIT 1`, assignmentID)
	d := &entity.Defense{}
	err := row.Scan(&d.ID, &d.AssignmentID, &d.JuryID, &d.ScheduledAt, &d.Room, &d.DefenseDeadline, &d.Status, &d.Result, &d.FinalGrade, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return d, nil
}

func (r *DefenseRepository) Update(d *entity.Defense) error {
	_, err := r.db.Exec(`UPDATE defenses SET scheduled_at = ?, room = ?, defense_deadline = ?, status = ?, result = ?, final_grade = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		d.ScheduledAt, d.Room, d.DefenseDeadline, d.Status, d.Result, d.FinalGrade, d.ID)
	return err
}

func (r *DefenseRepository) FindByJuryMember(teacherID int64) ([]*entity.Defense, error) {
	query := `SELECT d.id, d.assignment_id, d.jury_id, d.scheduled_at, d.room, d.defense_deadline, d.status, d.result, d.final_grade, d.created_at, d.updated_at
		FROM defenses d INNER JOIN defense_juries dj ON dj.id = d.jury_id
		WHERE dj.president_id = ? OR dj.member_id = ? ORDER BY d.scheduled_at DESC`
	rows, err := r.db.Query(query, teacherID, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var defenses []*entity.Defense
	for rows.Next() {
		d := &entity.Defense{}
		if err := rows.Scan(&d.ID, &d.AssignmentID, &d.JuryID, &d.ScheduledAt, &d.Room, &d.DefenseDeadline, &d.Status, &d.Result, &d.FinalGrade, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		defenses = append(defenses, d)
	}
	return defenses, rows.Err()
}
