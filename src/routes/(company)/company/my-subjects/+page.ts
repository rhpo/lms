import { company } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const subjects = await company.listSubjects();
    return { subjects: (subjects as any[]) ?? [] };
  } catch {
    return { subjects: [] };
  }
}
