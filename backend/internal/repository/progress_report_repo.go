package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// ProgressReportRepository gère les opérations base de données pour les comptes rendus d'avancement.
type ProgressReportRepository struct {
	db *sql.DB
}

// NewProgressReportRepository crée un nouveau ProgressReportRepository.
func NewProgressReportRepository(db *sql.DB) *ProgressReportRepository {
	return &ProgressReportRepository{db: db}
}

// FindByID cherche un compte rendu par son ID.
func (r *ProgressReportRepository) FindByID(id string) (*entity.PfeProgressReport, error) {
	row := r.db.QueryRow(`SELECT id, assignment_id, meeting_date, duration, meeting_type, topics, status, observation, created_at, updated_at
		FROM pfe_progress_reports WHERE id = ?`, id)
	p := &entity.PfeProgressReport{}
	err := row.Scan(&p.ID, &p.AssignmentID, &p.MeetingDate, &p.Duration, &p.MeetingType, &p.Topics, &p.Status, &p.Observation, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

// FindByAssignment retourne les comptes rendus d'une assignation.
func (r *ProgressReportRepository) FindByAssignment(assignmentID string) ([]*entity.PfeProgressReport, error) {
	rows, err := r.db.Query(`SELECT id, assignment_id, meeting_date, duration, meeting_type, topics, status, observation, created_at, updated_at
		FROM pfe_progress_reports WHERE assignment_id = ? ORDER BY meeting_date DESC`, assignmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*entity.PfeProgressReport
	for rows.Next() {
		p := &entity.PfeProgressReport{}
		if err := rows.Scan(&p.ID, &p.AssignmentID, &p.MeetingDate, &p.Duration, &p.MeetingType, &p.Topics, &p.Status, &p.Observation, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		reports = append(reports, p)
	}
	return reports, rows.Err()
}

// Insert crée un nouveau compte rendu.
func (r *ProgressReportRepository) Insert(p *entity.PfeProgressReport) error {
	_, err := r.db.Exec(`INSERT INTO pfe_progress_reports (id, assignment_id, meeting_date, duration, meeting_type, topics, status, observation)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, p.ID, p.AssignmentID, p.MeetingDate, p.Duration, p.MeetingType, p.Topics, p.Status, p.Observation)
	return err
}

// Update met à jour un compte rendu.
func (r *ProgressReportRepository) Update(p *entity.PfeProgressReport) error {
	_, err := r.db.Exec(`UPDATE pfe_progress_reports SET meeting_date = ?, duration = ?, meeting_type = ?, topics = ?, status = ?, observation = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		p.MeetingDate, p.Duration, p.MeetingType, p.Topics, p.Status, p.Observation, p.ID)
	return err
}
