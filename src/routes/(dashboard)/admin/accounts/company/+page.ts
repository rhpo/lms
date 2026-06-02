import { admin } from '$lib/api';
import type { Company } from '$lib/types';

export async function load() {
  try {
    const companies = await admin.listCompanies();
    return { companies: (companies ?? []) as Company[] };
  } catch {
    return { companies: [] as Company[] };
  }
}
