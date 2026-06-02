import { auth } from '$lib/api';



export async function load() {
  try {
    const me = await auth.me();
    return { teacher: me?.teacher ?? { availability_status: 'disponible', unavailable_until: null } as any };
  } catch {
    return { teacher: { availability_status: 'disponible', unavailable_until: null } as any };
  }
}
