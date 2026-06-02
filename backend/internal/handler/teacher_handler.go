package handler

import (
	"database/sql"
	"fmt"
	"time"

	"pfe-backend/internal/entity"
	"pfe-backend/internal/service"
	"pfe-backend/internal/shared/middleware"
	"pfe-backend/internal/shared/notify"
	"pfe-backend/internal/shared/response"

	"github.com/gofiber/fiber/v3"
)

// TeacherHandler gère les endpoints enseignant.
type TeacherHandler struct {
	svc      *service.TeacherService
	notifier *notify.Notifier
}

// NewTeacherHandler crée un nouveau TeacherHandler.
func NewTeacherHandler(svc *service.TeacherService, notifier *notify.Notifier) *TeacherHandler {
	return &TeacherHandler{svc: svc, notifier: notifier}
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
		Title       string  `json:"title"`
		Description string  `json:"description"`
		GroupType   string  `json:"group_type"`
		DomainIDs   []int64 `json:"domain_ids"`
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

	// Notify admins: new subject awaiting validation
	go h.notifier.NotifyAdmins(notify.TypeValidationRequise,
		fmt.Sprintf("Un nouveau sujet « %s » a été proposé par un enseignant et attend votre validation.", req.Title))

	return response.Created(c, subject)
}

// GetProposedSubject retourne un sujet proposé.
func (h *TeacherHandler) GetProposedSubject(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
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
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
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

// ResubmitSubject permet à l'auteur de modifier et resoumettre un sujet refusé ou accepté sous réserve.
func (h *TeacherHandler) ResubmitSubject(c fiber.Ctx) error {
	userID := middleware.GetProfileID(c)
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
	var req struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		GroupType   string  `json:"group_type"`
		DomainIDs   []int64 `json:"domain_ids"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.ResubmitSubject(userID, id, req.Title, req.Description, req.GroupType, req.DomainIDs); err != nil {
		return response.Error(c, err)
	}

	// Notify admins that a subject has been resubmitted
	title := req.Title
	if title == "" {
		title = h.svc.GetSubjectTitle(id)
	}
	go h.notifier.NotifyAdmins(notify.TypeValidationRequise,
		fmt.Sprintf("Le sujet « %s » a été modifié et resoumis pour validation par un enseignant.", title))

	return response.OK(c, map[string]string{"message": "Sujet resoumis pour validation"})
}

// ListCandidats liste les candidats pour un sujet.
func (h *TeacherHandler) ListCandidats(c fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
	candidats, err := h.svc.ListCandidats(id)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, candidats)
}

// AcceptCandidat accepte les étudiants sélectionnés pour un sujet (binôme/trinôme supporté).
func (h *TeacherHandler) AcceptCandidat(c fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
	var req struct {
		StudentIDs []int64 `json:"student_ids"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if len(req.StudentIDs) == 0 {
		return response.ValidationError(c, "Aucun étudiant sélectionné")
	}
	assignment, err := h.svc.AcceptCandidats(id, req.StudentIDs)
	if err != nil {
		return response.Error(c, err)
	}
	// Notify each student (needs profile ID, not entity ID) with subject title
	subjectTitle := h.svc.GetSubjectTitle(id)
	for _, studentID := range req.StudentIDs {
		go func(sID int64) {
			if profileID, err := h.svc.GetStudentProfileID(sID); err == nil {
				h.notifier.Send(profileID, notify.TypeAffectation,
					fmt.Sprintf("Votre candidature pour le sujet « %s » a été acceptée. Votre code PFE est %s.", subjectTitle, assignment.PfeCode))
			}
		}(studentID)
	}
	return response.OK(c, map[string]string{"message": "Candidats acceptés", "pfe_code": assignment.PfeCode})
}

// RejectCandidat refuse un étudiant pour un sujet.
func (h *TeacherHandler) RejectCandidat(c fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
	var req struct {
		StudentID int64 `json:"student_id"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if req.StudentID == 0 {
		return response.ValidationError(c, "L'ID de l'étudiant est requis")
	}
	if err := h.svc.RejectCandidat(id, req.StudentID); err != nil {
		return response.Error(c, err)
	}
	go func() {
		if profileID, err := h.svc.GetStudentProfileID(req.StudentID); err == nil {
			subjectTitle := h.svc.GetSubjectTitle(id)
			h.notifier.Send(profileID, notify.TypeAffectation,
				fmt.Sprintf("Votre candidature pour le sujet « %s » a été refusée par l'enseignant.", subjectTitle))
		}
	}()
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
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
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
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
	var req struct {
		Decision string `json:"decision"`
		Comment  string `json:"comment"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	// Load subject title before validation for the notification
	subjectTitle := fmt.Sprintf("sujet #%d", id) // fallback
	if sub, err := h.svc.GetSubjectToValidate(userID, id); err == nil && sub != nil {
		subjectTitle = sub.Title
	}

	if err := h.svc.ValidateSubject(userID, id, req.Decision, req.Comment); err != nil {
		return response.Error(c, err)
	}

	// Notify admins — human-readable decision label, subject title (never raw IDs)
	decisionLabels := map[string]string{
		"valide":               "validé",
		"accepte_sous_reserve": "accepté sous réserve",
		"refuse":               "refusé",
	}
	label := decisionLabels[req.Decision]
	if label == "" {
		label = req.Decision
	}
	msg := fmt.Sprintf("Un enseignant a %s le sujet « %s ».", label, subjectTitle)
	go h.notifier.NotifyAdmins(notify.TypeValidationRequise, msg)

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
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
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
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
	var req struct {
		MeetingDate string `json:"meeting_date"`
		Duration    int    `json:"duration"`
		MeetingType string `json:"meeting_type"`
		Topics      string `json:"topics"`
		Status      string `json:"status"`
		Observation string `json:"observation"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if req.MeetingDate == "" {
		return response.ValidationError(c, "La date est requise")
	}
	if req.MeetingType == "" {
		return response.ValidationError(c, "Le type de réunion est requis")
	}
	if req.Duration == 0 {
		return response.ValidationError(c, "La durée est requise")
	}

	// Parse date — support date-only and datetime-local formats
	var meetingDate time.Time
	var parseErr error
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04", "2006-01-02"} {
		meetingDate, parseErr = time.Parse(layout, req.MeetingDate)
		if parseErr == nil {
			break
		}
	}
	if parseErr != nil {
		return response.ValidationError(c, "Format de date invalide (attendu YYYY-MM-DD)")
	}

	status := req.Status
	if status == "" {
		status = "a_faire"
	}

	report := &entity.PfeProgressReport{
		AssignmentID: id,
		MeetingDate:  meetingDate,
		Duration:     req.Duration,
		MeetingType:  req.MeetingType,
		Topics:       req.Topics,
		Status:       status,
	}
	if req.Observation != "" {
		report.Observation = entity.NullString{NullString: sql.NullString{String: req.Observation, Valid: true}}
	}

	if err := h.svc.AddMeeting(report); err != nil {
		return response.Error(c, err)
	}

	// Notify student about the new meeting (with subject title and date)
	go func() {
		assignment, err := h.svc.GetSupervisedPFE(id)
		if err != nil || assignment == nil {
			return
		}
		subjectTitle := "votre PFE"
		if assignment.Subject != nil && assignment.Subject.Title != "" {
			subjectTitle = fmt.Sprintf("« %s »", assignment.Subject.Title)
		}
		dateStr := meetingDate.Format("02/01/2006")
		if assignment.Student != nil {
			h.notifier.Send(assignment.Student.ProfileID, notify.TypeDisponibilite,
				fmt.Sprintf("Une réunion de suivi pour le sujet %s a été planifiée le %s par votre encadrant.", subjectTitle, dateStr))
		}
	}()

	return response.Created(c, report)
}

// ListMeetings liste les réunions de suivi d'un PFE encadré.
func (h *TeacherHandler) ListMeetings(c fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
	meetings, err := h.svc.ListMeetings(id)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, meetings)
}

// GetEvaluation retourne l'évaluation d'un PFE encadré.
func (h *TeacherHandler) GetEvaluation(c fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
	eval, err := h.svc.GetEvaluation(id)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, eval)
}

// SubmitEvaluation soumet l'évaluation de l'encadrant.
func (h *TeacherHandler) SubmitEvaluation(c fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
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

	// Notify student and admins that the supervisor evaluation has been submitted
	go func() {
		assignment, err := h.svc.GetSupervisedPFE(id)
		if err != nil || assignment == nil {
			return
		}
		subjectTitle := "votre PFE"
		if assignment.Subject != nil && assignment.Subject.Title != "" {
			subjectTitle = fmt.Sprintf("« %s »", assignment.Subject.Title)
		}
		if assignment.Student != nil {
			h.notifier.Send(assignment.Student.ProfileID, notify.TypeSujet,
				fmt.Sprintf("L'évaluation de votre encadrant pour le sujet %s a été soumise.", subjectTitle))
		}
		h.notifier.NotifyAdmins(notify.TypeValidationRequise,
			fmt.Sprintf("L'évaluation de l'encadrant pour le sujet %s a été soumise.", subjectTitle))
	}()

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
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
	duty, err := h.svc.GetJuryDuty(id)
	if err != nil {
		return response.Error(c, err)
	}
	if duty == nil {
		return response.NotFound(c, "Obligation jury introuvable")
	}
	return response.OK(c, duty)
}

// GetGradeContext returns the grading context for a defense.
func (h *TeacherHandler) GetGradeContext(c fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
	callerID := middleware.GetProfileID(c)
	ctx, err := h.svc.GetGradeContext(id, callerID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, ctx)
}

// SubmitJuryGrade soumet la note du membre (examinateur) pour une soutenance.
func (h *TeacherHandler) SubmitJuryGrade(c fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
	callerID := middleware.GetProfileID(c)
	var req struct {
		Criterion1      float64 `json:"criterion1"`
		Criterion2      float64 `json:"criterion2"`
		Criterion3      float64 `json:"criterion3"`
		Criterion4      float64 `json:"criterion4"`
		ArchiveDecision string  `json:"archive_decision"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.SubmitJuryGrade(id, callerID, req.Criterion1, req.Criterion2, req.Criterion3, req.Criterion4, req.ArchiveDecision); err != nil {
		return response.Error(c, err)
	}

	// Notify admins that the member has submitted their grade
	go func() {
		duty, err := h.svc.GetJuryDuty(id)
		if err != nil || duty == nil {
			return
		}
		subjectTitle := "un PFE"
		if duty.Assignment != nil && duty.Assignment.Subject != nil {
			subjectTitle = fmt.Sprintf("« %s »", duty.Assignment.Subject.Title)
		}
		h.notifier.NotifyAdmins(notify.TypeJury,
			fmt.Sprintf("L'examinateur a soumis son évaluation pour le sujet %s.", subjectTitle))
	}()

	return response.OK(c, map[string]string{"message": "Évaluation soumise"})
}

// SubmitFinalGrade soumet la note finale par le président du jury.
func (h *TeacherHandler) SubmitFinalGrade(c fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return response.ValidationError(c, "ID invalide")
	}
	callerID := middleware.GetProfileID(c)
	var req struct {
		Choice          string  `json:"choice"`           // "member" or "new"
		Criterion1      float64 `json:"criterion1"`
		Criterion2      float64 `json:"criterion2"`
		Criterion3      float64 `json:"criterion3"`
		Criterion4      float64 `json:"criterion4"`
		ArchiveDecision string  `json:"archive_decision"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.ValidationError(c, "Données invalides")
	}
	if err := h.svc.SubmitFinalGrade(id, callerID, req.Choice, req.Criterion1, req.Criterion2, req.Criterion3, req.Criterion4, req.ArchiveDecision); err != nil {
		return response.Error(c, err)
	}

	// Notify admins that the president has finalized the grade
	go func() {
		duty, err := h.svc.GetJuryDuty(id)
		if err != nil || duty == nil {
			return
		}
		subjectTitle := "un PFE"
		if duty.Assignment != nil && duty.Assignment.Subject != nil {
			subjectTitle = fmt.Sprintf("« %s »", duty.Assignment.Subject.Title)
		}
		h.notifier.NotifyAdmins(notify.TypeJury,
			fmt.Sprintf("Le président du jury a finalisé la note pour le sujet %s.", subjectTitle))
	}()

	return response.OK(c, map[string]string{"message": "Note finale soumise"})
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
