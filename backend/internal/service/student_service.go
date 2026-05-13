package service

import (
	"pfe-backend/internal/entity"
	"pfe-backend/internal/repository"
	"pfe-backend/internal/shared/apperror"
)

// StudentService gère la logique métier des étudiants.
type StudentService struct {
	profileRepo       *repository.ProfileRepository
	studentRepo       *repository.StudentRepository
	pfeSubjectRepo    *repository.PfeSubjectRepository
	wishRepo          *repository.WishRepository
	pfeAssignmentRepo *repository.PfeAssignmentRepository
	progressRepo      *repository.ProgressReportRepository
	defenseRepo       *repository.DefenseRepository
	defenseJuryRepo   *repository.DefenseJuryRepository
	notificationRepo  *repository.NotificationRepository
	academicYearRepo  *repository.AcademicYearRepository
}

// NewStudentService crée un nouveau StudentService.
func NewStudentService(
	profileRepo *repository.ProfileRepository,
	studentRepo *repository.StudentRepository,
	pfeSubjectRepo *repository.PfeSubjectRepository,
	wishRepo *repository.WishRepository,
	pfeAssignmentRepo *repository.PfeAssignmentRepository,
	progressRepo *repository.ProgressReportRepository,
	defenseRepo *repository.DefenseRepository,
	defenseJuryRepo *repository.DefenseJuryRepository,
	notificationRepo *repository.NotificationRepository,
	academicYearRepo *repository.AcademicYearRepository,
) *StudentService {
	return &StudentService{
		profileRepo:       profileRepo,
		studentRepo:       studentRepo,
		pfeSubjectRepo:    pfeSubjectRepo,
		wishRepo:          wishRepo,
		pfeAssignmentRepo: pfeAssignmentRepo,
		progressRepo:      progressRepo,
		defenseRepo:       defenseRepo,
		defenseJuryRepo:   defenseJuryRepo,
		notificationRepo:  notificationRepo,
		academicYearRepo:  academicYearRepo,
	}
}

// getActiveAcademicYear récupère l'année académique active.
func (s *StudentService) getActiveAcademicYear() (string, error) {
	year, err := s.academicYearRepo.FindActive()
	if err != nil {
		return "", err
	}
	if year == nil {
		return "", apperror.Internal("Aucune année académique active")
	}
	return year.ID, nil
}

// Dashboard retourne les statistiques du tableau de bord étudiant.
func (s *StudentService) Dashboard(userID string) (map[string]any, error) {
	academicYearID, err := s.getActiveAcademicYear()
	if err != nil {
		return nil, err
	}

	assignment, _ := s.pfeAssignmentRepo.FindByStudent(userID, academicYearID)
	wishes, _ := s.wishRepo.FindByStudent(userID, academicYearID)

	result := map[string]any{
		"wishes_count": len(wishes),
	}
	if assignment != nil {
		result["has_pfe"] = true
		result["pfe_status"] = assignment.Status
	} else {
		result["has_pfe"] = false
	}
	return result, nil
}

// ListCatalogue liste tous les sujets disponibles (validés) pour l'étudiant.
func (s *StudentService) ListCatalogue() ([]*entity.PfeSubject, error) {
	return s.pfeSubjectRepo.FindByStatus("valide")
}

// GetCatalogueSubject retourne un sujet du catalogue.
func (s *StudentService) GetCatalogueSubject(id string) (*entity.PfeSubject, error) {
	return s.pfeSubjectRepo.FindByID(id)
}

// ListWishes liste les voeux de l'étudiant.
func (s *StudentService) ListWishes(userID string) ([]*entity.Wish, error) {
	academicYearID, err := s.getActiveAcademicYear()
	if err != nil {
		return nil, err
	}
	return s.wishRepo.FindByStudent(userID, academicYearID)
}

// CreateWish crée un voeu pour l'étudiant.
func (s *StudentService) CreateWish(userID, subjectID string) error {
	// Vérifier que le sujet existe
	subject, err := s.pfeSubjectRepo.FindByID(subjectID)
	if err != nil {
		return err
	}
	if subject == nil {
		return apperror.NotFound("Sujet introuvable")
	}
	if subject.Status != "valide" {
		return apperror.BadRequest("Ce sujet n'est pas disponible")
	}

	// Vérifier que l'étudiant n'a pas déjà un voeu pour ce sujet
	academicYearID, err := s.getActiveAcademicYear()
	if err != nil {
		return err
	}

	wishes, err := s.wishRepo.FindByStudent(userID, academicYearID)
	if err != nil {
		return err
	}
	for _, w := range wishes {
		if w.SubjectID == subjectID {
			return apperror.Conflict("Vous avez déjà un voeu pour ce sujet")
		}
	}

	wish := &entity.Wish{
		ID:             generateID(),
		StudentID:      userID,
		SubjectID:      subjectID,
		AcademicYearID: academicYearID,
		Status:         "en_attente",
	}
	return s.wishRepo.Insert(wish)
}

// DeleteWish supprime un voeu.
func (s *StudentService) DeleteWish(userID, wishID string) error {
	wish, err := s.wishRepo.FindByID(wishID)
	if err != nil {
		return err
	}
	if wish == nil {
		return apperror.NotFound("Voeu introuvable")
	}
	if wish.StudentID != userID {
		return apperror.Forbidden("Accès non autorisé à ce voeu")
	}
	return s.wishRepo.Delete(wishID)
}

// GetMyPFE retourne le PFE de l'étudiant.
func (s *StudentService) GetMyPFE(userID string) (*entity.PfeAssignment, error) {
	academicYearID, err := s.getActiveAcademicYear()
	if err != nil {
		return nil, err
	}
	return s.pfeAssignmentRepo.FindByStudent(userID, academicYearID)
}

// ListMyMeetings liste les meetings de suivi du PFE de l'étudiant.
func (s *StudentService) ListMyMeetings(assignmentID string) ([]*entity.PfeProgressReport, error) {
	return s.progressRepo.FindByAssignment(assignmentID)
}

// AddMyMeeting ajoute un meeting de suivi pour le PFE de l'étudiant.
func (s *StudentService) AddMyMeeting(userID string, report *entity.PfeProgressReport) error {
	academicYearID, err := s.getActiveAcademicYear()
	if err != nil {
		return err
	}
	assignment, err := s.pfeAssignmentRepo.FindByStudent(userID, academicYearID)
	if err != nil {
		return err
	}
	if assignment == nil {
		return apperror.NotFound("aucun PFE assigné")
	}
	report.ID = generateID()
	report.AssignmentID = assignment.ID
	if report.Status == "" {
		report.Status = "en_cours"
	}
	return s.progressRepo.Insert(report)
}

// SubmitMemoire soumet le mémoire PDF.
func (s *StudentService) SubmitMemoire(assignmentID, memoireURL string) error {
	return s.pfeAssignmentRepo.UpdateMemoire(assignmentID, memoireURL)
}

// GetSoutenance retourne les informations de soutenance de l'étudiant.
func (s *StudentService) GetSoutenance(userID string) (map[string]any, error) {
	academicYearID, err := s.getActiveAcademicYear()
	if err != nil {
		return nil, err
	}
	assignment, err := s.pfeAssignmentRepo.FindByStudent(userID, academicYearID)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, apperror.NotFound("Aucun PFE assigné")
	}

	defense, err := s.defenseRepo.FindByAssignment(assignment.ID)
	if err != nil {
		return nil, err
	}
	if defense == nil {
		return map[string]any{"has_soutenance": false}, nil
	}

	// Récupérer les infos du jury
	var jury *entity.DefenseJury
	if defense.JuryID != "" {
		jury, _ = s.defenseJuryRepo.FindByID(defense.JuryID)
	}

	return map[string]any{
		"has_soutenance": true,
		"defense":        defense,
		"jury":           jury,
	}, nil
}

// ListNotifications liste les notifications de l'étudiant.
func (s *StudentService) ListNotifications(userID string) ([]*entity.Notification, error) {
	return s.notificationRepo.FindByRecipient(userID)
}
