import { company } from '$lib/api';
import type { PfeSubject, PfeAssignment, Notification } from '$lib/types';

export async function load() {
  try {
    const [subjects, pfes, notifications] = await Promise.all([
      company.listSubjects(),
      company.listSupervisedPFEs(),
      company.listNotifications(),
    ]);
    return {
      subjects: (subjects ?? []) as PfeSubject[],
      pfes: (pfes ?? []) as PfeAssignment[],
      notifications: (notifications ?? []) as Notification[],
    };
  } catch {
    return {
      subjects: [] as PfeSubject[],
      pfes: [] as PfeAssignment[],
      notifications: [] as Notification[],
    };
  }
}
