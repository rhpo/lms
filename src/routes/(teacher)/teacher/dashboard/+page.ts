import { teacher } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const data = await teacher.dashboard() as Record<string, unknown>;
    return {
      supervisedCount: (data?.supervised_count as number) ?? 0,
      pendingValidationCount: (data?.pending_validation_count as number) ?? 0,
      proposedCount: (data?.proposed_count as number) ?? 0,
      upcomingJuryDuties: (data?.upcoming_jury_duties as any[]) ?? [],
      availabilityStatus: (data?.availability_status as string) ?? 'disponible',
      unavailableUntil: (data?.unavailable_until as string) ?? null,
      unreadCount: (data?.unread_count as number) ?? 0,
    };
  } catch {
    return {
      supervisedCount: 0, pendingValidationCount: 0, proposedCount: 0,
      upcomingJuryDuties: [], availabilityStatus: 'disponible',
      unavailableUntil: null, unreadCount: 0,
    };
  }
}
