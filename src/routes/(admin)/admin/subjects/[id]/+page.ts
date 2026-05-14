import { admin } from '$lib/api';
import type { LoadEvent } from '@sveltejs/kit';

export const ssr = false;
export const prerender = false;

export async function load({ params }: LoadEvent) {
  try {
    const [subject, teachers] = await Promise.all([
      admin.getSubject(params.id!) as Promise<Record<string, unknown>>,
      admin.listUsers() as Promise<Record<string, unknown>>,
    ]);
    return {
      subject: subject ?? null,
      teachers: ((teachers as Record<string, unknown>)?.teachers as any[]) ?? [],
    };
  } catch {
    return { subject: null, teachers: [] };
  }
}
