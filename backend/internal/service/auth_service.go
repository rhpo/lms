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
	teacherRepo *repository.TeacherRepository
	studentRepo *repository.StudentRepository
	companyRepo *repository.CompanyRepository
	cfg         *config.Config
}

// NewAuthService crée un nouveau AuthService.
func NewAuthService(profileRepo *repository.ProfileRepository, teacherRepo *repository.TeacherRepository, studentRepo *repository.StudentRepository, companyRepo *repository.CompanyRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		profileRepo: profileRepo,
		teacherRepo: teacherRepo,
		studentRepo: studentRepo,
		companyRepo: companyRepo,
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

	if profile.Role == "teacher" || profile.Role == "admin" {
		t, _ := s.teacherRepo.FindByProfileID(profile.ID)
		profile.Teacher = t
	} else if profile.Role == "student" {
		st, _ := s.studentRepo.FindByProfileID(profile.ID)
		profile.Student = st
	} else if profile.Role == "company" {
		c, _ := s.companyRepo.FindByProfileID(profile.ID)
		profile.Company = c
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
func (s *AuthService) GetProfile(id int64) (*entity.Profile, error) {
	profile, err := s.profileRepo.FindByID(id)
	if err != nil {
		return nil, apperror.Internal("Erreur base de données")
	}
	if profile == nil {
		return nil, apperror.NotFound("Profil introuvable")
	}

	if profile.Role == "teacher" || profile.Role == "admin" {
		t, _ := s.teacherRepo.FindByProfileID(profile.ID)
		profile.Teacher = t
	} else if profile.Role == "student" {
		st, _ := s.studentRepo.FindByProfileID(profile.ID)
		profile.Student = st
	} else if profile.Role == "company" {
		c, _ := s.companyRepo.FindByProfileID(profile.ID)
		profile.Company = c
	}

	return profile, nil
}

// RegisterCompanyRequest est la requête d'inscription d'un employé d'entreprise.
type RegisterCompanyRequest struct {
	// Employee info
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Position string `json:"position"`
	Phone    string `json:"phone"`

	// Existing company (mutually exclusive with NewCompany)
	CompanyID int64 `json:"company_id"`

	// New company (when CompanyID == 0)
	CompanyName  string `json:"company_name"`
	Sector       string `json:"sector"`
	Description  string `json:"description"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`
}

// RegisterCompanyEmployee crée un profil entreprise + company record.
// Si company_id est fourni, l'employé rejoint une entreprise existante (vérifiée).
// Sinon, une nouvelle entreprise est créée (non vérifiée, en attente de validation admin).
func (s *AuthService) RegisterCompanyEmployee(req *RegisterCompanyRequest) (*DevLoginResponse, error) {
	// Check email uniqueness
	existing, _ := s.profileRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, apperror.Conflict("Un compte existe déjà avec cet email")
	}

	// Validate
	if req.FullName == "" || req.Email == "" {
		return nil, apperror.BadRequest("Le nom complet et l'email sont obligatoires")
	}

	// Create profile
	profile := &entity.Profile{
		Role:     "company",
		FullName: req.FullName,
		Email:    req.Email,
		IsActive: true,
	}
	if err := s.profileRepo.Insert(profile); err != nil {
		return nil, apperror.Internal("Erreur création du profil")
	}

	var company *entity.Company

	if req.CompanyID > 0 {
		// Join existing verified company
		var err error
		company, err = s.companyRepo.FindByID(req.CompanyID)
		if err != nil || company == nil {
			return nil, apperror.NotFound("Entreprise introuvable")
		}
		if !company.IsVerified {
			return nil, apperror.Forbidden("Cette entreprise n'est pas encore vérifiée")
		}
		// Create a company record linked to this profile, copying company info
		newCompany := &entity.Company{
			ProfileID:    profile.ID,
			CompanyName:  company.CompanyName,
			Sector:       company.Sector,
			Description:  company.Description,
			LogoURL:      company.LogoURL,
			ContactEmail: company.ContactEmail,
			ContactPhone: company.ContactPhone,
			Website:      company.Website,
			IsVerified:   true, // auto-verified since the company is already verified
		}
		if err := s.companyRepo.Insert(newCompany); err != nil {
			return nil, apperror.Internal("Erreur création du compte entreprise")
		}
		company = newCompany
	} else {
		// Create new company (pending admin verification)
		if req.CompanyName == "" {
			return nil, apperror.BadRequest("Le nom de l'entreprise est obligatoire")
		}
		companyName := req.CompanyName
		sector := req.Sector
		description := req.Description
		contactEmail := req.ContactEmail
		contactPhone := req.ContactPhone

		company = &entity.Company{
			ProfileID:    profile.ID,
			CompanyName:  &companyName,
			Sector:       &sector,
			Description:  &description,
			ContactEmail: &contactEmail,
			ContactPhone: &contactPhone,
			IsVerified:   false,
		}
		if err := s.companyRepo.Insert(company); err != nil {
			return nil, apperror.Internal("Erreur création de l'entreprise")
		}
	}

	profile.Company = company

	token, err := s.generateToken(profile)
	if err != nil {
		return nil, apperror.Internal("Erreur génération du token")
	}

	return &DevLoginResponse{
		Token:   token,
		Profile: profile,
	}, nil
}

// ListVerifiedCompanies retourne les entreprises vérifiées (pour l'autocomplete d'inscription).
func (s *AuthService) ListVerifiedCompanies() ([]*entity.Company, error) {
	return s.companyRepo.FindAllVerified()
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
