import { admin, shared } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const [usersData, specialities] = await Promise.all([
      admin.listUsers() as Promise<Record<string, unknown>>,
      shared.specialities(),
    ]);
    return {
      formattedTeachers: (usersData?.teachers as unknown[]) ?? [],
      formattedStudents: (usersData?.students as unknown[]) ?? [],
      formattedCompanies: (usersData?.companies as unknown[]) ?? [],
      specialities: (specialities as unknown[]) ?? [],
    };
  } catch {
    return { formattedTeachers: [], formattedStudents: [], formattedCompanies: [], specialities: [] };
  }
}
