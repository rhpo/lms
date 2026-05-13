package tests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"pfe-backend/internal/config"
	"pfe-backend/internal/entity"
	"pfe-backend/internal/handler"
	"pfe-backend/internal/repository"
	"pfe-backend/internal/service"
	"pfe-backend/internal/shared/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/golang-jwt/jwt/v5"
	_ "modernc.org/sqlite"
)

// TestHelper fournit tout le nécessaire pour les tests d'intégration HTTP.
type TestHelper struct {
	App   *fiber.App
	DB    *sql.DB
	Cfg   *config.Config
	Admin string // token admin
}

// NewTestHelper initialise un serveur de test avec une base SQLite en mémoire.
func NewTestHelper() *TestHelper {
	os.Setenv("ENV", "development")
	os.Setenv("SUPABASE_URL", "https://test.supabase.co")
	os.Setenv("SUPABASE_SERVICE_ROLE_KEY", "test-key")
	os.Setenv("RESEND_API_KEY", "test-resend-key")
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	os.Setenv("DATABASE_PATH", ":memory:")
	os.Setenv("PORT", "9099")

	cfg := config.Load()

	db, err := sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		panic(fmt.Sprintf("Erreur ouverture DB test: %v", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Erreur ping DB test: %v", err))
	}

	if err := runTestMigrations(db); err != nil {
		panic(fmt.Sprintf("Erreur migration test: %v", err))
	}

	if err := runTestSeed(db); err != nil {
		panic(fmt.Sprintf("Erreur seed test: %v", err))
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(map[string]any{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	}))

	// Repositories
	profileRepo := repository.NewProfileRepository(db)
	teacherRepo := repository.NewTeacherRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	domainRepo := repository.NewDomainRepository(db)
	specialityRepo := repository.NewSpecialityRepository(db)
	promotionRepo := repository.NewPromotionRepository(db)
	academicYearRepo := repository.NewAcademicYearRepository(db)
	pfeSubjectRepo := repository.NewPfeSubjectRepository(db)
	wishRepo := repository.NewWishRepository(db)
	pfeAssignmentRepo := repository.NewPfeAssignmentRepository(db)
	progressRepo := repository.NewProgressReportRepository(db)
	defenseJuryRepo := repository.NewDefenseJuryRepository(db)
	defenseRepo := repository.NewDefenseRepository(db)
	juryGradeRepo := repository.NewJuryGradeRepository(db)
	supEvalRepo := repository.NewSupervisorEvaluationRepository(db)
	companyReportRepo := repository.NewCompanyReportRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)

	// Services
	authService := service.NewAuthService(profileRepo, cfg)
	adminService := service.NewAdminService(
		profileRepo, teacherRepo, studentRepo, companyRepo,
		domainRepo, specialityRepo, promotionRepo, academicYearRepo,
		pfeSubjectRepo, wishRepo, pfeAssignmentRepo,
		progressRepo, defenseJuryRepo, defenseRepo,
		juryGradeRepo, supEvalRepo, companyReportRepo,
		notificationRepo, auditLogRepo,
	)
	teacherService := service.NewTeacherService(
		profileRepo, teacherRepo, pfeSubjectRepo, wishRepo,
		pfeAssignmentRepo, progressRepo, defenseJuryRepo,
		defenseRepo, supEvalRepo, notificationRepo, academicYearRepo,
	)
	studentService := service.NewStudentService(
		profileRepo, studentRepo, pfeSubjectRepo, wishRepo,
		pfeAssignmentRepo, progressRepo, defenseRepo,
		defenseJuryRepo, notificationRepo, academicYearRepo,
	)
	companyService := service.NewCompanyService(
		profileRepo, companyRepo, pfeSubjectRepo, wishRepo,
		pfeAssignmentRepo, progressRepo, supEvalRepo,
		companyReportRepo, notificationRepo, academicYearRepo,
	)

	// Handlers
	authHandler := handler.NewAuthHandler(authService, cfg)
	adminHandler := handler.NewAdminHandler(adminService)
	teacherHandler := handler.NewTeacherHandler(teacherService)
	studentHandler := handler.NewStudentHandler(studentService)
	companyHandler := handler.NewCompanyHandler(companyService)
	uploadHandler := handler.NewUploadHandler(profileRepo, companyRepo, "./uploads")

	app.Static("/uploads", "./uploads")

	api := app.Group("/api")

	api.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(map[string]any{"success": true, "data": map[string]any{"status": "ok"}})
	})

	ref := api.Group("/ref", middleware.AuthRequired(cfg))
	ref.Get("/domains", func(c fiber.Ctx) error {
		domains, err := domainRepo.FindAll()
		if err != nil {
			return c.Status(500).JSON(map[string]any{"success": false, "error": "Erreur serveur"})
		}
		if domains == nil {
			domains = []*entity.Domain{}
		}
		return c.JSON(map[string]any{"success": true, "data": domains})
	})
	ref.Get("/specialities", func(c fiber.Ctx) error {
		specialities, err := specialityRepo.FindAll()
		if err != nil {
			return c.Status(500).JSON(map[string]any{"success": false, "error": "Erreur serveur"})
		}
		if specialities == nil {
			specialities = []*entity.Speciality{}
		}
		return c.JSON(map[string]any{"success": true, "data": specialities})
	})

	api.Post("/auth/dev-login", authHandler.DevLogin)
	auth := api.Group("/auth")
	auth.Use(middleware.AuthRequired(cfg))
	auth.Get("/me", authHandler.Me)
	auth.Post("/logout", authHandler.Logout)

	api.Post("/profile/avatar", middleware.AuthRequired(cfg), uploadHandler.UploadAvatar)

	upload := api.Group("/upload", middleware.AuthRequired(cfg))
	upload.Post("/company-logo", uploadHandler.UploadCompanyLogo)
	upload.Post("/memoire", uploadHandler.UploadMemoire)
	upload.Post("/avatar", uploadHandler.UploadAvatar)

	admin := api.Group("/admin", middleware.AuthRequired(cfg), middleware.RequireRole("admin"))
	admin.Get("/dashboard", adminHandler.Dashboard)
	admin.Get("/accounts/users", adminHandler.ListUsers)
	admin.Post("/accounts/users", adminHandler.CreateUser)
	admin.Get("/accounts/users/:id", adminHandler.GetUser)
	admin.Patch("/accounts/users/:id", adminHandler.UpdateUser)
	admin.Post("/accounts/users/:id/action", adminHandler.UserAction)
	admin.Post("/accounts/users/import-csv", adminHandler.ImportUsersCSV)
	admin.Get("/accounts/companies", adminHandler.ListCompanies)
	admin.Post("/accounts/companies/:id/action", adminHandler.CompanyAction)
	admin.Get("/reports", adminHandler.ListReports)
	admin.Post("/reports/:id/action", adminHandler.ReportAction)
	admin.Get("/subjects", adminHandler.ListSubjects)
	admin.Get("/subjects/:id", adminHandler.GetSubject)
	admin.Post("/subjects/:id/action", adminHandler.SubjectAction)
	admin.Get("/pfe", adminHandler.ListAssignments)
	admin.Get("/pfe/:id", adminHandler.GetAssignment)
	admin.Get("/defenses", adminHandler.ListDefenses)
	admin.Post("/defenses", adminHandler.CreateDefense)
	admin.Get("/defenses/recommend-jury", adminHandler.RecommendJury)
	admin.Get("/defenses/:id", adminHandler.GetDefense)
	admin.Post("/defenses/:id/submit-grade", adminHandler.SubmitGrade)
	admin.Post("/defenses/:id/resolve-grade", adminHandler.ResolveGrade)
	admin.Post("/defenses/:id/confirm-jury", adminHandler.ConfirmJury)
	admin.Post("/defenses/:id/decline-jury", adminHandler.DeclineJury)
	admin.Get("/settings/deadlines", adminHandler.ListDeadlines)
	admin.Post("/settings/deadlines", adminHandler.UpdateDeadlines)
	admin.Get("/settings/specialities", adminHandler.ListSpecialities)
	admin.Post("/settings/specialities", adminHandler.CreateSpeciality)
	admin.Delete("/settings/specialities/:id", adminHandler.DeleteSpeciality)
	admin.Get("/settings/domains", adminHandler.ListDomains)
	admin.Post("/settings/domains", adminHandler.CreateDomain)
	admin.Delete("/settings/domains/:id", adminHandler.DeleteDomain)
	admin.Get("/settings/promotions", adminHandler.ListPromotions)
	admin.Post("/settings/promotions", adminHandler.CreatePromotion)
	admin.Delete("/settings/promotions/:id", adminHandler.DeletePromotion)
	admin.Get("/settings/academic-years", adminHandler.ListAcademicYears)
	admin.Post("/settings/academic-years", adminHandler.CreateAcademicYear)
	admin.Post("/settings/academic-years/:id/close", adminHandler.CloseAcademicYear)
	admin.Get("/statistics", adminHandler.Statistics)
	admin.Get("/audit-log", adminHandler.AuditLog)
	admin.Get("/exports/affectations", adminHandler.ExportAffectations)
	admin.Get("/exports/plannings", adminHandler.ExportPlannings)
	admin.Get("/exports/statistiques", adminHandler.ExportStatistics)

	teacher := api.Group("/teacher", middleware.AuthRequired(cfg), middleware.RequireRole("teacher", "admin"))
	teacher.Get("/dashboard", teacherHandler.Dashboard)
	teacher.Get("/proposed-subjects", teacherHandler.ListProposedSubjects)
	teacher.Post("/proposed-subjects", teacherHandler.CreateProposedSubject)
	teacher.Get("/proposed-subjects/:id", teacherHandler.GetProposedSubject)
	teacher.Patch("/proposed-subjects/:id", teacherHandler.UpdateProposedSubject)
	teacher.Get("/proposed-subjects/:id/candidats", teacherHandler.ListCandidats)
	teacher.Post("/proposed-subjects/:id/candidats", teacherHandler.AcceptCandidat)
	teacher.Get("/subjects-to-validate", teacherHandler.ListSubjectsToValidate)
	teacher.Get("/subjects-to-validate/:id", teacherHandler.GetSubjectToValidate)
	teacher.Post("/subjects-to-validate/:id", teacherHandler.ValidateSubject)
	teacher.Get("/supervised-pfes", teacherHandler.ListSupervisedPFEs)
	teacher.Get("/supervised-pfes/:id", teacherHandler.GetSupervisedPFE)
	teacher.Post("/supervised-pfes/:id/meetings", teacherHandler.AddMeeting)
	teacher.Post("/supervised-pfes/:id/evaluation", teacherHandler.SubmitEvaluation)
	teacher.Get("/jury-duties", teacherHandler.ListJuryDuties)
	teacher.Get("/jury-duties/:id", teacherHandler.GetJuryDuty)
	teacher.Post("/availability", teacherHandler.UpdateAvailability)
	teacher.Get("/notifications", teacherHandler.ListNotifications)

	student := api.Group("/student", middleware.AuthRequired(cfg), middleware.RequireRole("student"))
	student.Get("/dashboard", studentHandler.Dashboard)
	student.Get("/catalogue", studentHandler.ListCatalogue)
	student.Get("/catalogue/:id", studentHandler.GetCatalogueSubject)
	student.Get("/wishes", studentHandler.ListWishes)
	student.Post("/wishes", studentHandler.CreateWish)
	student.Delete("/wishes/:id", studentHandler.DeleteWish)
	student.Get("/my-pfe", studentHandler.GetMyPFE)
	student.Get("/my-pfe/meetings", studentHandler.ListMyMeetings)
	student.Post("/my-pfe/meetings", studentHandler.AddMyMeeting)
	student.Post("/my-pfe/memoire", studentHandler.SubmitMemoire)
	student.Get("/soutenance", studentHandler.GetSoutenance)
	student.Get("/notifications", studentHandler.ListNotifications)

	company := api.Group("/company", middleware.AuthRequired(cfg), middleware.RequireRole("company"))
	company.Get("/dashboard", companyHandler.Dashboard)
	company.Get("/subjects", companyHandler.ListSubjects)
	company.Post("/subjects", companyHandler.CreateSubject)
	company.Get("/subjects/:id", companyHandler.GetSubject)
	company.Patch("/subjects/:id", companyHandler.UpdateSubject)
	company.Get("/subjects/:id/candidats", companyHandler.ListCandidats)
	company.Post("/subjects/:id/candidats", companyHandler.AcceptCandidat)
	company.Get("/supervised-pfes", companyHandler.ListSupervisedPFEs)
	company.Get("/supervised-pfes/:id", companyHandler.GetSupervisedPFE)
	company.Post("/supervised-pfes/:id/meetings", companyHandler.AddMeeting)
	company.Post("/supervised-pfes/:id/evaluation", companyHandler.SubmitEvaluation)
	company.Get("/reports", companyHandler.ListReports)
	company.Post("/reports", companyHandler.CreateReport)
	company.Get("/notifications", companyHandler.ListNotifications)

	notifs := api.Group("/notifications", middleware.AuthRequired(cfg))
	notifs.Get("/", func(c fiber.Ctx) error {
		role := middleware.GetRole(c)
		switch role {
		case "admin", "teacher":
			return teacherHandler.ListNotifications(c)
		case "student":
			return studentHandler.ListNotifications(c)
		case "company":
			return companyHandler.ListNotifications(c)
		default:
			return c.Status(403).JSON(map[string]any{"success": false, "error": "Rôle inconnu"})
		}
	})
	notifs.Post("/:id/read", func(c fiber.Ctx) error {
		return c.JSON(map[string]any{"success": true, "data": map[string]string{"message": "Notification marquée comme lue"}})
	})
	notifs.Post("/read-all", func(c fiber.Ctx) error {
		return c.JSON(map[string]any{"success": true, "data": map[string]string{"message": "Toutes les notifications lues"}})
	})

	adminToken := generateToken(cfg.JWTSecret, "seed-admin-001", "admin")

	return &TestHelper{
		App:   app,
		DB:    db,
		Cfg:   cfg,
		Admin: adminToken,
	}
}

func (h *TestHelper) Close() {
	h.DB.Close()
}

func (h *TestHelper) AuthHeader(profileID, role string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + generateToken(h.Cfg.JWTSecret, profileID, role),
	}
}

func (h *TestHelper) AuthToken(profileID, role string) string {
	return generateToken(h.Cfg.JWTSecret, profileID, role)
}

func generateToken(secret, profileID, role string) string {
	claims := jwt.MapClaims{
		"sub":  profileID,
		"role": role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func ParseResponse(resp *http.Response) (map[string]any, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture body: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("erreur unmarshal JSON: %w (body: %s)", err, string(body))
	}
	return result, nil
}

func MustParseResponse(resp *http.Response) map[string]any {
	result, err := ParseResponse(resp)
	if err != nil {
		panic(err)
	}
	return result
}

// AssertSuccess vérifie que la réponse a success=true.
func AssertSuccess(t TestingT, result map[string]any) {
	t.Helper()
	success, ok := result["success"].(bool)
	if !ok || !success {
		t.Fatalf("❌ Échec: success=false, response: %+v", result)
	}
}

// AssertError vérifie que la réponse a success=false.
func AssertError(t TestingT, result map[string]any) {
	t.Helper()
	success, ok := result["success"].(bool)
	if !ok || success {
		t.Fatalf("❌ Échec: attendu erreur mais success=true, response: %+v", result)
	}
}

// AssertErrorContains vérifie que l'erreur contient le texte attendu.
func AssertErrorContains(t TestingT, result map[string]any, expected string) {
	t.Helper()
	AssertError(t, result)
	errMsg, ok := result["error"].(string)
	if !ok {
		t.Fatalf("❌ Échec: champ error manquant ou non string, response: %+v", result)
	}
	if !strings.Contains(errMsg, expected) {
		t.Fatalf("❌ Échec: erreur attendue contenant %q, obtenu %q", expected, errMsg)
	}
}

// TestingT est une interface réduite pour compatibilité avec *testing.T et *testing.B.
type TestingT interface {
	Fatalf(format string, args ...any)
	Helper()
}

// ---- helpers de migration et seed ----

func runTestMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS profiles (
			id TEXT PRIMARY KEY,
			role TEXT NOT NULL CHECK(role IN ('admin','teacher','student','company')),
			full_name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			avatar_url TEXT DEFAULT '',
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS domains (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS specialities (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			code TEXT NOT NULL UNIQUE,
			year_type TEXT NOT NULL CHECK(year_type IN ('licence','master')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS academic_years (
			id TEXT PRIMARY KEY,
			label TEXT NOT NULL UNIQUE,
			status TEXT NOT NULL CHECK(status IN ('active','cloturee')),
			submission_open_at DATETIME,
			submission_close_at DATETIME,
			max_wishes INTEGER DEFAULT 5,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS promotions (
			id TEXT PRIMARY KEY,
			label TEXT NOT NULL,
			academic_year_id TEXT NOT NULL REFERENCES academic_years(id),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS teachers (
			id TEXT PRIMARY KEY,
			profile_id TEXT UNIQUE NOT NULL REFERENCES profiles(id),
			grade TEXT NOT NULL CHECK(grade IN ('assistant','mab','maa','mcb','mca','professeur')),
			department TEXT,
			availability_status TEXT DEFAULT 'disponible' CHECK(availability_status IN ('disponible','indisponible','indisponible_jusqu_au')),
			unavailable_until DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS teacher_domains (
			teacher_id TEXT NOT NULL REFERENCES teachers(id),
			domain_id TEXT NOT NULL REFERENCES domains(id),
			PRIMARY KEY (teacher_id, domain_id)
		)`,
		`CREATE TABLE IF NOT EXISTS students (
			id TEXT PRIMARY KEY,
			profile_id TEXT UNIQUE NOT NULL REFERENCES profiles(id),
			student_number TEXT UNIQUE NOT NULL,
			speciality_id TEXT REFERENCES specialities(id),
			level TEXT,
			promotion_id TEXT REFERENCES promotions(id),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS companies (
			id TEXT PRIMARY KEY,
			profile_id TEXT UNIQUE NOT NULL REFERENCES profiles(id),
			company_name TEXT NOT NULL,
			sector TEXT,
			description TEXT,
			logo_url TEXT DEFAULT '',
			contact_email TEXT,
			contact_phone TEXT,
			website TEXT,
			is_verified BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS pfe_subjects (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			group_type TEXT DEFAULT 'binome' CHECK(group_type IN ('monome','binome','trinome')),
			proposer_id TEXT NOT NULL REFERENCES profiles(id),
			proposer_role TEXT NOT NULL CHECK(proposer_role IN ('teacher','company')),
			company_id TEXT REFERENCES companies(id),
			academic_year_id TEXT NOT NULL REFERENCES academic_years(id),
			validator1_id TEXT REFERENCES teachers(id),
			validator2_id TEXT REFERENCES teachers(id),
			validator1_decision TEXT CHECK(validator1_decision IN ('valide','accepte_sous_reserve','refuse')),
			validator2_decision TEXT CHECK(validator2_decision IN ('valide','accepte_sous_reserve','refuse')),
			validator1_comment TEXT,
			validator2_comment TEXT,
			status TEXT DEFAULT 'en_attente' CHECK(status IN ('en_attente','valide','accepte_sous_reserve','refuse','expire')),
			co_supervisor_id TEXT REFERENCES teachers(id),
			pre_assigned_student_ids TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS subject_domains (
			subject_id TEXT NOT NULL REFERENCES pfe_subjects(id),
			domain_id TEXT NOT NULL REFERENCES domains(id),
			PRIMARY KEY (subject_id, domain_id)
		)`,
		`CREATE TABLE IF NOT EXISTS wishes (
			id TEXT PRIMARY KEY,
			student_id TEXT NOT NULL REFERENCES students(id),
			subject_id TEXT NOT NULL REFERENCES pfe_subjects(id),
			academic_year_id TEXT NOT NULL REFERENCES academic_years(id),
			status TEXT DEFAULT 'en_attente' CHECK(status IN ('en_attente','accepte','refuse')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS pfe_assignments (
			id TEXT PRIMARY KEY,
			pfe_code TEXT UNIQUE NOT NULL,
			subject_id TEXT NOT NULL REFERENCES pfe_subjects(id),
			academic_year_id TEXT NOT NULL REFERENCES academic_years(id),
			student_id TEXT NOT NULL REFERENCES students(id),
			student2_id TEXT REFERENCES students(id),
			student3_id TEXT REFERENCES students(id),
			supervisor_id TEXT NOT NULL REFERENCES teachers(id),
			co_supervisor_id TEXT REFERENCES teachers(id),
			memoire_url TEXT,
			status TEXT DEFAULT 'en_cours' CHECK(status IN ('en_cours','memoire_soumis','soutenance_planifiee','valide','refuse')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS pfe_progress_reports (
			id TEXT PRIMARY KEY,
			assignment_id TEXT NOT NULL REFERENCES pfe_assignments(id),
			meeting_date DATETIME NOT NULL,
			duration INTEGER NOT NULL,
			meeting_type TEXT NOT NULL CHECK(meeting_type IN ('presentiel','visio')),
			topics TEXT,
			status TEXT DEFAULT 'en_cours' CHECK(status IN ('en_cours','termine')),
			observation TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS defense_juries (
			id TEXT PRIMARY KEY,
			assignment_id TEXT NOT NULL REFERENCES pfe_assignments(id),
			president_id TEXT NOT NULL REFERENCES teachers(id),
			member_id TEXT NOT NULL REFERENCES teachers(id),
			president_confirmed BOOLEAN DEFAULT 0,
			member_confirmed BOOLEAN DEFAULT 0,
			president_wants_printed BOOLEAN DEFAULT 0,
			member_wants_printed BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS defenses (
			id TEXT PRIMARY KEY,
			assignment_id TEXT NOT NULL REFERENCES pfe_assignments(id),
			jury_id TEXT NOT NULL REFERENCES defense_juries(id),
			scheduled_at DATETIME,
			room TEXT,
			defense_deadline DATETIME,
			status TEXT DEFAULT 'scheduled' CHECK(status IN ('scheduled','done','postponed')),
			result TEXT CHECK(result IN ('admitted','corrections_required','not_admitted')),
			final_grade REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS jury_grades (
			id TEXT PRIMARY KEY,
			defense_id TEXT NOT NULL REFERENCES defenses(id),
			jury_member_id TEXT NOT NULL REFERENCES teachers(id),
			criterion1 REAL,
			criterion2 REAL,
			criterion3 REAL,
			criterion4 REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS supervisor_evaluations (
			id TEXT PRIMARY KEY,
			pfe_assignment_id TEXT UNIQUE NOT NULL REFERENCES pfe_assignments(id),
			evaluator_id TEXT NOT NULL REFERENCES teachers(id),
			criterion5 REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS company_reports (
			id TEXT PRIMARY KEY,
			company_id TEXT NOT NULL REFERENCES companies(id),
			submitted_by TEXT NOT NULL,
			correction_type TEXT NOT NULL,
			description TEXT,
			requested_value TEXT,
			status TEXT DEFAULT 'en_attente' CHECK(status IN ('en_attente','resolu','rejete')),
			resolved_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS notifications (
			id TEXT PRIMARY KEY,
			recipient_id TEXT NOT NULL REFERENCES profiles(id),
			type TEXT NOT NULL,
			payload TEXT DEFAULT '{}',
			read_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id TEXT PRIMARY KEY,
			actor_id TEXT NOT NULL REFERENCES profiles(id),
			action TEXT NOT NULL,
			entity TEXT NOT NULL,
			entity_id TEXT,
			metadata TEXT DEFAULT '{}',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("migration test échouée: %w\nSQL: %s", err, m[:50])
		}
	}
	return nil
}

func runTestSeed(db *sql.DB) error {
	seeds := []string{
		`INSERT OR IGNORE INTO domains (id, name) VALUES ('seed-domain-ia', 'Intelligence Artificielle')`,
		`INSERT OR IGNORE INTO domains (id, name) VALUES ('seed-domain-web', 'Développement Web')`,
		`INSERT OR IGNORE INTO domains (id, name) VALUES ('seed-domain-reseau', 'Réseaux et Sécurité')`,
		`INSERT OR IGNORE INTO domains (id, name) VALUES ('seed-domain-data', 'Data Science')`,
		`INSERT OR IGNORE INTO domains (id, name) VALUES ('seed-domain-emb', 'Systèmes Embarqués')`,
		`INSERT OR IGNORE INTO domains (id, name) VALUES ('seed-domain-mobile', 'Développement Mobile')`,
		`INSERT OR IGNORE INTO domains (id, name) VALUES ('seed-domain-cloud', 'Cloud Computing')`,
		`INSERT OR IGNORE INTO domains (id, name) VALUES ('seed-domain-biologie', 'Bio-Informatique')`,

		`INSERT OR IGNORE INTO specialities (id, name, code, year_type) VALUES ('seed-spec-isil', 'ISIL', 'ISIL', 'master')`,
		`INSERT OR IGNORE INTO specialities (id, name, code, year_type) VALUES ('seed-spec-chim', 'Chimie', 'CHIM', 'licence')`,
		`INSERT OR IGNORE INTO specialities (id, name, code, year_type) VALUES ('seed-spec-elec', 'Électrotechnique', 'ELEC', 'master')`,

		`INSERT OR IGNORE INTO academic_years (id, label, status) VALUES ('seed-ay-2324', '2023-2024', 'cloturee')`,
		`INSERT OR IGNORE INTO academic_years (id, label, status, submission_open_at, submission_close_at, max_wishes) 
		 VALUES ('seed-ay-2425', '2024-2025', 'active', datetime('now', '-30 days'), datetime('now', '+30 days'), 5)`,

		`INSERT OR IGNORE INTO promotions (id, label, academic_year_id) VALUES ('seed-promo-isil', 'ISIL 2024-2025', 'seed-ay-2425')`,
		`INSERT OR IGNORE INTO promotions (id, label, academic_year_id) VALUES ('seed-promo-chim', 'CHIM 2024-2025', 'seed-ay-2425')`,

		`INSERT OR IGNORE INTO profiles (id, role, full_name, email, is_active) VALUES ('seed-admin-001', 'admin', 'Admin Test', 'admin@test.dz', 1)`,
		`INSERT OR IGNORE INTO profiles (id, role, full_name, email, is_active) VALUES ('seed-teacher-isil-001', 'teacher', 'Dr. ISIL Teacher', 'teacher.isil@test.dz', 1)`,
		`INSERT OR IGNORE INTO profiles (id, role, full_name, email, is_active) VALUES ('seed-teacher-isil-002', 'teacher', 'Pr. ISIL Validator', 'teacher.isil2@test.dz', 1)`,
		`INSERT OR IGNORE INTO profiles (id, role, full_name, email, is_active) VALUES ('seed-teacher-chim-001', 'teacher', 'Dr. CHIM Teacher', 'teacher.chim@test.dz', 1)`,
		`INSERT OR IGNORE INTO profiles (id, role, full_name, email, is_active) VALUES ('seed-student-isil-001', 'student', 'Étudiant ISIL 1', 'student.isil1@test.dz', 1)`,
		`INSERT OR IGNORE INTO profiles (id, role, full_name, email, is_active) VALUES ('seed-student-isil-002', 'student', 'Étudiant ISIL 2', 'student.isil2@test.dz', 1)`,
		`INSERT OR IGNORE INTO profiles (id, role, full_name, email, is_active) VALUES ('seed-student-chim-001', 'student', 'Étudiant CHIM 1', 'student.chim1@test.dz', 1)`,
		`INSERT OR IGNORE INTO profiles (id, role, full_name, email, is_active) VALUES ('seed-student-isil-004', 'student', 'Étudiant ISIL 4', 'student.isil4@test.dz', 1)`,
		`INSERT OR IGNORE INTO profiles (id, role, full_name, email, is_active) VALUES ('seed-company-001', 'company', 'TechCorp Algérie', 'contact@techcorp.dz', 1)`,

		`INSERT OR IGNORE INTO teachers (id, profile_id, grade, department, availability_status) 
		 VALUES ('seed-teacher-isil-001', 'seed-teacher-isil-001', 'mca', 'Informatique', 'disponible')`,
		`INSERT OR IGNORE INTO teachers (id, profile_id, grade, department, availability_status) 
		 VALUES ('seed-teacher-isil-002', 'seed-teacher-isil-002', 'professeur', 'Informatique', 'disponible')`,
		`INSERT OR IGNORE INTO teachers (id, profile_id, grade, department, availability_status) 
		 VALUES ('seed-teacher-chim-001', 'seed-teacher-chim-001', 'mcb', 'Chimie', 'disponible')`,

		`INSERT OR IGNORE INTO teacher_domains (teacher_id, domain_id) VALUES ('seed-teacher-isil-001', 'seed-domain-ia')`,
		`INSERT OR IGNORE INTO teacher_domains (teacher_id, domain_id) VALUES ('seed-teacher-isil-001', 'seed-domain-web')`,
		`INSERT OR IGNORE INTO teacher_domains (teacher_id, domain_id) VALUES ('seed-teacher-isil-002', 'seed-domain-ia')`,
		`INSERT OR IGNORE INTO teacher_domains (teacher_id, domain_id) VALUES ('seed-teacher-chim-001', 'seed-domain-data')`,

		`INSERT OR IGNORE INTO students (id, profile_id, student_number, speciality_id, level, promotion_id)
		 VALUES ('seed-student-isil-001', 'seed-student-isil-001', '2024001', 'seed-spec-isil', 'M2', 'seed-promo-isil')`,
		`INSERT OR IGNORE INTO students (id, profile_id, student_number, speciality_id, level, promotion_id)
		 VALUES ('seed-student-isil-002', 'seed-student-isil-002', '2024002', 'seed-spec-isil', 'M2', 'seed-promo-isil')`,
		`INSERT OR IGNORE INTO students (id, profile_id, student_number, speciality_id, level, promotion_id)
		 VALUES ('seed-student-chim-001', 'seed-student-chim-001', '2024003', 'seed-spec-chim', 'L3', 'seed-promo-chim')`,
		`INSERT OR IGNORE INTO students (id, profile_id, student_number, speciality_id, level, promotion_id)
		 VALUES ('seed-student-isil-004', 'seed-student-isil-004', '2024004', 'seed-spec-isil', 'M2', 'seed-promo-isil')`,

		`INSERT OR IGNORE INTO companies (id, profile_id, company_name, sector, description, is_verified) 
		 VALUES ('seed-company-001', 'seed-company-001', 'TechCorp Algérie', 'Technologie', 'Entreprise tech', 1)`,

		`INSERT OR IGNORE INTO pfe_subjects (id, title, description, group_type, proposer_id, proposer_role, academic_year_id, status)
		 VALUES ('seed-subject-001', 'IA pour la santé', 'Sujet IA santé', 'binome', 'seed-teacher-isil-001', 'teacher', 'seed-ay-2425', 'en_attente')`,
		`INSERT OR IGNORE INTO pfe_subjects (id, title, description, group_type, proposer_id, proposer_role, academic_year_id, status)
		 VALUES ('seed-subject-002', 'Web App Sécurité', 'Sujet web sécurité', 'monome', 'seed-teacher-isil-001', 'teacher', 'seed-ay-2425', 'valide')`,
		`INSERT OR IGNORE INTO pfe_subjects (id, title, description, group_type, proposer_id, proposer_role, academic_year_id, status,
		    validator1_id, validator2_id, validator1_decision, validator2_decision)
		 VALUES ('seed-subject-003', 'Cloud Computing Avancé', 'Sujet cloud', 'binome', 'seed-teacher-isil-001', 'teacher', 'seed-ay-2425', 'valide',
		    'seed-teacher-isil-002', 'seed-teacher-chim-001', 'valide', 'valide')`,
		`INSERT OR IGNORE INTO pfe_subjects (id, title, description, group_type, proposer_id, proposer_role, academic_year_id, status,
		    validator1_id, validator1_decision)
		 VALUES ('seed-subject-004', 'Data Mining', 'Sujet data', 'binome', 'seed-teacher-chim-001', 'teacher', 'seed-ay-2425', 'accepte_sous_reserve',
		    'seed-teacher-isil-002', 'accepte_sous_reserve')`,
		`INSERT OR IGNORE INTO pfe_subjects (id, title, description, group_type, proposer_id, proposer_role, company_id, academic_year_id, status)
		 VALUES ('seed-subject-005', 'IoT Industriel', 'Sujet IoT', 'trinome', 'seed-company-001', 'company', 'seed-company-001', 'seed-ay-2425', 'valide')`,
		`INSERT OR IGNORE INTO pfe_subjects (id, title, description, group_type, proposer_id, proposer_role, academic_year_id, status)
		 VALUES ('seed-subject-006', 'Blockchain', 'Sujet blockchain', 'monome', 'seed-teacher-isil-001', 'teacher', 'seed-ay-2425', 'refuse')`,

		`INSERT OR IGNORE INTO subject_domains (subject_id, domain_id) VALUES ('seed-subject-001', 'seed-domain-ia')`,
		`INSERT OR IGNORE INTO subject_domains (subject_id, domain_id) VALUES ('seed-subject-002', 'seed-domain-web')`,
		`INSERT OR IGNORE INTO subject_domains (subject_id, domain_id) VALUES ('seed-subject-003', 'seed-domain-cloud')`,

		`INSERT OR IGNORE INTO wishes (id, student_id, subject_id, academic_year_id, status)
		 VALUES ('seed-wish-001', 'seed-student-isil-001', 'seed-subject-002', 'seed-ay-2425', 'en_attente')`,
		`INSERT OR IGNORE INTO wishes (id, student_id, subject_id, academic_year_id, status)
		 VALUES ('seed-wish-002', 'seed-student-isil-001', 'seed-subject-003', 'seed-ay-2425', 'en_attente')`,
		`INSERT OR IGNORE INTO wishes (id, student_id, subject_id, academic_year_id, status)
		 VALUES ('seed-wish-003', 'seed-student-isil-002', 'seed-subject-003', 'seed-ay-2425', 'accepte')`,

`INSERT OR IGNORE INTO pfe_assignments (id, pfe_code, subject_id, academic_year_id, student_id, student2_id, supervisor_id, status)
		 VALUES ('seed-assignment-001', 'PFE-ISIL-2025-001', 'seed-subject-003', 'seed-ay-2425', 'seed-student-isil-001', 'seed-student-isil-002', 'seed-teacher-isil-001', 'en_cours')`,
		`INSERT OR IGNORE INTO pfe_assignments (id, pfe_code, subject_id, academic_year_id, student_id, supervisor_id, status)
		 VALUES ('seed-assignment-002', 'PFE-ISIL-2025-002', 'seed-subject-005', 'seed-ay-2425', 'seed-student-isil-004', 'seed-teacher-isil-001', 'en_cours')`,

		`INSERT OR IGNORE INTO pfe_progress_reports (id, assignment_id, meeting_date, duration, meeting_type, topics, status)
		 VALUES ('seed-meeting-001', 'seed-assignment-001', datetime('now', '-14 days'), 60, 'presentiel', 'Introduction, planification', 'termine')`,
		`INSERT OR IGNORE INTO pfe_progress_reports (id, assignment_id, meeting_date, duration, meeting_type, topics, status)
		 VALUES ('seed-meeting-002', 'seed-assignment-001', datetime('now', '-7 days'), 45, 'visio', 'État d''avancement', 'termine')`,

		`INSERT OR IGNORE INTO supervisor_evaluations (id, pfe_assignment_id, evaluator_id, criterion5)
		 VALUES ('seed-sup-eval-001', 'seed-assignment-001', 'seed-teacher-isil-001', 3.5)`,

		`INSERT OR IGNORE INTO defense_juries (id, assignment_id, president_id, member_id, president_confirmed, member_confirmed)
		 VALUES ('seed-jury-001', 'seed-assignment-001', 'seed-teacher-isil-002', 'seed-teacher-chim-001', 1, 1)`,
		`INSERT OR IGNORE INTO defenses (id, assignment_id, jury_id, scheduled_at, room, status)
		 VALUES ('seed-defense-001', 'seed-assignment-001', 'seed-jury-001', datetime('now', '+14 days'), 'Salle A', 'scheduled')`,

		`INSERT OR IGNORE INTO notifications (id, recipient_id, type, payload)
		 VALUES ('seed-notif-001', 'seed-admin-001', 'nouveau_sujet', '{"subject_id":"seed-subject-001"}')`,
		`INSERT OR IGNORE INTO notifications (id, recipient_id, type, payload)
		 VALUES ('seed-notif-002', 'seed-teacher-isil-001', 'sujet_valide', '{"subject_id":"seed-subject-003"}')`,
	}

	for _, s := range seeds {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("seed test échoué: %w\nSQL: %s", err, s[:50])
		}
	}
	return nil
}

// CleanDB vide toutes les tables entre les tests.
func CleanDB(db *sql.DB) error {
	tables := []string{
		"audit_logs", "notifications", "company_reports", "supervisor_evaluations",
		"jury_grades", "defenses", "defense_juries", "pfe_progress_reports",
		"pfe_assignments", "wishes", "subject_domains", "pfe_subjects",
		"companies", "students", "teacher_domains", "teachers",
		"promotions", "academic_years", "specialities", "domains", "profiles",
	}
	for _, t := range tables {
		if _, err := db.Exec("DELETE FROM " + t); err != nil {
			return fmt.Errorf("erreur nettoyage %s: %w", t, err)
		}
	}
	return nil
}
