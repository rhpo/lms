import { admin } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const stats = await admin.statistics();
    return { stats: (stats ?? {}) as Record<string, unknown> };
  } catch {
    return { stats: {} };
  }
}
