





export type UserRole = 'admin' | 'teacher' | 'student' | 'company';

export type TeacherGrade = 'assistant' | 'mab' | 'maa' | 'mcb' | 'mca' | 'professeur';

export type AvailabilityStatus = 'disponible' | 'indisponible' | 'indisponible_jusqu_au';

export type YearType = 'licence' | 'master';

export type GroupType = 'monome' | 'binome' | 'trinome';

export type ProposerRole = 'teacher' | 'company';

export type SubjectStatus = 'en_attente' | 'valide' | 'accepte_sous_reserve' | 'refuse' | 'expire';

export type ReviewDecision = 'valide' | 'accepte_sous_reserve' | 'refuse';

export type WishStatus = 'en_attente' | 'accepte' | 'refuse';

export type PfeAssignmentStatus = 'en_cours' | 'memoire_soumis' | 'soutenance_planifiee' | 'valide' | 'refuse';

export type MeetingType = 'presentiel' | 'visio';

export type ProgressReportStatus = 'a_faire' | 'en_cours' | 'termine';

export type DefenseStatus = 'scheduled' | 'done' | 'postponed';

export type DefenseResult = 'admitted' | 'corrections_required' | 'not_admitted';

export type AcademicYearStatus = 'active' | 'cloturee';

export type CompanyReportStatus = 'en_attente' | 'resolu' | 'rejete';

export type TiebreakChoice = 'president' | 'member' | 'average';



export interface Profile {
  id: number;
  role: UserRole;
  full_name: string;
  email: string;
  avatar_url: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  teacher?: Teacher;
  student?: Student;
  company?: Company;
}

export interface Department {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Domain {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Speciality {
  id: number;
  name: string;
  code: string;
  year_type: YearType;
  department_id: number | null;
  created_at: string;
  updated_at: string;
  department?: Department;
}

export interface AcademicYear {
  id: number;
  label: string;
  status: AcademicYearStatus;
  submission_open_at: string | null;
  submission_close_at: string | null;
  max_wishes: number;
  created_at: string;
  updated_at: string;
}

export interface Promotion {
  id: number;
  label: string;
  academic_year_id: number;
  created_at: string;
  updated_at: string;
}

export interface Teacher {
  id: number;
  profile_id: number;
  grade: TeacherGrade | null;
  department_id: number | null;
  availability_status: AvailabilityStatus;
  unavailable_until: string | null;
  created_at: string;
  updated_at: string;

  profile?: Profile;
  department?: Department;
  domaines?: Domain[];
}

export interface Student {
  id: number;
  profile_id: number;
  student_number: string | null;
  speciality_id: number | null;
  level: string | null;
  promotion_id: number | null;
  created_at: string;
  updated_at: string;

  profile?: Profile;
  speciality?: Speciality;
  promotion?: Promotion;
}

export interface Company {
  id: number;
  profile_id: number;
  company_name: string | null;
  sector: string | null;
  description: string | null;
  logo_url: string | null;
  contact_email: string | null;
  contact_phone: string | null;
  website: string | null;
  is_verified: boolean;
  created_at: string;
  updated_at: string;

  profile?: Profile;
}

export interface PfeSubject {
  id: number;
  title: string;
  description: string;
  group_type: GroupType;
  proposer_id: number;
  proposer_role: ProposerRole;
  company_id: number | null;
  academic_year_id: number;
  validator1_id: number | null;
  validator2_id: number | null;
  validator1_decision: ReviewDecision | null;
  validator2_decision: ReviewDecision | null;
  validator1_comment: string | null;
  validator2_comment: string | null;
  status: SubjectStatus;
  co_supervisor_id: number | null;
  pre_assigned_student_ids: string | null;
  is_assigned: boolean;
  created_at: string;
  updated_at: string;

  proposer?: Profile;
  company?: Company;
  validator1?: Teacher;
  validator2?: Teacher;
  co_supervisor?: Teacher;
  domains?: Domain[];
}

export interface Wish {
  id: number;
  student_id: number;
  subject_id: number;
  academic_year_id: number;
  status: WishStatus;
  created_at: string;
  updated_at: string;

  student?: Student;
  subject?: PfeSubject;
  academic_year?: AcademicYear;

  student_name?: string;
  student_specialty?: string;
  subject_title?: string;
}

export interface PfeAssignment {
  id: number;
  pfe_code: string;
  subject_id: number;
  academic_year_id: number;
  student_id: number;
  student2_id: number | null;
  student3_id: number | null;
  supervisor_id: number;
  co_supervisor_id: number | null;
  memoire_url: string | null;
  status: PfeAssignmentStatus;
  created_at: string;
  updated_at: string;

  subject?: PfeSubject;
  academic_year?: AcademicYear;
  student?: Student;
  student2?: Student;
  student3?: Student;
  supervisor?: Teacher;
  co_supervisor?: Teacher;

  subject_title?: string;
}

export interface PfeProgressReport {
  id: number;
  assignment_id: number;
  meeting_date: string;
  duration: number;
  meeting_type: MeetingType;
  topics: string;
  status: ProgressReportStatus;
  observation: string | null;
  created_at: string;
  updated_at: string;

  assignment?: PfeAssignment;
}

export interface DefenseJury {
  id: number;
  assignment_id: number;
  president_id: number;
  member_id: number;
  president_confirmed: boolean;
  member_confirmed: boolean;
  president_wants_printed: boolean;
  member_wants_printed: boolean;
  created_at: string;
  updated_at: string;

  assignment?: PfeAssignment;
  president?: Teacher;
  member?: Teacher;
}

export interface Defense {
  id: number;
  assignment_id: number;
  jury_id: number;
  scheduled_at: string | null;
  room: string | null;
  defense_deadline: string | null;
  status: DefenseStatus;
  result: DefenseResult | null;
  final_grade: number | null;
  created_at: string;
  updated_at: string;

  assignment?: PfeAssignment;
  jury?: DefenseJury;
}

export type ArchiveDecision = 'archivable' | 'minor_corrections' | 'major_corrections';

export interface JuryGrade {
  id: number;
  defense_id: number;
  jury_member_id: number;
  criterion1: number | null;
  criterion2: number | null;
  criterion3: number | null;
  criterion4: number | null;
  archive_decision: ArchiveDecision | null;
  total?: number;
  created_at: string;
  updated_at: string;

  defense?: Defense;
  jury_member?: Teacher;
}

export interface SupervisorEvaluation {
  id: number;
  pfe_assignment_id: number;
  evaluator_id: number;
  criterion5: number | null;
  created_at: string;
  updated_at: string;

  assignment?: PfeAssignment;
  evaluator?: Teacher;
}

export interface CompanyReport {
  id: number;
  company_id: number;
  submitted_by: number;
  correction_type: string;
  description: string;
  requested_value: string;
  status: CompanyReportStatus;
  resolved_at: string | null;
  created_at: string;
  updated_at: string;

  company?: Company;
}

export interface Notification {
  id: number;
  recipient_id: number;
  type: string;
  payload: string;
  message?: string;
  read_at: string | null;
  created_at: string;
}

export interface AuditLog {
  id: number;
  actor_id: number;
  action: string;
  entity: string;
  entity_id: number | null;
  metadata: string | null;
  created_at: string;
}



export interface SessionUser {
  id: number;
  email: string;
  role: UserRole;
  full_name: string;
  avatar_url: string | null;
}

export interface AuthResult {
  token: string;
  profile: Profile;
}



export interface AdminDashboard {
  total_users: number;
  total_teachers: number;
  total_students: number;
  total_companies: number;
  total_subjects: number;
  total_pfes: number;
  pending_subjects: number;
  validated_subjects: number;
  refused_subjects: number;
  rejected_subjects: number;
  assigned_subjects: number;
  active_pfes: number;
  defended_pfes: number;
  total_assignments: number;
  total_defenses: number;
  scheduled_defenses: number;
  done_defenses: number;
  total_reports: number;
  timeline: {
    labels: string[];
    soumis_memoire: number[];
    avec_sujet: number[];
    sans_sujet: number[];
  };
}



export interface ValidatorRecommendation {
  teacher: Teacher;
  score: number;
  matching_domains: Domain[];
}




export type IconComponent = new (...args: any[]) => any;
