import { auth, teacher } from '$lib/api';
import type { Defense } from '$lib/types';

export async function load() {
  try {
    const [proposed, supervised, toValidate, juryDuties, notifications, me] = await Promise.all([
      teacher.listProposedSubjects(),
      teacher.listSupervisedPFEs(),
      teacher.listSubjectsToValidate(),
      teacher.listJuryDuties(),
      teacher.listNotifications(),
      auth.me(),
    ]);

    return {
      supervisedCount: supervised?.length ?? 0,
      pendingValidationCount: toValidate?.length ?? 0,
      proposedCount: proposed?.length ?? 0,
      upcomingJuryDuties: (juryDuties ?? []) as Defense[],
      availabilityStatus: (me as any)?.status ?? 'disponible',
      unavailableUntil: (me as any)?.unavailable_until ?? null,
      unreadCount: notifications?.filter((n) => !n.read_at).length ?? 0,
    };
  } catch {
    return {
      supervisedCount: 0, pendingValidationCount: 0, proposedCount: 0,
      upcomingJuryDuties: [] as Defense[], availabilityStatus: 'disponible',
      unavailableUntil: null as string | null, unreadCount: 0,
    };
  }
}
