import { student, shared } from '$lib/api';
import type { LoadEvent } from '@sveltejs/kit';

export const ssr = false;
export const prerender = false;

export async function load({ url }: LoadEvent) {
  try {
    const params = Object.fromEntries(url.searchParams.entries());
    const [subjects, specialities] = await Promise.all([
      student.listCatalogue(params) as Promise<unknown[]>,
      shared.specialities() as Promise<unknown[]>,
    ]);
    return {
      subjects: subjects ?? [],
      specialities: specialities ?? [],
    };
  } catch {
    return { subjects: [], specialities: [] };
  }
}
