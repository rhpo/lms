import { teacher } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const subjects = await teacher.listProposedSubjects();
    return { subjects: (subjects as any[]) ?? [] };
  } catch {
    return { subjects: [] };
  }
}
