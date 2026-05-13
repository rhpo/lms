package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// CompanyReportRepository gère les opérations base de données pour les signalements d'entreprises.
type CompanyReportRepository struct {
	db *sql.DB
}

// NewCompanyReportRepository crée un nouveau CompanyReportRepository.
func NewCompanyReportRepository(db *sql.DB) *CompanyReportRepository {
	return &CompanyReportRepository{db: db}
}

// FindByID cherche un signalement par son ID.
func (r *CompanyReportRepository) FindByID(id string) (*entity.CompanyReport, error) {
	row := r.db.QueryRow(`SELECT id, company_id, submitted_by, correction_type, description, requested_value, status, resolved_at, created_at, updated_at
		FROM company_reports WHERE id = ?`, id)
	cr := &entity.CompanyReport{}
	err := row.Scan(&cr.ID, &cr.CompanyID, &cr.SubmittedBy, &cr.CorrectionType, &cr.Description, &cr.RequestedValue, &cr.Status, &cr.ResolvedAt, &cr.CreatedAt, &cr.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return cr, nil
}

// FindByCompany retourne les signalements d'une entreprise.
func (r *CompanyReportRepository) FindByCompany(companyID string) ([]*entity.CompanyReport, error) {
	rows, err := r.db.Query(`SELECT id, company_id, submitted_by, correction_type, description, requested_value, status, resolved_at, created_at, updated_at
		FROM company_reports WHERE company_id = ? ORDER BY created_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReports(rows)
}

// FindAll retourne tous les signalements.
func (r *CompanyReportRepository) FindAll() ([]*entity.CompanyReport, error) {
	rows, err := r.db.Query(`SELECT id, company_id, submitted_by, correction_type, description, requested_value, status, resolved_at, created_at, updated_at
		FROM company_reports ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReports(rows)
}

// FindByStatus retourne les signalements filtrés par statut.
func (r *CompanyReportRepository) FindByStatus(status string) ([]*entity.CompanyReport, error) {
	rows, err := r.db.Query(`SELECT id, company_id, submitted_by, correction_type, description, requested_value, status, resolved_at, created_at, updated_at
		FROM company_reports WHERE status = ? ORDER BY created_at DESC`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanReports(rows)
}

// Insert crée un nouveau signalement.
func (r *CompanyReportRepository) Insert(cr *entity.CompanyReport) error {
	_, err := r.db.Exec(`INSERT INTO company_reports (id, company_id, submitted_by, correction_type, description, requested_value, status) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		cr.ID, cr.CompanyID, cr.SubmittedBy, cr.CorrectionType, cr.Description, cr.RequestedValue, cr.Status)
	return err
}

// UpdateStatus met à jour le statut d'un signalement.
func (r *CompanyReportRepository) UpdateStatus(id, status string) error {
	if status == "resolu" {
		_, err := r.db.Exec(`UPDATE company_reports SET status = ?, resolved_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
		return err
	}
	_, err := r.db.Exec(`UPDATE company_reports SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
	return err
}

func (r *CompanyReportRepository) scanReports(rows *sql.Rows) ([]*entity.CompanyReport, error) {
	var reports []*entity.CompanyReport
	for rows.Next() {
		cr := &entity.CompanyReport{}
		if err := rows.Scan(&cr.ID, &cr.CompanyID, &cr.SubmittedBy, &cr.CorrectionType, &cr.Description, &cr.RequestedValue, &cr.Status, &cr.ResolvedAt, &cr.CreatedAt, &cr.UpdatedAt); err != nil {
			return nil, err
		}
		reports = append(reports, cr)
	}
	return reports, rows.Err()
}
