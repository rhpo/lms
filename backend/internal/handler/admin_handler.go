package handler

import (
	"pfe-backend/internal/entity"
	"pfe-backend/internal/service"
	"pfe-backend/internal/shared/middleware"
	"pfe-backend/internal/shared/response"

	"github.com/gofiber/fiber/v3"
)

// AdminHandler gère les endpoints admin.
type AdminHandler struct {
	svc *service.AdminService
}

// NewAdminHandler crée un nouveau AdminHandler.
func NewAdminHandler(svc *service.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

// Dashboard retourne les statistiques du tableau de bord.
func (h *AdminHandler) Dashboard(c fiber.Ctx) error {
	data, err := h.svc.Dashboard()
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, data)
}

// ListUsers liste tous les utilisateurs.
func (h *AdminHandler) ListUsers(c fiber.Ctx) error {
	users, err := h.svc.ListUsers()
	if err != nil {
		return response.Error(c, err)
	}
	if users == nil {
		users = []*entity.Profile{}
	}
	return response.OK(c, users)
}

// GetUser retourne un utilisateur par son ID.
func (h *AdminHandler) GetUser(c fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.svc.GetUser(id)
	if err != nil {
		return response.Error(c, err)
	}
	if user == nil {
		return response.NotFound(c, "Utilisateur introuvable")
	}
	return response.OK(c, user)
}

// CreateUser crée un nouvel utilisateur. L'ID est généré côté serveur.
func (h *AdminHandler) CreateUser(c fiber.Ctx) error {
	var req struct {
		Role     string `json:"role"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if req.Role == "" || req.FullName == "" || req.Email == "" {
		return response.ValidationError(c, "role, full_name et email sont requis")
	}

	profile := &entity.Profile{
		Role:     req.Role,
		FullName: req.FullName,
		Email:    req.Email,
		IsActive: true,
	}
	if err := h.svc.CreateUser(profile); err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, profile)
}

// UpdateUser met à jour un utilisateur.
func (h *AdminHandler) UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")
	var req entity.Profile
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	req.ID = id
	if err := h.svc.UpdateUser(&req); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, req)
}

// UserAction gère les actions sur un utilisateur : deactivate, reactivate, transfer-admin.
func (h *AdminHandler) UserAction(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Action string `json:"action"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.UserAction(id, req.Action); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Action effectuée"})
}

// ImportUsersCSV importe des utilisateurs en masse depuis un CSV.
func (h *AdminHandler) ImportUsersCSV(c fiber.Ctx) error {
	var req struct {
		CSVData string `json:"csv_data"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.ImportUsersCSV(req.CSVData); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Import effectué"})
}

// ListCompanies liste toutes les entreprises.
func (h *AdminHandler) ListCompanies(c fiber.Ctx) error {
	companies, err := h.svc.ListCompanies()
	if err != nil {
		return response.Error(c, err)
	}
	if companies == nil {
		companies = []*entity.Company{}
	}
	return response.OK(c, companies)
}

// CompanyAction gère les actions sur une entreprise : validate, reject, update.
func (h *AdminHandler) CompanyAction(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Action string `json:"action"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.CompanyAction(id, req.Action); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Action effectuée"})
}

// ListReports liste tous les reports.
func (h *AdminHandler) ListReports(c fiber.Ctx) error {
	reports, err := h.svc.ListReports()
	if err != nil {
		return response.Error(c, err)
	}
	if reports == nil {
		reports = []*entity.CompanyReport{}
	}
	return response.OK(c, reports)
}

// ReportAction gère les actions sur un report : resolve, reject.
func (h *AdminHandler) ReportAction(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Action string `json:"action"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.ReportAction(id, req.Action); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Action effectuée"})
}

// ListSubjects liste tous les sujets PFE.
func (h *AdminHandler) ListSubjects(c fiber.Ctx) error {
	subjects, err := h.svc.ListSubjects()
	if err != nil {
		return response.Error(c, err)
	}
	if subjects == nil {
		subjects = []*entity.PfeSubject{}
	}
	return response.OK(c, subjects)
}

// GetSubject retourne un sujet par son ID.
func (h *AdminHandler) GetSubject(c fiber.Ctx) error {
	id := c.Params("id")
	subject, err := h.svc.GetSubject(id)
	if err != nil {
		return response.Error(c, err)
	}
	if subject == nil {
		return response.NotFound(c, "Sujet introuvable")
	}
	return response.OK(c, subject)
}

// SubjectAction gère les actions admin sur un sujet.
func (h *AdminHandler) SubjectAction(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Action    string `json:"action"`
		Validator string `json:"validator_id"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.SubjectAction(id, req.Action, req.Validator); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Action effectuée"})
}

// ListAssignments liste toutes les affectations PFE.
func (h *AdminHandler) ListAssignments(c fiber.Ctx) error {
	assignments, err := h.svc.ListAssignments()
	if err != nil {
		return response.Error(c, err)
	}
	if assignments == nil {
		assignments = []*entity.PfeAssignment{}
	}
	return response.OK(c, assignments)
}

// GetAssignment retourne une affectation par son ID.
func (h *AdminHandler) GetAssignment(c fiber.Ctx) error {
	id := c.Params("id")
	assignment, err := h.svc.GetAssignment(id)
	if err != nil {
		return response.Error(c, err)
	}
	if assignment == nil {
		return response.NotFound(c, "Affectation introuvable")
	}
	return response.OK(c, assignment)
}

// ListDefenses liste toutes les soutenances.
func (h *AdminHandler) ListDefenses(c fiber.Ctx) error {
	defenses, err := h.svc.ListDefenses()
	if err != nil {
		return response.Error(c, err)
	}
	if defenses == nil {
		defenses = []*entity.Defense{}
	}
	return response.OK(c, defenses)
}

// CreateDefense crée une nouvelle soutenance.
func (h *AdminHandler) CreateDefense(c fiber.Ctx) error {
	var req struct {
		AssignmentID string `json:"assignment_id"`
		PresidentID  string `json:"president_id"`
		MemberID     string `json:"member_id"`
		ScheduledAt  string `json:"scheduled_at"`
		Room         string `json:"room"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	defense, err := h.svc.CreateDefense(req.AssignmentID, req.PresidentID, req.MemberID, req.ScheduledAt, req.Room)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, defense)
}

// GetDefense retourne une soutenance par son ID.
func (h *AdminHandler) GetDefense(c fiber.Ctx) error {
	id := c.Params("id")
	defense, err := h.svc.GetDefense(id)
	if err != nil {
		return response.Error(c, err)
	}
	if defense == nil {
		return response.NotFound(c, "Soutenance introuvable")
	}
	return response.OK(c, defense)
}

// RecommendJury recommande un jury pour un PFE.
func (h *AdminHandler) RecommendJury(c fiber.Ctx) error {
	pfeID := c.Query("pfe_id")
	if pfeID == "" {
		return response.ValidationError(c, "pfe_id requis")
	}
	recommendation, err := h.svc.RecommendJury(pfeID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, recommendation)
}

// SubmitGrade soumet une note pour une soutenance.
func (h *AdminHandler) SubmitGrade(c fiber.Ctx) error {
	id := c.Params("id")
	callerID := middleware.GetProfileID(c)
	var req struct {
		Criterion1 float64 `json:"criterion1"`
		Criterion2 float64 `json:"criterion2"`
		Criterion3 float64 `json:"criterion3"`
		Criterion4 float64 `json:"criterion4"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.SubmitGrade(id, callerID, req.Criterion1, req.Criterion2, req.Criterion3, req.Criterion4); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Note soumise"})
}

// ResolveGrade résout la note finale d'une soutenance.
func (h *AdminHandler) ResolveGrade(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Choice     string             `json:"choice"`
		Criterion1 float64            `json:"criterion1"`
		Criterion2 float64            `json:"criterion2"`
		Criterion3 float64            `json:"criterion3"`
		Criterion4 float64            `json:"criterion4"`
		Grades     map[string]float64 `json:"grades"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	svcReq := service.ResolveGradeRequest{
		Choice:     req.Choice,
		Criterion1: req.Criterion1,
		Criterion2: req.Criterion2,
		Criterion3: req.Criterion3,
		Criterion4: req.Criterion4,
		Grades:     req.Grades,
	}
	if err := h.svc.ResolveGrade(id, svcReq); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Note résolue"})
}

// ConfirmJury confirme la participation d'un jury.
func (h *AdminHandler) ConfirmJury(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.ConfirmJury(id); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Jury confirmé"})
}

// DeclineJury décline la participation d'un jury.
func (h *AdminHandler) DeclineJury(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.DeclineJury(id); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Jury décliné"})
}

// ListDeadlines liste les délais configurés.
func (h *AdminHandler) ListDeadlines(c fiber.Ctx) error {
	return h.ListAcademicYears(c)
}

// UpdateDeadlines met à jour les délais.
func (h *AdminHandler) UpdateDeadlines(c fiber.Ctx) error {
	var req struct {
		OpenAt  string `json:"submission_open_at"`
		CloseAt string `json:"submission_close_at"`
		MaxWish int    `json:"max_wishes"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.UpdateDeadlines(req.OpenAt, req.CloseAt, req.MaxWish); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Délais mis à jour"})
}

// ListSpecialities liste toutes les spécialités.
func (h *AdminHandler) ListSpecialities(c fiber.Ctx) error {
	specialities, err := h.svc.ListSpecialities()
	if err != nil {
		return response.Error(c, err)
	}
	if specialities == nil {
		specialities = []*entity.Speciality{}
	}
	return response.OK(c, specialities)
}

// CreateSpeciality crée une spécialité.
func (h *AdminHandler) CreateSpeciality(c fiber.Ctx) error {
	var req entity.Speciality
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.CreateSpeciality(&req); err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, req)
}

// DeleteSpeciality supprime une spécialité.
func (h *AdminHandler) DeleteSpeciality(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.DeleteSpeciality(id); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Spécialité supprimée"})
}

// ListDomains liste tous les domaines.
func (h *AdminHandler) ListDomains(c fiber.Ctx) error {
	domains, err := h.svc.ListDomains()
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, domains)
}

// CreateDomain crée un domaine.
func (h *AdminHandler) CreateDomain(c fiber.Ctx) error {
	var req entity.Domain
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.CreateDomain(&req); err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, req)
}

// DeleteDomain supprime un domaine.
func (h *AdminHandler) DeleteDomain(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.DeleteDomain(id); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Domaine supprimé"})
}

// ListPromotions liste toutes les promotions.
func (h *AdminHandler) ListPromotions(c fiber.Ctx) error {
	promotions, err := h.svc.ListPromotions()
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, promotions)
}

// CreatePromotion crée une promotion.
func (h *AdminHandler) CreatePromotion(c fiber.Ctx) error {
	var req entity.Promotion
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.CreatePromotion(&req); err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, req)
}

// DeletePromotion supprime une promotion.
func (h *AdminHandler) DeletePromotion(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.DeletePromotion(id); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Promotion supprimée"})
}

// GetStatistics retourne les statistiques globales.
func (h *AdminHandler) Statistics(c fiber.Ctx) error {
	stats, err := h.svc.GetStatistics()
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, stats)
}

// AuditLog retourne les logs d'audit.
func (h *AdminHandler) AuditLog(c fiber.Ctx) error {
	logs, err := h.svc.AuditLog()
	if err != nil {
		return response.Error(c, err)
	}
	if logs == nil {
		logs = []*entity.AuditLog{}
	}
	return response.OK(c, logs)
}

// ExportAffectations exporte la liste des affectations.
func (h *AdminHandler) ExportAffectations(c fiber.Ctx) error {
	affectations, err := h.svc.ListAssignments()
	if err != nil {
		return response.Error(c, err)
	}
	if affectations == nil {
		affectations = []*entity.PfeAssignment{}
	}
	return response.OK(c, affectations)
}

// ExportPlannings exporte la liste des plannings de soutenance.
func (h *AdminHandler) ExportPlannings(c fiber.Ctx) error {
	defenses, err := h.svc.ListDefenses()
	if err != nil {
		return response.Error(c, err)
	}
	if defenses == nil {
		defenses = []*entity.Defense{}
	}
	return response.OK(c, defenses)
}

// ExportStatistics exporte les statistiques.
func (h *AdminHandler) ExportStatistics(c fiber.Ctx) error {
	stats, err := h.svc.GetStatistics()
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, stats)
}

// ListAcademicYears liste toutes les années académiques.
func (h *AdminHandler) ListAcademicYears(c fiber.Ctx) error {
	years, err := h.svc.ListAcademicYears()
	if err != nil {
		return response.Error(c, err)
	}
	if years == nil {
		years = []*entity.AcademicYear{}
	}
	return response.OK(c, years)
}

// CreateAcademicYear crée une année académique.
func (h *AdminHandler) CreateAcademicYear(c fiber.Ctx) error {
	var req entity.AcademicYear
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.CreateAcademicYear(&req); err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, req)
}

// CloseAcademicYear ferme une année académique.
func (h *AdminHandler) CloseAcademicYear(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.CloseAcademicYear(id); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Année académique clôturée"})
}
