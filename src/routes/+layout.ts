import { authStore } from '$lib/stores/auth';

export const ssr = false;

export async function load({ depends }: any) {
  depends('auth:profile');
  await authStore.init();
  return {
    profile: authStore.profile
  };
}
