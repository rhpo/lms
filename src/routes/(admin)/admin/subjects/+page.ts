import { admin } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const [subjects, specialities, teachers] = await Promise.all([
      admin.listSubjects() as Promise<unknown[]>,
      (await import('$lib/api')).shared.specialities() as Promise<unknown[]>,
      admin.listUsers() as Promise<Record<string, unknown>>,
    ]);
    return {
      subjects: subjects ?? [],
      specialities: specialities ?? [],
      teachers: ((teachers as Record<string, unknown>)?.teachers as any[]) ?? [],
    };
  } catch {
    return { subjects: [], specialities: [], teachers: [] };
  }
}
