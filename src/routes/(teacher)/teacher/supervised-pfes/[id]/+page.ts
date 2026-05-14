import { teacher } from '$lib/api';
import type { LoadEvent } from '@sveltejs/kit';

export const ssr = false;
export const prerender = false;

export async function load({ params }: LoadEvent) {
  try {
    const pfe = await teacher.getSupervisedPFE(params.id!);
    return { pfe: pfe ?? null };
  } catch {
    return { pfe: null };
  }
}
