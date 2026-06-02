package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type CompanyRepository struct {
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) FindByID(id int64) (*entity.Company, error) {
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

func (r *CompanyRepository) FindByProfileID(profileID int64) (*entity.Company, error) {
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

func (r *CompanyRepository) FindAllVerified() ([]*entity.Company, error) {
	query := `SELECT id, profile_id, company_name, sector, description, logo_url, contact_email, contact_phone, website, is_verified, created_at, updated_at
		FROM companies WHERE is_verified = 1 ORDER BY company_name ASC`
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

func (r *CompanyRepository) FindByName(name string) (*entity.Company, error) {
	query := `SELECT id, profile_id, company_name, sector, description, logo_url, contact_email, contact_phone, website, is_verified, created_at, updated_at
		FROM companies WHERE company_name = ? LIMIT 1`
	row := r.db.QueryRow(query, name)
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

func (r *CompanyRepository) FindAllByName(name string) ([]*entity.Company, error) {
	query := `SELECT id, profile_id, company_name, sector, description, logo_url, contact_email, contact_phone, website, is_verified, created_at, updated_at
		FROM companies WHERE company_name = ?`
	rows, err := r.db.Query(query, name)
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

func (r *CompanyRepository) Insert(c *entity.Company) error {
	query := `INSERT INTO companies (profile_id, company_name, sector, description, logo_url, contact_email, contact_phone, website, is_verified)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, c.ProfileID, c.CompanyName, c.Sector, c.Description, c.LogoURL, c.ContactEmail, c.ContactPhone, c.Website, c.IsVerified)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (r *CompanyRepository) Update(c *entity.Company) error {
	query := `UPDATE companies SET company_name = ?, sector = ?, description = ?, logo_url = ?, contact_email = ?, contact_phone = ?, website = ?, is_verified = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, c.CompanyName, c.Sector, c.Description, c.LogoURL, c.ContactEmail, c.ContactPhone, c.Website, c.IsVerified, c.ID)
	return err
}

func (r *CompanyRepository) UpdateVerification(id int64, isVerified bool) error {
	_, err := r.db.Exec(`UPDATE companies SET is_verified = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, isVerified, id)
	return err
}

func (r *CompanyRepository) UpdateLogoURLByProfileID(profileID int64, url string) error {
	_, err := r.db.Exec(`UPDATE companies SET logo_url = ?, updated_at = CURRENT_TIMESTAMP WHERE profile_id = ?`, url, profileID)
	return err
}

func (r *CompanyRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM companies WHERE id = ?`, id)
	return err
}
