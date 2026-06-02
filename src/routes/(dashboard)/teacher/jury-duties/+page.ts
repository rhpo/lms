import { teacher } from '$lib/api';
import type { Defense, JuryGrade, SupervisorEvaluation } from '$lib/types';

export type GradeContext = {
  my_role: 'president' | 'member';
  my_grade: JuryGrade | null;
  member_grade: JuryGrade | null;
  supervisor_eval: SupervisorEvaluation | null;
  member_submitted: boolean;
  supervisor_submitted: boolean;
  final_grade_set: boolean;
};

export type DutyWithContext = {
  defense: Defense;
  gradeCtx: GradeContext | null;
};

export async function load() {
  try {
    const list = await teacher.listJuryDuties();
    const defenses: Defense[] = list ?? [];


    const duties: DutyWithContext[] = await Promise.all(
      defenses.map(async (d) => {
        try {
          const gradeCtx = await teacher.getGradeContext(d.id);
          return { defense: d, gradeCtx };
        } catch {
          return { defense: d, gradeCtx: null };
        }
      })
    );

    return { duties };
  } catch {
    return { duties: [] as DutyWithContext[] };
  }
}
