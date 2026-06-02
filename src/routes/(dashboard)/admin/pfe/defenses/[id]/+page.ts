import { admin } from '$lib/api';
import type { Defense } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';

export async function load({ params }: LoadEvent) {
  const defenseId = Number(params.id);
  try {
    const defense = await admin.getDefense(defenseId);
    return {
      defense: defense ?? null,
    };
  } catch {
    return {
      defense: null as Defense | null,
    };
  }
}
