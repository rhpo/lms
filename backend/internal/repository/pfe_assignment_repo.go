package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// PfeAssignmentRepository gère les opérations base de données pour les assignations PFE.
type PfeAssignmentRepository struct {
	db *sql.DB
}

// NewPfeAssignmentRepository crée un nouveau PfeAssignmentRepository.
func NewPfeAssignmentRepository(db *sql.DB) *PfeAssignmentRepository {
	return &PfeAssignmentRepository{db: db}
}

// FindByID cherche une assignation par son ID.
func (r *PfeAssignmentRepository) FindByID(id string) (*entity.PfeAssignment, error) {
	query := `SELECT id, pfe_code, subject_id, academic_year_id, student_id, student2_id, student3_id,
		supervisor_id, co_supervisor_id, memoire_url, status, created_at, updated_at
		FROM pfe_assignments WHERE id = ?`
	row := r.db.QueryRow(query, id)
	a := &entity.PfeAssignment{}
	err := row.Scan(&a.ID, &a.PfeCode, &a.SubjectID, &a.AcademicYearID, &a.StudentID, &a.Student2ID, &a.Student3ID,
		&a.SupervisorID, &a.CoSupervisorID, &a.MemoireURL, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

// FindByStudent retourne l'assignation d'un étudiant pour une année universitaire.
func (r *PfeAssignmentRepository) FindByStudent(studentID, academicYearID string) (*entity.PfeAssignment, error) {
	query := `SELECT id, pfe_code, subject_id, academic_year_id, student_id, student2_id, student3_id,
		supervisor_id, co_supervisor_id, memoire_url, status, created_at, updated_at
		FROM pfe_assignments WHERE (student_id = ? OR student2_id = ? OR student3_id = ?) AND academic_year_id = ? LIMIT 1`
	row := r.db.QueryRow(query, studentID, studentID, studentID, academicYearID)
	a := &entity.PfeAssignment{}
	err := row.Scan(&a.ID, &a.PfeCode, &a.SubjectID, &a.AcademicYearID, &a.StudentID, &a.Student2ID, &a.Student3ID,
		&a.SupervisorID, &a.CoSupervisorID, &a.MemoireURL, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

// FindBySupervisor retourne les assignations supervisées par un enseignant.
func (r *PfeAssignmentRepository) FindBySupervisor(supervisorID string) ([]*entity.PfeAssignment, error) {
	query := `SELECT id, pfe_code, subject_id, academic_year_id, student_id, student2_id, student3_id,
		supervisor_id, co_supervisor_id, memoire_url, status, created_at, updated_at
		FROM pfe_assignments WHERE supervisor_id = ? OR co_supervisor_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, supervisorID, supervisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAssignments(rows)
}

// FindAll retourne toutes les assignations.
func (r *PfeAssignmentRepository) FindAll() ([]*entity.PfeAssignment, error) {
	rows, err := r.db.Query(`SELECT id, pfe_code, subject_id, academic_year_id, student_id, student2_id, student3_id,
		supervisor_id, co_supervisor_id, memoire_url, status, created_at, updated_at
		FROM pfe_assignments ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAssignments(rows)
}

// FindByAcademicYear retourne les assignations d'une année universitaire.
func (r *PfeAssignmentRepository) FindByAcademicYear(academicYearID string) ([]*entity.PfeAssignment, error) {
	query := `SELECT id, pfe_code, subject_id, academic_year_id, student_id, student2_id, student3_id,
		supervisor_id, co_supervisor_id, memoire_url, status, created_at, updated_at
		FROM pfe_assignments WHERE academic_year_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, academicYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAssignments(rows)
}

// FindByCompanySubject retourne les assignations dont le sujet appartient à une entreprise.
func (r *PfeAssignmentRepository) FindByCompanySubject(companyID string) ([]*entity.PfeAssignment, error) {
	query := `SELECT pa.id, pa.pfe_code, pa.subject_id, pa.academic_year_id, pa.student_id, pa.student2_id, pa.student3_id,
		pa.supervisor_id, pa.co_supervisor_id, pa.memoire_url, pa.status, pa.created_at, pa.updated_at
		FROM pfe_assignments pa
		INNER JOIN pfe_subjects ps ON ps.id = pa.subject_id
		WHERE ps.company_id = ? OR (ps.proposer_id = ? AND ps.proposer_role = 'company')
		ORDER BY pa.created_at DESC`
	rows, err := r.db.Query(query, companyID, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAssignments(rows)
}

// Insert crée une nouvelle assignation.
func (r *PfeAssignmentRepository) Insert(a *entity.PfeAssignment) error {
	query := `INSERT INTO pfe_assignments (id, pfe_code, subject_id, academic_year_id, student_id, student2_id, student3_id,
		supervisor_id, co_supervisor_id, memoire_url, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, a.ID, a.PfeCode, a.SubjectID, a.AcademicYearID, a.StudentID, a.Student2ID, a.Student3ID,
		a.SupervisorID, a.CoSupervisorID, a.MemoireURL, a.Status)
	return err
}

// UpdateStatus met à jour le statut d'une assignation.
func (r *PfeAssignmentRepository) UpdateStatus(id, status string) error {
	_, err := r.db.Exec(`UPDATE pfe_assignments SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
	return err
}

// UpdateMemoire met à jour l'URL du mémoire.
func (r *PfeAssignmentRepository) UpdateMemoire(id, memoireURL string) error {
	_, err := r.db.Exec(`UPDATE pfe_assignments SET memoire_url = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, memoireURL, id)
	return err
}

// Update met à jour une assignation complète.
func (r *PfeAssignmentRepository) Update(a *entity.PfeAssignment) error {
	query := `UPDATE pfe_assignments SET pfe_code = ?, subject_id = ?, student_id = ?, student2_id = ?, student3_id = ?,
		supervisor_id = ?, co_supervisor_id = ?, memoire_url = ?, status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, a.PfeCode, a.SubjectID, a.StudentID, a.Student2ID, a.Student3ID,
		a.SupervisorID, a.CoSupervisorID, a.MemoireURL, a.Status, a.ID)
	return err
}

func (r *PfeAssignmentRepository) scanAssignments(rows *sql.Rows) ([]*entity.PfeAssignment, error) {
	var assignments []*entity.PfeAssignment
	for rows.Next() {
		a := &entity.PfeAssignment{}
		if err := rows.Scan(&a.ID, &a.PfeCode, &a.SubjectID, &a.AcademicYearID, &a.StudentID, &a.Student2ID, &a.Student3ID,
			&a.SupervisorID, &a.CoSupervisorID, &a.MemoireURL, &a.Status, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, rows.Err()
}
