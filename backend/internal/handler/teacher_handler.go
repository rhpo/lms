package handler

import (
	"fmt"
	"time"

	"pfe-backend/internal/entity"
	"pfe-backend/internal/service"
	"pfe-backend/internal/shared/middleware"
	"pfe-backend/internal/shared/response"

	"github.com/gofiber/fiber/v3"
)

// TeacherHandler gère les endpoints enseignant.
type TeacherHandler struct {
	svc *service.TeacherService
}

// NewTeacherHandler crée un nouveau TeacherHandler.
func NewTeacherHandler(svc *service.TeacherService) *TeacherHandler {
	return &TeacherHandler{svc: svc}
}

// Dashboard retourne le tableau de bord enseignant.
func (h *TeacherHandler) Dashboard(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	data, err := h.svc.Dashboard(userID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, data)
}

// ListProposedSubjects liste les sujets proposés par l'enseignant.
func (h *TeacherHandler) ListProposedSubjects(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	subjects, err := h.svc.ListProposedSubjects(userID)
	if err != nil {
		return response.Error(c, err)
	}
	if subjects == nil {
		subjects = []*entity.PfeSubject{}
	}
	return response.OK(c, subjects)
}

// CreateProposedSubject crée un nouveau sujet proposé.
func (h *TeacherHandler) CreateProposedSubject(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	var req struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		GroupType   string   `json:"group_type"`
		DomainIDs   []string `json:"domain_ids"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if req.Title == "" || req.Description == "" {
		return response.ValidationError(c, "Titre et description requis")
	}

	subject := &entity.PfeSubject{
		Title:        req.Title,
		Description:  req.Description,
		GroupType:    req.GroupType,
		ProposerID:   userID,
		ProposerRole: "teacher",
		Status:       "en_attente",
	}
	if err := h.svc.CreateProposedSubject(subject, req.DomainIDs); err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, subject)
}

// GetProposedSubject retourne un sujet proposé.
func (h *TeacherHandler) GetProposedSubject(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	id := c.Params("id")
	subject, err := h.svc.GetProposedSubject(userID, id)
	if err != nil {
		return response.Error(c, err)
	}
	if subject == nil {
		return response.NotFound(c, "Sujet introuvable")
	}
	return response.OK(c, subject)
}

// UpdateProposedSubject met à jour un sujet proposé.
func (h *TeacherHandler) UpdateProposedSubject(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	id := c.Params("id")
	var req entity.PfeSubject
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	req.ID = id
	if err := h.svc.UpdateProposedSubject(userID, &req); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Sujet mis à jour"})
}

// ListCandidats liste les candidats pour un sujet.
func (h *TeacherHandler) ListCandidats(c fiber.Ctx) error {
	id := c.Params("id")
	candidats, err := h.svc.ListCandidats(id)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, candidats)
}

// AcceptCandidat accepte ou refuse un étudiant pour un sujet.
func (h *TeacherHandler) AcceptCandidat(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		StudentID string `json:"student_id"`
		Action    string `json:"action"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if req.Action == "accept" {
		if err := h.svc.AcceptCandidat(id, req.StudentID); err != nil {
			return response.Error(c, err)
		}
		return response.OK(c, map[string]string{"message": "Candidat accepté"})
	}
	if err := h.svc.RejectCandidat(id, req.StudentID); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Candidat refusé"})
}

// ListSubjectsToValidate liste les sujets à valider.
func (h *TeacherHandler) ListSubjectsToValidate(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	subjects, err := h.svc.ListSubjectsToValidate(userID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, subjects)
}

// GetSubjectToValidate retourne un sujet à valider.
func (h *TeacherHandler) GetSubjectToValidate(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	id := c.Params("id")
	subject, err := h.svc.GetSubjectToValidate(userID, id)
	if err != nil {
		return response.Error(c, err)
	}
	if subject == nil {
		return response.NotFound(c, "Sujet introuvable")
	}
	return response.OK(c, subject)
}

// ValidateSubject valide ou refuse un sujet.
func (h *TeacherHandler) ValidateSubject(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	id := c.Params("id")
	var req struct {
		Decision string `json:"decision"`
		Comment  string `json:"comment"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.ValidateSubject(userID, id, req.Decision, req.Comment); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Validation enregistrée"})
}

// ListSupervisedPFEs liste les PFE encadrés.
func (h *TeacherHandler) ListSupervisedPFEs(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	assignments, err := h.svc.ListSupervisedPFEs(userID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, assignments)
}

// GetSupervisedPFE retourne un PFE encadré.
func (h *TeacherHandler) GetSupervisedPFE(c fiber.Ctx) error {
	id := c.Params("id")
	assignment, err := h.svc.GetSupervisedPFE(id)
	if err != nil {
		return response.Error(c, err)
	}
	if assignment == nil {
		return response.NotFound(c, "PFE introuvable")
	}
	return response.OK(c, assignment)
}

// AddMeeting ajoute un meeting de suivi à un PFE encadré.
func (h *TeacherHandler) AddMeeting(c fiber.Ctx) error {
	id := c.Params("id")
	var req entity.PfeProgressReport
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if req.MeetingType == "" || req.Duration == 0 {
		return response.ValidationError(c, "meeting_type et duration sont requis")
	}
	req.ID = fmt.Sprintf("report-%d", time.Now().UnixNano())
	req.AssignmentID = id
	if req.Status == "" {
		req.Status = "en_cours"
	}
	if err := h.svc.AddMeeting(&req); err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, req)
}

// SubmitEvaluation soumet l'évaluation de l'encadrant.
func (h *TeacherHandler) SubmitEvaluation(c fiber.Ctx) error {
	id := c.Params("id")
	userID := middleware.GetProfileID(c)
	var req struct {
		Criterion5 float64 `json:"criterion5"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.SubmitEvaluation(id, userID, req.Criterion5); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Évaluation soumise"})
}

// ListJuryDuties liste les obligations de jury.
func (h *TeacherHandler) ListJuryDuties(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	duties, err := h.svc.ListJuryDuties(userID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, duties)
}

// GetJuryDuty retourne une obligation de jury.
func (h *TeacherHandler) GetJuryDuty(c fiber.Ctx) error {
	id := c.Params("id")
	duty, err := h.svc.GetJuryDuty(id)
	if err != nil {
		return response.Error(c, err)
	}
	if duty == nil {
		return response.NotFound(c, "Obligation jury introuvable")
	}
	return response.OK(c, duty)
}

// UpdateAvailability met à jour la disponibilité de l'enseignant.
func (h *TeacherHandler) UpdateAvailability(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	var req struct {
		Availability     string `json:"availability_status"`
		UnavailableUntil string `json:"unavailable_until"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.UpdateAvailability(userID, req.Availability, req.UnavailableUntil); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Disponibilité mise à jour"})
}

// ListNotifications liste les notifications de l'enseignant.
func (h *TeacherHandler) ListNotifications(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	notifications, err := h.svc.ListNotifications(userID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, notifications)
}
