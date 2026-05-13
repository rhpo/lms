package handler

import (
	"pfe-backend/internal/entity"
	"pfe-backend/internal/service"
	"pfe-backend/internal/shared/middleware"
	"pfe-backend/internal/shared/response"

	"github.com/gofiber/fiber/v3"
)

// StudentHandler gère les endpoints étudiant.
type StudentHandler struct {
	svc *service.StudentService
}

// NewStudentHandler crée un nouveau StudentHandler.
func NewStudentHandler(svc *service.StudentService) *StudentHandler {
	return &StudentHandler{svc: svc}
}

// Dashboard retourne le tableau de bord étudiant.
func (h *StudentHandler) Dashboard(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	data, err := h.svc.Dashboard(userID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, data)
}

// ListCatalogue liste tous les sujets disponibles.
func (h *StudentHandler) ListCatalogue(c fiber.Ctx) error {
	subjects, err := h.svc.ListCatalogue()
	if err != nil {
		return response.Error(c, err)
	}
	if subjects == nil {
		subjects = []*entity.PfeSubject{}
	}
	return response.OK(c, subjects)
}

// GetCatalogueSubject retourne un sujet du catalogue.
func (h *StudentHandler) GetCatalogueSubject(c fiber.Ctx) error {
	id := c.Params("id")
	subject, err := h.svc.GetCatalogueSubject(id)
	if err != nil {
		return response.Error(c, err)
	}
	if subject == nil {
		return response.NotFound(c, "Sujet introuvable")
	}
	return response.OK(c, subject)
}

// ListWishes liste les voeux de l'étudiant.
func (h *StudentHandler) ListWishes(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	wishes, err := h.svc.ListWishes(userID)
	if err != nil {
		return response.Error(c, err)
	}
	if wishes == nil {
		wishes = []*entity.Wish{}
	}
	return response.OK(c, wishes)
}

// CreateWish crée un voeu pour l'étudiant.
func (h *StudentHandler) CreateWish(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	var req struct {
		SubjectID string `json:"subject_id"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if req.SubjectID == "" {
		return response.ValidationError(c, "L'ID du sujet est requis")
	}
	if err := h.svc.CreateWish(userID, req.SubjectID); err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, map[string]string{"message": "Voeu créé"})
}

// DeleteWish supprime un voeu.
func (h *StudentHandler) DeleteWish(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	id := c.Params("id")
	if err := h.svc.DeleteWish(userID, id); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Voeu supprimé"})
}

// GetMyPFE retourne le PFE de l'étudiant.
func (h *StudentHandler) GetMyPFE(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	assignment, err := h.svc.GetMyPFE(userID)
	if err != nil {
		return response.Error(c, err)
	}
	if assignment == nil {
		return response.NotFound(c, "aucun PFE assigné")
	}
	return response.OK(c, assignment)
}

// ListMyMeetings liste les meetings de suivi.
func (h *StudentHandler) ListMyMeetings(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	assignment, err := h.svc.GetMyPFE(userID)
	if err != nil {
		return response.Error(c, err)
	}
	if assignment == nil {
		return response.NotFound(c, "aucun PFE assigné")
	}
	meetings, err := h.svc.ListMyMeetings(assignment.ID)
	if err != nil {
		return response.Error(c, err)
	}
	if meetings == nil {
		meetings = []*entity.PfeProgressReport{}
	}
	return response.OK(c, meetings)
}

// AddMyMeeting ajoute un meeting de suivi.
func (h *StudentHandler) AddMyMeeting(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	var req entity.PfeProgressReport
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if req.MeetingType == "" || req.Duration == 0 {
		return response.ValidationError(c, "meeting_type et duration sont requis")
	}
	if err := h.svc.AddMyMeeting(userID, &req); err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, req)
}

// SubmitMemoire soumet le mémoire PDF.
func (h *StudentHandler) SubmitMemoire(c fiber.Ctx) error {
	var req struct {
		MemoireURL string `json:"memoire_url"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if req.MemoireURL == "" {
		return response.ValidationError(c, "URL du mémoire requis")
	}

	userID := middleware.GetProfileID(c)
	assignment, err := h.svc.GetMyPFE(userID)
	if err != nil {
		return response.Error(c, err)
	}
	if assignment == nil {
		return response.NotFound(c, "aucun PFE assigné")
	}
	if err := h.svc.SubmitMemoire(assignment.ID, req.MemoireURL); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, map[string]string{"message": "Mémoire soumis"})
}

// GetSoutenance retourne les infos de soutenance.
func (h *StudentHandler) GetSoutenance(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	data, err := h.svc.GetSoutenance(userID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, data)
}

// ListNotifications liste les notifications de l'étudiant.
func (h *StudentHandler) ListNotifications(c fiber.Ctx) error {
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
