import { teacher } from '$lib/api';
import type { LoadEvent } from '@sveltejs/kit';

export const ssr = false;
export const prerender = false;

export async function load({ params }: LoadEvent) {
  try {
    const subject = await teacher.getSubjectToValidate(params.id!);
    return { subject: subject ?? null };
  } catch {
    return { subject: null };
  }
}
