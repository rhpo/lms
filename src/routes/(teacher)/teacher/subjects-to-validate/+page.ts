import { teacher } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const subjects = await teacher.listSubjectsToValidate();
    return { subjects: (subjects as unknown[]) ?? [] };
  } catch {
    return { subjects: [] };
  }
}
