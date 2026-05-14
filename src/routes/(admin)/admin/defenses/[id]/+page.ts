import { admin } from '$lib/api';
import type { LoadEvent } from '@sveltejs/kit';

export const ssr = false;
export const prerender = false;

export async function load({ params }: LoadEvent) {
  try {
    const defense = await admin.getDefense(params.id!);
    return { defense: defense ?? null };
  } catch {
    return { defense: null };
  }
}
