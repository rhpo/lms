/**
 * Domain Types
 *
 * Single source of truth for all business types in the application.
 * All types are string unions (never TypeScript enums) so they can be
 * used safely on both server and client without tree-shaking issues.
 *
 * Naming convention: English identifiers, French UI labels defined separately.
 */

// === Role & Auth ===============

export type UserRole = "student" | "teacher" | "admin" | "company";

// === Academic ==

export type Level = "licence" | "master" | "ingenieur";

export type AcademicYearStatus = "active" | "cloturee";

// === Teacher ===

/**
 * Teacher grades in ascending order:
 * assistant < mab < maa < mcb < mca < professeur
 */
export type TeacherGrade =
  | "assistant"
  | "mab"
  | "maa"
  | "mcb"
  | "mca"
  | "professeur";

export type AvailabilityStatus =
  | "disponible"
  | "indisponible"
  | "indisponible_jusqu_au"; // requires unavailable_until date

// === Subject ===

export type GroupType = "monome" | "binome" | "trinome";

/**
 * Subject lifecycle statuses (visible to end-users):
 * - en_attente : submitted, waiting for both validators
 * - valide     : both validators approved
 * - refuse     : at least one validator refused
 * - expire     : academic year closed while subject was pending
 */
export type SubjectStatus = "en_attente" | "valide" | "refuse" | "expire";

/**
 * Decision rendered by a validator teacher on a subject.
 */
export type ReviewDecision =
  | "accepte"
  | "accepte_sous_reserve"
  | "refuse";

// === Wish (Voeu) ==============

export type WishStatus = "en_attente" | "accepte" | "refuse";

// === PFE Assignment ===========

/**
 * PFE (assignment) lifecycle statuses:
 * - en_cours             : work in progress
 * - soutenance_planifiee : defense has been scheduled
 * - valide               : admitted after defense
 * - refuse               : not admitted after defense
 */
export type PfeStatus =
  | "en_cours"
  | "soutenance_planifiee"
  | "valide"
  | "refuse";

// === Defense ===

export type DefenseStatus = "scheduled" | "done" | "postponed";

export type DefenseResult =
  | "admitted"
  | "corrections_required"
  | "not_admitted";

/**
 * How the jury president resolves a grade disagreement (gap > 1 point).
 */
export type TiebreakChoice = "president" | "member" | "average";

// === Audit =====

export type AuditActionType =
  | "SUBJECT_SUBMITTED"
  | "SUBJECT_VALIDATED"
  | "SUBJECT_REFUSED"
  | "STUDENT_ASSIGNED"
  | "JURY_ASSIGNED"
  | "DEFENSE_SCHEDULED"
  | "GRADE_SUBMITTED"
  | "GRADE_VALIDATED"
  | "COMPANY_REGISTERED"
  | "WISH_SUBMITTED"
  | "NOTIFICATION_SENT"
  | "CSV_IMPORTED"
  | "ROLE_TRANSFERRED";

// === Interfaces ================

export interface Profile {
  id: string;
  role: UserRole;
  full_name: string;
  avatar_url: string | null;
  created_at: string;
  updated_at: string;
}

export interface Student {
  id: string;
  profile_id: string | null;
  email: string | null;
  student_number: string;
  specialty: string;
  level: Level;
  promotion_year: number;
  promotion_id: string | null;
  created_at: string;
  updated_at: string;
}

export interface Teacher {
  id: string;
  profile_id: string | null;
  email: string | null;
  grade: TeacherGrade;
  department: string;
  specialties: string[];
  availability_status: AvailabilityStatus;
  unavailable_until: string | null; // ISO date, only relevant when status = indisponible_jusqu_au
  created_at: string;
  updated_at: string;
}

export interface Company {
  profile_id: string;
  company_name: string;
  address: string;
  sector: string;
  website: string | null;
  contact_phone: string;
  is_verified: boolean;
  created_at: string;
  updated_at: string;
}

export interface AcademicYear {
  id: string;
  label: string;                // ex. "2024-2025"
  status: AcademicYearStatus;
  submission_open_at: string;   // ISO datetime
  submission_close_at: string;  // ISO datetime — subjects expire after this
  max_wishes: number;           // max voeux per student (configured by admin)
  created_at: string;
  updated_at: string;
}

export interface Promotion {
  id: string;
  label: string;                // ex. "Promotion ISIL 2024-2025"
  specialty: string;
  academic_year_id: string;
  created_at: string;
  updated_at: string;
}

export interface PFESubject {
  id: string;
  code: string | null;
  title: string;
  description: string;
  group_type: GroupType;
  specialty: string;
  level: Level;
  proposer_id: string;
  proposer_role: UserRole;
  company_id: string | null;
  academic_year_id: string | null;
  // Validators assigned by admin
  validator1_id: string | null;
  validator2_id: string | null;
  validator1_decision: ReviewDecision | null;
  validator2_decision: ReviewDecision | null;
  validator1_comment: string | null;
  validator2_comment: string | null;
  status: SubjectStatus;
  created_at: string;
  updated_at: string;
}

/**
 * A student wish (voeu) for a specific subject in a given academic year.
 * No priority order — the supervisor chooses among all applicants.
 */
export interface Wish {
  id: string;
  student_id: string;          // profile_id of the student
  subject_id: string;
  academic_year_id: string;
  status: WishStatus;
  created_at: string;
  updated_at: string;
}

/**
 * A PFE assignment links a validated subject to 1–3 students and a supervisor.
 * student2_id and student3_id are null for monome, student3_id is null for binome.
 */
export interface PFEAssignment {
  id: string;
  subject_id: string;
  academic_year_id: string;
  student_id: string;
  student2_id: string | null;
  student3_id: string | null;
  supervisor_id: string;
  co_supervisor_id: string | null; // company co-supervisor
  status: PfeStatus;
  assigned_at: string;
  created_at: string;
  updated_at: string;
}

export interface PFEProgressReport {
  id: string;
  assignment_id: string;
  meeting_date: string;          // ISO date "YYYY-MM-DD"
  student_notes: string;
  teacher_feedback: string | null;
  attachments: string[];
  signed_by_teacher: boolean;
  created_at: string;
  updated_at: string;
}

/**
 * A defense jury has exactly two members: a president and one member.
 */
export interface DefenseJury {
  id: string;
  assignment_id: string;
  president_id: string;
  member_id: string;
  created_at: string;
  updated_at: string;
}

/**
 * One jury member's grade for a defense.
 * Each criterion is scored 0–5; total = sum (max 20).
 */
export interface JuryGrade {
  id: string;
  defense_id: string;
  jury_member_id: string;   // profile_id of president or member
  criterion1: number;       // /5
  criterion2: number;       // /5
  criterion3: number;       // /5
  criterion4: number;       // /5
  total: number;            // computed: criterion1+2+3+4 (0–20)
  submitted_at: string;
  created_at: string;
  updated_at: string;
}

export interface Defense {
  id: string;
  assignment_id: string;
  jury_id: string;
  scheduled_at: string;
  room: string;
  status: DefenseStatus;
  result: DefenseResult | null;
  final_grade: number | null;              // resolved after grading or tiebreak
  grade_gap: number | null;               // |president_total - member_total|
  tiebreak_choice: TiebreakChoice | null; // only set when gap > 1
  comment: string | null;
  created_at: string;
  updated_at: string;
}

export interface Notification {
  id: string;
  recipient_id: string;
  type: string;
  payload: Record<string, unknown>;
  read_at: string | null;
  created_at: string;
}

export interface AuditLog {
  id: string;
  actor_id: string;
  action: AuditActionType;
  entity: string;
  entity_id: string;
  metadata: Record<string, unknown>;
  created_at: string;
}

export interface PFEStatusHistory {
  id: string;
  subject_id: string;
  from_status: SubjectStatus | null;
  to_status: SubjectStatus;
  actor_id: string;
  reason: string | null;
  created_at: string;
}

export interface SessionUser {
  id: string;
  email: string;
  role: UserRole;
  full_name: string;
  avatar_url: string | null;
  provider?: string;
  created_at?: string;
}

/**
 * Type for Svelte 4 class-based components used as icons (e.g. lucide-svelte).
 * These extend SvelteComponentTyped and accept size/color props.
 *
 * WARNING: This is a deliberately loose type to work around the Svelte 4/5
 * incompatibility with lucide-svelte (which exports Svelte 4 classes while
 * the project uses Svelte 5's Component type). DO NOT reuse this type for
 * anything other than icon props. For proper Svelte 5 component typing,
 * use `Component` from 'svelte' instead.
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type IconComponent = new (...args: any[]) => any;
