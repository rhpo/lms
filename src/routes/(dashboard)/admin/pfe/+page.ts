import { admin } from '$lib/api';
import type { PfeAssignment } from '$lib/types';

export async function load() {
  try {
    const pfe = await admin.listAssignments();
    return { assignments: pfe ?? [] };
  } catch {
    return { assignments: [] as PfeAssignment[] };
  }
}
