import { teacher } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const duties = await teacher.listJuryDuties();
    return { juryDuties: (duties as unknown[]) ?? [] };
  } catch {
    return { juryDuties: [] };
  }
}
