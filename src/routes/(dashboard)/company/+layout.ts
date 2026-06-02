import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const ssr = false;

export const load: LayoutLoad = async ({ parent }) => {
  const { profile } = await parent();

  if (profile && profile.role !== 'company') {
    throw redirect(302, `/${profile.role}/dashboard`);
  }

  return { profile };
};
