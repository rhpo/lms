import { student } from '$lib/api';
import type { Defense, JuryGrade } from '$lib/types';

export async function load() {
  try {
    const raw = await student.getSoutenance();
    if (!raw || !raw.has_soutenance || !raw.defense) {
      return {
        defense: null as Defense | null,
        grades: [] as JuryGrade[],
        supervisorNote: null as number | null,
        finalGradeBreakdown: null as null | {
          criterion1: number; criterion2: number;
          criterion3: number; criterion4: number; criterion5: number;
        },
      };
    }


    const defense: Defense = { ...raw.defense, jury: raw.jury };


    const supEval = (raw as any).supervisor_note ?? null;
    const supervisorNote: number | null = supEval?.criterion5 ?? null;





    return {
      defense,
      grades: [] as JuryGrade[],
      supervisorNote,
      finalGradeBreakdown: null,
    };
  } catch {
    return {
      defense: null as Defense | null,
      grades: [] as JuryGrade[],
      supervisorNote: null as number | null,
      finalGradeBreakdown: null,
    };
  }
}
