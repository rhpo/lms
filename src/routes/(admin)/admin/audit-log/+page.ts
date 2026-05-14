import { admin } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const logs = await admin.auditLog();
    return { auditLogs: (logs as any[]) ?? [] };
  } catch {
    return { auditLogs: [] };
  }
}
