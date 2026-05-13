package entity

import "time"

// Profile représente un profil utilisateur.
type Profile struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	AvatarURL *string   `json:"avatar_url"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
