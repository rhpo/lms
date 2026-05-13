package entity

import "time"

// Student représente un étudiant.
type Student struct {
	ID            string    `json:"id"`
	ProfileID     string    `json:"profile_id"`
	StudentNumber *string   `json:"student_number"`
	SpecialityID  *string   `json:"specialty_id"`
	Level         *string   `json:"level"`
	PromotionID   *string   `json:"promotion_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Profile    *Profile    `json:"profile,omitempty"`
	Speciality *Speciality `json:"speciality,omitempty"`
	Promotion  *Promotion  `json:"promotion,omitempty"`
}
