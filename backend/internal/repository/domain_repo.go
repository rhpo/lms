package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type DomainRepository struct {
	db *sql.DB
}

func NewDomainRepository(db *sql.DB) *DomainRepository {
	return &DomainRepository{db: db}
}

func (r *DomainRepository) FindByID(id int64) (*entity.Domain, error) {
	row := r.db.QueryRow(`SELECT id, name, created_at, updated_at FROM domains WHERE id = ?`, id)
	d := &entity.Domain{}
	err := row.Scan(&d.ID, &d.Name, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return d, nil
}

func (r *DomainRepository) FindAll() ([]*entity.Domain, error) {
	rows, err := r.db.Query(`SELECT id, name, created_at, updated_at FROM domains ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var domains []*entity.Domain
	for rows.Next() {
		d := &entity.Domain{}
		if err := rows.Scan(&d.ID, &d.Name, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}
	return domains, rows.Err()
}

func (r *DomainRepository) Insert(d *entity.Domain) error {
	result, err := r.db.Exec(`INSERT INTO domains (name) VALUES (?)`, d.Name)
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

func (r *DomainRepository) Update(d *entity.Domain) error {
	_, err := r.db.Exec(`UPDATE domains SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, d.Name, d.ID)
	return err
}

func (r *DomainRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM domains WHERE id = ?`, id)
	return err
}
