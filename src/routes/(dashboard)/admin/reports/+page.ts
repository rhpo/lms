import { admin } from '$lib/api';
import type { CompanyReport } from '$lib/types';

export async function load() {
  try {
    const reports = await admin.listReports();
    return { reports: reports ?? [] };
  } catch {
    return { reports: [] as CompanyReport[] };
  }
}
