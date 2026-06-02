package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"pfe-backend/internal/config"
	"pfe-backend/internal/entity"
	"pfe-backend/internal/handler"
	"pfe-backend/internal/repository"
	"pfe-backend/internal/service"
	"pfe-backend/internal/shared/middleware"
	"pfe-backend/internal/shared/notify"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"

	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Erreur connexion DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Erreur ping DB: %v", err)
	}

	db.SetMaxOpenConns(1)

	if err := runMigrations(db); err != nil {
		log.Fatalf("Erreur migration: %v", err)
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
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	}))


	profileRepo := repository.NewProfileRepository(db)
	teacherRepo := repository.NewTeacherRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	departmentRepo := repository.NewDepartmentRepository(db)
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


	authService := service.NewAuthService(profileRepo, teacherRepo, studentRepo, companyRepo, cfg)


	notifier := notify.New(notificationRepo, profileRepo, cfg.ResendAPIKey)

	adminService := service.NewAdminService(
		profileRepo, teacherRepo, studentRepo, companyRepo, departmentRepo,
		domainRepo, specialityRepo, promotionRepo, academicYearRepo,
		pfeSubjectRepo, wishRepo, pfeAssignmentRepo,
		progressRepo, defenseJuryRepo, defenseRepo,
		juryGradeRepo, supEvalRepo, companyReportRepo,
		notificationRepo, auditLogRepo, notifier, "./uploads",
	)

	teacherService := service.NewTeacherService(
		profileRepo, teacherRepo, studentRepo, companyRepo, specialityRepo,
		pfeSubjectRepo, wishRepo, pfeAssignmentRepo, progressRepo,
		defenseJuryRepo, defenseRepo, juryGradeRepo, supEvalRepo, notificationRepo, academicYearRepo, notifier,
	)

	studentService := service.NewStudentService(
		profileRepo, studentRepo, teacherRepo, companyRepo, specialityRepo,
		pfeSubjectRepo, wishRepo, pfeAssignmentRepo, progressRepo,
		defenseRepo, defenseJuryRepo, supEvalRepo, notificationRepo, academicYearRepo,
	)

	companyService := service.NewCompanyService(
		profileRepo, companyRepo, teacherRepo, studentRepo, specialityRepo,
		pfeSubjectRepo, wishRepo, pfeAssignmentRepo, progressRepo,
		supEvalRepo, companyReportRepo, notificationRepo, academicYearRepo, notifier,
	)


	authHandler := handler.NewAuthHandler(authService, cfg, notifier)
	adminHandler := handler.NewAdminHandler(adminService, notifier)
	teacherHandler := handler.NewTeacherHandler(teacherService, notifier)
	studentHandler := handler.NewStudentHandler(studentService, notifier)
	companyHandler := handler.NewCompanyHandler(companyService, notifier)
	uploadHandler := handler.NewUploadHandler(profileRepo, companyRepo, "./uploads")

	api := app.Group("/api")


	api.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(map[string]any{
			"success": true,
			"data": map[string]any{
				"status": "ok",
			},
		})
	})


	authLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Trop de tentatives. Veuillez réessayer dans une minute.",
			})
		},
	})

	api.Post("/auth/dev-login", authLimiter, authHandler.DevLogin)
	api.Get("/auth/companies", authHandler.ListVerifiedCompanies)
	api.Post("/auth/register-company", authLimiter, authHandler.RegisterCompany)

	api.Get("/accounts/users", adminHandler.ListUsers)


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
	ref.Get("/departments", func(c fiber.Ctx) error {
		departments, err := departmentRepo.FindAll()
		if err != nil {
			return c.Status(500).JSON(map[string]any{"success": false, "error": "Erreur serveur"})
		}
		if departments == nil {
			departments = []*entity.Department{}
		}
		return c.JSON(map[string]any{"success": true, "data": departments})
	})


	auth := api.Group("/auth")
	auth.Use(middleware.AuthRequired(cfg))
	auth.Get("/me", authHandler.Me)
	auth.Post("/logout", authHandler.Logout)



	app.Use("/uploads", static.New("./uploads"))


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
	admin.Post("/accounts/users/:id/avatar", adminHandler.UpdateUserAvatar)
	admin.Post("/accounts/users/:id/action", adminHandler.UserAction)
	admin.Post("/accounts/users/enseignant", adminHandler.CreateTeacher)
	admin.Post("/accounts/users/etudiant", adminHandler.CreateStudent)
	admin.Post("/accounts/users/import-csv", adminHandler.ImportUsersCSV)
	admin.Get("/accounts/companies", adminHandler.ListCompanies)
	admin.Post("/accounts/companies/:id/action", adminHandler.CompanyAction)
	admin.Get("/reports", adminHandler.ListReports)
	admin.Post("/reports/:id/action", adminHandler.ReportAction)
	admin.Get("/subjects", adminHandler.ListSubjects)
	admin.Get("/subjects/:id", adminHandler.GetSubject)
	admin.Post("/subjects/:id/action", adminHandler.SubjectAction)
	admin.Get("/pfe", adminHandler.ListAssignments)
	admin.Get("/pfe/recommend-co-supervisor", adminHandler.RecommendCoSupervisorHandler)
	admin.Get("/pfe/:id", adminHandler.GetAssignment)
	admin.Post("/pfe/:id/action", adminHandler.AssignmentAction)
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
	admin.Get("/settings/departments", adminHandler.ListDepartments)
	admin.Post("/settings/departments", adminHandler.CreateDepartment)
	admin.Delete("/settings/departments/:id", adminHandler.DeleteDepartment)
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
	teacher.Post("/proposed-subjects/:id/resubmit", teacherHandler.ResubmitSubject)
	teacher.Get("/proposed-subjects/:id/candidats", teacherHandler.ListCandidats)
	teacher.Post("/proposed-subjects/:id/candidats", teacherHandler.AcceptCandidat)
	teacher.Get("/subjects-to-validate", teacherHandler.ListSubjectsToValidate)
	teacher.Get("/subjects-to-validate/:id", teacherHandler.GetSubjectToValidate)
	teacher.Post("/subjects-to-validate/:id", teacherHandler.ValidateSubject)
	teacher.Get("/supervised-pfes", teacherHandler.ListSupervisedPFEs)
	teacher.Get("/supervised-pfes/:id", teacherHandler.GetSupervisedPFE)
	teacher.Get("/supervised-pfes/:id/meetings", teacherHandler.ListMeetings)
	teacher.Post("/supervised-pfes/:id/meetings", teacherHandler.AddMeeting)
	teacher.Get("/supervised-pfes/:id/evaluation", teacherHandler.GetEvaluation)
	teacher.Post("/supervised-pfes/:id/evaluation", teacherHandler.SubmitEvaluation)
	teacher.Get("/jury-duties", teacherHandler.ListJuryDuties)
	teacher.Get("/jury-duties/:id", teacherHandler.GetJuryDuty)
	teacher.Get("/jury-duties/:id/grade-context", teacherHandler.GetGradeContext)
	teacher.Post("/jury-duties/:id/grade", teacherHandler.SubmitJuryGrade)
	teacher.Post("/jury-duties/:id/final-grade", teacherHandler.SubmitFinalGrade)
	teacher.Post("/availability", teacherHandler.UpdateAvailability)
	teacher.Get("/notifications", teacherHandler.ListNotifications)


	student := api.Group("/student", middleware.AuthRequired(cfg), middleware.RequireRole("student"))
	student.Get("/settings", studentHandler.GetSettings)
	student.Get("/dashboard", studentHandler.Dashboard)
	student.Get("/catalogue", studentHandler.ListCatalogue)
	student.Get("/catalogue/:id", studentHandler.GetCatalogueSubject)
	student.Get("/wishes", studentHandler.ListWishes)
	student.Post("/wishes", studentHandler.CreateWish)
	student.Delete("/wishes/:id", studentHandler.DeleteWish)
	student.Get("/my-pfe", studentHandler.GetMyPFE)
	student.Get("/my-pfe/meetings", studentHandler.ListMyMeetings)
	student.Post("/my-pfe/meetings", studentHandler.AddMyMeeting)
	student.Patch("/my-pfe/meetings/:id", studentHandler.UpdateMyMeeting)
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
	company.Get("/supervised-pfes/:id/meetings", companyHandler.ListMeetings)
	company.Post("/supervised-pfes/:id/meetings", companyHandler.AddMeeting)
	company.Get("/supervised-pfes/:id/evaluation", companyHandler.GetEvaluation)
	company.Post("/supervised-pfes/:id/evaluation", companyHandler.SubmitEvaluation)
	company.Get("/reports", companyHandler.ListReports)
	company.Post("/reports", companyHandler.CreateReport)
	company.Get("/notifications", companyHandler.ListNotifications)


	notifs := api.Group("/notifications", middleware.AuthRequired(cfg))
	notifs.Get("/", func(c fiber.Ctx) error {
		profileID := middleware.GetProfileID(c)
		list, err := notificationRepo.FindByRecipient(profileID)
		if err != nil {
			return c.Status(500).JSON(map[string]any{"success": false, "error": err.Error()})
		}
		if list == nil {
			list = []*entity.Notification{}
		}
		return c.JSON(map[string]any{"success": true, "data": list})
	})
	notifs.Get("/unread-count", func(c fiber.Ctx) error {
		profileID := middleware.GetProfileID(c)
		count, err := notificationRepo.CountUnread(profileID)
		if err != nil {
			return c.Status(500).JSON(map[string]any{"success": false, "error": err.Error()})
		}
		return c.JSON(map[string]any{"success": true, "data": count})
	})
	notifs.Post("/:id/read", func(c fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(400).JSON(map[string]any{"success": false, "error": "ID invalide"})
		}
		if err := notificationRepo.MarkAsRead(id); err != nil {
			return c.Status(500).JSON(map[string]any{"success": false, "error": err.Error()})
		}
		return c.JSON(map[string]any{"success": true, "data": map[string]string{"message": "Notification marquée comme lue"}})
	})
	notifs.Post("/read-all", func(c fiber.Ctx) error {
		profileID := middleware.GetProfileID(c)
		if err := notificationRepo.MarkAllAsRead(profileID); err != nil {
			return c.Status(500).JSON(map[string]any{"success": false, "error": err.Error()})
		}
		return c.JSON(map[string]any{"success": true, "data": map[string]string{"message": "Toutes les notifications lues"}})
	})

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	go func() {
		fmt.Printf("Serveur démarré sur le port %s (env: %s)\n", port, cfg.Env)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Erreur serveur: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Arrêt du serveur...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Erreur shutdown: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			role TEXT NOT NULL CHECK(role IN ('admin','teacher','student','company')),
			full_name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			avatar_url TEXT DEFAULT '',
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS domains (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS departments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS specialities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			code TEXT NOT NULL UNIQUE,
			year_type TEXT NOT NULL CHECK(year_type IN ('licence','master')),
			department_id INTEGER REFERENCES departments(id),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS academic_years (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			label TEXT NOT NULL UNIQUE,
			status TEXT NOT NULL CHECK(status IN ('active','cloturee')),
			submission_open_at DATETIME,
			submission_close_at DATETIME,
			max_wishes INTEGER DEFAULT 5,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS promotions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			label TEXT NOT NULL,
			academic_year_id INTEGER NOT NULL REFERENCES academic_years(id),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS teachers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			profile_id INTEGER UNIQUE NOT NULL REFERENCES profiles(id),
			grade TEXT NOT NULL CHECK(grade IN ('assistant','mab','maa','mcb','mca','professeur')),
			department_id INTEGER REFERENCES departments(id),
			availability_status TEXT DEFAULT 'disponible' CHECK(availability_status IN ('disponible','indisponible','indisponible_jusqu_au')),
			unavailable_until DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS teacher_domains (
			teacher_id INTEGER NOT NULL REFERENCES teachers(id),
			domain_id INTEGER NOT NULL REFERENCES domains(id),
			PRIMARY KEY (teacher_id, domain_id)
		)`,
		`CREATE TABLE IF NOT EXISTS students (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			profile_id INTEGER UNIQUE NOT NULL REFERENCES profiles(id),
			student_number TEXT UNIQUE NOT NULL,
			speciality_id INTEGER REFERENCES specialities(id),
			level TEXT,
			promotion_id INTEGER REFERENCES promotions(id),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS companies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			profile_id INTEGER UNIQUE NOT NULL REFERENCES profiles(id),
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
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			group_type TEXT DEFAULT 'binome' CHECK(group_type IN ('monome','binome','trinome')),
			proposer_id INTEGER NOT NULL REFERENCES profiles(id),
			proposer_role TEXT NOT NULL CHECK(proposer_role IN ('teacher','company')),
			company_id INTEGER REFERENCES companies(id),
			academic_year_id INTEGER NOT NULL REFERENCES academic_years(id),
			validator1_id INTEGER REFERENCES teachers(id),
			validator2_id INTEGER REFERENCES teachers(id),
			validator1_decision TEXT CHECK(validator1_decision IN ('valide','accepte_sous_reserve','refuse')),
			validator2_decision TEXT CHECK(validator2_decision IN ('valide','accepte_sous_reserve','refuse')),
			validator1_comment TEXT,
			validator2_comment TEXT,
			status TEXT DEFAULT 'en_attente' CHECK(status IN ('en_attente','valide','accepte_sous_reserve','refuse','expire')),
			co_supervisor_id INTEGER REFERENCES teachers(id),
			pre_assigned_student_ids TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS subject_domains (
			subject_id INTEGER NOT NULL REFERENCES pfe_subjects(id),
			domain_id INTEGER NOT NULL REFERENCES domains(id),
			PRIMARY KEY (subject_id, domain_id)
		)`,
		`CREATE TABLE IF NOT EXISTS wishes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			student_id INTEGER NOT NULL REFERENCES students(id),
			subject_id INTEGER NOT NULL REFERENCES pfe_subjects(id),
			academic_year_id INTEGER NOT NULL REFERENCES academic_years(id),
			status TEXT DEFAULT 'en_attente' CHECK(status IN ('en_attente','accepte','refuse')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS pfe_assignments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			pfe_code TEXT UNIQUE NOT NULL,
			subject_id INTEGER NOT NULL REFERENCES pfe_subjects(id),
			academic_year_id INTEGER NOT NULL REFERENCES academic_years(id),
			student_id INTEGER NOT NULL REFERENCES students(id),
			student2_id INTEGER REFERENCES students(id),
			student3_id INTEGER REFERENCES students(id),
			supervisor_id INTEGER NOT NULL REFERENCES teachers(id),
			co_supervisor_id INTEGER REFERENCES teachers(id),
			memoire_url TEXT,
			status TEXT DEFAULT 'en_cours' CHECK(status IN ('en_cours','memoire_soumis','soutenance_planifiee','valide','refuse')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS pfe_progress_reports (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			assignment_id INTEGER NOT NULL REFERENCES pfe_assignments(id),
			meeting_date DATETIME NOT NULL,
			duration INTEGER NOT NULL,
			meeting_type TEXT NOT NULL CHECK(meeting_type IN ('presentiel','visio')),
			topics TEXT,
			status TEXT DEFAULT 'a_faire' CHECK(status IN ('a_faire','en_cours','termine')),
			observation TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS defense_juries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			assignment_id INTEGER NOT NULL REFERENCES pfe_assignments(id),
			president_id INTEGER NOT NULL REFERENCES teachers(id),
			member_id INTEGER NOT NULL REFERENCES teachers(id),
			president_confirmed BOOLEAN DEFAULT 0,
			member_confirmed BOOLEAN DEFAULT 0,
			president_wants_printed BOOLEAN DEFAULT 0,
			member_wants_printed BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS defenses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			assignment_id INTEGER NOT NULL REFERENCES pfe_assignments(id),
			jury_id INTEGER NOT NULL REFERENCES defense_juries(id),
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
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			defense_id INTEGER NOT NULL REFERENCES defenses(id),
			jury_member_id INTEGER NOT NULL REFERENCES teachers(id),
			criterion1 REAL,
			criterion2 REAL,
			criterion3 REAL,
			criterion4 REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS supervisor_evaluations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			pfe_assignment_id INTEGER UNIQUE NOT NULL REFERENCES pfe_assignments(id),
			evaluator_id INTEGER NOT NULL REFERENCES teachers(id),
			criterion5 REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS company_reports (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			company_id INTEGER NOT NULL REFERENCES companies(id),
			submitted_by INTEGER NOT NULL,
			correction_type TEXT NOT NULL,
			description TEXT,
			requested_value TEXT,
			status TEXT DEFAULT 'en_attente' CHECK(status IN ('en_attente','resolu','rejete')),
			resolved_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			recipient_id INTEGER NOT NULL REFERENCES profiles(id),
			type TEXT NOT NULL,
			payload TEXT DEFAULT '{}',
			read_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			actor_id INTEGER NOT NULL REFERENCES profiles(id),
			action TEXT NOT NULL,
			entity TEXT NOT NULL,
			entity_id INTEGER,
			metadata TEXT DEFAULT '{}',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("migration échouée: %w\nSQL: %s", err, m[:50])
		}
	}


	_, _ = db.Exec(`ALTER TABLE jury_grades ADD COLUMN archive_decision TEXT CHECK(archive_decision IN ('archivable','minor_corrections','major_corrections'))`)


	var progressSchema string
	_ = db.QueryRow(`SELECT sql FROM sqlite_master WHERE type='table' AND name='pfe_progress_reports'`).Scan(&progressSchema)
	if progressSchema != "" && !strings.Contains(progressSchema, "a_faire") {
		if _, err := db.Exec(`PRAGMA foreign_keys=OFF`); err != nil {
			return fmt.Errorf("pragma foreign_keys off: %w", err)
		}
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin migration tx: %w", err)
		}
		steps := []string{
			`CREATE TABLE pfe_progress_reports_v2 (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					assignment_id INTEGER NOT NULL REFERENCES pfe_assignments(id),
					meeting_date DATETIME NOT NULL,
					duration INTEGER NOT NULL,
					meeting_type TEXT NOT NULL CHECK(meeting_type IN ('presentiel','visio')),
					topics TEXT,
					status TEXT DEFAULT 'a_faire' CHECK(status IN ('a_faire','en_cours','termine')),
					observation TEXT,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				)`,
			`INSERT INTO pfe_progress_reports_v2 SELECT * FROM pfe_progress_reports`,
			`DROP TABLE pfe_progress_reports`,
			`ALTER TABLE pfe_progress_reports_v2 RENAME TO pfe_progress_reports`,
		}
		for _, s := range steps {
			if _, err := tx.Exec(s); err != nil {
				_ = tx.Rollback()
				_, _ = db.Exec(`PRAGMA foreign_keys=ON`)
				return fmt.Errorf("migration pfe_progress_reports step: %w", err)
			}
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration: %w", err)
		}
		if _, err := db.Exec(`PRAGMA foreign_keys=ON`); err != nil {
			return fmt.Errorf("pragma foreign_keys on: %w", err)
		}
	}


	defaultDepts := []string{"Informatique", "Electronique", "Chimie", "Finance"}
	for _, name := range defaultDepts {
		_, _ = db.Exec(`INSERT OR IGNORE INTO departments (name) VALUES (?)`, name)
	}


	defaultDomains := []string{
		"Intelligence Artificielle", "Développement Web", "Sécurité Informatique",
		"Big Data", "Réseaux & Télécommunications", "Systèmes Embarqués",
		"Finance & Comptabilité", "Electronique", "Chimie Industrielle",
	}
	for _, name := range defaultDomains {
		_, _ = db.Exec(`INSERT OR IGNORE INTO domains (name) VALUES (?)`, name)
	}


	deptIDs := map[string]int64{}
	for _, name := range defaultDepts {
		var id int64
		_ = db.QueryRow(`SELECT id FROM departments WHERE name = ?`, name).Scan(&id)
		deptIDs[name] = id
	}


	defaultSpecs := []struct {
		name, code, yearType, dept string
	}{
		{"Ingénierie des Systèmes d'Information et Logiciels", "ISIL", "licence", "Informatique"},
		{"Sécurité des Systèmes d'Information", "SSI", "master", "Informatique"},
		{"Big Data Technologies", "BDT", "master", "Informatique"},
		{"Electronique", "ELEC", "licence", "Electronique"},
		{"Systèmes Embarqués", "SEM", "master", "Electronique"},
		{"Chimie", "CHIM", "licence", "Chimie"},
		{"Chimie Pharmaceutique", "CF", "master", "Chimie"},
		{"Business & Finances", "BIZ", "licence", "Finance"},
	}
	for _, s := range defaultSpecs {
		_, _ = db.Exec(`INSERT OR IGNORE INTO specialities (name, code, year_type, department_id) VALUES (?, ?, ?, ?)`,
			s.name, s.code, s.yearType, deptIDs[s.dept])
	}



	var count int
	_ = db.QueryRow(`SELECT COUNT(*) FROM profiles WHERE role = 'admin'`).Scan(&count)
	if count == 0 {
		res, _ := db.Exec(`INSERT INTO profiles (role, full_name, email, is_active) VALUES ('admin', 'DJOUAMAA Amir', 'djouamaa.amir@codiha.com', 1)`)
		adminID, _ := res.LastInsertId()
		_, _ = db.Exec(`INSERT INTO teachers (profile_id, grade, department_id, availability_status) VALUES (?, 'mca', ?, 'disponible')`, adminID, deptIDs["Informatique"])
	}


	var ayCount int
	_ = db.QueryRow(`SELECT COUNT(*) FROM academic_years WHERE status = 'active'`).Scan(&ayCount)
	if ayCount == 0 {
		now := time.Now()
		year := now.Year()

		if now.Month() < 9 {
			year--
		}
		label := fmt.Sprintf("%d-%d", year, year+1)
		_, _ = db.Exec(`INSERT INTO academic_years (label, status) VALUES (?, 'active')`, label)
	}

	return nil
}
