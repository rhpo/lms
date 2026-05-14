import { student } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const data = await student.dashboard() as Record<string, unknown>;
    return {
      currentPfe: data?.current_pfe ?? null,
      wishes: (data?.wishes as unknown[]) ?? [],
      notifications: (data?.notifications as unknown[]) ?? [],
      yearId: (data?.active_year_id as string) ?? null,
    };
  } catch {
    return { currentPfe: null, wishes: [], notifications: [], yearId: null };
  }
}
