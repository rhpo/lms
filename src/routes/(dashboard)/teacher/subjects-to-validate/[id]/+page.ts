import { teacher } from '$lib/api';
import type { LoadEvent } from '@sveltejs/kit';



export async function load({ params }: LoadEvent) {
  try {
    const subject = await teacher.getSubjectToValidate(Number(params.id));
    return { subject: subject ?? null };
  } catch {
    return { subject: null };
  }
}
