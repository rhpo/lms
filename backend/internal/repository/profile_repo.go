package repository

import (
	"database/sql"

	"pfe-backend/internal/entity"
	"pfe-backend/internal/shared/convert"
)

// ProfileRepository gère les opérations base de données pour les profils.
type ProfileRepository struct {
	db *sql.DB
}

// NewProfileRepository crée un nouveau ProfileRepository.
func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// FindByEmail cherche un profil par son email.
func (r *ProfileRepository) FindByEmail(email string) (*entity.Profile, error) {
	query := `SELECT id, role, full_name, email, avatar_url, is_active, created_at, updated_at 
		FROM profiles WHERE email = ?`
	row := r.db.QueryRow(query, email)

	var avatarURL sql.NullString
	profile := &entity.Profile{}
	err := row.Scan(
		&profile.ID,
		&profile.Role,
		&profile.FullName,
		&profile.Email,
		&avatarURL,
		&profile.IsActive,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	profile.AvatarURL = convert.StringPtr(avatarURL)
	return profile, nil
}

// FindByID cherche un profil par son ID.
func (r *ProfileRepository) FindByID(id string) (*entity.Profile, error) {
	query := `SELECT id, role, full_name, email, avatar_url, is_active, created_at, updated_at 
		FROM profiles WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var avatarURL sql.NullString
	profile := &entity.Profile{}
	err := row.Scan(
		&profile.ID,
		&profile.Role,
		&profile.FullName,
		&profile.Email,
		&avatarURL,
		&profile.IsActive,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	profile.AvatarURL = convert.StringPtr(avatarURL)
	return profile, nil
}

// FindAll retourne tous les profils.
func (r *ProfileRepository) FindAll() ([]*entity.Profile, error) {
	query := `SELECT id, role, full_name, email, avatar_url, is_active, created_at, updated_at 
		FROM profiles ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*entity.Profile
	for rows.Next() {
		var avatarURL sql.NullString
		profile := &entity.Profile{}
		err := rows.Scan(
			&profile.ID,
			&profile.Role,
			&profile.FullName,
			&profile.Email,
			&avatarURL,
			&profile.IsActive,
			&profile.CreatedAt,
			&profile.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		profile.AvatarURL = convert.StringPtr(avatarURL)
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

// Insert crée un nouveau profil.
func (r *ProfileRepository) Insert(p *entity.Profile) error {
	query := `INSERT INTO profiles (id, role, full_name, email, avatar_url, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	_, err := r.db.Exec(query, p.ID, p.Role, p.FullName, p.Email, convert.NullString(p.AvatarURL), p.IsActive)
	return err
}

// Update met à jour un profil.
func (r *ProfileRepository) Update(p *entity.Profile) error {
	query := `UPDATE profiles SET role = ?, full_name = ?, email = ?, avatar_url = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, p.Role, p.FullName, p.Email, convert.NullString(p.AvatarURL), p.IsActive, p.ID)
	return err
}

// UpdateAvatarURL met à jour l'avatar d'un profil.
func (r *ProfileRepository) UpdateAvatarURL(id, url string) error {
	_, err := r.db.Exec(`UPDATE profiles SET avatar_url = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, url, id)
	return err
}

// Delete supprime un profil.
func (r *ProfileRepository) Delete(id string) error {
	query := `DELETE FROM profiles WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
