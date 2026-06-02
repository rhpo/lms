import { admin } from '$lib/api';
import type { PfeSubject, Profile } from '$lib/types';
import type { LoadEvent } from '@sveltejs/kit';

export async function load({ params }: LoadEvent) {
  try {
    const [subject, users] = await Promise.all([
      admin.getSubject(Number(params.id)),
      admin.listUsers(),
    ]);
    const teachers = (users ?? []).filter((u: Profile) => u.role === 'teacher');
    return {
      subject: subject ?? null,
      teachers: teachers as Profile[],
    };
  } catch {
    return {
      subject: null as PfeSubject | null,
      teachers: [] as Profile[],
    };
  }
}
