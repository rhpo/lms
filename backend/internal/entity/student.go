package entity

import "time"

// Student représente un étudiant.
type Student struct {
	ID            int64     `json:"id"`
	ProfileID     int64     `json:"profile_id"`
	StudentNumber *string   `json:"student_number"`
	SpecialityID  *int64    `json:"speciality_id"`
	Level         *string   `json:"level"`
	PromotionID   *int64    `json:"promotion_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Profile    *Profile    `json:"profile,omitempty"`
	Speciality *Speciality `json:"speciality,omitempty"`
	Promotion  *Promotion  `json:"promotion,omitempty"`
}
