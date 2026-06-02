import { teacher } from '$lib/api';
import type { LoadEvent } from '@sveltejs/kit';
import type { PfeAssignment, PfeProgressReport, SupervisorEvaluation } from '$lib/types';

export async function load({ params }: LoadEvent) {
  const id = Number(params.id);
  try {
    const [pfe, progressReports, supervisorEval] = await Promise.all([
      teacher.getSupervisedPFE(id),
      teacher.listMeetings(id).catch(() => [] as PfeProgressReport[]),
      teacher.getEvaluation(id).catch(() => null),
    ]);
    return {
      pfe: pfe ?? null,
      progressReports: progressReports ?? [],
      supervisorEval: supervisorEval ?? null,
    };
  } catch {
    return {
      pfe: null as PfeAssignment | null,
      progressReports: [] as PfeProgressReport[],
      supervisorEval: null as SupervisorEvaluation | null,
    };
  }
}
