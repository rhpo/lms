import { student } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const data = await student.listWishes() as Record<string, unknown>;
    return {
      wishes: (data?.wishes as unknown[]) ?? [],
      maxWishes: (data?.max_wishes as number) ?? 5,
      canStillSubmit: !!data?.can_still_submit,
    };
  } catch {
    return { wishes: [], maxWishes: 5, canStillSubmit: false };
  }
}
