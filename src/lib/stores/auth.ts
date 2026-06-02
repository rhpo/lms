import { writable, get } from 'svelte/store';
import { auth as authApi, getToken, setToken, clearToken } from '$lib/api';
import type { Profile } from '$lib/types';
import { goto, invalidate } from '$app/navigation';

const _profile = writable<Profile | null>(null);
const _loading = writable(false);
let _initialized = false;

export const authStore = {
  get profile() { return get(_profile); },
  get loading() { return get(_loading); },
  get initialized() { return _initialized; },
  get isAuthenticated() { return get(_profile) !== null; },

  /** Subscribe to profile changes (for reactive use in components) */
  subscribe: _profile.subscribe,

  /**
   * Call once on app mount to restore the session from localStorage.
   * Always re-checks the token on the client, even if called before during SSR.
   */
  async init() {

    if (typeof window === 'undefined') return;


    if (_initialized) return;
    _initialized = true;

    const token = getToken();
    if (!token) return;

    _loading.set(true);
    try {
      const profile = await authApi.me();
      _profile.set(profile);
    } catch {
      clearToken();
      _profile.set(null);
    } finally {
      _loading.set(false);
    }
  },

  /** Re-fetch profile from backend (e.g. after avatar upload). */
  async refreshProfile(): Promise<void> {
    try {
      const profile = await authApi.me();
      _profile.set(profile);
      await invalidate('auth:profile');
    } catch {

    }
  },

  async devLogin(email: string): Promise<void> {
    _loading.set(true);
    try {
      const result = await authApi.devLogin(email);
      setToken(result.token);
      _profile.set(result.profile);
      await invalidate('auth:profile');
      const redirects: Record<string, string> = {
        admin: '/admin/dashboard',
        teacher: '/teacher/dashboard',
        student: '/student/dashboard',
        company: '/company/dashboard',
      };
      await goto(redirects[result.profile.role] ?? '/', { invalidateAll: true });
    } finally {
      _loading.set(false);
    }
  },

  async logout(): Promise<void> {
    try { await authApi.logout(); } catch { /* ignore */ }
    clearToken();
    _profile.set(null);
    _initialized = false;
    await invalidate('auth:profile');
    await goto('/accounts/login', { invalidateAll: true });
  },
};
