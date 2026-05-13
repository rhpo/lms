package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// PfeSubjectRepository gère les opérations base de données pour les sujets PFE.
type PfeSubjectRepository struct {
	db *sql.DB
}

// NewPfeSubjectRepository crée un nouveau PfeSubjectRepository.
func NewPfeSubjectRepository(db *sql.DB) *PfeSubjectRepository {
	return &PfeSubjectRepository{db: db}
}

// FindByID cherche un sujet par son ID.
func (r *PfeSubjectRepository) FindByID(id string) (*entity.PfeSubject, error) {
	query := `SELECT id, title, description, group_type, proposer_id, proposer_role, company_id, academic_year_id,
		validator1_id, validator2_id, validator1_decision, validator2_decision,
		validator1_comment, validator2_comment, status, co_supervisor_id, pre_assigned_student_ids, created_at, updated_at
		FROM pfe_subjects WHERE id = ?`
	row := r.db.QueryRow(query, id)
	s := &entity.PfeSubject{}
	err := row.Scan(
		&s.ID, &s.Title, &s.Description, &s.GroupType, &s.ProposerID, &s.ProposerRole, &s.CompanyID, &s.AcademicYearID,
		&s.Validator1ID, &s.Validator2ID, &s.Validator1Decision, &s.Validator2Decision,
		&s.Validator1Comment, &s.Validator2Comment, &s.Status, &s.CoSupervisorID, &s.PreAssignedStudentIDs, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

// FindAll retourne tous les sujets PFE.
func (r *PfeSubjectRepository) FindAll() ([]*entity.PfeSubject, error) {
	query := `SELECT id, title, description, group_type, proposer_id, proposer_role, company_id, academic_year_id,
		validator1_id, validator2_id, validator1_decision, validator2_decision,
		validator1_comment, validator2_comment, status, co_supervisor_id, pre_assigned_student_ids, created_at, updated_at
		FROM pfe_subjects ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSubjects(rows)
}

// FindByProposer retourne les sujets proposés par un utilisateur.
func (r *PfeSubjectRepository) FindByProposer(proposerID string) ([]*entity.PfeSubject, error) {
	query := `SELECT id, title, description, group_type, proposer_id, proposer_role, company_id, academic_year_id,
		validator1_id, validator2_id, validator1_decision, validator2_decision,
		validator1_comment, validator2_comment, status, co_supervisor_id, pre_assigned_student_ids, created_at, updated_at
		FROM pfe_subjects WHERE proposer_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, proposerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSubjects(rows)
}

// FindByStatus retourne les sujets filtrés par statut.
func (r *PfeSubjectRepository) FindByStatus(status string) ([]*entity.PfeSubject, error) {
	query := `SELECT id, title, description, group_type, proposer_id, proposer_role, company_id, academic_year_id,
		validator1_id, validator2_id, validator1_decision, validator2_decision,
		validator1_comment, validator2_comment, status, co_supervisor_id, pre_assigned_student_ids, created_at, updated_at
		FROM pfe_subjects WHERE status = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSubjects(rows)
}

// FindPendingValidation retourne les sujets en attente de validation.
func (r *PfeSubjectRepository) FindPendingValidation(validatorID string) ([]*entity.PfeSubject, error) {
	query := `SELECT id, title, description, group_type, proposer_id, proposer_role, company_id, academic_year_id,
		validator1_id, validator2_id, validator1_decision, validator2_decision,
		validator1_comment, validator2_comment, status, co_supervisor_id, pre_assigned_student_ids, created_at, updated_at
		FROM pfe_subjects WHERE (validator1_id = ? OR validator2_id = ?) AND status = 'en_attente' ORDER BY created_at DESC`
	rows, err := r.db.Query(query, validatorID, validatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSubjects(rows)
}

// FindByAcademicYear retourne les sujets d'une année universitaire.
func (r *PfeSubjectRepository) FindByAcademicYear(academicYearID string) ([]*entity.PfeSubject, error) {
	query := `SELECT id, title, description, group_type, proposer_id, proposer_role, company_id, academic_year_id,
		validator1_id, validator2_id, validator1_decision, validator2_decision,
		validator1_comment, validator2_comment, status, co_supervisor_id, pre_assigned_student_ids, created_at, updated_at
		FROM pfe_subjects WHERE academic_year_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, academicYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSubjects(rows)
}

// FindByCompany retourne les sujets proposés par une entreprise.
func (r *PfeSubjectRepository) FindByCompany(companyID string) ([]*entity.PfeSubject, error) {
	query := `SELECT id, title, description, group_type, proposer_id, proposer_role, company_id, academic_year_id,
		validator1_id, validator2_id, validator1_decision, validator2_decision,
		validator1_comment, validator2_comment, status, co_supervisor_id, pre_assigned_student_ids, created_at, updated_at
		FROM pfe_subjects WHERE company_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSubjects(rows)
}

// FindAvailable retourne les sujets disponibles (validés) pour les étudiants.
func (r *PfeSubjectRepository) FindAvailable(academicYearID string) ([]*entity.PfeSubject, error) {
	query := `SELECT id, title, description, group_type, proposer_id, proposer_role, company_id, academic_year_id,
		validator1_id, validator2_id, validator1_decision, validator2_decision,
		validator1_comment, validator2_comment, status, co_supervisor_id, pre_assigned_student_ids, created_at, updated_at
		FROM pfe_subjects WHERE status = 'valide' AND academic_year_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, academicYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSubjects(rows)
}

// Insert crée un nouveau sujet PFE.
func (r *PfeSubjectRepository) Insert(s *entity.PfeSubject) error {
	query := `INSERT INTO pfe_subjects (id, title, description, group_type, proposer_id, proposer_role, company_id, academic_year_id,
		validator1_id, validator2_id, validator1_decision, validator2_decision,
		validator1_comment, validator2_comment, status, co_supervisor_id, pre_assigned_student_ids)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, s.ID, s.Title, s.Description, s.GroupType, s.ProposerID, s.ProposerRole, s.CompanyID,
		s.AcademicYearID, s.Validator1ID, s.Validator2ID, s.Validator1Decision, s.Validator2Decision,
		s.Validator1Comment, s.Validator2Comment, s.Status, s.CoSupervisorID, s.PreAssignedStudentIDs)
	return err
}

// Update met à jour un sujet PFE (sans le statut, utiliser UpdateStatus pour changer le statut).
func (r *PfeSubjectRepository) Update(s *entity.PfeSubject) error {
	query := `UPDATE pfe_subjects SET title = ?, description = ?, group_type = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, s.Title, s.Description, s.GroupType, s.ID)
	return err
}

// UpdateStatus met à jour le statut d'un sujet.
func (r *PfeSubjectRepository) UpdateStatus(id, status string) error {
	_, err := r.db.Exec(`UPDATE pfe_subjects SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
	return err
}

// UpdateValidation met à jour la décision d'un validateur.
func (r *PfeSubjectRepository) UpdateValidation(id, validatorField, decision, comment string) error {
	var query string
	if validatorField == "validator1" {
		query = `UPDATE pfe_subjects SET validator1_decision = ?, validator1_comment = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	} else {
		query = `UPDATE pfe_subjects SET validator2_decision = ?, validator2_comment = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	}
	_, err := r.db.Exec(query, decision, comment, id)
	return err
}

// AssignValidators assigne les validateurs à un sujet.
func (r *PfeSubjectRepository) AssignValidators(id, validator1ID, validator2ID string) error {
	_, err := r.db.Exec(`UPDATE pfe_subjects SET validator1_id = ?, validator2_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		validator1ID, validator2ID, id)
	return err
}

// AssignCoSupervisor assigne un co-encadrant.
func (r *PfeSubjectRepository) AssignCoSupervisor(id, coSupervisorID string) error {
	_, err := r.db.Exec(`UPDATE pfe_subjects SET co_supervisor_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		coSupervisorID, id)
	return err
}

// AddDomain ajoute un domaine à un sujet.
func (r *PfeSubjectRepository) AddDomain(subjectID, domainID string) error {
	_, err := r.db.Exec(`INSERT OR IGNORE INTO subject_domains (subject_id, domain_id) VALUES (?, ?)`, subjectID, domainID)
	return err
}

// RemoveDomain retire un domaine d'un sujet.
func (r *PfeSubjectRepository) RemoveDomain(subjectID, domainID string) error {
	_, err := r.db.Exec(`DELETE FROM subject_domains WHERE subject_id = ? AND domain_id = ?`, subjectID, domainID)
	return err
}

// GetDomains retourne les domaines d'un sujet.
func (r *PfeSubjectRepository) GetDomains(subjectID string) ([]*entity.Domain, error) {
	query := `SELECT d.id, d.name, d.created_at, d.updated_at
		FROM domains d INNER JOIN subject_domains sd ON sd.domain_id = d.id WHERE sd.subject_id = ?`
	rows, err := r.db.Query(query, subjectID)
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

// scanSubjects scanne les résultats d'une requête de sujets.
func (r *PfeSubjectRepository) scanSubjects(rows *sql.Rows) ([]*entity.PfeSubject, error) {
	var subjects []*entity.PfeSubject
	for rows.Next() {
		s := &entity.PfeSubject{}
		if err := rows.Scan(
			&s.ID, &s.Title, &s.Description, &s.GroupType, &s.ProposerID, &s.ProposerRole, &s.CompanyID, &s.AcademicYearID,
			&s.Validator1ID, &s.Validator2ID, &s.Validator1Decision, &s.Validator2Decision,
			&s.Validator1Comment, &s.Validator2Comment, &s.Status, &s.CoSupervisorID, &s.PreAssignedStudentIDs, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}
	return subjects, rows.Err()
}
