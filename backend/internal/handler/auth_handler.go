package handler

import (
	"pfe-backend/internal/config"
	"pfe-backend/internal/service"
	"pfe-backend/internal/shared/middleware"
	"pfe-backend/internal/shared/response"

	"github.com/gofiber/fiber/v3"
)

// AuthHandler gère les endpoints d'authentification.
type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config
}

// NewAuthHandler crée un nouveau AuthHandler.
func NewAuthHandler(authService *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

// devLoginRequest est la structure de la requête de login.
type devLoginRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// DevLogin gère POST /api/auth/dev-login.
func (h *AuthHandler) DevLogin(c fiber.Ctx) error {
	if !h.cfg.IsDevelopment() {
		return response.NotFound(c, "Endpoint non disponible")
	}

	var req devLoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if req.Email == "" {
		return response.ValidationError(c, "L'email est requis")
	}

	result, err := h.authService.DevLogin(req.Email)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, result)
}

// Me gère GET /api/auth/me.
func (h *AuthHandler) Me(c fiber.Ctx) error {
	profileID := middleware.GetProfileID(c)
	if profileID == "" {
		return response.Unauthorized(c, "Non authentifié")
	}

	profile, err := h.authService.GetProfile(profileID)
	if err != nil {
		return response.NotFound(c, "Profil introuvable")
	}

	return response.OK(c, profile)
}

// Logout gère POST /api/auth/logout.
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	// En JWT, la déconnexion est gérée côté client (suppression du token).
	// On pourrait ajouter une blacklist de tokens si nécessaire.
	return response.OK(c, map[string]string{"message": "Déconnexion réussie"})
}
