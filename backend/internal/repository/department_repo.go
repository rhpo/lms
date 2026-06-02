package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type DepartmentRepository struct {
	db *sql.DB
}

func NewDepartmentRepository(db *sql.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) FindByID(id int64) (*entity.Department, error) {
	row := r.db.QueryRow(`SELECT id, name, created_at, updated_at FROM departments WHERE id = ?`, id)
	d := &entity.Department{}
	err := row.Scan(&d.ID, &d.Name, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return d, nil
}

func (r *DepartmentRepository) FindAll() ([]*entity.Department, error) {
	rows, err := r.db.Query(`SELECT id, name, created_at, updated_at FROM departments ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var departments []*entity.Department
	for rows.Next() {
		d := &entity.Department{}
		if err := rows.Scan(&d.ID, &d.Name, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		departments = append(departments, d)
	}
	return departments, rows.Err()
}

func (r *DepartmentRepository) Insert(d *entity.Department) error {
	result, err := r.db.Exec(`INSERT INTO departments (name) VALUES (?)`, d.Name)
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

func (r *DepartmentRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM departments WHERE id = ?`, id)
	return err
}
