package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) FindByID(id int64) (*entity.Student, error) {
	query := `SELECT id, profile_id, student_number, speciality_id, level, promotion_id, created_at, updated_at
		FROM students WHERE id = ?`
	row := r.db.QueryRow(query, id)
	s := &entity.Student{}
	err := row.Scan(&s.ID, &s.ProfileID, &s.StudentNumber, &s.SpecialityID, &s.Level, &s.PromotionID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func (r *StudentRepository) FindByProfileID(profileID int64) (*entity.Student, error) {
	query := `SELECT id, profile_id, student_number, speciality_id, level, promotion_id, created_at, updated_at
		FROM students WHERE profile_id = ?`
	row := r.db.QueryRow(query, profileID)
	s := &entity.Student{}
	err := row.Scan(&s.ID, &s.ProfileID, &s.StudentNumber, &s.SpecialityID, &s.Level, &s.PromotionID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func (r *StudentRepository) FindByStudentNumber(number string) (*entity.Student, error) {
	query := `SELECT id, profile_id, student_number, speciality_id, level, promotion_id, created_at, updated_at
		FROM students WHERE student_number = ?`
	row := r.db.QueryRow(query, number)
	s := &entity.Student{}
	err := row.Scan(&s.ID, &s.ProfileID, &s.StudentNumber, &s.SpecialityID, &s.Level, &s.PromotionID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func (r *StudentRepository) FindAll() ([]*entity.Student, error) {
	query := `SELECT id, profile_id, student_number, speciality_id, level, promotion_id, created_at, updated_at
		FROM students ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var students []*entity.Student
	for rows.Next() {
		s := &entity.Student{}
		if err := rows.Scan(&s.ID, &s.ProfileID, &s.StudentNumber, &s.SpecialityID, &s.Level, &s.PromotionID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, rows.Err()
}

func (r *StudentRepository) Insert(s *entity.Student) error {
	query := `INSERT INTO students (profile_id, student_number, speciality_id, level, promotion_id)
		VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, s.ProfileID, s.StudentNumber, s.SpecialityID, s.Level, s.PromotionID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}

func (r *StudentRepository) Update(s *entity.Student) error {
	query := `UPDATE students SET student_number = ?, speciality_id = ?, level = ?, promotion_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, s.StudentNumber, s.SpecialityID, s.Level, s.PromotionID, s.ID)
	return err
}

func (r *StudentRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM students WHERE id = ?`, id)
	return err
}
