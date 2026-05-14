import { admin } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const [deadlines, specialities, domains, promotions, academicYears] = await Promise.all([
      admin.listDeadlines(),
      admin.listSpecialities(),
      admin.listDomains(),
      admin.listPromotions(),
      admin.listAcademicYears(),
    ]);
    return {
      deadlines: deadlines ?? null,
      specialities: (specialities as any[]) ?? [],
      domains: (domains as any[]) ?? [],
      promotions: (promotions as any[]) ?? [],
      academicYears: (academicYears as any[]) ?? [],
    };
  } catch {
    return { deadlines: null, specialities: [], domains: [], promotions: [], academicYears: [] };
  }
}
