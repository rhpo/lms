import { student } from '$lib/api';
import type { PfeSubject } from '$lib/types';

export async function load() {
  const [subjects, pfe] = await Promise.all([
    student.listCatalogue().catch(() => [] as PfeSubject[]),
    student.getMyPFE().catch(() => null),
  ]);
  return {
    subjects: (subjects ?? []) as PfeSubject[],
    hasPfe: pfe != null,
  };
}
