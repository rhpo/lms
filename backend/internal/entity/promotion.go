package entity

import "time"

// Promotion représente une promotion.
type Promotion struct {
	ID           string    `json:"id"`
	Label        string    `json:"label"`
	AcademicYearID string    `json:"academic_year_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
