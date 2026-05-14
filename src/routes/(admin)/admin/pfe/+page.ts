import { admin } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const pfe = await admin.listAssignments();
    return { assignments: (pfe as unknown[]) ?? [] };
  } catch {
    return { assignments: [] };
  }
}
