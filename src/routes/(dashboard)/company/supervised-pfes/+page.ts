import { company } from '$lib/api';
import type { PfeAssignment } from '$lib/types';

export async function load() {
  try {
    const pfes = await company.listSupervisedPFEs();
    return { supervisedPfes: pfes ?? [] };
  } catch {
    return { supervisedPfes: [] as PfeAssignment[] };
  }
}
