import type {
  Profile, Teacher, Student, Company, Department, Domain, Speciality, AcademicYear,
  Promotion, PfeSubject, Wish, PfeAssignment, PfeProgressReport,
  DefenseJury, Defense, JuryGrade, SupervisorEvaluation, CompanyReport,
  Notification, AuditLog, SessionUser, AuthResult, AdminDashboard,
  ValidatorRecommendation,
  TeacherGrade, ReviewDecision, GroupType, MeetingType,
} from './types';



const TOKEN_KEY = 'pfe_token';
const API_URL = '';

export function getToken(): string | null {
  if (typeof localStorage === 'undefined') return null;
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY);
}



const API_BASE = '/api';

interface ApiResponse<T = unknown> {
  success: boolean;
  data?: T;
  error?: string;
}

async function request<T = unknown>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = headers['Content-Type'] ?? 'application/json';
  }

  const res = await fetch(`${API_URL}${API_BASE}${path}`, { ...options, headers });

  if (!res.ok) {
    let errMsg = `HTTP ${res.status}`;
    try {
      const json: ApiResponse = await res.json();
      errMsg = json.error ?? errMsg;
    } catch { /* ignore parse error */ }
    throw new Error(errMsg);
  }

  const json: ApiResponse<T> = await res.json();
  return json.data as T;
}

function get<T>(path: string) {
  return request<T>(path, { method: 'GET' });
}

function post<T>(path: string, body?: unknown) {
  return request<T>(path, {
    method: 'POST',
    body: body instanceof FormData ? body : JSON.stringify(body),
  });
}

function patch<T>(path: string, body?: unknown) {
  return request<T>(path, { method: 'PATCH', body: JSON.stringify(body) });
}

function del<T = void>(path: string) {
  return request<T>(path, { method: 'DELETE' });
}

export async function downloadBlob(path: string): Promise<Blob> {
  const token = getToken();
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(`${API_BASE}${path}`, { headers });
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return await res.blob();
}



export const auth = {
  devLogin: (email: string) => post<AuthResult>('/auth/dev-login', { email }),
  me: () => get<Profile>('/auth/me'),
  logout: () => post<void>('/auth/logout'),
  listVerifiedCompanies: () => get<Company[]>('/auth/companies'),
  registerCompany: (body: {
    full_name: string;
    email: string;
    position: string;
    phone: string;
    company_id?: number;
    company_name?: string;
    sector?: string;
    description?: string;
    contact_email?: string;
    contact_phone?: string;
  }) => post<AuthResult>('/auth/register-company', body),
};



export const ref = {
  domains: () => get<Domain[]>('/ref/domains'),
  specialities: () => get<Speciality[]>('/ref/specialities'),
  departments: () => get<Department[]>('/ref/departments'),
  allUsers: () => get<Profile[]>('/accounts/users'),
};



export const admin = {

  dashboard: () => get<AdminDashboard>('/admin/dashboard'),


  listUsers: () => get<Profile[]>('/admin/accounts/users'),
  createUser: (body: { role: string; full_name: string; email: string }) =>
    post<Profile>('/admin/accounts/users', body),
  createTeacher: (body: { full_name: string; email: string; grade: string; department_id: number }) =>
    post<Profile>('/admin/accounts/users/enseignant', body),
  createStudent: (body: { full_name: string; email: string; student_number: string; speciality_id?: number; level?: string; promotion_id?: number }) =>
    post<Profile>('/admin/accounts/users/etudiant', body),
  getUser: (id: number) => get<Profile>(`/admin/accounts/users/${id}`),
  updateUser: (id: number, body: {
    full_name?: string; email?: string;
    grade?: string; department_id?: number | null;
    domain_ids?: number[];
    student_number?: string; level?: string;
    speciality_id?: number | null; promotion_id?: number | null;
  }) => patch<Profile>(`/admin/accounts/users/${id}`, body),
  updateUserAvatar: (id: number, formData: FormData) =>
    request<{ url: string }>(`/admin/accounts/users/${id}/avatar`, { method: 'POST', body: formData }),
  userAction: (id: number, action: string) =>
    post<{ message: string }>(`/admin/accounts/users/${id}/action`, { action }),
  importUsersCSV: (csvData: string, csvType: string, replaceMode: boolean) =>
    post<{ message: string }>('/admin/accounts/users/import-csv', { csv_data: csvData, csv_type: csvType, replace: replaceMode }),


  listCompanies: () => get<Company[]>('/admin/accounts/companies'),
  companyAction: (id: number, action: string, payload?: Record<string, unknown>) =>
    post<{ message: string }>(`/admin/accounts/companies/${id}/action`, { action, ...payload }),


  listReports: () => get<CompanyReport[]>('/admin/reports'),
  reportAction: (id: number, action: string) =>
    post<{ message: string }>(`/admin/reports/${id}/action`, { action }),


  listSubjects: () => get<PfeSubject[]>('/admin/subjects'),
  getSubject: (id: number) => get<PfeSubject>(`/admin/subjects/${id}`),
  subjectAction: (id: number, action: string, payload?: Record<string, unknown>) =>
    post<{ message: string }>(`/admin/subjects/${id}/action`, { action, ...payload }),


  listAssignments: () => get<PfeAssignment[]>('/admin/pfe'),
  getAssignment: (id: number) => get<PfeAssignment>(`/admin/pfe/${id}`),
  assignmentAction: (id: number, action: string, body?: Record<string, unknown>) =>
    post<{ message: string }>(`/admin/pfe/${id}/action`, { action, ...body }),
  recommendCoSupervisor: (assignmentId: number) =>
    get<{ recommended: ValidatorRecommendation[]; assignment_id: number; subject_domains: Domain[] }>(
      `/admin/pfe/recommend-co-supervisor?assignment_id=${assignmentId}`
    ),


  listDefenses: () => get<Defense[]>('/admin/defenses'),
  createDefense: (body: {
    assignment_id: number;
    president_id: number;
    member_id: number;
    scheduled_at: string;
    room: string;
  }) => post<Defense>('/admin/defenses', body),
  getDefense: (id: number) => get<Defense>(`/admin/defenses/${id}`),
  recommendJury: (pfeId: number) =>
    get<{ recommended: ValidatorRecommendation[]; pfe_id: number; subject_domains: Domain[] }>(`/admin/defenses/recommend-jury?pfe_id=${pfeId}`),
  submitGrade: (defenseId: number, body: {
    criterion1: number; criterion2: number; criterion3: number; criterion4: number;
  }) => post<{ message: string }>(`/admin/defenses/${defenseId}/submit-grade`, body),
  resolveGrade: (defenseId: number, body: {
    choice: string;
    criterion1?: number; criterion2?: number; criterion3?: number; criterion4?: number;
    grades?: Record<string, number>;
  }) => post<{ message: string }>(`/admin/defenses/${defenseId}/resolve-grade`, body),
  confirmJury: (defenseId: number) =>
    post<{ message: string }>(`/admin/defenses/${defenseId}/confirm-jury`),
  declineJury: (defenseId: number) =>
    post<{ message: string }>(`/admin/defenses/${defenseId}/decline-jury`),


  listDeadlines: () => get<AcademicYear[]>('/admin/settings/deadlines'),
  updateDeadlines: (body: {
    submission_open_at: string; submission_close_at: string; max_wishes: number;
  }) => post<{ message: string }>('/admin/settings/deadlines', body),


  listSpecialities: () => get<Speciality[]>('/admin/settings/specialities'),
  createSpeciality: (body: { name: string; code: string; year_type: string; department_id?: number }) =>
    post<Speciality>('/admin/settings/specialities', body),
  deleteSpeciality: (id: number) => del(`/admin/settings/specialities/${id}`),


  listDomains: () => get<Domain[]>('/admin/settings/domains'),
  createDomain: (body: { name: string }) =>
    post<Domain>('/admin/settings/domains', body),
  deleteDomain: (id: number) => del(`/admin/settings/domains/${id}`),


  listDepartments: () => get<Department[]>('/admin/settings/departments'),
  createDepartment: (body: { name: string }) =>
    post<Department>('/admin/settings/departments', body),
  deleteDepartment: (id: number) => del(`/admin/settings/departments/${id}`),


  listPromotions: () => get<Promotion[]>('/admin/settings/promotions'),
  createPromotion: (body: { label: string; academic_year_id: number }) =>
    post<Promotion>('/admin/settings/promotions', body),
  deletePromotion: (id: number) => del(`/admin/settings/promotions/${id}`),


  listAcademicYears: () => get<AcademicYear[]>('/admin/settings/academic-years'),
  createAcademicYear: (body: { label: string; status: string; max_wishes?: number }) =>
    post<AcademicYear>('/admin/settings/academic-years', body),
  closeAcademicYear: (id: number) =>
    post<{ message: string }>(`/admin/settings/academic-years/${id}/close`),




  statistics: () => get<AdminDashboard>('/admin/statistics'),
  auditLog: () => get<AuditLog[]>('/admin/audit-log'),


  exportAffectations: () => get<PfeAssignment[]>('/admin/exports/affectations'),
  exportPlannings: () => get<Defense[]>('/admin/exports/plannings'),
  exportStatistics: () => get<AdminDashboard>('/admin/exports/statistiques'),
};



export const teacher = {
  dashboard: () => get<Record<string, unknown>>('/teacher/dashboard'),


  listProposedSubjects: () => get<PfeSubject[]>('/teacher/proposed-subjects'),
  createProposedSubject: (body: {
    title: string; description: string; group_type: GroupType;
    domain_ids?: number[];
  }) => post<PfeSubject>('/teacher/proposed-subjects', body),
  getProposedSubject: (id: number) => get<PfeSubject>(`/teacher/proposed-subjects/${id}`),
  updateProposedSubject: (id: number, body: Partial<PfeSubject>) =>
    patch<PfeSubject>(`/teacher/proposed-subjects/${id}`, body),
  resubmitProposedSubject: (id: number, body: {
    title: string; description: string; group_type: GroupType; domain_ids?: number[];
  }) => post<{ message: string }>(`/teacher/proposed-subjects/${id}/resubmit`, body),


  listCandidats: (subjectId: number) =>
    get<Wish[]>(`/teacher/proposed-subjects/${subjectId}/candidats`),
  acceptCandidat: (subjectId: number, body: { student_ids: number[] }) =>
    post<{ message: string }>(`/teacher/proposed-subjects/${subjectId}/candidats`, body),


  listSubjectsToValidate: () => get<PfeSubject[]>('/teacher/subjects-to-validate'),
  getSubjectToValidate: (id: number) => get<PfeSubject>(`/teacher/subjects-to-validate/${id}`),
  validateSubject: (id: number, body: { decision: ReviewDecision; comment?: string }) =>
    post<{ message: string }>(`/teacher/subjects-to-validate/${id}`, body),


  listSupervisedPFEs: () => get<PfeAssignment[]>('/teacher/supervised-pfes'),
  getSupervisedPFE: (id: number) => get<PfeAssignment>(`/teacher/supervised-pfes/${id}`),
  listMeetings: (id: number) => get<PfeProgressReport[]>(`/teacher/supervised-pfes/${id}/meetings`),
  addMeeting: (id: number, body: {
    meeting_date: string; duration: number; meeting_type: MeetingType;
    topics: string; status?: string; observation?: string;
  }) => post<PfeProgressReport>(`/teacher/supervised-pfes/${id}/meetings`, body),
  getEvaluation: (id: number) => get<SupervisorEvaluation | null>(`/teacher/supervised-pfes/${id}/evaluation`),
  submitEvaluation: (id: number, body: { criterion5: number }) =>
    post<{ message: string }>(`/teacher/supervised-pfes/${id}/evaluation`, body),


  listJuryDuties: () => get<Defense[]>('/teacher/jury-duties'),
  getJuryDuty: (id: number) => get<Defense>(`/teacher/jury-duties/${id}`),


  updateAvailability: (body: {
    availability_status: string; unavailable_until?: string;
  }) => post<{ message: string }>('/teacher/availability', body),


  listNotifications: () => get<Notification[]>('/teacher/notifications'),


  getGradeContext: (defenseId: number) =>
    get<{
      my_role: 'president' | 'member';
      my_grade: JuryGrade | null;
      member_grade: JuryGrade | null;
      supervisor_eval: SupervisorEvaluation | null;
      member_submitted: boolean;
      supervisor_submitted: boolean;
      final_grade_set: boolean;
    }>(`/teacher/jury-duties/${defenseId}/grade-context`),
  submitGrade: (defenseId: number, body: {
    criterion1: number; criterion2: number; criterion3: number; criterion4: number;
    archive_decision: string;
  }) => post<{ message: string }>(`/teacher/jury-duties/${defenseId}/grade`, body),
  submitFinalGrade: (defenseId: number, body: {
    choice: 'member' | 'new';
    criterion1?: number; criterion2?: number; criterion3?: number; criterion4?: number;
    archive_decision: string;
  }) => post<{ message: string }>(`/teacher/jury-duties/${defenseId}/final-grade`, body),
};



export const student = {
  settings: () => get<{ max_wishes: number; submission_open_at: string | null; submission_close_at: string | null }>('/student/settings'),
  dashboard: () => get<Record<string, unknown>>('/student/dashboard'),
  listCatalogue: () => get<PfeSubject[]>('/student/catalogue'),
  getCatalogueSubject: (id: number) => get<PfeSubject>(`/student/catalogue/${id}`),
  listWishes: () => get<Wish[]>('/student/wishes'),
  createWish: (body: { subject_id: number }) => post<Wish>('/student/wishes', body),
  deleteWish: (id: number) => del(`/student/wishes/${id}`),
  getMyPFE: () => get<PfeAssignment | null>('/student/my-pfe'),
  listMyMeetings: () => get<PfeProgressReport[]>('/student/my-pfe/meetings'),
  addMyMeeting: (body: {
    meeting_date: string; duration: number; meeting_type: MeetingType;
    topics: string; status?: string; observation?: string;
  }) => post<PfeProgressReport>('/student/my-pfe/meetings', body),
  updateMyMeeting: (id: number, body: { status: string }) =>
    patch<{ message: string }>(`/student/my-pfe/meetings/${id}`, body),
  submitMemoire: (body: { memoire_url: string }) =>
    post<{ message: string }>('/student/my-pfe/memoire', body),
  getSoutenance: () => get<{ has_soutenance: boolean; defense?: Defense; jury?: DefenseJury } | null>('/student/soutenance'),
  listNotifications: () => get<Notification[]>('/student/notifications'),
};



export const company = {
  dashboard: () => get<Record<string, unknown>>('/company/dashboard'),
  listSubjects: () => get<PfeSubject[]>('/company/subjects'),
  createSubject: (body: {
    title: string; description: string; group_type: GroupType;
    domain_ids?: number[];
  }) => post<PfeSubject>('/company/subjects', body),
  getSubject: (id: number) => get<PfeSubject>(`/company/subjects/${id}`),
  updateSubject: (id: number, body: Partial<PfeSubject>) =>
    patch<PfeSubject>(`/company/subjects/${id}`, body),
  listCandidats: (subjectId: number) =>
    get<Wish[]>(`/company/subjects/${subjectId}/candidats`),
  acceptCandidat: (subjectId: number, body: { student_ids: number[] }) =>
    post<{ message: string }>(`/company/subjects/${subjectId}/candidats`, body),
  listSupervisedPFEs: () => get<PfeAssignment[]>('/company/supervised-pfes'),
  getSupervisedPFE: (id: number) => get<PfeAssignment>(`/company/supervised-pfes/${id}`),
  listMeetings: (id: number) => get<PfeProgressReport[]>(`/company/supervised-pfes/${id}/meetings`),
  addMeeting: (id: number, body: {
    meeting_date: string; duration: number; meeting_type: MeetingType;
    topics: string; status?: string; observation?: string;
  }) => post<PfeProgressReport>(`/company/supervised-pfes/${id}/meetings`, body),
  getEvaluation: (id: number) => get<SupervisorEvaluation | null>(`/company/supervised-pfes/${id}/evaluation`),
  submitEvaluation: (id: number, body: { criterion5: number }) =>
    post<{ message: string }>(`/company/supervised-pfes/${id}/evaluation`, body),
  listReports: () => get<CompanyReport[]>('/company/reports'),
  createReport: (body: {
    correction_type: string; description: string; requested_value: string;
  }) => post<CompanyReport>('/company/reports', body),
  listNotifications: () => get<Notification[]>('/company/notifications'),
};



export const notifications = {
  list: () => get<Notification[]>('/notifications'),
  unreadCount: () => get<number>('/notifications/unread-count'),
  markRead: (id: number) => post<{ message: string }>(`/notifications/${id}/read`),
  markAllRead: () => post<{ message: string }>('/notifications/read-all'),
};



export const upload = {
  avatar: (formData: FormData) =>
    request<{ url: string }>('/upload/avatar', { method: 'POST', body: formData }),
  companyLogo: (formData: FormData) =>
    request<{ url: string }>('/upload/company-logo', { method: 'POST', body: formData }),
  memoire: (formData: FormData) =>
    request<{ url: string }>('/upload/memoire', { method: 'POST', body: formData }),
};



export const shared = {
  domains: ref.domains,
  specialities: ref.specialities,
  accounts: ref.allUsers,
  markNotificationRead: notifications.markRead,
  markAllNotificationsRead: notifications.markAllRead,
};
