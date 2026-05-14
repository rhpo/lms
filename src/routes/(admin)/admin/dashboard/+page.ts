import { admin } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const data = await admin.dashboard() as Record<string, unknown>;
    return {
      stats: (data?.stats as Record<string, number>) ?? {},
      activeYear: data?.active_year ?? null,
      recentAuditLogs: (data?.recent_audit_logs as unknown[]) ?? [],
    };
  } catch {
    return { stats: {}, activeYear: null, recentAuditLogs: [] };
  }
}
