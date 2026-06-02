import { authStore } from '$lib/stores/auth';



export async function load({ parent }) {
  const { profile } = await parent();
  return { profile: profile ?? null };
}
