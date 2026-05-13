package entity

import (
	"time"
)

// Wish représente un vœu d'étudiant pour un sujet PFE.
type Wish struct {
	ID             string    `json:"id"`
	StudentID      string    `json:"student_id"`
	SubjectID      string    `json:"subject_id"`
	AcademicYearID string    `json:"academic_year_id"`
	Status         string    `json:"status"` // en_attente/accepte/refuse
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations
	Student      *Student      `json:"student,omitempty"`
	Subject      *PfeSubject   `json:"subject,omitempty"`
	AcademicYear *AcademicYear `json:"academic_year,omitempty"`
}
