import { admin } from '$lib/api';
import type { AuditLog } from '$lib/types';

export async function load({ url }: { url: URL }) {
  const actionTypeFilter = url.searchParams.get('action_type');
  const actorFilter = url.searchParams.get('actor');
  const dateFrom = url.searchParams.get('date_from');
  const dateTo = url.searchParams.get('date_to');
  const search = url.searchParams.get('search');

  try {
    const list = await admin.auditLog();
    return {
      logs: list ?? [],
      actionTypeFilter,
      actorFilter,
      dateFrom,
      dateTo,
      search,
    };
  } catch {
    return {
      logs: [] as AuditLog[],
      actionTypeFilter: null,
      actorFilter: null,
      dateFrom: null,
      dateTo: null,
      search: null,
    };
  }
}
