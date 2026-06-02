import { admin } from '$lib/api';
import type { Profile, Company } from '$lib/types';

export async function load() {
  try {
    const [profiles, companies] = await Promise.all([
      admin.listUsers(),
      admin.listCompanies(),
    ]);
    return {
      profiles: profiles ?? [],
      companies: companies ?? [],
    };
  } catch {
    return { profiles: [] as Profile[], companies: [] as Company[] };
  }
}
