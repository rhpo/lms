import { teacher } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const list = await teacher.listJuryDuties() as any;
    return { duties: list ?? [] };
  } catch {
    return { duties: [] };
  }
}
