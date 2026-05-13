package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// AcademicYearRepository gère les opérations base de données pour les années universitaires.
type AcademicYearRepository struct {
	db *sql.DB
}

// NewAcademicYearRepository crée un nouveau AcademicYearRepository.
func NewAcademicYearRepository(db *sql.DB) *AcademicYearRepository {
	return &AcademicYearRepository{db: db}
}

// FindByID cherche une année universitaire par son ID.
func (r *AcademicYearRepository) FindByID(id string) (*entity.AcademicYear, error) {
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

// FindActive retourne l'année universitaire active.
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

// FindAll retourne toutes les années universitaires.
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

// Insert crée une nouvelle année universitaire.
func (r *AcademicYearRepository) Insert(a *entity.AcademicYear) error {
	_, err := r.db.Exec(`INSERT INTO academic_years (id, label, status, submission_open_at, submission_close_at, max_wishes) VALUES (?, ?, ?, ?, ?, ?)`,
		a.ID, a.Label, a.Status, a.SubmissionOpenAt, a.SubmissionCloseAt, a.MaxWishes)
	return err
}

// Close ferme une année universitaire (passe son statut à 'cloturee').
func (r *AcademicYearRepository) Close(id string) error {
	_, err := r.db.Exec(`UPDATE academic_years SET status = 'cloturee', updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	return err
}

// Update met à jour une année universitaire.
func (r *AcademicYearRepository) Update(a *entity.AcademicYear) error {
	_, err := r.db.Exec(`UPDATE academic_years SET label = ?, status = ?, submission_open_at = ?, submission_close_at = ?, max_wishes = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		a.Label, a.Status, a.SubmissionOpenAt, a.SubmissionCloseAt, a.MaxWishes, a.ID)
	return err
}
