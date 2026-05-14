/**
 * Centralized API client for the Go backend.
 * All requests go through /api (proxied to http://localhost:8080 by Vite in dev).
 */
import type {
  PFESubject, Wish, PFEAssignment, PFEProgressReport, Notification, SessionUser
} from './types';

const API_BASE = '/api';
const TOKEN_KEY = 'pfe_token';

// ─── Token management ────────────────────────────────────────────────────────

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

// ─── Core fetch wrapper ──────────────────────────────────────────────────────

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

  // Don't set Content-Type for FormData (browser sets it with boundary)
  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = headers['Content-Type'] ?? 'application/json';
  }

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers });

  if (!res.ok) {
    let errMsg = `HTTP ${res.status}`;
    try {
      const json: ApiResponse = await res.json();
      errMsg = json.error ?? errMsg;
    } catch {
      // ignore parse error
    }
    throw new Error(errMsg);
  }

  const json: ApiResponse<T> = await res.json();
  return json.data as T;
}

export async function downloadBlob(path: string): Promise<Blob> {
  const token = getToken();
  const headers: Record<string, string> = {};
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  const res = await fetch(`${API_BASE}${path}`, { headers });
  if (!res.ok) {
    throw new Error(`HTTP ${res.status}`);
  }
  return await res.blob();
}

function get<T = any>(path: string) {
  return request<T>(path, { method: 'GET' });
}

function post<T = any>(path: string, body?: unknown) {
  return request<T>(path, {
    method: 'POST',
    body: body instanceof FormData ? body : JSON.stringify(body),
  });
}

function patch<T = any>(path: string, body?: unknown) {
  return request<T>(path, { method: 'PATCH', body: JSON.stringify(body) });
}

function del<T = void>(path: string) {
  return request<T>(path, { method: 'DELETE' });
}

// ─── Auth ────────────────────────────────────────────────────────────────────

export interface AuthResult {
  token: string;
  profile: SessionUser;
}

export const auth = {
  devLogin: (email: string) => post<AuthResult>('/auth/dev-login', { email }),
  me: () => get<SessionUser>('/auth/me'),
  logout: () => post('/auth/logout'),
};

// ─── Admin ───────────────────────────────────────────────────────────────────

export const admin = {
  dashboard: () => get('/admin/dashboard'),
  listUsers: (params?: Record<string, string>) => {
    const qs = params ? '?' + new URLSearchParams(params).toString() : '';
    return get(`/admin/accounts/users${qs}`);
  },
  createUser: (body: unknown) => post('/admin/accounts/users', body),
  getUser: (id: string) => get(`/admin/accounts/users/${id}`),
  updateUser: (id: string, body: unknown) => patch(`/admin/accounts/users/${id}`, body),
  userAction: (id: string, action: string, payload?: unknown) =>
    post(`/admin/accounts/users/${id}/action`, { action, ...payload }),
  importUsersCSV: (formData: FormData) =>
    request('/admin/accounts/users/import-csv', { method: 'POST', body: formData }),
  listCompanies: () => get('/admin/accounts/companies'),
  companyAction: (id: string, action: string) =>
    post(`/admin/accounts/companies/${id}/action`, { action }),
  listReports: () => get('/admin/reports'),
  reportAction: (id: string, action: string) =>
    post(`/admin/reports/${id}/action`, { action }),
  listSubjects: (params?: Record<string, string>) => {
    const qs = params ? '?' + new URLSearchParams(params).toString() : '';
    return get(`/admin/subjects${qs}`);
  },
  getSubject: (id: string) => get(`/admin/subjects/${id}`),
  subjectAction: (id: string, action: string, payload?: unknown) =>
    post(`/admin/subjects/${id}/action`, { action, ...payload }),
  listAssignments: () => get('/admin/pfe'),
  getAssignment: (id: string) => get(`/admin/pfe/${id}`),
  listDefenses: () => get('/admin/defenses'),
  createDefense: (body: unknown) => post('/admin/defenses', body),
  recommendJury: (assignmentId: string) =>
    get(`/admin/defenses/recommend-jury?assignment_id=${assignmentId}`),
  getDefense: (id: string) => get(`/admin/defenses/${id}`),
  submitGrade: (id: string, body: unknown) => post(`/admin/defenses/${id}/submit-grade`, body),
  resolveGrade: (id: string, body: unknown) => post(`/admin/defenses/${id}/resolve-grade`, body),
  confirmJury: (id: string) => post(`/admin/defenses/${id}/confirm-jury`),
  declineJury: (id: string) => post(`/admin/defenses/${id}/decline-jury`),
  listDeadlines: () => get('/admin/settings/deadlines'),
  updateDeadlines: (body: unknown) => post('/admin/settings/deadlines', body),
  listSpecialities: () => get('/admin/settings/specialities'),
  createSpeciality: (body: unknown) => post('/admin/settings/specialities', body),
  deleteSpeciality: (id: string) => del(`/admin/settings/specialities/${id}`),
  listDomains: () => get('/admin/settings/domains'),
  createDomain: (body: unknown) => post('/admin/settings/domains', body),
  deleteDomain: (id: string) => del(`/admin/settings/domains/${id}`),
  listPromotions: () => get('/admin/settings/promotions'),
  createPromotion: (body: unknown) => post('/admin/settings/promotions', body),
  deletePromotion: (id: string) => del(`/admin/settings/promotions/${id}`),
  listAcademicYears: () => get('/admin/settings/academic-years'),
  createAcademicYear: (body: unknown) => post('/admin/settings/academic-years', body),
  closeAcademicYear: (id: string) => post(`/admin/settings/academic-years/${id}/close`),
  statistics: () => get('/admin/statistics'),
  auditLog: () => get('/admin/audit-log'),
  exportAffectations: () => get('/admin/exports/affectations'),
  exportPlannings: () => get('/admin/exports/plannings'),
  exportStatistics: () => get('/admin/exports/statistiques'),
  sendEmail: (body: unknown) => post('/admin/send-email', body),
};

// ─── Teacher ─────────────────────────────────────────────────────────────────

export const teacher = {
  dashboard: () => get('/teacher/dashboard'),
  listProposedSubjects: () => get('/teacher/proposed-subjects'),
  createProposedSubject: (body: unknown) => post('/teacher/proposed-subjects', body),
  getProposedSubject: (id: string) => get(`/teacher/proposed-subjects/${id}`),
  updateProposedSubject: (id: string, body: unknown) =>
    patch(`/teacher/proposed-subjects/${id}`, body),
  listCandidats: (subjectId: string) =>
    get(`/teacher/proposed-subjects/${subjectId}/candidats`),
  acceptCandidat: (subjectId: string, body: unknown) =>
    post(`/teacher/proposed-subjects/${subjectId}/candidats`, body),
  listSubjectsToValidate: () => get('/teacher/subjects-to-validate'),
  getSubjectToValidate: (id: string) => get(`/teacher/subjects-to-validate/${id}`),
  validateSubject: (id: string, body: unknown) =>
    post(`/teacher/subjects-to-validate/${id}`, body),
  listSupervisedPFEs: () => get('/teacher/supervised-pfes'),
  getSupervisedPFE: (id: string) => get(`/teacher/supervised-pfes/${id}`),
  addMeeting: (id: string, body: unknown) =>
    post(`/teacher/supervised-pfes/${id}/meetings`, body),
  submitEvaluation: (id: string, body: unknown) =>
    post(`/teacher/supervised-pfes/${id}/evaluation`, body),
  listJuryDuties: () => get('/teacher/jury-duties'),
  getJuryDuty: (id: string) => get(`/teacher/jury-duties/${id}`),
  updateAvailability: (body: unknown) => post('/teacher/availability', body),
  listNotifications: () => get('/teacher/notifications'),
  submitGrade: (defenseId: string, body: unknown) =>
    post(`/admin/defenses/${defenseId}/submit-grade`, body),
};

// ─── Student ─────────────────────────────────────────────────────────────────

export const student = {
  dashboard: () => get<any>('/student/dashboard'),
  listCatalogue: (params?: Record<string, string>) => {
    const qs = params ? '?' + new URLSearchParams(params).toString() : '';
    return get<PFESubject[]>(`/student/catalogue${qs}`);
  },
  getCatalogueSubject: (id: string) => get<PFESubject>(`/student/catalogue/${id}`),
  listWishes: () => get<Wish[]>('/student/wishes'),
  createWish: (body: unknown) => post<Wish>('/student/wishes', body),
  deleteWish: (id: string) => del<void>(`/student/wishes/${id}`),
  getMyPFE: () => get<PFEAssignment | null>('/student/my-pfe'),
  listMyMeetings: () => get<PFEProgressReport[]>('/student/my-pfe/meetings'),
  addMyMeeting: (body: unknown) => post<PFEProgressReport>('/student/my-pfe/meetings', body),
  submitMemoire: (formData: FormData) =>
    request('/student/my-pfe/memoire', { method: 'POST', body: formData }),
  getSoutenance: () => get('/student/soutenance'),
  listNotifications: () => get('/student/notifications'),
};

// ─── Company ─────────────────────────────────────────────────────────────────

export const company = {
  dashboard: () => get('/company/dashboard'),
  listSubjects: () => get('/company/subjects'),
  createSubject: (body: unknown) => post('/company/subjects', body),
  getSubject: (id: string) => get(`/company/subjects/${id}`),
  updateSubject: (id: string, body: unknown) => patch(`/company/subjects/${id}`, body),
  listCandidats: (subjectId: string) => get(`/company/subjects/${subjectId}/candidats`),
  acceptCandidat: (subjectId: string, body: unknown) =>
    post(`/company/subjects/${subjectId}/candidats`, body),
  listSupervisedPFEs: () => get('/company/supervised-pfes'),
  getSupervisedPFE: (id: string) => get(`/company/supervised-pfes/${id}`),
  addMeeting: (id: string, body: unknown) =>
    post(`/company/supervised-pfes/${id}/meetings`, body),
  submitEvaluation: (id: string, body: unknown) =>
    post(`/company/supervised-pfes/${id}/evaluation`, body),
  listReports: () => get('/company/reports'),
  createReport: (body: unknown) => post('/company/reports', body),
  listNotifications: () => get('/company/notifications'),
};

// ─── Shared (all authenticated roles) ────────────────────────────────────────

export const shared = {
  domains: () => get('/ref/domains'),
  specialities: () => get('/ref/specialities'),
  markNotificationRead: (id: string) => post(`/notifications/${id}/read`),
  markAllNotificationsRead: () => post('/notifications/read-all'),
};

// ─── Upload ──────────────────────────────────────────────────────────────────

export const upload = {
  avatar: (formData: FormData) =>
    request('/upload/avatar', { method: 'POST', body: formData }),
  companyLogo: (formData: FormData) =>
    request('/upload/company-logo', { method: 'POST', body: formData }),
  memoire: (formData: FormData) =>
    request('/upload/memoire', { method: 'POST', body: formData }),
};
