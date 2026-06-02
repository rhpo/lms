import { admin } from '$lib/api';
import type { PfeAssignment, Profile } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';

export async function load({ params }: LoadEvent) {
  const pfeId = Number(params.id);
  try {
    const [assignment, users] = await Promise.all([
      admin.getAssignment(pfeId),
      admin.listUsers(),
    ]);
    const teachers = (users ?? []).filter((u: Profile) => u.role === 'teacher' || u.role === 'admin');
    return {
      assignment: assignment ?? null,
      teachers: teachers as Profile[],
    };
  } catch {
    return {
      assignment: null as PfeAssignment | null,
      teachers: [] as Profile[],
    };
  }
}
