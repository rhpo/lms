package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type WishRepository struct {
	db *sql.DB
}

func NewWishRepository(db *sql.DB) *WishRepository {
	return &WishRepository{db: db}
}

func (r *WishRepository) FindByID(id int64) (*entity.Wish, error) {
	row := r.db.QueryRow(`SELECT id, student_id, subject_id, academic_year_id, status, created_at, updated_at FROM wishes WHERE id = ?`, id)
	w := &entity.Wish{}
	err := row.Scan(&w.ID, &w.StudentID, &w.SubjectID, &w.AcademicYearID, &w.Status, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return w, nil
}

func (r *WishRepository) FindByStudent(studentID, academicYearID int64) ([]*entity.Wish, error) {
	query := `SELECT id, student_id, subject_id, academic_year_id, status, created_at, updated_at
		FROM wishes WHERE student_id = ? AND academic_year_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, studentID, academicYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var wishes []*entity.Wish
	for rows.Next() {
		w := &entity.Wish{}
		if err := rows.Scan(&w.ID, &w.StudentID, &w.SubjectID, &w.AcademicYearID, &w.Status, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		wishes = append(wishes, w)
	}
	return wishes, rows.Err()
}

func (r *WishRepository) FindBySubject(subjectID int64) ([]*entity.Wish, error) {
	query := `SELECT id, student_id, subject_id, academic_year_id, status, created_at, updated_at
		FROM wishes WHERE subject_id = ? ORDER BY created_at ASC`
	rows, err := r.db.Query(query, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var wishes []*entity.Wish
	for rows.Next() {
		w := &entity.Wish{}
		if err := rows.Scan(&w.ID, &w.StudentID, &w.SubjectID, &w.AcademicYearID, &w.Status, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		wishes = append(wishes, w)
	}
	return wishes, rows.Err()
}

func (r *WishRepository) CountByStudent(studentID, academicYearID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM wishes WHERE student_id = ? AND academic_year_id = ?`, studentID, academicYearID).Scan(&count)
	return count, err
}

func (r *WishRepository) Insert(w *entity.Wish) error {
	result, err := r.db.Exec(`INSERT INTO wishes (student_id, subject_id, academic_year_id, status) VALUES (?, ?, ?, ?)`,
		w.StudentID, w.SubjectID, w.AcademicYearID, w.Status)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	w.ID = id
	return nil
}

func (r *WishRepository) UpdateStatus(id int64, status string) error {
	_, err := r.db.Exec(`UPDATE wishes SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
	return err
}

func (r *WishRepository) Update(w *entity.Wish) error {
	_, err := r.db.Exec(`UPDATE wishes SET status = ?, student_id = ?, subject_id = ?, academic_year_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		w.Status, w.StudentID, w.SubjectID, w.AcademicYearID, w.ID)
	return err
}

func (r *WishRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM wishes WHERE id = ?`, id)
	return err
}
