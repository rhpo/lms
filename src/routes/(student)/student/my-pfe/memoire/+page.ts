import { student } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const data = await student.getMyPFE() as Record<string, unknown>;
    return {
      pfe: data?.pfe ?? null,
    };
  } catch {
    return { pfe: null };
  }
}
