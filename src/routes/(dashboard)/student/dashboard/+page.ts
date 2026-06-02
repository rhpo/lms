import { student } from '$lib/api';
import type { PfeAssignment, Wish, Notification } from '$lib/types';

export async function load() {
  try {
    const [pfe, wishes, notifications] = await Promise.all([
      student.getMyPFE(),
      student.listWishes(),
      student.listNotifications(),
    ]);
    return {
      currentPfe: pfe ?? null,
      wishes: wishes ?? [],
      notifications: notifications ?? [],
      yearId: null as string | null,
    };
  } catch {
    return {
      currentPfe: null as PfeAssignment | null,
      wishes: [] as Wish[],
      notifications: [] as Notification[],
      yearId: null as string | null,
    };
  }
}
