package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// DefenseJuryRepository gère les opérations base de données pour les jurys de soutenance.
type DefenseJuryRepository struct {
	db *sql.DB
}

// NewDefenseJuryRepository crée un nouveau DefenseJuryRepository.
func NewDefenseJuryRepository(db *sql.DB) *DefenseJuryRepository {
	return &DefenseJuryRepository{db: db}
}

// FindByID cherche un jury par son ID.
func (r *DefenseJuryRepository) FindByID(id string) (*entity.DefenseJury, error) {
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

// FindByAssignment retourne le jury d'une assignation.
func (r *DefenseJuryRepository) FindByAssignment(assignmentID string) (*entity.DefenseJury, error) {
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

// Insert crée un nouveau jury.
func (r *DefenseJuryRepository) Insert(j *entity.DefenseJury) error {
	_, err := r.db.Exec(`INSERT INTO defense_juries (id, assignment_id, president_id, member_id, president_confirmed, member_confirmed,
		president_wants_printed, member_wants_printed) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		j.ID, j.AssignmentID, j.PresidentID, j.MemberID, j.PresidentConfirmed, j.MemberConfirmed,
		j.PresidentWantsPrinted, j.MemberWantsPrinted)
	return err
}

// ConfirmPresident confirme la participation du président.
func (r *DefenseJuryRepository) ConfirmPresident(id string) error {
	_, err := r.db.Exec(`UPDATE defense_juries SET president_confirmed = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	return err
}

// ConfirmMember confirme la participation du membre.
func (r *DefenseJuryRepository) ConfirmMember(id string) error {
	_, err := r.db.Exec(`UPDATE defense_juries SET member_confirmed = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	return err
}

// SetPresidentWantsPrinted définit si le président veut une version papier.
func (r *DefenseJuryRepository) SetPresidentWantsPrinted(id string, wants bool) error {
	_, err := r.db.Exec(`UPDATE defense_juries SET president_wants_printed = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, wants, id)
	return err
}

// SetMemberWantsPrinted définit si le membre veut une version papier.
func (r *DefenseJuryRepository) SetMemberWantsPrinted(id string, wants bool) error {
	_, err := r.db.Exec(`UPDATE defense_juries SET member_wants_printed = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, wants, id)
	return err
}

// Delete supprime un jury.
func (r *DefenseJuryRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM defense_juries WHERE id = ?`, id)
	return err
}
