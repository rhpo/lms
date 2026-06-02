import { company } from '$lib/api';
import type { PfeAssignment, PfeProgressReport, SupervisorEvaluation } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';

export async function load({ params }: LoadEvent) {
  const id = Number(params.id);
  try {
    const [pfe, progressReports, supervisorEval] = await Promise.all([
      company.getSupervisedPFE(id),
      company.listMeetings(id).catch(() => [] as PfeProgressReport[]),
      company.getEvaluation(id).catch(() => null),
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
