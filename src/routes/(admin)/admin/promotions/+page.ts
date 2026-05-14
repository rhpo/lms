import { admin } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const promotions = await admin.listPromotions();
    return { promotions: (promotions as any[]) ?? [] };
  } catch {
    return { promotions: [] };
  }
}
