import { student } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const soutenance = await student.getSoutenance();
    return { soutenance: soutenance ?? null };
  } catch {
    return { soutenance: null };
  }
}
