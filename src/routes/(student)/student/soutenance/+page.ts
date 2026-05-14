import { student } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const data = await student.getSoutenance() as any;
    return { 
      defense: data?.defense ?? data ?? null, 
      grades: data?.grades ?? [] 
    };
  } catch {
    return { defense: null, grades: [] };
  }
}
