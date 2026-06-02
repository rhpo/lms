import { admin } from '$lib/api';
import type { AcademicYear } from '$lib/types';

export async function load() {
  try {
    const academicYears = await admin.listAcademicYears();
    return { academicYears: academicYears ?? [] };
  } catch {
    return { academicYears: [] as AcademicYear[] };
  }
}
