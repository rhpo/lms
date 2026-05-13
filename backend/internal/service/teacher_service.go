package service

import (
	"database/sql"
	"time"

	"pfe-backend/internal/entity"
	"pfe-backend/internal/repository"
	"pfe-backend/internal/shared/apperror"
)

// TeacherService gère la logique métier des enseignants.
type TeacherService struct {
	profileRepo       *repository.ProfileRepository
	teacherRepo       *repository.TeacherRepository
	pfeSubjectRepo    *repository.PfeSubjectRepository
	wishRepo          *repository.WishRepository
	pfeAssignmentRepo *repository.PfeAssignmentRepository
	progressRepo      *repository.ProgressReportRepository
	defenseJuryRepo   *repository.DefenseJuryRepository
	defenseRepo       *repository.DefenseRepository
	supEvalRepo       *repository.SupervisorEvaluationRepository
	notificationRepo  *repository.NotificationRepository
	academicYearRepo  *repository.AcademicYearRepository
}

// NewTeacherService crée un nouveau TeacherService.
func NewTeacherService(
	profileRepo *repository.ProfileRepository,
	teacherRepo *repository.TeacherRepository,
	pfeSubjectRepo *repository.PfeSubjectRepository,
	wishRepo *repository.WishRepository,
	pfeAssignmentRepo *repository.PfeAssignmentRepository,
	progressRepo *repository.ProgressReportRepository,
	defenseJuryRepo *repository.DefenseJuryRepository,
	defenseRepo *repository.DefenseRepository,
	supEvalRepo *repository.SupervisorEvaluationRepository,
	notificationRepo *repository.NotificationRepository,
	academicYearRepo *repository.AcademicYearRepository,
) *TeacherService {
	return &TeacherService{
		profileRepo:       profileRepo,
		teacherRepo:       teacherRepo,
		pfeSubjectRepo:    pfeSubjectRepo,
		wishRepo:          wishRepo,
		pfeAssignmentRepo: pfeAssignmentRepo,
		progressRepo:      progressRepo,
		defenseJuryRepo:   defenseJuryRepo,
		defenseRepo:       defenseRepo,
		supEvalRepo:       supEvalRepo,
		notificationRepo:  notificationRepo,
		academicYearRepo:  academicYearRepo,
	}
}

// Dashboard retourne les statistiques du tableau de bord enseignant.
func (s *TeacherService) Dashboard(userID string) (map[string]any, error) {
	subjects, _ := s.pfeSubjectRepo.FindByProposer(userID)
	supervised, _ := s.pfeAssignmentRepo.FindBySupervisor(userID)

	return map[string]any{
		"proposed_subjects": len(subjects),
		"supervised_pfes":   len(supervised),
	}, nil
}

// ListProposedSubjects liste les sujets proposés par l'enseignant.
func (s *TeacherService) ListProposedSubjects(userID string) ([]*entity.PfeSubject, error) {
	return s.pfeSubjectRepo.FindByProposer(userID)
}

// CreateProposedSubject crée un sujet proposé par l'enseignant.
func (s *TeacherService) CreateProposedSubject(subject *entity.PfeSubject, domainIDs []string) error {
	if err := s.pfeSubjectRepo.Insert(subject); err != nil {
		return err
	}
	for _, domainID := range domainIDs {
		if err := s.pfeSubjectRepo.AddDomain(subject.ID, domainID); err != nil {
			return err
		}
	}
	return nil
}

// GetProposedSubject retourne un sujet proposé par l'enseignant.
func (s *TeacherService) GetProposedSubject(userID, id string) (*entity.PfeSubject, error) {
	subject, err := s.pfeSubjectRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if subject == nil {
		return nil, nil
	}
	if subject.ProposerID != userID {
		return nil, apperror.Forbidden("Vous n'êtes pas l'auteur de ce sujet")
	}
	return subject, nil
}

// UpdateProposedSubject met à jour un sujet proposé par l'enseignant.
func (s *TeacherService) UpdateProposedSubject(userID string, subject *entity.PfeSubject) error {
	existing, err := s.pfeSubjectRepo.FindByID(subject.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return apperror.NotFound("Sujet introuvable")
	}
	if existing.ProposerID != userID {
		return apperror.Forbidden("Vous n'êtes pas l'auteur de ce sujet")
	}
	if subject.Title != "" {
		existing.Title = subject.Title
	}
	if subject.Description != "" {
		existing.Description = subject.Description
	}
	if subject.GroupType != "" {
		existing.GroupType = subject.GroupType
	}
	return s.pfeSubjectRepo.Update(existing)
}

// ListCandidats liste les candidats pour un sujet.
func (s *TeacherService) ListCandidats(subjectID string) ([]*entity.Wish, error) {
	return s.wishRepo.FindBySubject(subjectID)
}

// AcceptCandidat accepte un étudiant pour un sujet.
func (s *TeacherService) AcceptCandidat(subjectID, studentID string) error {
	wishes, err := s.wishRepo.FindBySubject(subjectID)
	if err != nil {
		return err
	}
	for _, w := range wishes {
		if w.StudentID == studentID {
			w.Status = "accepte"
			return s.wishRepo.Update(w)
		}
	}
	ay, err := s.academicYearRepo.FindActive()
	if err != nil {
		return err
	}
	ayID := ""
	if ay != nil {
		ayID = ay.ID
	}
	wish := &entity.Wish{
		ID:             generateID(),
		StudentID:      studentID,
		SubjectID:      subjectID,
		AcademicYearID: ayID,
		Status:         "accepte",
	}
	return s.wishRepo.Insert(wish)
}

// RejectCandidat refuse un étudiant pour un sujet.
func (s *TeacherService) RejectCandidat(subjectID, studentID string) error {
	wishes, err := s.wishRepo.FindBySubject(subjectID)
	if err != nil {
		return err
	}
	for _, w := range wishes {
		if w.StudentID == studentID {
			w.Status = "refuse"
			return s.wishRepo.Update(w)
		}
	}
	return apperror.NotFound("Candidature introuvable")
}

// ListSupervisedPFEs liste les PFE encadrés par l'enseignant.
func (s *TeacherService) ListSupervisedPFEs(userID string) ([]*entity.PfeAssignment, error) {
	return s.pfeAssignmentRepo.FindBySupervisor(userID)
}

// GetSupervisedPFE retourne un PFE encadré.
func (s *TeacherService) GetSupervisedPFE(id string) (*entity.PfeAssignment, error) {
	return s.pfeAssignmentRepo.FindByID(id)
}

// AddMeeting ajoute un meeting de suivi à un PFE.
func (s *TeacherService) AddMeeting(report *entity.PfeProgressReport) error {
	return s.progressRepo.Insert(report)
}

// SubmitEvaluation soumet l'évaluation de l'encadrant.
func (s *TeacherService) SubmitEvaluation(assignmentID, evaluatorID string, criterion5 float64) error {
	if criterion5 < 0 || criterion5 > 4 {
		return apperror.BadRequest("Le critère 5 doit être entre 0 et 4")
	}
	// Vérifier si une évaluation existe déjà
	existing, err := s.supEvalRepo.FindByAssignment(assignmentID)
	if err != nil {
		return err
	}
	if existing != nil {
		existing.EvaluatorID = evaluatorID
		existing.Criterion5 = sql.NullFloat64{Float64: criterion5, Valid: true}
		return s.supEvalRepo.Update(existing)
	}
	eval := &entity.SupervisorEvaluation{
		ID:              generateID(),
		PfeAssignmentID: assignmentID,
		EvaluatorID:     evaluatorID,
		Criterion5:      sql.NullFloat64{Float64: criterion5, Valid: true},
	}
	return s.supEvalRepo.Insert(eval)
}

// UpdateAvailability met à jour la disponibilité de l'enseignant.
func (s *TeacherService) UpdateAvailability(userID, status string, unavailableUntilStr string) error {
	validStatuses := map[string]bool{
		"disponible":             true,
		"indisponible":           true,
		"indisponible_jusqu_au":  true,
	}
	if !validStatuses[status] {
		return apperror.BadRequest("Statut invalide: utilisez disponible, indisponible ou indisponible_jusqu_au")
	}
	var unavailableUntil *sql.NullTime
	if unavailableUntilStr != "" {
		t, err := time.Parse("2006-01-02", unavailableUntilStr)
		if err != nil {
			return apperror.BadRequest("Format de date invalide, utilisez YYYY-MM-DD")
		}
		unavailableUntil = &sql.NullTime{Time: t, Valid: true}
	}
	return s.teacherRepo.UpdateAvailability(userID, status, unavailableUntil)
}

// ListSubjectsToValidate liste les sujets à valider par l'enseignant.
func (s *TeacherService) ListSubjectsToValidate(userID string) ([]*entity.PfeSubject, error) {
	return s.pfeSubjectRepo.FindPendingValidation(userID)
}

// GetSubjectToValidate retourne un sujet à valider.
func (s *TeacherService) GetSubjectToValidate(userID, id string) (*entity.PfeSubject, error) {
	subject, err := s.pfeSubjectRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if subject == nil {
		return nil, nil
	}
	if !isTeacherValidator(userID, subject) {
		return nil, apperror.Forbidden("Vous n'êtes pas validateur de ce sujet")
	}
	return subject, nil
}

// ValidateSubject valide ou refuse un sujet.
func (s *TeacherService) ValidateSubject(userID, id, decision, comment string) error {
	validDecisions := map[string]bool{"valide": true, "accepte_sous_reserve": true, "refuse": true}
	if !validDecisions[decision] {
		return apperror.BadRequest("Décision invalide: utilisez valide, accepte_sous_reserve ou refuse")
	}
	subject, err := s.pfeSubjectRepo.FindByID(id)
	if err != nil {
		return err
	}
	if subject == nil {
		return apperror.NotFound("Sujet introuvable")
	}
	if !isTeacherValidator(userID, subject) {
		return apperror.Forbidden("Vous n'êtes pas validateur de ce sujet")
	}

	// Déterminer quel validateur (validator1 ou validator2)
	var validatorField string
	if subject.Validator1ID.Valid && subject.Validator1ID.String == userID {
		validatorField = "validator1"
	} else if subject.Validator2ID.Valid && subject.Validator2ID.String == userID {
		validatorField = "validator2"
	} else {
		return apperror.Forbidden("Vous n'êtes pas validateur de ce sujet")
	}

	// Enregistrer la décision
	if err := s.pfeSubjectRepo.UpdateValidation(id, validatorField, decision, comment); err != nil {
		return err
	}

	// Recharger pour calculer le nouveau statut
	subject, err = s.pfeSubjectRepo.FindByID(id)
	if err != nil {
		return err
	}
	newStatus := computeSubjectStatus(subject, decision)
	return s.pfeSubjectRepo.UpdateStatus(id, newStatus)
}

// ListJuryDuties liste les obligations de jury de l'enseignant.
func (s *TeacherService) ListJuryDuties(userID string) ([]*entity.Defense, error) {
	return s.defenseRepo.FindByJuryMember(userID)
}

// GetJuryDuty retourne une obligation de jury spécifique.
func (s *TeacherService) GetJuryDuty(id string) (*entity.Defense, error) {
	return s.defenseRepo.FindByID(id)
}

// ListNotifications liste les notifications de l'enseignant.
func (s *TeacherService) ListNotifications(userID string) ([]*entity.Notification, error) {
	return s.notificationRepo.FindByRecipient(userID)
}

// isTeacherValidator vérifie si un enseignant est validateur d'un sujet donné.
func isTeacherValidator(userID string, subject *entity.PfeSubject) bool {
	return (subject.Validator1ID.Valid && subject.Validator1ID.String == userID) ||
		(subject.Validator2ID.Valid && subject.Validator2ID.String == userID)
}

// setValidatorDecision définit la décision d'un validateur.
func setValidatorDecision(subject *entity.PfeSubject, userID, decision, comment string) {
	if subject.Validator1ID.Valid && subject.Validator1ID.String == userID {
		subject.Validator1Decision = sql.NullString{String: decision, Valid: true}
		subject.Validator1Comment = sql.NullString{String: comment, Valid: true}
	} else if subject.Validator2ID.Valid && subject.Validator2ID.String == userID {
		subject.Validator2Decision = sql.NullString{String: decision, Valid: true}
		subject.Validator2Comment = sql.NullString{String: comment, Valid: true}
	}
}

// computeSubjectStatus calcule le nouveau statut après validation.
func computeSubjectStatus(subject *entity.PfeSubject, decision string) string {
	bothValid := subject.Validator1Decision.Valid && subject.Validator2Decision.Valid &&
		subject.Validator1Decision.String == "valide" && subject.Validator2Decision.String == "valide"
	bothRefused := subject.Validator1Decision.Valid && subject.Validator1Decision.String == "refuse" &&
		subject.Validator2Decision.Valid && subject.Validator2Decision.String == "refuse"

	if bothValid {
		return "valide"
	}
	if bothRefused {
		return "refuse"
	}
	if decision == "refuse" {
		return "refuse"
	}
	if decision == "accepte_sous_reserve" {
		return "accepte_sous_reserve"
	}
	return subject.Status
}
