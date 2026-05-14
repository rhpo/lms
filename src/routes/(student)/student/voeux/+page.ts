import { student } from '$lib/api';

export const ssr = false;
export const prerender = false;

export async function load() {
  try {
    const wishes = await student.listWishes();
    return {
      wishes,
      maxWishes: 5,
      canStillSubmit: true,
    };
  } catch {
    return { wishes: [], maxWishes: 5, canStillSubmit: false };
  }
}
