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
      formattedTeachers: (usersData?.teachers as any[]) ?? [],
      formattedStudents: (usersData?.students as any[]) ?? [],
      formattedCompanies: (usersData?.companies as any[]) ?? [],
      specialities: (specialities as any[]) ?? [],
    };
  } catch {
    return { formattedTeachers: [], formattedStudents: [], formattedCompanies: [], specialities: [] };
  }
}
