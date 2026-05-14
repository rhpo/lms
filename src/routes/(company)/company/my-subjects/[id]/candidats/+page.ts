import { company } from '$lib/api';
import type { LoadEvent } from '@sveltejs/kit';

export const ssr = false;
export const prerender = false;

export async function load({ params }: LoadEvent) {
  try {
    const [subject, candidats] = await Promise.all([
      company.getSubject(params.id!) as Promise<Record<string, unknown>>,
      company.listCandidats(params.id!) as Promise<unknown[]>,
    ]);
    return { subject: subject ?? null, candidats: candidats ?? [] };
  } catch {
    return { subject: null, candidats: [] };
  }
}
