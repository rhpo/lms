import { student } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const notifications = await student.listNotifications();
    return { notifications: (notifications as unknown[]) ?? [] };
  } catch {
    return { notifications: [] };
  }
}
