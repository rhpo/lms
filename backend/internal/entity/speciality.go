package entity

import "time"

// Speciality représente une spécialité.
type Speciality struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	YearType  string    `json:"year_type"` // licence/master
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
