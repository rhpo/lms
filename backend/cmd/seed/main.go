package main

import (
	"database/sql"
	"fmt"
	"log"

	"pfe-backend/internal/config"

	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Erreur connexion DB: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		log.Fatalf("Erreur migration: %v", err)
	}

	fmt.Println("🔧 Seeding database...")

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Erreur début transaction: %v", err)
	}
	defer tx.Rollback()

	// 1. Domaines
	domains := []struct{ id, name string }{
		{"seed-domain-001", "Intelligence Artificielle"},
		{"seed-domain-002", "Génie Logiciel"},
		{"seed-domain-003", "Réseaux & Télécommunications"},
		{"seed-domain-004", "Systèmes Embarqués"},
		{"seed-domain-005", "Cybersécurité"},
		{"seed-domain-006", "Base de Données & Big Data"},
		{"seed-domain-007", "Cloud Computing"},
		{"seed-domain-008", "Systèmes d'Information"},
	}
	for _, d := range domains {
		_, err := tx.Exec(`INSERT OR IGNORE INTO domains (id, name) VALUES (?, ?)`, d.id, d.name)
		if err != nil {
			log.Fatalf("Erreur insertion domaine %s: %v", d.name, err)
		}
	}
	fmt.Println("  ✓ Domaines (8)")

	// 2. Spécialités
	specialities := []struct{ id, name, code, yearType string }{
		{"seed-spec-001", "Ingénierie des Systèmes Informatiques et Logiciels", "ISIL", "master"},
		{"seed-spec-002", "Chimie", "CHIM", "licence"},
		{"seed-spec-003", "Électronique", "ELEC", "master"},
	}
	for _, s := range specialities {
		_, err := tx.Exec(`INSERT OR IGNORE INTO specialities (id, name, code, year_type) VALUES (?, ?, ?, ?)`, s.id, s.name, s.code, s.yearType)
		if err != nil {
			log.Fatalf("Erreur insertion spécialité %s: %v", s.name, err)
		}
	}
	fmt.Println("  ✓ Spécialités (3)")

	// 3. Années universitaires
	academicYears := []struct{ id, label, status string }{
		{"seed-year-001", "2023-2024", "cloturee"},
		{"seed-year-002", "2024-2025", "active"},
	}
	for _, y := range academicYears {
		_, err := tx.Exec(`INSERT OR IGNORE INTO academic_years (id, label, status) VALUES (?, ?, ?)`, y.id, y.label, y.status)
		if err != nil {
			log.Fatalf("Erreur insertion année %s: %v", y.label, err)
		}
	}
	fmt.Println("  ✓ Années universitaires (2)")

	// 4. Promotions
	promotions := []struct{ id, label, academicYearID string }{
		{"seed-promo-001", "ISIL 2024-2025", "seed-year-002"},
		{"seed-promo-002", "CHIM 2024-2025", "seed-year-002"},
	}
	for _, p := range promotions {
		_, err := tx.Exec(`INSERT OR IGNORE INTO promotions (id, label, academic_year_id) VALUES (?, ?, ?)`, p.id, p.label, p.academicYearID)
		if err != nil {
			log.Fatalf("Erreur insertion promotion %s: %v", p.label, err)
		}
	}
	fmt.Println("  ✓ Promotions (2)")

	// 5. Profiles
	profiles := []struct{ id, role, fullName, email string }{
		{"seed-admin-001", "admin", "Admin PFE", "admin@pfe.dz"},
		{"seed-teacher-isil-001", "teacher", "Dr. Mohamed Benali", "m.benali@pfe.dz"},
		{"seed-teacher-isil-002", "teacher", "Pr. Fatima Zohra", "f.zohra@pfe.dz"},
		{"seed-teacher-isil-003", "teacher", "Dr. Ahmed Mansour", "a.mansour@pfe.dz"},
		{"seed-teacher-chim-001", "teacher", "Dr. Leila Meziane", "l.meziane@pfe.dz"},
		{"seed-teacher-chim-002", "teacher", "Pr. Omar Boudiaf", "o.boudiaf@pfe.dz"},
		{"seed-teacher-elec-001", "teacher", "Dr. Nadia Khelifi", "n.khelifi@pfe.dz"},
		{"seed-teacher-elec-002", "teacher", "Dr. Redouane Amrane", "r.amrane@pfe.dz"},
		{"seed-teacher-info-001", "teacher", "Pr. Samira Belkacem", "s.belkacem@pfe.dz"},
		{"seed-teacher-info-002", "teacher", "Dr. Karim Oussama", "k.oussama@pfe.dz"},
		{"seed-student-001", "student", "Anis Djerbi", "a.djerbi@pfe.dz"},
		{"seed-student-002", "student", "Meriem Tabet", "m.tabet@pfe.dz"},
		{"seed-student-003", "student", "Yacine Bouaziz", "y.bouaziz@pfe.dz"},
		{"seed-student-004", "student", "Lyna Chettah", "l.chettah@pfe.dz"},
		{"seed-student-005", "student", "Rayan Haddad", "r.haddad@pfe.dz"},
		{"seed-company-001", "company", "TechCorp Algérie", "contact@techcorp.dz"},
	}
	for _, p := range profiles {
		_, err := tx.Exec(`INSERT OR IGNORE INTO profiles (id, role, full_name, email) VALUES (?, ?, ?, ?)`, p.id, p.role, p.fullName, p.email)
		if err != nil {
			log.Fatalf("Erreur insertion profil %s: %v", p.email, err)
		}
	}
	fmt.Println("  ✓ Profiles (16)")

	// 6. Teachers
	teachers := []struct{ id, profileID, grade, department string }{
		{"seed-teacher-isil-001", "seed-teacher-isil-001", "mca", "ISIL"},
		{"seed-teacher-isil-002", "seed-teacher-isil-002", "professeur", "ISIL"},
		{"seed-teacher-isil-003", "seed-teacher-isil-003", "mcb", "ISIL"},
		{"seed-teacher-chim-001", "seed-teacher-chim-001", "mca", "CHIM"},
		{"seed-teacher-chim-002", "seed-teacher-chim-002", "professeur", "CHIM"},
		{"seed-teacher-elec-001", "seed-teacher-elec-001", "mca", "ELEC"},
		{"seed-teacher-elec-002", "seed-teacher-elec-002", "mcb", "ELEC"},
		{"seed-teacher-info-001", "seed-teacher-info-001", "professeur", "ISIL"},
		{"seed-teacher-info-002", "seed-teacher-info-002", "maa", "ISIL"},
	}
	for _, t := range teachers {
		_, err := tx.Exec(`INSERT OR IGNORE INTO teachers (id, profile_id, grade, department) VALUES (?, ?, ?, ?)`, t.id, t.profileID, t.grade, t.department)
		if err != nil {
			log.Fatalf("Erreur insertion enseignant %s: %v", t.id, err)
		}
	}
	fmt.Println("  ✓ Enseignants (9)")

	// Teacher domains
	teacherDomains := []struct{ teacherID, domainID string }{
		{"seed-teacher-isil-001", "seed-domain-001"},
		{"seed-teacher-isil-001", "seed-domain-002"},
		{"seed-teacher-isil-002", "seed-domain-002"},
		{"seed-teacher-isil-002", "seed-domain-006"},
		{"seed-teacher-isil-003", "seed-domain-003"},
		{"seed-teacher-isil-003", "seed-domain-005"},
		{"seed-teacher-chim-001", "seed-domain-008"},
		{"seed-teacher-chim-002", "seed-domain-008"},
		{"seed-teacher-elec-001", "seed-domain-004"},
		{"seed-teacher-elec-002", "seed-domain-004"},
	}
	for _, td := range teacherDomains {
		_, err := tx.Exec(`INSERT OR IGNORE INTO teacher_domains (teacher_id, domain_id) VALUES (?, ?)`, td.teacherID, td.domainID)
		if err != nil {
			log.Fatalf("Erreur insertion teacher_domain %s-%s: %v", td.teacherID, td.domainID, err)
		}
	}
	fmt.Println("  ✓ Domaines enseignants")

	// 7. Students
	students := []struct{ id, profileID, studentNumber, specialityID, level, promotionID string }{
		{"seed-student-001", "seed-student-001", "202001001", "seed-spec-001", "M2", "seed-promo-001"},
		{"seed-student-002", "seed-student-002", "202001002", "seed-spec-001", "M2", "seed-promo-001"},
		{"seed-student-003", "seed-student-003", "202001003", "seed-spec-001", "M2", "seed-promo-001"},
		{"seed-student-004", "seed-student-004", "202002001", "seed-spec-002", "L3", "seed-promo-002"},
		{"seed-student-005", "seed-student-005", "202001004", "seed-spec-001", "M2", "seed-promo-001"},
	}
	for _, s := range students {
		_, err := tx.Exec(`INSERT OR IGNORE INTO students (id, profile_id, student_number, speciality_id, level, promotion_id) VALUES (?, ?, ?, ?, ?, ?)`,
			s.id, s.profileID, s.studentNumber, s.specialityID, s.level, s.promotionID)
		if err != nil {
			log.Fatalf("Erreur insertion étudiant %s: %v", s.id, err)
		}
	}
	fmt.Println("  ✓ Étudiants (5)")

	// 8. Company
	_, err = tx.Exec(`INSERT OR IGNORE INTO companies (id, profile_id, company_name, sector, description, is_verified) VALUES (?, ?, ?, ?, ?, ?)`,
		"seed-company-001", "seed-company-001", "TechCorp Algérie", "Technologies", "Entreprise de services numériques", 1)
	if err != nil {
		log.Fatalf("Erreur insertion entreprise: %v", err)
	}
	fmt.Println("  ✓ Entreprise (1)")

	// 9. Sujets PFE
	subjects := []struct {
		id, title, description, groupType, proposerID, proposerRole, companyID, academicYearID, status string
		validator1ID, validator2ID                                                                     string
	}{
		{"seed-subject-001", "Chatbot IA pour la scolarité", "Développement d'un chatbot basé sur l'IA pour répondre aux questions des étudiants", "binome", "seed-teacher-isil-001", "teacher", "", "seed-year-002", "en_attente", "", ""},
		{"seed-subject-002", "Plateforme d'évaluation en ligne", "Conception et développement d'une plateforme pour les examens en ligne", "binome", "seed-teacher-isil-002", "teacher", "", "seed-year-002", "valide", "seed-teacher-isil-001", "seed-teacher-isil-003"},
		{"seed-subject-003", "Application IoT pour la gestion d'énergie", "Système de monitoring énergétique basé sur l'IoT", "trinome", "seed-teacher-elec-001", "teacher", "", "seed-year-002", "valide", "seed-teacher-isil-001", "seed-teacher-info-001"},
		{"seed-subject-004", "Analyse des processus chimiques", "Optimisation des processus chimiques avec l'IA", "monome", "seed-teacher-chim-001", "teacher", "", "seed-year-002", "refuse", "seed-teacher-chim-002", "seed-teacher-info-001"},
		{"seed-subject-005", "Migration cloud de l'infrastructure", "Étude et migration de l'infrastructure IT vers le cloud", "binome", "seed-company-001", "company", "seed-company-001", "seed-year-002", "en_attente", "", ""},
		{"seed-subject-006", "Application de gestion des stages", "Développement d'une application web pour la gestion des stages en entreprise", "binome", "seed-teacher-isil-003", "teacher", "", "seed-year-002", "valide", "seed-teacher-info-001", "seed-teacher-isil-002"},
	}
	for _, s := range subjects {
		_, err := tx.Exec(`INSERT OR IGNORE INTO pfe_subjects (id, title, description, group_type, proposer_id, proposer_role, company_id, academic_year_id, status, validator1_id, validator2_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			s.id, s.title, s.description, s.groupType, s.proposerID, s.proposerRole, s.companyID, s.academicYearID, s.status, s.validator1ID, s.validator2ID)
		if err != nil {
			log.Fatalf("Erreur insertion sujet %s: %v", s.id, err)
		}
	}
	fmt.Println("  ✓ Sujets PFE (6)")

	// 10. Voeux
	wishes := []struct{ id, studentID, subjectID, academicYearID, status string }{
		{"seed-wish-001", "seed-student-001", "seed-subject-002", "seed-year-002", "accepte"},
		{"seed-wish-002", "seed-student-002", "seed-subject-002", "seed-year-002", "accepte"},
		{"seed-wish-003", "seed-student-003", "seed-subject-003", "seed-year-002", "accepte"},
		{"seed-wish-004", "seed-student-001", "seed-subject-006", "seed-year-002", "en_attente"},
		{"seed-wish-005", "seed-student-005", "seed-subject-003", "seed-year-002", "en_attente"},
	}
	for _, w := range wishes {
		_, err := tx.Exec(`INSERT OR IGNORE INTO wishes (id, student_id, subject_id, academic_year_id, status) VALUES (?, ?, ?, ?, ?)`,
			w.id, w.studentID, w.subjectID, w.academicYearID, w.status)
		if err != nil {
			log.Fatalf("Erreur insertion voeu %s: %v", w.id, err)
		}
	}
	fmt.Println("  ✓ Voeux (5)")

	// 11. Affectations PFE
	_, err = tx.Exec(`INSERT OR IGNORE INTO pfe_assignments (id, pfe_code, subject_id, academic_year_id, student_id, student2_id, supervisor_id, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		"seed-assign-001", "PFE-ISIL-2025-001", "seed-subject-002", "seed-year-002", "seed-student-001", "seed-student-002", "seed-teacher-isil-002", "en_cours")
	if err != nil {
		log.Fatalf("Erreur insertion assignation 1: %v", err)
	}
	_, err = tx.Exec(`INSERT OR IGNORE INTO pfe_assignments (id, pfe_code, subject_id, academic_year_id, student_id, student2_id, student3_id, supervisor_id, co_supervisor_id, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"seed-assign-002", "PFE-ELEC-2025-001", "seed-subject-003", "seed-year-002", "seed-student-003", "seed-student-005", "", "seed-teacher-elec-001", "seed-teacher-isil-001", "en_cours")
	if err != nil {
		log.Fatalf("Erreur insertion assignation 2: %v", err)
	}
	fmt.Println("  ✓ Affectations PFE (2)")

	// 12. Meetings pour assignation 1
	meetings := []struct{ id, assignmentID, date, meetingType, topics, status string }{
		{"seed-meeting-001", "seed-assign-001", "2025-02-10 10:00:00", "presentiel", "Introduction, cahier des charges", "termine"},
		{"seed-meeting-002", "seed-assign-001", "2025-02-24 10:00:00", "visio", "Architecture, technologies", "termine"},
		{"seed-meeting-003", "seed-assign-001", "2025-03-10 10:00:00", "presentiel", "Avancement module auth", "termine"},
		{"seed-meeting-004", "seed-assign-001", "2025-03-24 10:00:00", "presentiel", "Tests, déploiement", "termine"},
		{"seed-meeting-005", "seed-assign-001", "2025-04-07 10:00:00", "visio", "Revue finale", "en_cours"},
	}
	for _, m := range meetings {
		_, err := tx.Exec(`INSERT OR IGNORE INTO pfe_progress_reports (id, assignment_id, meeting_date, duration, meeting_type, topics, status) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			m.id, m.assignmentID, m.date, 60, m.meetingType, m.topics, m.status)
		if err != nil {
			log.Fatalf("Erreur insertion meeting %s: %v", m.id, err)
		}
	}
	fmt.Println("  ✓ Meetings (5)")

	// 13. Évaluation encadrant
	_, err = tx.Exec(`INSERT OR IGNORE INTO supervisor_evaluations (id, pfe_assignment_id, evaluator_id, criterion5) VALUES (?, ?, ?, ?)`,
		"seed-eval-001", "seed-assign-001", "seed-teacher-isil-002", 3.5)
	if err != nil {
		log.Fatalf("Erreur insertion évaluation: %v", err)
	}
	fmt.Println("  ✓ Évaluation encadrant (1)")

	// 14. Soutenance + Jury
	_, err = tx.Exec(`INSERT OR IGNORE INTO defense_juries (id, assignment_id, president_id, member_id) VALUES (?, ?, ?, ?)`,
		"seed-jury-001", "seed-assign-001", "seed-teacher-isil-001", "seed-teacher-isil-003")
	if err != nil {
		log.Fatalf("Erreur insertion jury: %v", err)
	}
	_, err = tx.Exec(`INSERT OR IGNORE INTO defenses (id, assignment_id, jury_id, scheduled_at, room, status) VALUES (?, ?, ?, ?, ?, ?)`,
		"seed-defense-001", "seed-assign-001", "seed-jury-001", "2025-06-15 09:00:00", "Salle A101", "scheduled")
	if err != nil {
		log.Fatalf("Erreur insertion soutenance: %v", err)
	}
	fmt.Println("  ✓ Soutenance + Jury (1)")

	if err := tx.Commit(); err != nil {
		log.Fatalf("Erreur commit: %v", err)
	}

	fmt.Println("\n✅ Seed terminé avec succès !")
}

func runMigrations(db *sql.DB) error {
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
			return fmt.Errorf("migration échouée: %w\nSQL: %s", err, m[:50])
		}
	}
	return nil
}
