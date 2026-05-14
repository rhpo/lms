import { company } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const data = await company.dashboard() as Record<string, unknown>;
    return {
      supervisedCount: (data?.supervised_count as number) ?? 0,
      proposedCount: (data?.proposed_count as number) ?? 0,
      activeReportsCount: (data?.active_reports_count as number) ?? 0,
      unreadCount: (data?.unread_count as number) ?? 0,
      companyStatus: (data?.company_status as string) ?? 'en_attente',
    };
  } catch {
    return {
      supervisedCount: 0, proposedCount: 0, activeReportsCount: 0,
      unreadCount: 0, companyStatus: 'en_attente',
    };
  }
}
