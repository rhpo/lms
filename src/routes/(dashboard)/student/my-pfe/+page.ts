import { student } from '$lib/api';
import type { PfeAssignment, PfeProgressReport, Defense } from '$lib/types';

export async function load() {
  try {
    const [pfe, reports, soutenanceRaw] = await Promise.all([
      student.getMyPFE(),
      student.listMyMeetings(),
      student.getSoutenance().catch(() => null),
    ]);


    let defense: Defense | null = null;
    let supervisorNote: number | null = null;
    if (soutenanceRaw?.has_soutenance && soutenanceRaw.defense) {
      defense = { ...soutenanceRaw.defense, jury: soutenanceRaw.jury } as Defense;
      const supEval = (soutenanceRaw as any).supervisor_note ?? null;
      supervisorNote = supEval?.criterion5 ?? null;
    }

    return {
      pfe: pfe ?? null,
      progressReports: reports ?? [],
      defense,
      supervisorNote,
    };
  } catch {
    return {
      pfe: null as PfeAssignment | null,
      progressReports: [] as PfeProgressReport[],
      defense: null as Defense | null,
      supervisorNote: null as number | null,
    };
  }
}
