package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

// TeacherRepository gère les opérations base de données pour les enseignants.
type TeacherRepository struct {
	db *sql.DB
}

// NewTeacherRepository crée un nouveau TeacherRepository.
func NewTeacherRepository(db *sql.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

// FindByID cherche un enseignant par son ID avec son profil et ses domaines.
func (r *TeacherRepository) FindByID(id string) (*entity.Teacher, error) {
	query := `SELECT id, profile_id, grade, department, availability_status, unavailable_until, created_at, updated_at
		FROM teachers WHERE id = ?`
	row := r.db.QueryRow(query, id)

	t := &entity.Teacher{}
	err := row.Scan(
		&t.ID, &t.ProfileID, &t.Grade, &t.Department,
		&t.AvailabilityStatus, &t.UnavailableUntil,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

// FindByProfileID cherche un enseignant par son profile_id.
func (r *TeacherRepository) FindByProfileID(profileID string) (*entity.Teacher, error) {
	query := `SELECT id, profile_id, grade, department, availability_status, unavailable_until, created_at, updated_at
		FROM teachers WHERE profile_id = ?`
	row := r.db.QueryRow(query, profileID)

	t := &entity.Teacher{}
	err := row.Scan(
		&t.ID, &t.ProfileID, &t.Grade, &t.Department,
		&t.AvailabilityStatus, &t.UnavailableUntil,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

// FindAll retourne tous les enseignants.
func (r *TeacherRepository) FindAll() ([]*entity.Teacher, error) {
	query := `SELECT id, profile_id, grade, department, availability_status, unavailable_until, created_at, updated_at
		FROM teachers ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []*entity.Teacher
	for rows.Next() {
		t := &entity.Teacher{}
		if err := rows.Scan(
			&t.ID, &t.ProfileID, &t.Grade, &t.Department,
			&t.AvailabilityStatus, &t.UnavailableUntil,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}
	return teachers, rows.Err()
}

// UpdateAvailability met à jour le statut de disponibilité d'un enseignant.
func (r *TeacherRepository) UpdateAvailability(id string, status string, unavailableUntil *sql.NullTime) error {
	query := `UPDATE teachers SET availability_status = ?, unavailable_until = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, status, unavailableUntil, id)
	return err
}

// AddDomain ajoute un domaine à un enseignant.
func (r *TeacherRepository) AddDomain(teacherID, domainID string) error {
	_, err := r.db.Exec(`INSERT OR IGNORE INTO teacher_domains (teacher_id, domain_id) VALUES (?, ?)`, teacherID, domainID)
	return err
}

// RemoveDomain retire un domaine d'un enseignant.
func (r *TeacherRepository) RemoveDomain(teacherID, domainID string) error {
	_, err := r.db.Exec(`DELETE FROM teacher_domains WHERE teacher_id = ? AND domain_id = ?`, teacherID, domainID)
	return err
}

// GetDomains retourne les domaines d'un enseignant.
func (r *TeacherRepository) GetDomains(teacherID string) ([]*entity.Domain, error) {
	query := `SELECT d.id, d.name, d.created_at, d.updated_at
		FROM domains d
		INNER JOIN teacher_domains td ON td.domain_id = d.id
		WHERE td.teacher_id = ?`
	rows, err := r.db.Query(query, teacherID)
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

// FindAvailableTeachers retourne les enseignants disponibles pour une date donnée.
func (r *TeacherRepository) FindAvailableTeachers() ([]*entity.Teacher, error) {
	query := `SELECT id, profile_id, grade, department, availability_status, unavailable_until, created_at, updated_at
		FROM teachers WHERE availability_status = 'disponible' ORDER BY grade DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []*entity.Teacher
	for rows.Next() {
		t := &entity.Teacher{}
		if err := rows.Scan(
			&t.ID, &t.ProfileID, &t.Grade, &t.Department,
			&t.AvailabilityStatus, &t.UnavailableUntil,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}
	return teachers, rows.Err()
}

// Insert crée un nouvel enseignant.
func (r *TeacherRepository) Insert(t *entity.Teacher) error {
	query := `INSERT INTO teachers (id, profile_id, grade, department, availability_status, unavailable_until) 
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, t.ID, t.ProfileID, t.Grade, t.Department, t.AvailabilityStatus, t.UnavailableUntil)
	return err
}

// Update met à jour les informations d'un enseignant.
func (r *TeacherRepository) Update(t *entity.Teacher) error {
	query := `UPDATE teachers SET grade = ?, department = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, t.Grade, t.Department, t.ID)
	return err
}

// Delete supprime un enseignant.
func (r *TeacherRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM teacher_domains WHERE teacher_id = ?`, id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`DELETE FROM teachers WHERE id = ?`, id)
	return err
}
