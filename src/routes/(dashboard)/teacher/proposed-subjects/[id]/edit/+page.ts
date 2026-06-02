import { teacher } from '$lib/api';
import type { LoadEvent } from '@sveltejs/kit';
import type { PfeSubject } from '$lib/types';

export async function load({ params }: LoadEvent) {
  const id = Number(params.id);
  const subject = await teacher.getProposedSubject(id).catch(() => null);
  return {
    subject: subject as PfeSubject | null,
  };
}
