import { admin } from '$lib/api';
import type { PfeSubject, Profile } from '$lib/types';

export async function load() {
  try {
    const [subjects, users] = await Promise.all([
      admin.listSubjects(),
      admin.listUsers(),
    ]);

    const teachers = (users ?? []).filter((u) => u.role === 'teacher');
    return {
      subjects: subjects ?? [],
      teachers,
    };
  } catch {
    return {
      subjects: [] as PfeSubject[],
      teachers: [] as Profile[],
    };
  }
}
