import { company } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const pfes = await company.listSupervisedPFEs();
    return { supervisedPfes: (pfes as unknown[]) ?? [] };
  } catch {
    return { supervisedPfes: [] };
  }
}
