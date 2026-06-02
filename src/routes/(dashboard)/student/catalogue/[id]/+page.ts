import { student } from '$lib/api';
import type { PfeSubject, Wish } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';

export async function load({ params }: LoadEvent) {
  const id = Number(params.id);

  // Fetch subject and wishes independently — a wishes failure must not kill the subject load
  const [subject, wishes] = await Promise.all([
    student.getCatalogueSubject(id).catch(() => null),
    student.listWishes().catch(() => [] as Wish[]),
  ]);

  const alreadyWished = wishes?.some((w: Wish) => w.subject_id === id) ?? false;
  const pfe = await student.getMyPFE().catch(() => null);

  return {
    subject: subject ?? null,
    alreadyWished,
    alreadyAssigned: pfe != null,
    subjectTaken: subject?.is_assigned ?? false,
    wishesCount: wishes?.length ?? 0,
  };
}
