import { student } from '$lib/api';
import type { PfeAssignment, PfeProgressReport } from '$lib/types';

export async function load() {
  const [meetings, pfe] = await Promise.all([
    student.listMyMeetings().catch(() => [] as PfeProgressReport[]),
    student.getMyPFE().catch(() => null),
  ]);
  return {
    meetings: meetings ?? [],
    pfe: pfe ?? null,
  };
}
