import { admin } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const companies = await admin.listCompanies();
    return { companies: (companies as any[]) ?? [] };
  } catch {
    return { companies: [] };
  }
}
