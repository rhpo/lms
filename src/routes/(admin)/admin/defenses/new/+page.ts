import { admin } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const [assignments, teachers] = await Promise.all([
      admin.listAssignments() as Promise<unknown[]>,
      admin.listUsers() as Promise<Record<string, unknown>>,
    ]);
    return {
      assignments: assignments ?? [],
      teachers: ((teachers as Record<string, unknown>)?.teachers as unknown[]) ?? [],
    };
  } catch {
    return { assignments: [], teachers: [] };
  }
}
