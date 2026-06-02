import { company } from '$lib/api';
import type { PfeSubject, Wish } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';

export async function load({ params }: LoadEvent) {
  try {
    const [subject, wishes] = await Promise.all([
      company.getSubject(Number(params.id)),
      company.listCandidats(Number(params.id)),
    ]);
    const wishesList = (wishes ?? []) as Wish[];

    const isLocked = wishesList.some((w) => w.status === 'accepte');
    return {
      subject: subject ?? null,
      wishes: wishesList,
      isLocked,
    };
  } catch {
    return {
      subject: null as PfeSubject | null,
      wishes: [] as Wish[],
      isLocked: false,
    };
  }
}
