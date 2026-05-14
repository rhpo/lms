import { authStore } from '$lib/stores/auth';

export const ssr = false;
export const prerender = false;

export async function load() {
  // Wait for the auth store to initialize if it hasn't already
  if (!authStore.initialized) {
    await authStore.init();
  }
  return { profile: authStore.profile ?? null };
}
