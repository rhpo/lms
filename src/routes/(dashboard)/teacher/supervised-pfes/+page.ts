import { teacher } from '$lib/api';
import type { PfeAssignment } from '$lib/types';

export async function load() {
  try {
    const pfes = await teacher.listSupervisedPFEs();
    return { supervisedPfes: pfes ?? [] };
  } catch {
    return { supervisedPfes: [] as PfeAssignment[] };
  }
}
