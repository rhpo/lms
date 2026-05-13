package entity

import (
	"database/sql"
	"time"
)

// PfeAssignment représente l'affectation d'étudiants à un sujet PFE.
type PfeAssignment struct {
	ID             string         `json:"id"`
	PfeCode        string         `json:"pfe_code"` // PFE-ISIL-2025-001
	SubjectID      string         `json:"subject_id"`
	AcademicYearID string         `json:"academic_year_id"`
	StudentID      string         `json:"student_id"`
	Student2ID     sql.NullString `json:"student2_id"`
	Student3ID     sql.NullString `json:"student3_id"`
	SupervisorID   string         `json:"supervisor_id"`
	CoSupervisorID sql.NullString `json:"co_supervisor_id"`
	MemoireURL     sql.NullString `json:"memoire_url"`
	Status         string         `json:"status"` // en_cours/memoire_soumis/soutenance_planifiee/valide/refuse
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`

	// Relations
	Subject      *PfeSubject   `json:"subject,omitempty"`
	AcademicYear *AcademicYear `json:"academic_year,omitempty"`
	Student      *Student      `json:"student,omitempty"`
	Student2     *Student      `json:"student2,omitempty"`
	Student3     *Student      `json:"student3,omitempty"`
	Supervisor   *Teacher      `json:"supervisor,omitempty"`
	CoSupervisor *Teacher      `json:"co_supervisor,omitempty"`
}
