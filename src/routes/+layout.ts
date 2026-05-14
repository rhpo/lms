import { authStore } from '$lib/stores/auth';

export const ssr = false;
export const prerender = false;

export async function load() {
  if (typeof window !== 'undefined' && !authStore.initialized) {
    await authStore.init();
  }
  return {
    profile: authStore.profile
  };
}
