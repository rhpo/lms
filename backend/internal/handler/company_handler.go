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

// CompanyHandler gère les endpoints entreprise.
type CompanyHandler struct {
	svc *service.CompanyService
}

// NewCompanyHandler crée un nouveau CompanyHandler.
func NewCompanyHandler(svc *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{svc: svc}
}

// Dashboard retourne le tableau de bord entreprise.
func (h *CompanyHandler) Dashboard(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	data, err := h.svc.Dashboard(userID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, data)
}

// ListSubjects liste les sujets proposés par l'entreprise.
func (h *CompanyHandler) ListSubjects(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	subjects, err := h.svc.ListSubjects(userID)
	if err != nil {
		return response.Error(c, err)
	}
	if subjects == nil {
		subjects = []*entity.PfeSubject{}
	}
	return response.OK(c, subjects)
}

// CreateSubject crée un nouveau sujet proposé.
func (h *CompanyHandler) CreateSubject(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	var req entity.PfeSubject
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if req.Title == "" || req.Description == "" {
		return response.ValidationError(c, "Titre et description requis")
	}
	if err := h.svc.CreateSubject(userID, &req); err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, req)
}

// GetSubject retourne un sujet de l'entreprise.
func (h *CompanyHandler) GetSubject(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	id := c.Params("id")
	subject, err := h.svc.GetSubject(userID, id)
	if err != nil {
		return response.Error(c, err)
	}
	if subject == nil {
		return response.NotFound(c, "Sujet introuvable")
	}
	return response.OK(c, subject)
}

// UpdateSubject met à jour un sujet.
func (h *CompanyHandler) UpdateSubject(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	id := c.Params("id")
	var req entity.PfeSubject
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	req.ID = id
	if err := h.svc.UpdateSubject(userID, &req); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Sujet mis à jour"})
}

// ListCandidats liste les candidats pour un sujet.
func (h *CompanyHandler) ListCandidats(c fiber.Ctx) error {
	id := c.Params("id")
	candidats, err := h.svc.ListCandidats(id)
	if err != nil {
		return response.Error(c, err)
	}
	if candidats == nil {
		candidats = []*entity.Wish{}
	}
	return response.OK(c, candidats)
}

// AcceptCandidat accepte ou refuse un étudiant pour un sujet.
func (h *CompanyHandler) AcceptCandidat(c fiber.Ctx) error {
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

// ListSupervisedPFEs liste les PFE encadrés.
func (h *CompanyHandler) ListSupervisedPFEs(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	assignments, err := h.svc.ListSupervisedPFEs(userID)
	if err != nil {
		return response.Error(c, err)
	}
	if assignments == nil {
		assignments = []*entity.PfeAssignment{}
	}
	return response.OK(c, assignments)
}

// GetSupervisedPFE retourne un PFE encadré.
func (h *CompanyHandler) GetSupervisedPFE(c fiber.Ctx) error {
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

// AddMeeting ajoute un meeting de suivi.
func (h *CompanyHandler) AddMeeting(c fiber.Ctx) error {
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
func (h *CompanyHandler) SubmitEvaluation(c fiber.Ctx) error {
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

// ListReports liste les signalements.
func (h *CompanyHandler) ListReports(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	reports, err := h.svc.ListReports(userID)
	if err != nil {
		return response.Error(c, err)
	}
	if reports == nil {
		reports = []*entity.CompanyReport{}
	}
	return response.OK(c, reports)
}

// CreateReport crée un signalement.
func (h *CompanyHandler) CreateReport(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	var req entity.CompanyReport
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if req.Description == "" || req.CorrectionType == "" {
		return response.ValidationError(c, "Description et correction_type requis")
	}
	if err := h.svc.CreateReport(userID, &req); err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, req)
}

// ListNotifications liste les notifications.
func (h *CompanyHandler) ListNotifications(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	notifications, err := h.svc.ListNotifications(userID)
	if err != nil {
		return response.Error(c, err)
	}
	if notifications == nil {
		notifications = []*entity.Notification{}
	}
	return response.OK(c, notifications)
}
