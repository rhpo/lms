package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type DefenseJuryRepository struct {
	db *sql.DB
}

func NewDefenseJuryRepository(db *sql.DB) *DefenseJuryRepository {
	return &DefenseJuryRepository{db: db}
}

func (r *DefenseJuryRepository) FindByID(id int64) (*entity.DefenseJury, error) {
	row := r.db.QueryRow(`SELECT id, assignment_id, president_id, member_id, president_confirmed, member_confirmed,
		president_wants_printed, member_wants_printed, created_at, updated_at FROM defense_juries WHERE id = ?`, id)
	j := &entity.DefenseJury{}
	err := row.Scan(&j.ID, &j.AssignmentID, &j.PresidentID, &j.MemberID, &j.PresidentConfirmed, &j.MemberConfirmed,
		&j.PresidentWantsPrinted, &j.MemberWantsPrinted, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return j, nil
}

func (r *DefenseJuryRepository) FindByAssignment(assignmentID int64) (*entity.DefenseJury, error) {
	row := r.db.QueryRow(`SELECT id, assignment_id, president_id, member_id, president_confirmed, member_confirmed,
		president_wants_printed, member_wants_printed, created_at, updated_at FROM defense_juries WHERE assignment_id = ?`, assignmentID)
	j := &entity.DefenseJury{}
	err := row.Scan(&j.ID, &j.AssignmentID, &j.PresidentID, &j.MemberID, &j.PresidentConfirmed, &j.MemberConfirmed,
		&j.PresidentWantsPrinted, &j.MemberWantsPrinted, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return j, nil
}

func (r *DefenseJuryRepository) Insert(j *entity.DefenseJury) error {
	result, err := r.db.Exec(`INSERT INTO defense_juries (assignment_id, president_id, member_id, president_confirmed, member_confirmed,
		president_wants_printed, member_wants_printed) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		j.AssignmentID, j.PresidentID, j.MemberID, j.PresidentConfirmed, j.MemberConfirmed,
		j.PresidentWantsPrinted, j.MemberWantsPrinted)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	j.ID = id
	return nil
}

func (r *DefenseJuryRepository) ConfirmPresident(id int64) error {
	_, err := r.db.Exec(`UPDATE defense_juries SET president_confirmed = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	return err
}

func (r *DefenseJuryRepository) ConfirmMember(id int64) error {
	_, err := r.db.Exec(`UPDATE defense_juries SET member_confirmed = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	return err
}

func (r *DefenseJuryRepository) SetPresidentWantsPrinted(id int64, wants bool) error {
	_, err := r.db.Exec(`UPDATE defense_juries SET president_wants_printed = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, wants, id)
	return err
}

func (r *DefenseJuryRepository) SetMemberWantsPrinted(id int64, wants bool) error {
	_, err := r.db.Exec(`UPDATE defense_juries SET member_wants_printed = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, wants, id)
	return err
}

func (r *DefenseJuryRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM defense_juries WHERE id = ?`, id)
	return err
}
