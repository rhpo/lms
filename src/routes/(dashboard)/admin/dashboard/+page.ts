import { admin } from '$lib/api';
import type { AdminDashboard, AuditLog } from '$lib/types';

export async function load() {
  try {
    const data = await admin.dashboard();
    const s = (data ?? {}) as AdminDashboard & {
      active_year?: any;
      recent_audit_logs?: AuditLog[];
      pending_companies?: number;
    };
    return {
      stats: {
        totalStudents: s.total_students ?? 0,
        totalTeachers: s.total_teachers ?? 0,
        totalCompanies: s.total_companies ?? 0,
        pendingCompanies: s.pending_companies ?? 0,
        totalSubjects: s.total_subjects ?? 0,
        pendingSubjects: s.pending_subjects ?? 0,
        validatedSubjects: s.validated_subjects ?? 0,
        rejectedSubjects: s.rejected_subjects ?? 0,
        assignedSubjects: s.assigned_subjects ?? 0,
        totalAssignments: s.total_assignments ?? 0,
        totalDefenses: s.total_defenses ?? 0,
      },
      timeline: s.timeline ?? { labels: [], soumis_memoire: [], avec_sujet: [], sans_sujet: [] },
      activeYear: s.active_year ?? null,
      recentAuditLogs: s.recent_audit_logs ?? [],
    };
  } catch {
    return {
      stats: {} as Record<string, number>,
      activeYear: null,
      recentAuditLogs: [] as AuditLog[],
      timeline: { labels: [] as string[], soumis_memoire: [] as number[], avec_sujet: [] as number[], sans_sujet: [] as number[] },
    };
  }
}
