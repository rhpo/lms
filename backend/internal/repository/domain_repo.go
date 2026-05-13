package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// DomainRepository gère les opérations base de données pour les domaines.
type DomainRepository struct {
	db *sql.DB
}

// NewDomainRepository crée un nouveau DomainRepository.
func NewDomainRepository(db *sql.DB) *DomainRepository {
	return &DomainRepository{db: db}
}

// FindByID cherche un domaine par son ID.
func (r *DomainRepository) FindByID(id string) (*entity.Domain, error) {
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

// FindAll retourne tous les domaines.
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

// Insert crée un nouveau domaine.
func (r *DomainRepository) Insert(d *entity.Domain) error {
	_, err := r.db.Exec(`INSERT INTO domains (id, name) VALUES (?, ?)`, d.ID, d.Name)
	return err
}

// Update met à jour un domaine.
func (r *DomainRepository) Update(d *entity.Domain) error {
	_, err := r.db.Exec(`UPDATE domains SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, d.Name, d.ID)
	return err
}

// Delete supprime un domaine.
func (r *DomainRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM domains WHERE id = ?`, id)
	return err
}
