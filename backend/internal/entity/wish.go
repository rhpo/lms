package entity

import (
	"time"
)

// Wish représente un vœu d'étudiant pour un sujet PFE.
type Wish struct {
	ID             int64     `json:"id"`
	StudentID      int64     `json:"student_id"`
	SubjectID      int64     `json:"subject_id"`
	AcademicYearID int64     `json:"academic_year_id"`
	Status         string    `json:"status"` // en_attente/accepte/refuse
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations
	Student      *Student      `json:"student,omitempty"`
	Subject      *PfeSubject   `json:"subject,omitempty"`
	AcademicYear *AcademicYear `json:"academic_year,omitempty"`
}
