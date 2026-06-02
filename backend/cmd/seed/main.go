package main

import (
	"database/sql"
	"fmt"
	"log"

	"pfe-backend/internal/config"

	_ "modernc.org/sqlite"
)

func mustExec(tx *sql.Tx, query string, args ...any) sql.Result {
	res, err := tx.Exec(query, args...)
	if err != nil {
		log.Fatalf("SQL error:\n  query : %s\n  args  : %v\n  error : %v", query, args, err)
	}
	return res
}

func queryID(tx *sql.Tx, query string, args ...any) int64 {
	var id int64
	_ = tx.QueryRow(query, args...).Scan(&id)
	return id
}

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

	deptNames := []string{"Informatique", "Electronique", "Chimie", "Finance"}
	deptIDs := map[string]int64{}
	for _, name := range deptNames {
		mustExec(tx, `INSERT OR IGNORE INTO departments (name) VALUES (?)`, name)
		deptIDs[name] = queryID(tx, `SELECT id FROM departments WHERE name = ?`, name)
	}
	fmt.Println("[DONE] Départements (4)")

	domainNames := []string{
		"Intelligence Artificielle",
		"Génie Logiciel",
		"Systèmes Embarqués",
		"Internet des Objets (IoT)",
		"Réseaux & Sécurité",
		"Big Data & Science des Données",
		"Cloud Computing",
		"Cybersécurité",
	}
	domainIDs := map[string]int64{}
	for _, name := range domainNames {
		mustExec(tx, `INSERT OR IGNORE INTO domains (name) VALUES (?)`, name)
		domainIDs[name] = queryID(tx, `SELECT id FROM domains WHERE name = ?`, name)
	}
	fmt.Println("[DONE] Domaines (8)")

	type specDef struct{ name, code, yearType, dept string }
	specs := []specDef{
		{"Ingénierie des Systèmes d'Information et Logiciels", "ISIL", "licence", "Informatique"},
		{"Sécurité des Systèmes d'Information", "SSI", "master", "Informatique"},
		{"Big Data Technologies", "BDT", "master", "Informatique"},
		{"Electronique", "ELEC", "licence", "Electronique"},
		{"Systèmes Embarqués", "SEM", "master", "Electronique"},
		{"Chimie", "CHIM", "licence", "Chimie"},
		{"Chimie Pharmaceutique", "CF", "master", "Chimie"},
		{"Business & Finances", "BIZ", "licence", "Finance"},
	}
	specIDs := map[string]int64{}
	for _, s := range specs {
		mustExec(tx,
			`INSERT OR IGNORE INTO specialities (name, code, year_type, department_id) VALUES (?, ?, ?, ?)`,
			s.name, s.code, s.yearType, deptIDs[s.dept])
		specIDs[s.code] = queryID(tx, `SELECT id FROM specialities WHERE code = ?`, s.code)
	}
	fmt.Println("[DONE] Spécialités (8)")

	mustExec(tx, `INSERT OR IGNORE INTO academic_years (label, status) VALUES (?, ?)`, "2023-2024", "cloturee")
	mustExec(tx, `INSERT OR IGNORE INTO academic_years (label, status, submission_open_at, submission_close_at, max_wishes) VALUES (?, ?, ?, ?, ?)`,
		"2024-2025", "active", "2025-01-15 00:00:00", "2025-06-30 23:59:59", 5)
	yearID2024 := queryID(tx, `SELECT id FROM academic_years WHERE label = ?`, "2023-2024")
	yearID2025 := queryID(tx, `SELECT id FROM academic_years WHERE label = ?`, "2024-2025")
	_ = yearID2024
	fmt.Println("[DONE] Années académiques (2)")

	mustExec(tx, `INSERT OR IGNORE INTO promotions (label, academic_year_id) VALUES (?, ?)`, "ISIL M2 2024-2025", yearID2025)
	mustExec(tx, `INSERT OR IGNORE INTO promotions (label, academic_year_id) VALUES (?, ?)`, "ELEC M2 2024-2025", yearID2025)
	promoISIL := queryID(tx, `SELECT id FROM promotions WHERE label = ?`, "ISIL M2 2024-2025")
	promoELEC := queryID(tx, `SELECT id FROM promotions WHERE label = ?`, "ELEC M2 2024-2025")
	fmt.Println("[DONE] Promotions (2)")

	mustExec(tx,
		`INSERT OR IGNORE INTO profiles (role, full_name, email, is_active) VALUES ('admin', 'DJOUAMAA Amir', 'djouamaa.amir@codiha.com', 1)`)
	djouamaaProfileID := queryID(tx, `SELECT id FROM profiles WHERE email = ?`, "djouamaa.amir@codiha.com")

	mustExec(tx,
		`INSERT OR IGNORE INTO teachers (profile_id, grade, department_id, availability_status) VALUES (?, 'mca', ?, 'disponible')`,
		djouamaaProfileID, deptIDs["Informatique"])
	djouamaaTeacherID := queryID(tx, `SELECT id FROM teachers WHERE profile_id = ?`, djouamaaProfileID)

	for _, d := range []string{"Intelligence Artificielle", "Génie Logiciel", "Systèmes Embarqués"} {
		mustExec(tx, `INSERT OR IGNORE INTO teacher_domains (teacher_id, domain_id) VALUES (?, ?)`, djouamaaTeacherID, domainIDs[d])
	}

	mustExec(tx,
		`INSERT OR IGNORE INTO profiles (role, full_name, email, is_active) VALUES ('teacher', 'BOUYAKOUB Faycal', 'bouyakoub.faycal@codiha.com', 1)`)
	bouyakProfileID := queryID(tx, `SELECT id FROM profiles WHERE email = ?`, "bouyakoub.faycal@codiha.com")

	mustExec(tx,
		`INSERT OR IGNORE INTO teachers (profile_id, grade, department_id, availability_status) VALUES (?, 'professeur', ?, 'disponible')`,
		bouyakProfileID, deptIDs["Informatique"])
	bouyakTeacherID := queryID(tx, `SELECT id FROM teachers WHERE profile_id = ?`, bouyakProfileID)

	for _, d := range []string{"Intelligence Artificielle", "Systèmes Embarqués", "Génie Logiciel"} {
		mustExec(tx, `INSERT OR IGNORE INTO teacher_domains (teacher_id, domain_id) VALUES (?, ?)`, bouyakTeacherID, domainIDs[d])
	}

	mustExec(tx,
		`INSERT OR IGNORE INTO profiles (role, full_name, email, is_active) VALUES ('teacher', 'OUADAHI Nazim', 'ouadahi.nazim@codiha.com', 1)`)
	ouadahiProfileID := queryID(tx, `SELECT id FROM profiles WHERE email = ?`, "ouadahi.nazim@codiha.com")

	mustExec(tx,
		`INSERT OR IGNORE INTO teachers (profile_id, grade, department_id, availability_status) VALUES (?, 'mca', ?, 'disponible')`,
		ouadahiProfileID, deptIDs["Electronique"])
	ouadahiTeacherID := queryID(tx, `SELECT id FROM teachers WHERE profile_id = ?`, ouadahiProfileID)

	for _, d := range []string{"Internet des Objets (IoT)", "Systèmes Embarqués"} {
		mustExec(tx, `INSERT OR IGNORE INTO teacher_domains (teacher_id, domain_id) VALUES (?, ?)`, ouadahiTeacherID, domainIDs[d])
	}

	mustExec(tx,
		`INSERT OR IGNORE INTO profiles (role, full_name, email, is_active) VALUES ('teacher', 'MAHMOUDI Mohamed', 'mahmoudi.mohamed@codiha.com', 1)`)
	mahmoudiProfileID := queryID(tx, `SELECT id FROM profiles WHERE email = ?`, "mahmoudi.mohamed@codiha.com")

	mustExec(tx,
		`INSERT OR IGNORE INTO teachers (profile_id, grade, department_id, availability_status) VALUES (?, 'mca', ?, 'disponible')`,
		mahmoudiProfileID, deptIDs["Informatique"])
	mahmoudiTeacherID := queryID(tx, `SELECT id FROM teachers WHERE profile_id = ?`, mahmoudiProfileID)

	for _, d := range []string{"Intelligence Artificielle", "Internet des Objets (IoT)", "Génie Logiciel"} {
		mustExec(tx, `INSERT OR IGNORE INTO teacher_domains (teacher_id, domain_id) VALUES (?, ?)`, mahmoudiTeacherID, domainIDs[d])
	}

	fmt.Printf("[DONE] Enseignants (4) - IDs teacher: DJOUAMAA=%d, BOUYAKOUB=%d, OUADAHI=%d, MAHMOUDI=%d\n",
		djouamaaTeacherID, bouyakTeacherID, ouadahiTeacherID, mahmoudiTeacherID)

	type studentDef struct {
		fullName, email, number, specCode, level string
		promoID                                  int64
	}
	students := []studentDef{
		{"HADID Rami", "rami.hadid@codiha.com", "202331676584", "ISIL", "M2", promoISIL},
		{"GOUMIRI Samy", "samy.goumiri@codiha.com", "202331589234", "ISIL", "M2", promoISIL},
		{"BOUDRAA Mohamed Anis", "anis.boudraa@codiha.com", "202332145678", "ISIL", "M2", promoISIL},
		{"SAADI Saber", "saber.saadi@codiha.com", "202330987654", "ISIL", "M2", promoISIL},
		{"ABDERAHIM Arab", "arab.abderahim@codiha.com", "202331234567", "ISIL", "M2", promoISIL},
		{"BOUAFIA Racim", "racim.bouafia@codiha.com", "202332056789", "ELEC", "M2", promoELEC},
	}

	studentIDs := map[string]int64{}
	for _, s := range students {
		mustExec(tx,
			`INSERT OR IGNORE INTO profiles (role, full_name, email, is_active) VALUES ('student', ?, ?, 1)`,
			s.fullName, s.email)
		pid := queryID(tx, `SELECT id FROM profiles WHERE email = ?`, s.email)

		mustExec(tx,
			`INSERT OR IGNORE INTO students (profile_id, student_number, speciality_id, level, promotion_id) VALUES (?, ?, ?, ?, ?)`,
			pid, s.number, specIDs[s.specCode], s.level, s.promoID)
		studentIDs[s.email] = queryID(tx, `SELECT id FROM students WHERE profile_id = ?`, pid)
	}
	fmt.Println("[DONE] Étudiants (6)")

	mustExec(tx,
		`INSERT OR IGNORE INTO profiles (role, full_name, email, is_active) VALUES ('company', 'Salah Zeghdani', 'salah.zeghdani@sonatrach.dz', 1)`)
	zeghdaniProfileID := queryID(tx, `SELECT id FROM profiles WHERE email = ?`, "salah.zeghdani@sonatrach.dz")

	mustExec(tx, `INSERT OR IGNORE INTO companies
		(profile_id, company_name, sector, description, contact_email, contact_phone, is_verified)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		zeghdaniProfileID,
		"Sonatrach",
		"Énergie & Pétrole",
		"Société Nationale pour la Recherche, la Production, le Transport, la Transformation et la Commercialisation des Hydrocarbures.",
		"sonatrach@sonatrach.dz",
		"023 48 31 31",
		0,
	)
	sonatrachCompanyID := queryID(tx, `SELECT id FROM companies WHERE profile_id = ?`, zeghdaniProfileID)
	fmt.Printf("[DONE] Entreprise Sonatrach (en attente) - company_id=%d\n", sonatrachCompanyID)

	type subjectDef struct {
		title, description, groupType string
		proposerID                    int64
		proposerRole                  string
		companyID                     int64
		yearID                        int64
		status                        string
		v1ID, v2ID                    int64
		v1Decision, v2Decision        string
	}

	subjects := []subjectDef{
		{
			title:       "Plateforme intelligente de gestion des PFE",
			description: "Conception et développement d'une plateforme web complète pour la gestion du cycle de vie des Projets de Fin d'Études : dépôt de sujets, expression de vœux, affectation automatique, suivi d'avancement et organisation des soutenances.",
			groupType:   "binome",
			proposerID:  djouamaaProfileID, proposerRole: "teacher",
			yearID: yearID2025,
			status: "valide",
			v1ID:   bouyakTeacherID, v1Decision: "valide",
			v2ID: ouadahiTeacherID, v2Decision: "valide",
		},
		{
			title:       "Système de recommandation de sujets PFE par apprentissage automatique",
			description: "Développement d'un moteur de recommandation utilisant des algorithmes de machine learning pour suggérer des sujets PFE adaptés au profil académique et aux intérêts de chaque étudiant. Le système exploitera l'historique des affectations et les domaines de compétence des encadrants.",
			groupType:   "binome",
			proposerID:  bouyakProfileID, proposerRole: "teacher",
			yearID: yearID2025,
			status: "valide",
			v1ID:   djouamaaTeacherID, v1Decision: "valide",
			v2ID: mahmoudiTeacherID, v2Decision: "valide",
		},
		{
			title:       "Réseau de capteurs IoT pour la surveillance environnementale",
			description: "Conception d'une architecture IoT distribuée pour la surveillance en temps réel de paramètres environnementaux (température, humidité, qualité de l'air). Le système inclut une passerelle LoRaWAN, un tableau de bord web et des alertes intelligentes basées sur des seuils adaptatifs.",
			groupType:   "binome",
			proposerID:  ouadahiProfileID, proposerRole: "teacher",
			yearID: yearID2025,
			status: "valide",
			v1ID:   djouamaaTeacherID, v1Decision: "valide",
			v2ID: bouyakTeacherID, v2Decision: "valide",
		},
		{
			title:       "Monitoring des pipelines par IoT et Intelligence Artificielle",
			description: "Développement d'un système de surveillance en temps réel des pipelines de transport d'hydrocarbures de Sonatrach. Le projet intègre des capteurs IoT pour la détection de fuites et d'anomalies de pression, couplés à des modèles d'IA pour la maintenance prédictive.",
			groupType:   "binome",
			proposerID:  zeghdaniProfileID, proposerRole: "company",
			companyID: sonatrachCompanyID,
			yearID:    yearID2025,
			status:    "en_attente",
		},
	}

	subjectIDs := make([]int64, len(subjects))
	for i, s := range subjects {
		mustExec(tx,
			`INSERT OR IGNORE INTO pfe_subjects
				(title, description, group_type, proposer_id, proposer_role, company_id,
				 academic_year_id, status, validator1_id, validator2_id,
				 validator1_decision, validator2_decision)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			s.title, s.description, s.groupType, s.proposerID, s.proposerRole,
			nullInt(s.companyID), s.yearID, s.status,
			nullInt(s.v1ID), nullInt(s.v2ID),
			nullStr(s.v1Decision), nullStr(s.v2Decision),
		)
		subjectIDs[i] = queryID(tx, `SELECT id FROM pfe_subjects WHERE title = ?`, s.title)
	}
	fmt.Println("[DONE] Sujets PFE (4)")

	type wishDef struct {
		studentEmail string
		subjectIdx   int
		status       string
	}
	wishes := []wishDef{
		{"rami.hadid@codiha.com", 1, "accepte"},
		{"samy.goumiri@codiha.com", 1, "accepte"},
		{"anis.boudraa@codiha.com", 0, "accepte"},
		{"saber.saadi@codiha.com", 0, "accepte"},
		{"arab.abderahim@codiha.com", 1, "en_attente"},
		{"racim.bouafia@codiha.com", 2, "en_attente"},
	}
	for _, w := range wishes {
		sid := studentIDs[w.studentEmail]
		subID := subjectIDs[w.subjectIdx]
		mustExec(tx,
			`INSERT OR IGNORE INTO wishes (student_id, subject_id, academic_year_id, status) VALUES (?, ?, ?, ?)`,
			sid, subID, yearID2025, w.status)
	}
	fmt.Println("[DONE] Vœux (6)")

	mustExec(tx,
		`INSERT OR IGNORE INTO pfe_assignments
			(pfe_code, subject_id, academic_year_id, student_id, student2_id, supervisor_id, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"PFE-ISIL-2025-001",
		subjectIDs[1], yearID2025,
		studentIDs["rami.hadid@codiha.com"],
		studentIDs["samy.goumiri@codiha.com"],
		bouyakTeacherID,
		"soutenance_planifiee",
	)
	assign1ID := queryID(tx, `SELECT id FROM pfe_assignments WHERE pfe_code = ?`, "PFE-ISIL-2025-001")

	mustExec(tx,
		`INSERT OR IGNORE INTO pfe_assignments
			(pfe_code, subject_id, academic_year_id, student_id, student2_id, supervisor_id, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"PFE-ISIL-2025-002",
		subjectIDs[0], yearID2025,
		studentIDs["anis.boudraa@codiha.com"],
		studentIDs["saber.saadi@codiha.com"],
		djouamaaTeacherID,
		"en_cours",
	)
	assign2ID := queryID(tx, `SELECT id FROM pfe_assignments WHERE pfe_code = ?`, "PFE-ISIL-2025-002")
	_ = assign2ID

	fmt.Println("[DONE] Affectations PFE (2)")

	type meetingDef struct {
		date, meetingType, topics, observation, status string
		duration                                       int
	}
	meetings := []meetingDef{
		{
			date: "2025-02-17 10:00:00", duration: 60, meetingType: "presentiel",
			topics:      "Définition du sujet, périmètre fonctionnel, cahier des charges initial",
			observation: "Bonne compréhension du sujet. Les étudiants ont proposé des fonctionnalités pertinentes.",
			status:      "termine",
		},
		{
			date: "2025-03-03 14:00:00", duration: 90, meetingType: "visio",
			topics:      "Architecture logicielle, choix technologiques (SvelteKit, Go Fiber, SQLite), modèle de données",
			observation: "Architecture validée. Découpage en modules bien pensé.",
			status:      "termine",
		},
		{
			date: "2025-03-24 10:00:00", duration: 60, meetingType: "presentiel",
			topics:      "Avancement module authentification et gestion des rôles, démonstration partielle",
			observation: "Retard sur le module de notification. À rattraper d'ici la prochaine séance.",
			status:      "termine",
		},
		{
			date: "2025-04-14 15:00:00", duration: 90, meetingType: "visio",
			topics:      "Revue mi-parcours : module voeux étudiant, affectation admin, journal de suivi encadrant",
			observation: "",
			status:      "en_cours",
		},
		{
			date: "2025-05-12 10:00:00", duration: 60, meetingType: "presentiel",
			topics:      "Finalisation des fonctionnalités, préparation du rapport et de la soutenance",
			observation: "",
			status:      "a_faire",
		},
	}
	for _, m := range meetings {
		obs := nullStr(m.observation)
		mustExec(tx,
			`INSERT OR IGNORE INTO pfe_progress_reports
				(assignment_id, meeting_date, duration, meeting_type, topics, observation, status)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			assign1ID, m.date, m.duration, m.meetingType, m.topics, obs, m.status,
		)
	}
	fmt.Println("[DONE] Journal de suivi - 5 réunions (assignment 1)")

	mustExec(tx,
		`INSERT OR IGNORE INTO supervisor_evaluations (pfe_assignment_id, evaluator_id, criterion5) VALUES (?, ?, ?)`,
		assign1ID, bouyakTeacherID, 3.5,
	)
	fmt.Println("[DONE] Évaluation encadrant - 3.5/4 (BOUYAKOUB -> HADID/GOUMIRI)")

	mustExec(tx,
		`INSERT OR IGNORE INTO defense_juries
			(assignment_id, president_id, member_id)
		VALUES (?, ?, ?)`,
		assign1ID, djouamaaTeacherID, mahmoudiTeacherID,
	)
	juryID := queryID(tx, `SELECT id FROM defense_juries WHERE assignment_id = ?`, assign1ID)

	mustExec(tx,
		`INSERT OR IGNORE INTO defenses
			(assignment_id, jury_id, scheduled_at, room, status)
		VALUES (?, ?, ?, ?, ?)`,
		assign1ID, juryID, "2025-06-16 09:00:00", "Salle B205", "scheduled",
	)
	fmt.Println("[DONE] Jury + Soutenance planifiée - 16 juin 2025, Salle B205")

	if err := tx.Commit(); err != nil {
		log.Fatalf("Erreur commit: %v", err)
	}

	fmt.Println("\n✅ Seed terminé avec succès !")
	fmt.Println("\n📋 Comptes disponibles :")
	fmt.Println("  Admin/Enseignant : djouamaa.amir@codiha.com")
	fmt.Println("  Enseignants      : bouyakoub.faycal@codiha.com | ouadahi.nazim@codiha.com | mahmoudi.mohamed@codiha.com")
	fmt.Println("  Étudiants ISIL   : rami.hadid@codiha.com | samy.goumiri@codiha.com | anis.boudraa@codiha.com")
	fmt.Println("                     saber.saadi@codiha.com | arab.abderahim@codiha.com")
	fmt.Println("  Étudiant ELEC    : racim.bouafia@codiha.com")
	fmt.Println("  Entreprise       : salah.zeghdani@sonatrach.dz  (Sonatrach - en attente de validation)")
}

func nullInt(v int64) any {
	if v == 0 {
		return nil
	}
	return v
}

func nullStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func runMigrations(db *sql.DB) error {
	stmts := []string{
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
		`CREATE TABLE IF NOT EXISTS departments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS domains (
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
			entity_id TEXT,
			metadata TEXT DEFAULT '{}',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("migration échouée: %w", err)
		}
	}
	return nil
}
