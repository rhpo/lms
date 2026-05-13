package service

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"pfe-backend/internal/entity"
	"pfe-backend/internal/repository"
	"pfe-backend/internal/shared/apperror"
	"strings"
	"time"
)

// AdminService gère la logique métier de l'admin.
type AdminService struct {
	profileRepo       *repository.ProfileRepository
	teacherRepo       *repository.TeacherRepository
	studentRepo       *repository.StudentRepository
	companyRepo       *repository.CompanyRepository
	domainRepo        *repository.DomainRepository
	specialityRepo    *repository.SpecialityRepository
	promotionRepo     *repository.PromotionRepository
	academicYearRepo  *repository.AcademicYearRepository
	pfeSubjectRepo    *repository.PfeSubjectRepository
	wishRepo          *repository.WishRepository
	pfeAssignmentRepo *repository.PfeAssignmentRepository
	progressRepo      *repository.ProgressReportRepository
	defenseJuryRepo   *repository.DefenseJuryRepository
	defenseRepo       *repository.DefenseRepository
	juryGradeRepo     *repository.JuryGradeRepository
	supEvalRepo       *repository.SupervisorEvaluationRepository
	companyReportRepo *repository.CompanyReportRepository
	notificationRepo  *repository.NotificationRepository
	auditLogRepo      *repository.AuditLogRepository
}

// NewAdminService crée un nouveau AdminService.
func NewAdminService(
	profileRepo *repository.ProfileRepository,
	teacherRepo *repository.TeacherRepository,
	studentRepo *repository.StudentRepository,
	companyRepo *repository.CompanyRepository,
	domainRepo *repository.DomainRepository,
	specialityRepo *repository.SpecialityRepository,
	promotionRepo *repository.PromotionRepository,
	academicYearRepo *repository.AcademicYearRepository,
	pfeSubjectRepo *repository.PfeSubjectRepository,
	wishRepo *repository.WishRepository,
	pfeAssignmentRepo *repository.PfeAssignmentRepository,
	progressRepo *repository.ProgressReportRepository,
	defenseJuryRepo *repository.DefenseJuryRepository,
	defenseRepo *repository.DefenseRepository,
	juryGradeRepo *repository.JuryGradeRepository,
	supEvalRepo *repository.SupervisorEvaluationRepository,
	companyReportRepo *repository.CompanyReportRepository,
	notificationRepo *repository.NotificationRepository,
	auditLogRepo *repository.AuditLogRepository,
) *AdminService {
	return &AdminService{
		profileRepo:       profileRepo,
		teacherRepo:       teacherRepo,
		studentRepo:       studentRepo,
		companyRepo:       companyRepo,
		domainRepo:        domainRepo,
		specialityRepo:    specialityRepo,
		promotionRepo:     promotionRepo,
		academicYearRepo:  academicYearRepo,
		pfeSubjectRepo:    pfeSubjectRepo,
		wishRepo:          wishRepo,
		pfeAssignmentRepo: pfeAssignmentRepo,
		progressRepo:      progressRepo,
		defenseJuryRepo:   defenseJuryRepo,
		defenseRepo:       defenseRepo,
		juryGradeRepo:     juryGradeRepo,
		supEvalRepo:       supEvalRepo,
		companyReportRepo: companyReportRepo,
		notificationRepo:  notificationRepo,
		auditLogRepo:      auditLogRepo,
	}
}

// Dashboard retourne les statistiques du tableau de bord admin.
func (s *AdminService) Dashboard() (map[string]any, error) {
	profiles, err := s.profileRepo.FindAll()
	if err != nil {
		return nil, err
	}
	teachers, _ := s.teacherRepo.FindAll()
	students, _ := s.studentRepo.FindAll()
	companies, _ := s.companyRepo.FindAll()
	subjects, _ := s.pfeSubjectRepo.FindAll()
	assignments, _ := s.pfeAssignmentRepo.FindAll()
	defenses, _ := s.defenseRepo.FindAll()
	reports, _ := s.companyReportRepo.FindAll()

	return map[string]any{
		"total_users":       len(profiles),
		"total_teachers":    len(teachers),
		"total_students":    len(students),
		"total_companies":   len(companies),
		"total_subjects":    len(subjects),
		"total_assignments": len(assignments),
		"total_defenses":    len(defenses),
		"total_reports":     len(reports),
	}, nil
}

// ListUsers retourne tous les profils utilisateurs.
func (s *AdminService) ListUsers() ([]*entity.Profile, error) {
	return s.profileRepo.FindAll()
}

// GetUser retourne un profil par son ID.
func (s *AdminService) GetUser(id string) (*entity.Profile, error) {
	return s.profileRepo.FindByID(id)
}

// CreateUser crée un nouveau profil utilisateur. Génère l'ID si absent.
func (s *AdminService) CreateUser(profile *entity.Profile) error {
	if profile.ID == "" {
		profile.ID = generateID()
	}
	return s.profileRepo.Insert(profile)
}

// UpdateUser met à jour un profil utilisateur en fusionnant les champs non-vides.
func (s *AdminService) UpdateUser(profile *entity.Profile) error {
	existing, err := s.profileRepo.FindByID(profile.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return apperror.NotFound("Utilisateur introuvable")
	}
	if profile.FullName != "" {
		existing.FullName = profile.FullName
	}
	if profile.Email != "" {
		existing.Email = profile.Email
	}
	if profile.Role != "" {
		existing.Role = profile.Role
	}
	return s.profileRepo.Update(existing)
}

// DeactivateUser désactive un utilisateur.
func (s *AdminService) DeactivateUser(id string) error {
	profile, err := s.profileRepo.FindByID(id)
	if err != nil {
		return err
	}
	if profile == nil {
		return apperror.NotFound("Utilisateur introuvable")
	}
	profile.IsActive = false
	return s.profileRepo.Update(profile)
}

// ReactivateUser réactive un utilisateur.
func (s *AdminService) ReactivateUser(id string) error {
	profile, err := s.profileRepo.FindByID(id)
	if err != nil {
		return err
	}
	if profile == nil {
		return apperror.NotFound("Utilisateur introuvable")
	}
	profile.IsActive = true
	return s.profileRepo.Update(profile)
}

// ListCompanies retourne toutes les entreprises.
func (s *AdminService) ListCompanies() ([]*entity.Company, error) {
	return s.companyRepo.FindAll()
}

// ValidateCompany valide une entreprise.
func (s *AdminService) ValidateCompany(id string) error {
	company, err := s.companyRepo.FindByID(id)
	if err != nil {
		return err
	}
	if company == nil {
		return apperror.NotFound("Entreprise introuvable")
	}
	return s.companyRepo.UpdateVerification(id, true)
}

// RejectCompany rejette une entreprise.
func (s *AdminService) RejectCompany(id string) error {
	company, err := s.companyRepo.FindByID(id)
	if err != nil {
		return err
	}
	if company == nil {
		return apperror.NotFound("Entreprise introuvable")
	}
	return s.companyRepo.UpdateVerification(id, false)
}

// ListReports retourne tous les reports entreprises.
func (s *AdminService) ListReports() ([]*entity.CompanyReport, error) {
	return s.companyReportRepo.FindAll()
}

// ResolveReport résout un report.
func (s *AdminService) ResolveReport(id string) error {
	report, err := s.companyReportRepo.FindByID(id)
	if err != nil {
		return err
	}
	if report == nil {
		return apperror.NotFound("Report introuvable")
	}
	return s.companyReportRepo.UpdateStatus(id, "resolu")
}

// RejectReport rejette un report.
func (s *AdminService) RejectReport(id string) error {
	report, err := s.companyReportRepo.FindByID(id)
	if err != nil {
		return err
	}
	if report == nil {
		return apperror.NotFound("Report introuvable")
	}
	return s.companyReportRepo.UpdateStatus(id, "rejete")
}

// ListSubjects retourne tous les sujets PFE.
func (s *AdminService) ListSubjects() ([]*entity.PfeSubject, error) {
	return s.pfeSubjectRepo.FindAll()
}

// GetSubject retourne un sujet par son ID.
func (s *AdminService) GetSubject(id string) (*entity.PfeSubject, error) {
	return s.pfeSubjectRepo.FindByID(id)
}

// ListAssignments retourne toutes les affectations PFE.
func (s *AdminService) ListAssignments() ([]*entity.PfeAssignment, error) {
	return s.pfeAssignmentRepo.FindAll()
}

// GetAssignment retourne une affectation par son ID.
func (s *AdminService) GetAssignment(id string) (*entity.PfeAssignment, error) {
	return s.pfeAssignmentRepo.FindByID(id)
}

// ListDefenses retourne toutes les soutenances.
func (s *AdminService) ListDefenses() ([]*entity.Defense, error) {
	return s.defenseRepo.FindAll()
}

// GetDefense retourne une soutenance par son ID.
func (s *AdminService) GetDefense(id string) (*entity.Defense, error) {
	return s.defenseRepo.FindByID(id)
}

// ListAcademicYears retourne toutes les années académiques.
func (s *AdminService) ListAcademicYears() ([]*entity.AcademicYear, error) {
	return s.academicYearRepo.FindAll()
}

// CreateAcademicYear crée une année académique.
func (s *AdminService) CreateAcademicYear(ay *entity.AcademicYear) error {
	return s.academicYearRepo.Insert(ay)
}

// CloseAcademicYear ferme une année académique.
func (s *AdminService) CloseAcademicYear(id string) error {
	ay, err := s.academicYearRepo.FindByID(id)
	if err != nil {
		return err
	}
	if ay == nil {
		return apperror.NotFound("Année académique introuvable")
	}
	return s.academicYearRepo.Close(id)
}

// ListSpecialities retourne toutes les spécialités.
func (s *AdminService) ListSpecialities() ([]*entity.Speciality, error) {
	return s.specialityRepo.FindAll()
}

// CreateSpeciality crée une spécialité.
func (s *AdminService) CreateSpeciality(sp *entity.Speciality) error {
	return s.specialityRepo.Insert(sp)
}

// DeleteSpeciality supprime une spécialité.
func (s *AdminService) DeleteSpeciality(id string) error {
	return s.specialityRepo.Delete(id)
}

// ListDomains retourne tous les domaines.
func (s *AdminService) ListDomains() ([]*entity.Domain, error) {
	return s.domainRepo.FindAll()
}

// CreateDomain crée un domaine.
func (s *AdminService) CreateDomain(d *entity.Domain) error {
	return s.domainRepo.Insert(d)
}

// DeleteDomain supprime un domaine.
func (s *AdminService) DeleteDomain(id string) error {
	return s.domainRepo.Delete(id)
}

// ListPromotions retourne toutes les promotions.
func (s *AdminService) ListPromotions() ([]*entity.Promotion, error) {
	return s.promotionRepo.FindAll()
}

// CreatePromotion crée une promotion.
func (s *AdminService) CreatePromotion(p *entity.Promotion) error {
	return s.promotionRepo.Insert(p)
}

// DeletePromotion supprime une promotion.
func (s *AdminService) DeletePromotion(id string) error {
	return s.promotionRepo.Delete(id)
}

// AssignValidators assigne les validateurs à un sujet.
func (s *AdminService) AssignValidators(subjectID, validator1ID, validator2ID string) error {
	subject, err := s.pfeSubjectRepo.FindByID(subjectID)
	if err != nil {
		return err
	}
	if subject == nil {
		return apperror.NotFound("Sujet introuvable")
	}
	return s.pfeSubjectRepo.AssignValidators(subjectID, validator1ID, validator2ID)
}

// AssignCoSupervisor assigne un co-encadrant à un sujet.
func (s *AdminService) AssignCoSupervisor(subjectID, coSupervisorID string) error {
	subject, err := s.pfeSubjectRepo.FindByID(subjectID)
	if err != nil {
		return err
	}
	if subject == nil {
		return apperror.NotFound("Sujet introuvable")
	}
	return s.pfeSubjectRepo.AssignCoSupervisor(subjectID, coSupervisorID)
}

// GetStatistics retourne les statistiques globales.
func (s *AdminService) GetStatistics() (map[string]any, error) {
	return s.Dashboard()
}

// AuditLog retourne tous les logs d'audit.
func (s *AdminService) AuditLog() ([]*entity.AuditLog, error) {
	return s.auditLogRepo.FindAll()
}

// UserAction gère les actions sur un utilisateur.
func (s *AdminService) UserAction(id, action string) error {
	switch action {
	case "deactivate":
		return s.DeactivateUser(id)
	case "reactivate":
		return s.ReactivateUser(id)
	case "transfer-admin":
		return s.TransferAdmin(id)
	default:
		return apperror.BadRequest("Action non reconnue: " + action)
	}
}

// TransferAdmin transfère le rôle admin à un autre utilisateur.
func (s *AdminService) TransferAdmin(id string) error {
	profile, err := s.profileRepo.FindByID(id)
	if err != nil {
		return err
	}
	if profile == nil {
		return apperror.NotFound("Utilisateur introuvable")
	}
	if profile.Role != "teacher" {
		return apperror.BadRequest("Le transfert admin n'est possible que vers un enseignant")
	}
	// Désactiver l'admin actuel
	currentAdmin, err := s.FindAdmin()
	if err != nil {
		return err
	}
	if currentAdmin != nil {
		currentAdmin.Role = "teacher"
		if err := s.profileRepo.Update(currentAdmin); err != nil {
			return err
		}
	}
	profile.Role = "admin"
	return s.profileRepo.Update(profile)
}

// FindAdmin trouve le compte admin.
func (s *AdminService) FindAdmin() (*entity.Profile, error) {
	all, err := s.profileRepo.FindAll()
	if err != nil {
		return nil, err
	}
	for _, p := range all {
		if p.Role == "admin" {
			return p, nil
		}
	}
	return nil, nil
}

// ImportUsersCSV importe des utilisateurs depuis un CSV (colonnes: role,full_name,email).
func (s *AdminService) ImportUsersCSV(csvData string) error {
	r := csv.NewReader(strings.NewReader(csvData))
	records, err := r.ReadAll()
	if err != nil {
		return apperror.BadRequest("CSV invalide: " + err.Error())
	}
	if len(records) < 2 {
		return apperror.BadRequest("CSV vide ou sans données")
	}
	for i, row := range records[1:] {
		if len(row) < 3 {
			return apperror.BadRequest(fmt.Sprintf("Ligne %d invalide: 3 colonnes requises", i+2))
		}
		profile := &entity.Profile{
			ID:       generateID(),
			Role:     strings.TrimSpace(row[0]),
			FullName: strings.TrimSpace(row[1]),
			Email:    strings.TrimSpace(row[2]),
			IsActive: true,
		}
		if err := s.profileRepo.Insert(profile); err != nil {
			return fmt.Errorf("erreur import ligne %d: %w", i+2, err)
		}
	}
	return nil
}

// CompanyAction gère les actions sur une entreprise.
func (s *AdminService) CompanyAction(id, action string) error {
	switch action {
	case "validate":
		return s.ValidateCompany(id)
	case "reject":
		return s.RejectCompany(id)
	default:
		return apperror.BadRequest("Action non reconnue: " + action)
	}
}

// ReportAction gère les actions sur un report.
func (s *AdminService) ReportAction(id, action string) error {
	switch action {
	case "resolve":
		return s.ResolveReport(id)
	case "reject":
		return s.RejectReport(id)
	default:
		return apperror.BadRequest("Action non reconnue: " + action)
	}
}

// SubjectAction gère les actions admin sur un sujet.
func (s *AdminService) SubjectAction(id, action, validatorID string) error {
	switch action {
	case "assign-validators":
		if validatorID == "" {
			return apperror.BadRequest("validator_id requis")
		}
		return s.AssignValidators(id, validatorID, "")
	case "assign-co-supervisor":
		if validatorID == "" {
			return apperror.BadRequest("co_supervisor_id requis")
		}
		return s.AssignCoSupervisor(id, validatorID)
	case "unblock":
		return s.UnblockSubject(id)
	default:
		return apperror.BadRequest("Action non reconnue: " + action)
	}
}

// UnblockSubject débloque un sujet.
func (s *AdminService) UnblockSubject(id string) error {
	subject, err := s.pfeSubjectRepo.FindByID(id)
	if err != nil {
		return err
	}
	if subject == nil {
		return apperror.NotFound("Sujet introuvable")
	}
	return s.pfeSubjectRepo.UpdateStatus(id, "en_attente")
}

// CreateDefense crée une nouvelle soutenance avec son jury.
func (s *AdminService) CreateDefense(assignmentID, presidentID, memberID, scheduledAt, room string) (*entity.Defense, error) {
	if assignmentID == "" || presidentID == "" || memberID == "" {
		return nil, apperror.BadRequest("assignment_id, president_id et member_id sont requis")
	}
	if presidentID == memberID {
		return nil, apperror.BadRequest("Le président et le membre du jury doivent être différents")
	}
	// Vérifier que l'assignation existe
	assignment, err := s.pfeAssignmentRepo.FindByID(assignmentID)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, apperror.NotFound("Affectation introuvable")
	}

	// Créer le jury
	jury := &entity.DefenseJury{
		ID:           generateID(),
		AssignmentID: assignmentID,
		PresidentID:  presidentID,
		MemberID:     memberID,
	}
	if err := s.defenseJuryRepo.Insert(jury); err != nil {
		return nil, err
	}

	// Parsing du temps scheduledAt
	var scheduledAtTime sql.NullTime
	if t, err := time.Parse(time.RFC3339, scheduledAt); err == nil {
		scheduledAtTime = sql.NullTime{Time: t, Valid: true}
	} else {
		scheduledAtTime = sql.NullTime{Valid: false}
	}

	// Créer la soutenance
	defense := &entity.Defense{
		ID:           generateID(),
		AssignmentID: assignmentID,
		JuryID:       jury.ID,
		ScheduledAt:  scheduledAtTime,
		Room:         sql.NullString{String: room, Valid: room != ""},
		Status:       "scheduled",
	}
	if err := s.defenseRepo.Insert(defense); err != nil {
		return nil, err
	}

	// Mettre à jour le statut de l'assignation
	_ = s.pfeAssignmentRepo.UpdateStatus(assignmentID, "soutenance_planifiee")

	// Notifier les étudiants et le jury
	s.notifyDefenseScheduled(defense, jury)

	return defense, nil
}

// notifyDefenseScheduled envoie les notifications pour une soutenance planifiée.
func (s *AdminService) notifyDefenseScheduled(defense *entity.Defense, jury *entity.DefenseJury) {
	// Notification aux étudiants
	assignment, _ := s.pfeAssignmentRepo.FindByID(defense.AssignmentID)
	if assignment != nil {
		students := []string{assignment.StudentID}
		if assignment.Student2ID.Valid {
			students = append(students, assignment.Student2ID.String)
		}
		if assignment.Student3ID.Valid {
			students = append(students, assignment.Student3ID.String)
		}
		for _, studentID := range students {
			profile, _ := s.profileRepo.FindByID(studentID)
			if profile != nil {
				s.notificationRepo.Insert(&entity.Notification{
					ID:          generateID(),
					RecipientID: profile.ID,
					Type:        "soutenance_planifiee",
				})
			}
		}
	}
	// Notification au jury
	s.notificationRepo.Insert(&entity.Notification{
		ID:          generateID(),
		RecipientID: jury.PresidentID,
		Type:        "soutenance_planifiee",
	})
	s.notificationRepo.Insert(&entity.Notification{
		ID:          generateID(),
		RecipientID: jury.MemberID,
		Type:        "soutenance_planifiee",
	})
}

// RecommendJury recommande un jury pour un PFE.
func (s *AdminService) RecommendJury(pfeID string) (map[string]any, error) {
	// Trouver des enseignants disponibles
	teachers, err := s.teacherRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"recommended": teachers,
		"pfe_id":      pfeID,
	}, nil
}

// SubmitGrade soumet une note jury pour l'utilisateur authentifié.
func (s *AdminService) SubmitGrade(defenseID, callerID string, c1, c2, c3, c4 float64) error {
	// Valider les critères (0-4)
	for _, v := range []float64{c1, c2, c3, c4} {
		if v < 0 || v > 4 {
			return apperror.BadRequest("Les critères doivent être entre 0 et 4")
		}
	}

	// Vérifier la défense
	defense, err := s.defenseRepo.FindByID(defenseID)
	if err != nil {
		return err
	}
	if defense == nil {
		return apperror.NotFound("Soutenance introuvable")
	}

	// Upsert: vérifier si une note existe déjà pour ce membre
	existing, err := s.juryGradeRepo.FindByDefenseAndMember(defenseID, callerID)
	if err != nil {
		return err
	}

	if existing != nil {
		existing.Criterion1 = sql.NullFloat64{Float64: c1, Valid: true}
		existing.Criterion2 = sql.NullFloat64{Float64: c2, Valid: true}
		existing.Criterion3 = sql.NullFloat64{Float64: c3, Valid: true}
		existing.Criterion4 = sql.NullFloat64{Float64: c4, Valid: true}
		return s.juryGradeRepo.Update(existing)
	}

	grade := &entity.JuryGrade{
		ID:           generateID(),
		DefenseID:    defenseID,
		JuryMemberID: callerID,
		Criterion1:   sql.NullFloat64{Float64: c1, Valid: true},
		Criterion2:   sql.NullFloat64{Float64: c2, Valid: true},
		Criterion3:   sql.NullFloat64{Float64: c3, Valid: true},
		Criterion4:   sql.NullFloat64{Float64: c4, Valid: true},
	}
	return s.juryGradeRepo.Insert(grade)
}

// ResolveGradeRequest contient les paramètres de résolution de notes.
type ResolveGradeRequest struct {
	Choice     string             // "president", "member", "new", ou "" (ancien mode direct)
	Criterion1 float64
	Criterion2 float64
	Criterion3 float64
	Criterion4 float64
	Grades     map[string]float64 // pour choice="new"
}

// ResolveGrade résout la note finale d'une soutenance.
func (s *AdminService) ResolveGrade(defenseID string, req ResolveGradeRequest) error {
	defense, err := s.defenseRepo.FindByID(defenseID)
	if err != nil {
		return err
	}
	if defense == nil {
		return apperror.NotFound("Soutenance introuvable")
	}

	var c1, c2, c3, c4 float64

	switch req.Choice {
	case "president", "member":
		// Trouver le jury
		jury, err := s.defenseJuryRepo.FindByID(defense.JuryID)
		if err != nil {
			return err
		}
		if jury == nil {
			return apperror.NotFound("Jury introuvable")
		}

		var memberID string
		if req.Choice == "president" {
			memberID = jury.PresidentID
		} else {
			memberID = jury.MemberID
		}

		grade, err := s.juryGradeRepo.FindByDefenseAndMember(defenseID, memberID)
		if err != nil {
			return err
		}
		if grade == nil {
			return apperror.BadRequest("Aucune note soumise par ce membre du jury")
		}
		c1 = grade.Criterion1.Float64
		c2 = grade.Criterion2.Float64
		c3 = grade.Criterion3.Float64
		c4 = grade.Criterion4.Float64

	case "new":
		if req.Grades == nil {
			return apperror.BadRequest("Les notes sont requises pour le choix 'new'")
		}
		c1 = req.Grades["criterion1"]
		c2 = req.Grades["criterion2"]
		c3 = req.Grades["criterion3"]
		c4 = req.Grades["criterion4"]

	default:
		// Mode direct (rétrocompatibilité)
		c1 = req.Criterion1
		c2 = req.Criterion2
		c3 = req.Criterion3
		c4 = req.Criterion4
	}

	// Récupérer l'évaluation encadrant
	assignment, err := s.pfeAssignmentRepo.FindByID(defense.AssignmentID)
	if err != nil {
		return err
	}
	if assignment == nil {
		return apperror.NotFound("Affectation introuvable")
	}

	supEval, _ := s.supEvalRepo.FindByAssignment(assignment.ID)
	criterion5 := 0.0
	if supEval != nil {
		criterion5 = supEval.Criterion5.Float64
	}

	totalGrade := c1 + c2 + c3 + c4 + criterion5
	return s.defenseRepo.UpdateResult(defenseID, "admitted", totalGrade)
}

// findJuryByIDOrDefenseID cherche un jury directement ou via une soutenance.
func (s *AdminService) findJuryByIDOrDefenseID(id string) (*entity.DefenseJury, error) {
	jury, err := s.defenseJuryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if jury != nil {
		return jury, nil
	}
	// Essayer comme ID de soutenance
	defense, err := s.defenseRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if defense == nil {
		return nil, apperror.NotFound("Jury introuvable")
	}
	jury, err = s.defenseJuryRepo.FindByID(defense.JuryID)
	if err != nil {
		return nil, err
	}
	if jury == nil {
		return nil, apperror.NotFound("Jury introuvable")
	}
	return jury, nil
}

// ConfirmJury confirme la participation d'un jury (accepte ID jury ou ID soutenance).
func (s *AdminService) ConfirmJury(id string) error {
	jury, err := s.findJuryByIDOrDefenseID(id)
	if err != nil {
		return err
	}
	if err := s.defenseJuryRepo.ConfirmPresident(jury.ID); err != nil {
		return err
	}
	return s.defenseJuryRepo.ConfirmMember(jury.ID)
}

// DeclineJury décline la participation d'un jury (accepte ID jury ou ID soutenance).
func (s *AdminService) DeclineJury(id string) error {
	jury, err := s.findJuryByIDOrDefenseID(id)
	if err != nil {
		return err
	}
	return s.defenseJuryRepo.Delete(jury.ID)
}

// UpdateDeadlines met à jour les délais de soumission.
func (s *AdminService) UpdateDeadlines(openAt, closeAt string, maxWishes int) error {
	// Trouver l'année académique active
	years, err := s.academicYearRepo.FindAll()
	if err != nil {
		return err
	}
	for _, y := range years {
		if y.Status == "active" {
			// Parsing des dates
			if t, err := time.Parse(time.RFC3339, openAt); err == nil {
				y.SubmissionOpenAt = sql.NullTime{Time: t, Valid: true}
			}
			if t, err := time.Parse(time.RFC3339, closeAt); err == nil {
				y.SubmissionCloseAt = sql.NullTime{Time: t, Valid: true}
			}
			y.MaxWishes = maxWishes
			return s.academicYearRepo.Update(y)
		}
	}
	return apperror.NotFound("Aucune année académique active trouvée")
}

// idCounter is used to ensure unique IDs even within the same nanosecond.
var idCounter int64

// generateID génère un ID unique.
func generateID() string {
	idCounter++
	return fmt.Sprintf("id-%d-%d", time.Now().UnixNano(), idCounter)
}
