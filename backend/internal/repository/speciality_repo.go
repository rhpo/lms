package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type SpecialityRepository struct {
	db *sql.DB
}

func NewSpecialityRepository(db *sql.DB) *SpecialityRepository {
	return &SpecialityRepository{db: db}
}

func (r *SpecialityRepository) FindByID(id int64) (*entity.Speciality, error) {
	row := r.db.QueryRow(`SELECT id, name, code, year_type, department_id, created_at, updated_at FROM specialities WHERE id = ?`, id)
	s := &entity.Speciality{}
	err := row.Scan(&s.ID, &s.Name, &s.Code, &s.YearType, &s.DepartmentID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func (r *SpecialityRepository) FindByCode(code string) (*entity.Speciality, error) {
	row := r.db.QueryRow(`SELECT id, name, code, year_type, department_id, created_at, updated_at FROM specialities WHERE code = ?`, code)
	s := &entity.Speciality{}
	err := row.Scan(&s.ID, &s.Name, &s.Code, &s.YearType, &s.DepartmentID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func (r *SpecialityRepository) FindAll() ([]*entity.Speciality, error) {
	rows, err := r.db.Query(`SELECT id, name, code, year_type, department_id, created_at, updated_at FROM specialities ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var specialities []*entity.Speciality
	for rows.Next() {
		s := &entity.Speciality{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Code, &s.YearType, &s.DepartmentID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		specialities = append(specialities, s)
	}
	return specialities, rows.Err()
}

func (r *SpecialityRepository) Insert(s *entity.Speciality) error {
	result, err := r.db.Exec(`INSERT INTO specialities (name, code, year_type, department_id) VALUES (?, ?, ?, ?)`, s.Name, s.Code, s.YearType, s.DepartmentID)
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

func (r *SpecialityRepository) Update(s *entity.Speciality) error {
	_, err := r.db.Exec(`UPDATE specialities SET name = ?, code = ?, year_type = ?, department_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, s.Name, s.Code, s.YearType, s.DepartmentID, s.ID)
	return err
}

func (r *SpecialityRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM specialities WHERE id = ?`, id)
	return err
}
