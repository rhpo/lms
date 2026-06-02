import { teacher } from '$lib/api';
import type { PfeSubject, Wish } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';

export async function load({ params }: LoadEvent) {
  try {
    const [subject, wishes] = await Promise.all([
      teacher.getProposedSubject(Number(params.id)),
      teacher.listCandidats(Number(params.id)),
    ]);
    return {
      subject: subject ?? null,
      wishes: (wishes ?? []) as Wish[],
    };
  } catch {
    return {
      subject: null as PfeSubject | null,
      wishes: [] as Wish[],
    };
  }
}
