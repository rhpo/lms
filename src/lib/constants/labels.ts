import type {
  UserRole, TeacherGrade, AvailabilityStatus, YearType, GroupType,
  SubjectStatus, ReviewDecision, WishStatus, PfeAssignmentStatus,
  MeetingType, ProgressReportStatus, DefenseStatus, DefenseResult,
  AcademicYearStatus, CompanyReportStatus,
} from '$lib/types';



export const ROLE_LABELS: Record<UserRole, string> = {
  admin: 'Administrateur',
  teacher: 'Enseignant',
  student: 'Étudiant',
  company: 'Entreprise',
};

export const TEACHER_GRADE_LABELS: Record<TeacherGrade, string> = {
  assistant: 'Assistant',
  mab: 'MAB',
  maa: 'MAA',
  mcb: 'MCB',
  mca: 'MCA',
  professeur: 'Professeur',
};

export const AVAILABILITY_LABELS: Record<AvailabilityStatus, string> = {
  disponible: 'Disponible',
  indisponible: 'Indisponible',
  indisponible_jusqu_au: 'Indisponible (temporaire)',
};

export const YEAR_TYPE_LABELS: Record<YearType, string> = {
  licence: 'Licence',
  master: 'Master',
};

export const GROUP_TYPE_LABELS: Record<GroupType, string> = {
  monome: 'Monôme',
  binome: 'Binôme',
  trinome: 'Trinôme',
};

export const SUBJECT_STATUS_LABELS: Record<SubjectStatus, string> = {
  en_attente: 'En attente',
  valide: 'Validé',
  accepte_sous_reserve: 'Accepté sous réserve',
  refuse: 'Refusé',
  expire: 'Expiré',
};

export const REVIEW_DECISION_LABELS: Record<ReviewDecision, string> = {
  valide: 'Validé',
  accepte_sous_reserve: 'Accepté sous réserve',
  refuse: 'Refusé',
};

export const REVIEW_DECISION_VARIANTS: Record<ReviewDecision, 'info' | 'success' | 'warning' | 'danger'> = {
  valide: 'success',
  accepte_sous_reserve: 'warning',
  refuse: 'danger',
};

export const WISH_STATUS_LABELS: Record<WishStatus, string> = {
  en_attente: 'En attente',
  accepte: 'Accepté',
  refuse: 'Refusé',
};

export const PFE_ASSIGNMENT_STATUS_LABELS: Record<PfeAssignmentStatus, string> = {
  en_cours: 'En cours',
  memoire_soumis: 'Mémoire soumis',
  soutenance_planifiee: 'Soutenance planifiée',
  valide: 'Validé',
  refuse: 'Refusé',
};

export const MEETING_TYPE_LABELS: Record<MeetingType, string> = {
  presentiel: 'Présentiel',
  visio: 'Visioconférence',
};

export const PROGRESS_REPORT_STATUS_LABELS: Record<ProgressReportStatus, string> = {
  a_faire: 'À faire',
  en_cours: 'En cours',
  termine: 'Terminé',
};

export const DEFENSE_STATUS_LABELS: Record<DefenseStatus, string> = {
  scheduled: 'Planifiée',
  done: 'Passée',
  postponed: 'Reportée',
};

export const DEFENSE_RESULT_LABELS: Record<DefenseResult, string> = {
  admitted: 'Admis',
  corrections_required: 'Corrections requises',
  not_admitted: 'Non admis',
};

export const ACADEMIC_YEAR_STATUS_LABELS: Record<AcademicYearStatus, string> = {
  active: 'Active',
  cloturee: 'Clôturée',
};

export const COMPANY_REPORT_STATUS_LABELS: Record<CompanyReportStatus, string> = {
  en_attente: 'En attente',
  resolu: 'Résolu',
  rejete: 'Rejeté',
};



type BadgeVariant = 'info' | 'success' | 'warning' | 'danger';

export const SUBJECT_STATUS_VARIANTS: Record<SubjectStatus, BadgeVariant> = {
  en_attente: 'warning',
  valide: 'success',
  accepte_sous_reserve: 'info',
  refuse: 'danger',
  expire: 'danger',
};

export const DEFENSE_STATUS_VARIANTS: Record<DefenseStatus, BadgeVariant> = {
  scheduled: 'info',
  done: 'success',
  postponed: 'warning',
};

export const DEFENSE_RESULT_VARIANTS: Record<DefenseResult, BadgeVariant> = {
  admitted: 'success',
  corrections_required: 'warning',
  not_admitted: 'danger',
};

export const AVAILABILITY_VARIANTS: Record<AvailabilityStatus, BadgeVariant> = {
  disponible: 'success',
  indisponible: 'danger',
  indisponible_jusqu_au: 'warning',
};

export const PFE_ASSIGNMENT_STATUS_VARIANTS: Record<PfeAssignmentStatus, BadgeVariant> = {
  en_cours: 'info',
  memoire_soumis: 'warning',
  soutenance_planifiee: 'info',
  valide: 'success',
  refuse: 'danger',
};

export const WISH_STATUS_VARIANTS: Record<WishStatus, BadgeVariant> = {
  en_attente: 'warning',
  accepte: 'success',
  refuse: 'danger',
};

export const YEAR_TYPE_VARIANTS: Record<YearType, BadgeVariant> = {
  licence: 'info',
  master: 'warning',
};



export const YEAR_TYPE_OPTIONS = [
  { value: 'licence' as const, label: 'Licence' },
  { value: 'master' as const, label: 'Master' },
] as const;

export const TEACHER_GRADE_OPTIONS = [
  { value: 'assistant' as const, label: 'Assistant' },
  { value: 'mab' as const, label: 'MAB' },
  { value: 'maa' as const, label: 'MAA' },
  { value: 'mcb' as const, label: 'MCB' },
  { value: 'mca' as const, label: 'MCA' },
  { value: 'professeur' as const, label: 'Professeur' },
] as const;

export const GROUP_TYPE_OPTIONS = [
  { value: 'monome' as const, label: 'Monôme' },
  { value: 'binome' as const, label: 'Binôme' },
  { value: 'trinome' as const, label: 'Trinôme' },
] as const;

export const MEETING_TYPE_OPTIONS = [
  { value: 'presentiel' as const, label: 'Présentiel' },
  { value: 'visio' as const, label: 'Visioconférence' },
] as const;



export const NOTIFICATION_TYPE_LABELS: Record<string, string> = {
  validation_requise: "Validation requise",
  affectation: "Affectation",
  jury: "Jury",
  disponibilite: "Disponibilité",
};




export function formatDate(iso: string | null | undefined): string {
  if (!iso) return '-';
  return new Date(iso).toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  });
}

export function formatDateTime(iso: string | null | undefined): string {
  if (!iso) return '-';
  return new Date(iso).toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}
