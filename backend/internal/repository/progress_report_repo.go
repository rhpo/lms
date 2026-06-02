package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type ProgressReportRepository struct {
	db *sql.DB
}

func NewProgressReportRepository(db *sql.DB) *ProgressReportRepository {
	return &ProgressReportRepository{db: db}
}

func (r *ProgressReportRepository) FindByID(id int64) (*entity.PfeProgressReport, error) {
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

func (r *ProgressReportRepository) FindByAssignment(assignmentID int64) ([]*entity.PfeProgressReport, error) {
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

func (r *ProgressReportRepository) Insert(p *entity.PfeProgressReport) error {
	result, err := r.db.Exec(`INSERT INTO pfe_progress_reports (assignment_id, meeting_date, duration, meeting_type, topics, status, observation)
		VALUES (?, ?, ?, ?, ?, ?, ?)`, p.AssignmentID, p.MeetingDate, p.Duration, p.MeetingType, p.Topics, p.Status, p.Observation)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = id
	return nil
}

func (r *ProgressReportRepository) Update(p *entity.PfeProgressReport) error {
	_, err := r.db.Exec(`UPDATE pfe_progress_reports SET meeting_date = ?, duration = ?, meeting_type = ?, topics = ?, status = ?, observation = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		p.MeetingDate, p.Duration, p.MeetingType, p.Topics, p.Status, p.Observation, p.ID)
	return err
}
