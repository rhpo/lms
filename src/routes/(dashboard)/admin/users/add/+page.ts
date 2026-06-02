import { admin } from '$lib/api';
import type { Department, Domain, Speciality, Promotion } from '$lib/types';

export async function load({ url }) {
  const type = url.searchParams.get('type') ?? 'teacher';

  const [departments, domains, specialities, promotions] = await Promise.all([
    admin.listDepartments(),
    admin.listDomains(),
    admin.listSpecialities(),
    admin.listPromotions(),
  ]);

  return {
    type,
    departments: departments as Department[],
    domains: domains as Domain[],
    specialities: specialities as Speciality[],
    promotions: promotions as Promotion[],
  };
}
