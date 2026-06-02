import { redirect } from '@sveltejs/kit';
import { browser } from '$app/environment';
import type { LayoutLoad } from './$types';

export const ssr = false;

export const load: LayoutLoad = async ({ parent }) => {
  const { profile } = await parent();


  if (browser && !profile) {
    throw redirect(302, '/accounts/login');
  }

  return { profile };
};
