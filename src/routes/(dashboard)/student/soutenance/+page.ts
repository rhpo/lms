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

    // Merge jury into defense so the template can access defense.jury
    const defense: Defense = { ...raw.defense, jury: raw.jury };

    // supervisor_note from the API
    const supEval = (raw as any).supervisor_note ?? null;
    const supervisorNote: number | null = supEval?.criterion5 ?? null;

    // If final_grade is set, try to show breakdown from the president's jury grade
    // The defense object already carries final_grade — we pass it through.
    // The soutenance endpoint does not return the jury grades, so we just show the total.

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
