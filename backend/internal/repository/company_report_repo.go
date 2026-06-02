package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type CompanyReportRepository struct {
	db *sql.DB
}

func NewCompanyReportRepository(db *sql.DB) *CompanyReportRepository {
	return &CompanyReportRepository{db: db}
}

func (r *CompanyReportRepository) FindByID(id int64) (*entity.CompanyReport, error) {
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

func (r *CompanyReportRepository) FindByCompany(companyID int64) ([]*entity.CompanyReport, error) {
	rows, err := r.db.Query(`SELECT id, company_id, submitted_by, correction_type, description, requested_value, status, resolved_at, created_at, updated_at
		FROM company_reports WHERE company_id = ? ORDER BY created_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanReports(rows)
}

func (r *CompanyReportRepository) FindAll() ([]*entity.CompanyReport, error) {
	rows, err := r.db.Query(`SELECT id, company_id, submitted_by, correction_type, description, requested_value, status, resolved_at, created_at, updated_at
		FROM company_reports ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanReports(rows)
}

func (r *CompanyReportRepository) FindByStatus(status string) ([]*entity.CompanyReport, error) {
	rows, err := r.db.Query(`SELECT id, company_id, submitted_by, correction_type, description, requested_value, status, resolved_at, created_at, updated_at
		FROM company_reports WHERE status = ? ORDER BY created_at DESC`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanReports(rows)
}

func (r *CompanyReportRepository) Insert(cr *entity.CompanyReport) error {
	result, err := r.db.Exec(`INSERT INTO company_reports (company_id, submitted_by, correction_type, description, requested_value, status) VALUES (?, ?, ?, ?, ?, ?)`,
		cr.CompanyID, cr.SubmittedBy, cr.CorrectionType, cr.Description, cr.RequestedValue, cr.Status)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	cr.ID = id
	return nil
}

func (r *CompanyReportRepository) UpdateStatus(id int64, status string) error {
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
