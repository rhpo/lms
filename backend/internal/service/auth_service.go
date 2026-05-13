package service

import (
	"time"

	"pfe-backend/internal/config"
	"pfe-backend/internal/entity"
	"pfe-backend/internal/repository"
	"pfe-backend/internal/shared/apperror"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService gère la logique d'authentification.
type AuthService struct {
	profileRepo *repository.ProfileRepository
	cfg         *config.Config
}

// NewAuthService crée un nouveau AuthService.
func NewAuthService(profileRepo *repository.ProfileRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		profileRepo: profileRepo,
		cfg:         cfg,
	}
}

// DevLoginResponse est la réponse du login de développement.
type DevLoginResponse struct {
	Token   string          `json:"token"`
	Profile *entity.Profile `json:"profile"`
}

// DevLogin connecte un utilisateur par email en mode développement.
func (s *AuthService) DevLogin(email string) (*DevLoginResponse, error) {
	if !s.cfg.IsDevelopment() {
		return nil, apperror.NotFound("Endpoint non disponible")
	}

	profile, err := s.profileRepo.FindByEmail(email)
	if err != nil {
		return nil, apperror.Internal("Erreur base de données")
	}
	if profile == nil {
		return nil, apperror.NotFound("Aucun profil trouvé avec cet email")
	}
	if !profile.IsActive {
		return nil, apperror.Forbidden("Compte désactivé")
	}

	token, err := s.generateToken(profile)
	if err != nil {
		return nil, apperror.Internal("Erreur génération du token")
	}

	return &DevLoginResponse{
		Token:   token,
		Profile: profile,
	}, nil
}

// GetProfile récupère un profil par son ID.
func (s *AuthService) GetProfile(id string) (*entity.Profile, error) {
	profile, err := s.profileRepo.FindByID(id)
	if err != nil {
		return nil, apperror.Internal("Erreur base de données")
	}
	if profile == nil {
		return nil, apperror.NotFound("Profil introuvable")
	}
	return profile, nil
}

// generateToken génère un token JWT pour un profil.
func (s *AuthService) generateToken(profile *entity.Profile) (string, error) {
	claims := jwt.MapClaims{
		"sub":   profile.ID,
		"role":  profile.Role,
		"email": profile.Email,
		"name":  profile.FullName,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
