import { student } from '$lib/api';
import type { LoadEvent } from '@sveltejs/kit';

export const ssr = false;
export const prerender = false;

export async function load({ params }: LoadEvent) {
  try {
    const [subject, wishes] = await Promise.all([
      student.getCatalogueSubject(params.id!),
      student.listWishes()
    ]);
    const alreadyWished = wishes?.some((w: any) => w.subject_id === params.id) ?? false;
    return {
      subject: subject ?? null,
      alreadyWished,
      alreadyAssigned: false, // Could be fetched from listAssignments if needed
      wishesCount: wishes?.length ?? 0
    };
  } catch {
    return { subject: null, alreadyWished: false, alreadyAssigned: false, wishesCount: 0 };
  }
}
