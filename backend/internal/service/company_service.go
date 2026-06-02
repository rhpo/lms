package service

import (
	"database/sql"
	"fmt"

	"pfe-backend/internal/entity"
	"pfe-backend/internal/repository"
	"pfe-backend/internal/shared/apperror"
	"pfe-backend/internal/shared/notify"
	pfe_code "pfe-backend/internal/shared/pfe_code"
)

type CompanyService struct {
	profileRepo       *repository.ProfileRepository
	companyRepo       *repository.CompanyRepository
	teacherRepo       *repository.TeacherRepository
	studentRepo       *repository.StudentRepository
	specialityRepo    *repository.SpecialityRepository
	pfeSubjectRepo    *repository.PfeSubjectRepository
	wishRepo          *repository.WishRepository
	pfeAssignmentRepo *repository.PfeAssignmentRepository
	progressRepo      *repository.ProgressReportRepository
	supEvalRepo       *repository.SupervisorEvaluationRepository
	companyReportRepo *repository.CompanyReportRepository
	notificationRepo  *repository.NotificationRepository
	academicYearRepo  *repository.AcademicYearRepository
	notifier          *notify.Notifier
}

func NewCompanyService(
	profileRepo *repository.ProfileRepository,
	companyRepo *repository.CompanyRepository,
	teacherRepo *repository.TeacherRepository,
	studentRepo *repository.StudentRepository,
	specialityRepo *repository.SpecialityRepository,
	pfeSubjectRepo *repository.PfeSubjectRepository,
	wishRepo *repository.WishRepository,
	pfeAssignmentRepo *repository.PfeAssignmentRepository,
	progressRepo *repository.ProgressReportRepository,
	supEvalRepo *repository.SupervisorEvaluationRepository,
	companyReportRepo *repository.CompanyReportRepository,
	notificationRepo *repository.NotificationRepository,
	academicYearRepo *repository.AcademicYearRepository,
	notifier *notify.Notifier,
) *CompanyService {
	return &CompanyService{
		profileRepo:       profileRepo,
		companyRepo:       companyRepo,
		teacherRepo:       teacherRepo,
		studentRepo:       studentRepo,
		specialityRepo:    specialityRepo,
		pfeSubjectRepo:    pfeSubjectRepo,
		wishRepo:          wishRepo,
		pfeAssignmentRepo: pfeAssignmentRepo,
		progressRepo:      progressRepo,
		supEvalRepo:       supEvalRepo,
		companyReportRepo: companyReportRepo,
		notificationRepo:  notificationRepo,
		academicYearRepo:  academicYearRepo,
		notifier:          notifier,
	}
}

func (s *CompanyService) getCompanyByProfile(userID int64) (*entity.Company, error) {
	company, err := s.companyRepo.FindByProfileID(userID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, apperror.NotFound("Profil entreprise introuvable")
	}
	return company, nil
}

func (s *CompanyService) Dashboard(userID int64) (map[string]any, error) {
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

func (s *CompanyService) ListSubjects(userID int64) ([]*entity.PfeSubject, error) {
	company, err := s.getCompanyByProfile(userID)
	if err != nil {
		return nil, err
	}
	subjects, err := s.pfeSubjectRepo.FindByCompany(company.ID)
	if err != nil {
		return nil, err
	}
	for _, sub := range subjects {
		s.hydrateSubject(sub)
	}
	return subjects, nil
}

func (s *CompanyService) CreateSubject(userID int64, subject *entity.PfeSubject, domainIDs []int64) error {
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

	subject.CompanyID = entity.NullInt64{NullInt64: sql.NullInt64{Int64: company.ID, Valid: true}}
	subject.ProposerID = company.ProfileID
	subject.ProposerRole = "company"
	subject.AcademicYearID = ay.ID
	if subject.Status == "" {
		subject.Status = "en_attente"
	}

	if err := s.pfeSubjectRepo.Insert(subject); err != nil {
		return err
	}

	for _, did := range domainIDs {
		_ = s.pfeSubjectRepo.AddDomain(subject.ID, did)
	}

	go s.notifier.Send(company.ProfileID, notify.TypeSujet,
		fmt.Sprintf("Votre sujet « %s » a bien été soumis et est en attente de validation.", subject.Title))
	go s.notifier.NotifyAdmins(notify.TypeSujet,
		fmt.Sprintf("Nouveau sujet soumis par %s : « %s ».", company.CompanyName, subject.Title))

	return nil
}

func isCompanySubject(subject *entity.PfeSubject, companyID, companyProfileID int64) bool {
	if subject.CompanyID.Valid && subject.CompanyID.Int64 == companyID {
		return true
	}

	return subject.ProposerID == companyProfileID && subject.ProposerRole == "company"
}

func (s *CompanyService) GetSubject(userID, id int64) (*entity.PfeSubject, error) {
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
	if !isCompanySubject(subject, company.ID, company.ProfileID) {
		return nil, apperror.Forbidden("Accès non autorisé à ce sujet")
	}
	return subject, nil
}

func (s *CompanyService) UpdateSubject(userID int64, subject *entity.PfeSubject) error {
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
	if !isCompanySubject(existing, company.ID, company.ProfileID) {
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

func (s *CompanyService) ListCandidats(subjectID int64) ([]*entity.Wish, error) {
	wishes, err := s.wishRepo.FindBySubject(subjectID)
	if err != nil {
		return nil, err
	}
	for _, w := range wishes {
		w.Student = s.hydrateStudent(w.StudentID)
		sub, _ := s.pfeSubjectRepo.FindByID(w.SubjectID)
		if sub != nil {
			w.Subject = sub
		}
	}
	return wishes, nil
}

func (s *CompanyService) AcceptCandidats(subjectID int64, studentIDs []int64) error {

	existing, _ := s.pfeAssignmentRepo.FindBySubjectID(subjectID)
	if existing != nil {
		return apperror.Conflict("Les étudiants ont déjà été affectés à ce sujet")
	}

	subject, err := s.pfeSubjectRepo.FindByID(subjectID)
	if err != nil || subject == nil {
		return apperror.NotFound("Sujet introuvable")
	}
	if !subject.Validator1ID.Valid {
		return apperror.Conflict("Ce sujet n'a pas encore de validateur assigné - contactez l'administration")
	}

	ay, err := s.academicYearRepo.FindActive()
	if err != nil || ay == nil {
		return apperror.NotFound("Aucune année académique active")
	}

	wishes, err := s.wishRepo.FindBySubject(subjectID)
	if err != nil {
		return err
	}
	selectedSet := make(map[int64]bool, len(studentIDs))
	for _, id := range studentIDs {
		selectedSet[id] = true
	}
	for _, w := range wishes {
		if selectedSet[w.StudentID] {
			w.Status = "accepte"
			_ = s.wishRepo.Update(w)
		} else if w.Status == "en_attente" {
			w.Status = "refuse"
			_ = s.wishRepo.Update(w)
		}
	}

	specialityCode := "GEN"
	student1, err := s.studentRepo.FindByID(studentIDs[0])
	if err == nil && student1 != nil && student1.SpecialityID != nil {
		if sp, err := s.specialityRepo.FindByID(*student1.SpecialityID); err == nil && sp != nil {
			specialityCode = sp.Code
		}
	}
	seq, _ := s.pfeAssignmentRepo.CountBySpecialityAndYear(ay.ID, specialityCode)
	code := pfe_code.Generate(specialityCode, ay.Label, seq+1)

	assignment := &entity.PfeAssignment{
		PfeCode:        code,
		SubjectID:      subjectID,
		AcademicYearID: ay.ID,
		StudentID:      studentIDs[0],
		SupervisorID:   subject.Validator1ID.Int64,
		Status:         "en_cours",
	}
	if len(studentIDs) >= 2 {
		assignment.Student2ID = entity.NullInt64{NullInt64: sql.NullInt64{Int64: studentIDs[1], Valid: true}}
	}
	if len(studentIDs) >= 3 {
		assignment.Student3ID = entity.NullInt64{NullInt64: sql.NullInt64{Int64: studentIDs[2], Valid: true}}
	}
	return s.pfeAssignmentRepo.Insert(assignment)
}

func (s *CompanyService) RejectCandidat(subjectID, studentID int64) error {
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

func (s *CompanyService) ListSupervisedPFEs(userID int64) ([]*entity.PfeAssignment, error) {
	company, err := s.getCompanyByProfile(userID)
	if err != nil {
		return nil, err
	}
	assignments, err := s.pfeAssignmentRepo.FindByCompanySubject(company.ID)
	if err != nil {
		return nil, err
	}
	for _, a := range assignments {
		s.hydrateAssignment(a)
	}
	return assignments, nil
}

func (s *CompanyService) GetSupervisedPFE(id int64) (*entity.PfeAssignment, error) {
	a, err := s.pfeAssignmentRepo.FindByID(id)
	if err != nil || a == nil {
		return a, err
	}
	s.hydrateAssignment(a)
	return a, nil
}

func (s *CompanyService) AddMeeting(report *entity.PfeProgressReport) error {
	return s.progressRepo.Insert(report)
}

func (s *CompanyService) ListMeetings(assignmentID int64) ([]*entity.PfeProgressReport, error) {
	reports, err := s.progressRepo.FindByAssignment(assignmentID)
	if err != nil {
		return nil, err
	}
	if reports == nil {
		return []*entity.PfeProgressReport{}, nil
	}
	return reports, nil
}

func (s *CompanyService) GetEvaluation(assignmentID int64) (*entity.SupervisorEvaluation, error) {
	return s.supEvalRepo.FindByAssignment(assignmentID)
}

func (s *CompanyService) SubmitEvaluation(assignmentID, evaluatorID int64, criterion5 float64) error {
	if criterion5 < 0 || criterion5 > 4 {
		return apperror.BadRequest("Le critère 5 doit être entre 0 et 4")
	}

	assignment, err := s.pfeAssignmentRepo.FindByID(assignmentID)
	if err != nil || assignment == nil {
		return apperror.NotFound("Affectation introuvable")
	}
	existing, err := s.supEvalRepo.FindByAssignment(assignmentID)
	if err != nil {
		return err
	}
	if existing != nil {
		existing.EvaluatorID = assignment.SupervisorID
		existing.Criterion5 = entity.NullFloat64{NullFloat64: sql.NullFloat64{Float64: criterion5, Valid: true}}
		return s.supEvalRepo.Update(existing)
	}
	eval := &entity.SupervisorEvaluation{
		PfeAssignmentID: assignmentID,
		EvaluatorID:     assignment.SupervisorID,
		Criterion5:      entity.NullFloat64{NullFloat64: sql.NullFloat64{Float64: criterion5, Valid: true}},
	}
	return s.supEvalRepo.Insert(eval)
}

func (s *CompanyService) ListReports(userID int64) ([]*entity.CompanyReport, error) {
	company, err := s.getCompanyByProfile(userID)
	if err != nil {
		return nil, err
	}
	return s.companyReportRepo.FindByCompany(company.ID)
}

func (s *CompanyService) CreateReport(userID int64, report *entity.CompanyReport) error {
	company, err := s.getCompanyByProfile(userID)
	if err != nil {
		return err
	}
	report.CompanyID = company.ID
	report.SubmittedBy = userID
	if report.Status == "" {
		report.Status = "en_attente"
	}
	return s.companyReportRepo.Insert(report)
}

func (s *CompanyService) ListNotifications(userID int64) ([]*entity.Notification, error) {
	return s.notificationRepo.FindByRecipient(userID)
}

func (s *CompanyService) GetStudentProfileID(studentID int64) (int64, error) {
	st, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return 0, err
	}
	if st == nil {
		return 0, apperror.NotFound("Étudiant introuvable")
	}
	return st.ProfileID, nil
}

func (s *CompanyService) GetSubjectTitle(subjectID int64) string {
	sub, err := s.pfeSubjectRepo.FindByID(subjectID)
	if err != nil || sub == nil {
		return fmt.Sprintf("sujet #%d", subjectID)
	}
	return sub.Title
}

func (s *CompanyService) hydrateTeacher(id int64) *entity.Teacher {
	if id == 0 {
		return nil
	}
	t, _ := s.teacherRepo.FindByID(id)
	if t == nil {
		t, _ = s.teacherRepo.FindByProfileID(id)
	}
	if t != nil {
		t.Profile, _ = s.profileRepo.FindByID(t.ProfileID)
	}
	return t
}

func (s *CompanyService) hydrateStudent(id int64) *entity.Student {
	if id == 0 {
		return nil
	}
	st, _ := s.studentRepo.FindByID(id)
	if st == nil {
		st, _ = s.studentRepo.FindByProfileID(id)
	}
	if st != nil {
		st.Profile, _ = s.profileRepo.FindByID(st.ProfileID)
		if st.SpecialityID != nil {
			st.Speciality, _ = s.specialityRepo.FindByID(*st.SpecialityID)
		}
	}
	return st
}

func (s *CompanyService) hydrateSubject(sub *entity.PfeSubject) {
	sub.Proposer, _ = s.profileRepo.FindByID(sub.ProposerID)
	if sub.CompanyID.Valid {
		sub.Company, _ = s.companyRepo.FindByID(sub.CompanyID.Int64)
		if sub.Company == nil {
			sub.Company, _ = s.companyRepo.FindByProfileID(sub.CompanyID.Int64)
		}
	}
	if sub.Validator1ID.Valid {
		sub.Validator1 = s.hydrateTeacher(sub.Validator1ID.Int64)
	}
	if sub.Validator2ID.Valid {
		sub.Validator2 = s.hydrateTeacher(sub.Validator2ID.Int64)
	}
	if sub.CoSupervisorID.Valid {
		sub.CoSupervisor = s.hydrateTeacher(sub.CoSupervisorID.Int64)
	}
	sub.Domains, _ = s.pfeSubjectRepo.GetDomains(sub.ID)
}

func (s *CompanyService) hydrateAssignment(a *entity.PfeAssignment) {
	sub, _ := s.pfeSubjectRepo.FindByID(a.SubjectID)
	if sub != nil {
		s.hydrateSubject(sub)
		a.Subject = sub
	}
	a.Student = s.hydrateStudent(a.StudentID)
	if a.Student2ID.Valid {
		a.Student2 = s.hydrateStudent(a.Student2ID.Int64)
	}
	if a.Student3ID.Valid {
		a.Student3 = s.hydrateStudent(a.Student3ID.Int64)
	}
	a.Supervisor = s.hydrateTeacher(a.SupervisorID)
	if a.CoSupervisorID.Valid {
		a.CoSupervisor = s.hydrateTeacher(a.CoSupervisorID.Int64)
	}
	ay, _ := s.academicYearRepo.FindByID(a.AcademicYearID)
	a.AcademicYear = ay
}
