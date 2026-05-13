package service

import (
	"database/sql"

	"pfe-backend/internal/entity"
	"pfe-backend/internal/repository"
	"pfe-backend/internal/shared/apperror"
)

// CompanyService gère la logique métier des entreprises.
type CompanyService struct {
	profileRepo       *repository.ProfileRepository
	companyRepo       *repository.CompanyRepository
	pfeSubjectRepo    *repository.PfeSubjectRepository
	wishRepo          *repository.WishRepository
	pfeAssignmentRepo *repository.PfeAssignmentRepository
	progressRepo      *repository.ProgressReportRepository
	supEvalRepo       *repository.SupervisorEvaluationRepository
	companyReportRepo *repository.CompanyReportRepository
	notificationRepo  *repository.NotificationRepository
	academicYearRepo  *repository.AcademicYearRepository
}

// NewCompanyService crée un nouveau CompanyService.
func NewCompanyService(
	profileRepo *repository.ProfileRepository,
	companyRepo *repository.CompanyRepository,
	pfeSubjectRepo *repository.PfeSubjectRepository,
	wishRepo *repository.WishRepository,
	pfeAssignmentRepo *repository.PfeAssignmentRepository,
	progressRepo *repository.ProgressReportRepository,
	supEvalRepo *repository.SupervisorEvaluationRepository,
	companyReportRepo *repository.CompanyReportRepository,
	notificationRepo *repository.NotificationRepository,
	academicYearRepo *repository.AcademicYearRepository,
) *CompanyService {
	return &CompanyService{
		profileRepo:       profileRepo,
		companyRepo:       companyRepo,
		pfeSubjectRepo:    pfeSubjectRepo,
		wishRepo:          wishRepo,
		pfeAssignmentRepo: pfeAssignmentRepo,
		progressRepo:      progressRepo,
		supEvalRepo:       supEvalRepo,
		companyReportRepo: companyReportRepo,
		notificationRepo:  notificationRepo,
		academicYearRepo:  academicYearRepo,
	}
}

// getCompanyByProfile récupère l'entité Company liée au profil, ou erreur si introuvable.
func (s *CompanyService) getCompanyByProfile(userID string) (*entity.Company, error) {
	company, err := s.companyRepo.FindByProfileID(userID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, apperror.NotFound("Profil entreprise introuvable")
	}
	return company, nil
}

// Dashboard retourne les statistiques du tableau de bord entreprise.
func (s *CompanyService) Dashboard(userID string) (map[string]any, error) {
	company, err := s.getCompanyByProfile(userID)
	if err != nil {
		return nil, err
	}

	subjects, _ := s.pfeSubjectRepo.FindByProposer(company.ID)
	supervised, _ := s.pfeAssignmentRepo.FindBySupervisor(company.ID)
	reports, _ := s.companyReportRepo.FindByCompany(company.ID)

	return map[string]any{
		"subjects_count":  len(subjects),
		"supervised_pfes": len(supervised),
		"reports_count":   len(reports),
	}, nil
}

// ListSubjects liste les sujets proposés par l'entreprise.
func (s *CompanyService) ListSubjects(userID string) ([]*entity.PfeSubject, error) {
	company, err := s.getCompanyByProfile(userID)
	if err != nil {
		return nil, err
	}
	return s.pfeSubjectRepo.FindByCompany(company.ID)
}

// CreateSubject crée un sujet proposé par l'entreprise.
func (s *CompanyService) CreateSubject(userID string, subject *entity.PfeSubject) error {
	company, err := s.getCompanyByProfile(userID)
	if err != nil {
		return err
	}

	ay, err := s.academicYearRepo.FindActive()
	if err != nil {
		return err
	}
	if ay == nil {
		return apperror.Internal("Aucune année académique active")
	}

	subject.ID = generateID()
	subject.CompanyID = sql.NullString{String: company.ID, Valid: true}
	subject.ProposerID = company.ID
	subject.ProposerRole = "company"
	subject.AcademicYearID = ay.ID
	if subject.Status == "" {
		subject.Status = "en_attente"
	}
	return s.pfeSubjectRepo.Insert(subject)
}

// isCompanySubject vérifie si un sujet appartient à l'entreprise.
func isCompanySubject(subject *entity.PfeSubject, companyID string) bool {
	if subject.CompanyID.Valid && subject.CompanyID.String == companyID {
		return true
	}
	return subject.ProposerID == companyID && subject.ProposerRole == "company"
}

// GetSubject retourne un sujet de l'entreprise.
func (s *CompanyService) GetSubject(userID, id string) (*entity.PfeSubject, error) {
	company, err := s.getCompanyByProfile(userID)
	if err != nil {
		return nil, err
	}

	subject, err := s.pfeSubjectRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if subject == nil {
		return nil, nil
	}
	if !isCompanySubject(subject, company.ID) {
		return nil, apperror.Forbidden("Accès non autorisé à ce sujet")
	}
	return subject, nil
}

// UpdateSubject met à jour un sujet de l'entreprise.
func (s *CompanyService) UpdateSubject(userID string, subject *entity.PfeSubject) error {
	company, err := s.getCompanyByProfile(userID)
	if err != nil {
		return err
	}

	existing, err := s.pfeSubjectRepo.FindByID(subject.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return apperror.NotFound("Sujet introuvable")
	}
	if !isCompanySubject(existing, company.ID) {
		return apperror.Forbidden("Accès non autorisé à ce sujet")
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

// ListCandidats liste les candidats pour un sujet de l'entreprise.
func (s *CompanyService) ListCandidats(subjectID string) ([]*entity.Wish, error) {
	return s.wishRepo.FindBySubject(subjectID)
}

// AcceptCandidat accepte un étudiant pour un sujet.
func (s *CompanyService) AcceptCandidat(subjectID, studentID string) error {
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
func (s *CompanyService) RejectCandidat(subjectID, studentID string) error {
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

// ListSupervisedPFEs liste les PFE encadrés par l'entreprise (sujets proposés par l'entreprise).
func (s *CompanyService) ListSupervisedPFEs(userID string) ([]*entity.PfeAssignment, error) {
	company, err := s.getCompanyByProfile(userID)
	if err != nil {
		return nil, err
	}
	return s.pfeAssignmentRepo.FindByCompanySubject(company.ID)
}

// GetSupervisedPFE retourne un PFE encadré.
func (s *CompanyService) GetSupervisedPFE(id string) (*entity.PfeAssignment, error) {
	return s.pfeAssignmentRepo.FindByID(id)
}

// AddMeeting ajoute un meeting de suivi à un PFE.
func (s *CompanyService) AddMeeting(report *entity.PfeProgressReport) error {
	return s.progressRepo.Insert(report)
}

// SubmitEvaluation soumet l'évaluation de l'encadrant entreprise.
func (s *CompanyService) SubmitEvaluation(assignmentID, evaluatorID string, criterion5 float64) error {
	if criterion5 < 0 || criterion5 > 4 {
		return apperror.BadRequest("Le critère 5 doit être entre 0 et 4")
	}
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

// ListReports liste les signalements de l'entreprise.
func (s *CompanyService) ListReports(userID string) ([]*entity.CompanyReport, error) {
	company, err := s.getCompanyByProfile(userID)
	if err != nil {
		return nil, err
	}
	return s.companyReportRepo.FindByCompany(company.ID)
}

// CreateReport crée un signalement.
func (s *CompanyService) CreateReport(userID string, report *entity.CompanyReport) error {
	company, err := s.getCompanyByProfile(userID)
	if err != nil {
		return err
	}
	report.ID = generateID()
	report.CompanyID = company.ID
	report.SubmittedBy = userID
	if report.Status == "" {
		report.Status = "en_attente"
	}
	return s.companyReportRepo.Insert(report)
}

// ListNotifications liste les notifications de l'entreprise.
func (s *CompanyService) ListNotifications(userID string) ([]*entity.Notification, error) {
	return s.notificationRepo.FindByRecipient(userID)
}
