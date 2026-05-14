import { student } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const meetings = await student.listMyMeetings();
    return { meetings: (meetings as any[]) ?? [] };
  } catch {
    return { meetings: [] };
  }
}
