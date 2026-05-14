import { company } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const reports = await company.listReports();
    return { reports: (reports as unknown[]) ?? [] };
  } catch {
    return { reports: [] };
  }
}
