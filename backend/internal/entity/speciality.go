package entity

import "time"

// Speciality représente une spécialité.
type Speciality struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	YearType     string    `json:"year_type"` // licence/master
	DepartmentID *int64    `json:"department_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	Department *Department `json:"department,omitempty"`
}
