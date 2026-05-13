package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// SpecialityRepository gère les opérations base de données pour les spécialités.
type SpecialityRepository struct {
	db *sql.DB
}

// NewSpecialityRepository crée un nouveau SpecialityRepository.
func NewSpecialityRepository(db *sql.DB) *SpecialityRepository {
	return &SpecialityRepository{db: db}
}

// FindByID cherche une spécialité par son ID.
func (r *SpecialityRepository) FindByID(id string) (*entity.Speciality, error) {
	row := r.db.QueryRow(`SELECT id, name, code, year_type, created_at, updated_at FROM specialities WHERE id = ?`, id)
	s := &entity.Speciality{}
	err := row.Scan(&s.ID, &s.Name, &s.Code, &s.YearType, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

// FindByCode cherche une spécialité par son code.
func (r *SpecialityRepository) FindByCode(code string) (*entity.Speciality, error) {
	row := r.db.QueryRow(`SELECT id, name, code, year_type, created_at, updated_at FROM specialities WHERE code = ?`, code)
	s := &entity.Speciality{}
	err := row.Scan(&s.ID, &s.Name, &s.Code, &s.YearType, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

// FindAll retourne toutes les spécialités.
func (r *SpecialityRepository) FindAll() ([]*entity.Speciality, error) {
	rows, err := r.db.Query(`SELECT id, name, code, year_type, created_at, updated_at FROM specialities ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specialities []*entity.Speciality
	for rows.Next() {
		s := &entity.Speciality{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Code, &s.YearType, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		specialities = append(specialities, s)
	}
	return specialities, rows.Err()
}

// Insert crée une nouvelle spécialité.
func (r *SpecialityRepository) Insert(s *entity.Speciality) error {
	_, err := r.db.Exec(`INSERT INTO specialities (id, name, code, year_type) VALUES (?, ?, ?, ?)`, s.ID, s.Name, s.Code, s.YearType)
	return err
}

// Update met à jour une spécialité.
func (r *SpecialityRepository) Update(s *entity.Speciality) error {
	_, err := r.db.Exec(`UPDATE specialities SET name = ?, code = ?, year_type = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, s.Name, s.Code, s.YearType, s.ID)
	return err
}

// Delete supprime une spécialité.
func (r *SpecialityRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM specialities WHERE id = ?`, id)
	return err
}
