import { admin } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const academicYears = await admin.listAcademicYears();
    return { academicYears: (academicYears as any[]) ?? [] };
  } catch {
    return { academicYears: [] };
  }
}
