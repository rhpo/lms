package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// StudentRepository gère les opérations base de données pour les étudiants.
type StudentRepository struct {
	db *sql.DB
}

// NewStudentRepository crée un nouveau StudentRepository.
func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

// FindByID cherche un étudiant par son ID.
func (r *StudentRepository) FindByID(id string) (*entity.Student, error) {
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

// FindByProfileID cherche un étudiant par son profile_id.
func (r *StudentRepository) FindByProfileID(profileID string) (*entity.Student, error) {
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

// FindByStudentNumber cherche un étudiant par son numéro étudiant.
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

// FindAll retourne tous les étudiants.
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

// Insert crée un nouvel étudiant.
func (r *StudentRepository) Insert(s *entity.Student) error {
	query := `INSERT INTO students (id, profile_id, student_number, speciality_id, level, promotion_id)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, s.ID, s.ProfileID, s.StudentNumber, s.SpecialityID, s.Level, s.PromotionID)
	return err
}

// Update met à jour les informations d'un étudiant.
func (r *StudentRepository) Update(s *entity.Student) error {
	query := `UPDATE students SET student_number = ?, speciality_id = ?, level = ?, promotion_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, s.StudentNumber, s.SpecialityID, s.Level, s.PromotionID, s.ID)
	return err
}

// Delete supprime un étudiant.
func (r *StudentRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM students WHERE id = ?`, id)
	return err
}
