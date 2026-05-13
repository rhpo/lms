package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// CompanyRepository gère les opérations base de données pour les entreprises.
type CompanyRepository struct {
	db *sql.DB
}

// NewCompanyRepository crée un nouveau CompanyRepository.
func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// FindByID cherche une entreprise par son ID.
func (r *CompanyRepository) FindByID(id string) (*entity.Company, error) {
	query := `SELECT id, profile_id, company_name, sector, description, logo_url, contact_email, contact_phone, website, is_verified, created_at, updated_at
		FROM companies WHERE id = ?`
	row := r.db.QueryRow(query, id)

	c := &entity.Company{}
	err := row.Scan(&c.ID, &c.ProfileID, &c.CompanyName, &c.Sector, &c.Description, &c.LogoURL, &c.ContactEmail, &c.ContactPhone, &c.Website, &c.IsVerified, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return c, nil
}

// FindByProfileID cherche une entreprise par son profile_id.
func (r *CompanyRepository) FindByProfileID(profileID string) (*entity.Company, error) {
	query := `SELECT id, profile_id, company_name, sector, description, logo_url, contact_email, contact_phone, website, is_verified, created_at, updated_at
		FROM companies WHERE profile_id = ?`
	row := r.db.QueryRow(query, profileID)

	c := &entity.Company{}
	err := row.Scan(&c.ID, &c.ProfileID, &c.CompanyName, &c.Sector, &c.Description, &c.LogoURL, &c.ContactEmail, &c.ContactPhone, &c.Website, &c.IsVerified, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return c, nil
}

// FindAll retourne toutes les entreprises.
func (r *CompanyRepository) FindAll() ([]*entity.Company, error) {
	query := `SELECT id, profile_id, company_name, sector, description, logo_url, contact_email, contact_phone, website, is_verified, created_at, updated_at
		FROM companies ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*entity.Company
	for rows.Next() {
		c := &entity.Company{}
		if err := rows.Scan(&c.ID, &c.ProfileID, &c.CompanyName, &c.Sector, &c.Description, &c.LogoURL, &c.ContactEmail, &c.ContactPhone, &c.Website, &c.IsVerified, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, rows.Err()
}

// Insert crée une nouvelle entreprise.
func (r *CompanyRepository) Insert(c *entity.Company) error {
	query := `INSERT INTO companies (id, profile_id, company_name, sector, description, logo_url, contact_email, contact_phone, website, is_verified)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, c.ID, c.ProfileID, c.CompanyName, c.Sector, c.Description, c.LogoURL, c.ContactEmail, c.ContactPhone, c.Website, c.IsVerified)
	return err
}

// Update met à jour les informations d'une entreprise.
func (r *CompanyRepository) Update(c *entity.Company) error {
	query := `UPDATE companies SET company_name = ?, sector = ?, description = ?, logo_url = ?, contact_email = ?, contact_phone = ?, website = ?, is_verified = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, c.CompanyName, c.Sector, c.Description, c.LogoURL, c.ContactEmail, c.ContactPhone, c.Website, c.IsVerified, c.ID)
	return err
}

// UpdateVerification met à jour le statut de vérification d'une entreprise.
func (r *CompanyRepository) UpdateVerification(id string, isVerified bool) error {
	_, err := r.db.Exec(`UPDATE companies SET is_verified = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, isVerified, id)
	return err
}

// UpdateLogoURLByProfileID met à jour le logo d'une entreprise via son profile_id.
func (r *CompanyRepository) UpdateLogoURLByProfileID(profileID, url string) error {
	_, err := r.db.Exec(`UPDATE companies SET logo_url = ?, updated_at = CURRENT_TIMESTAMP WHERE profile_id = ?`, url, profileID)
	return err
}

// Delete supprime une entreprise.
func (r *CompanyRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM companies WHERE id = ?`, id)
	return err
}
