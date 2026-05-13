package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// WishRepository gère les opérations base de données pour les vœux.
type WishRepository struct {
	db *sql.DB
}

// NewWishRepository crée un nouveau WishRepository.
func NewWishRepository(db *sql.DB) *WishRepository {
	return &WishRepository{db: db}
}

// FindByID cherche un vœu par son ID.
func (r *WishRepository) FindByID(id string) (*entity.Wish, error) {
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

// FindByStudent retourne les vœux d'un étudiant pour une année universitaire.
func (r *WishRepository) FindByStudent(studentID, academicYearID string) ([]*entity.Wish, error) {
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

// FindBySubject retourne les vœux pour un sujet.
func (r *WishRepository) FindBySubject(subjectID string) ([]*entity.Wish, error) {
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

// CountByStudent compte les vœux d'un étudiant pour une année universitaire.
func (r *WishRepository) CountByStudent(studentID, academicYearID string) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM wishes WHERE student_id = ? AND academic_year_id = ?`, studentID, academicYearID).Scan(&count)
	return count, err
}

// Insert crée un nouveau vœu.
func (r *WishRepository) Insert(w *entity.Wish) error {
	_, err := r.db.Exec(`INSERT INTO wishes (id, student_id, subject_id, academic_year_id, status) VALUES (?, ?, ?, ?, ?)`,
		w.ID, w.StudentID, w.SubjectID, w.AcademicYearID, w.Status)
	return err
}

// UpdateStatus met à jour le statut d'un vœu.
func (r *WishRepository) UpdateStatus(id, status string) error {
	_, err := r.db.Exec(`UPDATE wishes SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
	return err
}

// Update met à jour un vœu complet.
func (r *WishRepository) Update(w *entity.Wish) error {
	_, err := r.db.Exec(`UPDATE wishes SET status = ?, student_id = ?, subject_id = ?, academic_year_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		w.Status, w.StudentID, w.SubjectID, w.AcademicYearID, w.ID)
	return err
}

// Delete supprime un vœu.
func (r *WishRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM wishes WHERE id = ?`, id)
	return err
}
