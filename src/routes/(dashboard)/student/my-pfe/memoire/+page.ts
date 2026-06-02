import { student } from '$lib/api';
import type { PfeAssignment } from '$lib/types';

export async function load() {
  try {
    const pfe = await student.getMyPFE();
    return {
      pfe: pfe ?? null,
      deadlinePassed: false,
      memoireDeadline: '',
    };
  } catch {
    return {
      pfe: null as PfeAssignment | null,
      deadlinePassed: false,
      memoireDeadline: '',
    };
  }
}
