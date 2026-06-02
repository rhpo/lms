import { admin } from '$lib/api';
import type { AcademicYear, Speciality, Domain, Department, Promotion } from '$lib/types';

export async function load() {
  try {
    const [deadlines, specialities, domains, departments, promotions, academicYears] = await Promise.all([
      admin.listDeadlines(),
      admin.listSpecialities(),
      admin.listDomains(),
      admin.listDepartments(),
      admin.listPromotions(),
      admin.listAcademicYears(),
    ]);
    const years = deadlines ?? academicYears ?? [];
    const activeYear = years.find((y) => y.status === 'active') ?? years[0];

    return {
      settings: activeYear ?? null,
      specialities: specialities ?? [],
      domains: domains ?? [],
      departments: departments ?? [],
      promotions: promotions ?? [],
      academicYears: years,
    };
  } catch {
    return {
      settings: null as AcademicYear | null,
      specialities: [] as Speciality[],
      domains: [] as Domain[],
      departments: [] as Department[],
      promotions: [] as Promotion[],
      academicYears: [] as AcademicYear[],
    };
  }
}
