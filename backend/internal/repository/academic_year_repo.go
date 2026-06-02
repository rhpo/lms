package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type AcademicYearRepository struct {
	db *sql.DB
}

func NewAcademicYearRepository(db *sql.DB) *AcademicYearRepository {
	return &AcademicYearRepository{db: db}
}

func (r *AcademicYearRepository) FindByID(id int64) (*entity.AcademicYear, error) {
	row := r.db.QueryRow(`SELECT id, label, status, submission_open_at, submission_close_at, max_wishes, created_at, updated_at FROM academic_years WHERE id = ?`, id)
	a := &entity.AcademicYear{}
	err := row.Scan(&a.ID, &a.Label, &a.Status, &a.SubmissionOpenAt, &a.SubmissionCloseAt, &a.MaxWishes, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

func (r *AcademicYearRepository) FindActive() (*entity.AcademicYear, error) {
	row := r.db.QueryRow(`SELECT id, label, status, submission_open_at, submission_close_at, max_wishes, created_at, updated_at FROM academic_years WHERE status = 'active' LIMIT 1`)
	a := &entity.AcademicYear{}
	err := row.Scan(&a.ID, &a.Label, &a.Status, &a.SubmissionOpenAt, &a.SubmissionCloseAt, &a.MaxWishes, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

func (r *AcademicYearRepository) FindAll() ([]*entity.AcademicYear, error) {
	rows, err := r.db.Query(`SELECT id, label, status, submission_open_at, submission_close_at, max_wishes, created_at, updated_at FROM academic_years ORDER BY label DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var years []*entity.AcademicYear
	for rows.Next() {
		a := &entity.AcademicYear{}
		if err := rows.Scan(&a.ID, &a.Label, &a.Status, &a.SubmissionOpenAt, &a.SubmissionCloseAt, &a.MaxWishes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		years = append(years, a)
	}
	return years, rows.Err()
}

func (r *AcademicYearRepository) Insert(a *entity.AcademicYear) error {
	result, err := r.db.Exec(`INSERT INTO academic_years (label, status, submission_open_at, submission_close_at, max_wishes) VALUES (?, ?, ?, ?, ?)`,
		a.Label, a.Status, a.SubmissionOpenAt, a.SubmissionCloseAt, a.MaxWishes)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	a.ID = id
	return nil
}

func (r *AcademicYearRepository) Close(id int64) error {
	_, err := r.db.Exec(`UPDATE academic_years SET status = 'cloturee', updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	return err
}

func (r *AcademicYearRepository) Update(a *entity.AcademicYear) error {
	_, err := r.db.Exec(`UPDATE academic_years SET label = ?, status = ?, submission_open_at = ?, submission_close_at = ?, max_wishes = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		a.Label, a.Status, a.SubmissionOpenAt, a.SubmissionCloseAt, a.MaxWishes, a.ID)
	return err
}
