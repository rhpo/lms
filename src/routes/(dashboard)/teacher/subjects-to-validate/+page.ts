import { teacher } from '$lib/api';
import type { PfeSubject } from '$lib/types';

export async function load() {
  try {
    const subjects = await teacher.listSubjectsToValidate();
    return { subjects: subjects ?? [] };
  } catch {
    return { subjects: [] as PfeSubject[] };
  }
}
