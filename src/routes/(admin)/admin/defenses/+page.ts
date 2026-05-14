import { admin } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const defenses = await admin.listDefenses();
    return { defenses: (defenses as any[]) ?? [] };
  } catch {
    return { defenses: [] };
  }
}
