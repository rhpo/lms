import { admin } from '$lib/api';
import type { Promotion } from '$lib/types';

export async function load() {
  try {
    const promotions = await admin.listPromotions();
    return { promotions: promotions ?? [] };
  } catch {
    return { promotions: [] as Promotion[] };
  }
}
