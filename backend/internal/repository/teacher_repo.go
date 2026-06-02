package repository

import (
	"database/sql"
	"pfe-backend/internal/entity"
)

type TeacherRepository struct {
	db *sql.DB
}

func NewTeacherRepository(db *sql.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

func (r *TeacherRepository) FindByID(id int64) (*entity.Teacher, error) {
	query := `SELECT id, profile_id, grade, department_id, availability_status, unavailable_until, created_at, updated_at
		FROM teachers WHERE id = ?`
	row := r.db.QueryRow(query, id)
	t := &entity.Teacher{}
	err := row.Scan(&t.ID, &t.ProfileID, &t.Grade, &t.DepartmentID, &t.AvailabilityStatus, &t.UnavailableUntil, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

func (r *TeacherRepository) FindByProfileID(profileID int64) (*entity.Teacher, error) {
	query := `SELECT id, profile_id, grade, department_id, availability_status, unavailable_until, created_at, updated_at
		FROM teachers WHERE profile_id = ?`
	row := r.db.QueryRow(query, profileID)
	t := &entity.Teacher{}
	err := row.Scan(&t.ID, &t.ProfileID, &t.Grade, &t.DepartmentID, &t.AvailabilityStatus, &t.UnavailableUntil, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

func (r *TeacherRepository) FindAll() ([]*entity.Teacher, error) {
	query := `SELECT id, profile_id, grade, department_id, availability_status, unavailable_until, created_at, updated_at
		FROM teachers ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teachers []*entity.Teacher
	for rows.Next() {
		t := &entity.Teacher{}
		if err := rows.Scan(&t.ID, &t.ProfileID, &t.Grade, &t.DepartmentID, &t.AvailabilityStatus, &t.UnavailableUntil, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}
	return teachers, rows.Err()
}

func (r *TeacherRepository) UpdateAvailability(id int64, status string, unavailableUntil *sql.NullTime) error {
	query := `UPDATE teachers SET availability_status = ?, unavailable_until = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, status, unavailableUntil, id)
	return err
}

func (r *TeacherRepository) AddDomain(teacherID, domainID int64) error {
	_, err := r.db.Exec(`INSERT OR IGNORE INTO teacher_domains (teacher_id, domain_id) VALUES (?, ?)`, teacherID, domainID)
	return err
}

func (r *TeacherRepository) RemoveDomain(teacherID, domainID int64) error {
	_, err := r.db.Exec(`DELETE FROM teacher_domains WHERE teacher_id = ? AND domain_id = ?`, teacherID, domainID)
	return err
}

func (r *TeacherRepository) GetDomains(teacherID int64) ([]*entity.Domain, error) {
	query := `SELECT d.id, d.name, d.created_at, d.updated_at
		FROM domains d INNER JOIN teacher_domains td ON td.domain_id = d.id WHERE td.teacher_id = ?`
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

func (r *TeacherRepository) FindAvailableTeachers() ([]*entity.Teacher, error) {
	query := `SELECT id, profile_id, grade, department_id, availability_status, unavailable_until, created_at, updated_at
		FROM teachers WHERE availability_status = 'disponible' ORDER BY grade DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teachers []*entity.Teacher
	for rows.Next() {
		t := &entity.Teacher{}
		if err := rows.Scan(&t.ID, &t.ProfileID, &t.Grade, &t.DepartmentID, &t.AvailabilityStatus, &t.UnavailableUntil, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}
	return teachers, rows.Err()
}

func (r *TeacherRepository) Insert(t *entity.Teacher) error {
	query := `INSERT INTO teachers (profile_id, grade, department_id, availability_status, unavailable_until)
		VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, t.ProfileID, t.Grade, t.DepartmentID, t.AvailabilityStatus, t.UnavailableUntil)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = id
	return nil
}

func (r *TeacherRepository) Update(t *entity.Teacher) error {
	query := `UPDATE teachers SET grade = ?, department_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, t.Grade, t.DepartmentID, t.ID)
	return err
}

func (r *TeacherRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM teacher_domains WHERE teacher_id = ?`, id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`DELETE FROM teachers WHERE id = ?`, id)
	return err
}
