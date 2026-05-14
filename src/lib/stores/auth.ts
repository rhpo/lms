import { auth as authApi, getToken, setToken, clearToken } from '$lib/api';
import type { Profile } from '$lib/api';
import { goto } from '$app/navigation';

// ── Reactive state ──────────────────────────────────────────────────────────

let _profile = $state<Profile | null>(null);
let _loading = $state(false);
let _initialized = $state(false);

export const authStore = {
  get profile() { return _profile; },
  get loading() { return _loading; },
  get initialized() { return _initialized; },
  get isAuthenticated() { return _profile !== null; },

  /**
   * Call once on app mount to restore the session from localStorage.
   */
  async init() {
    if (_initialized) return;
    _initialized = true;
    const token = getToken();
    if (!token) return;
    _loading = true;
    try {
      _profile = await authApi.me();
    } catch {
      // Token expired or invalid
      clearToken();
      _profile = null;
    } finally {
      _loading = false;
    }
  },

  /**
   * Dev-login (only available in development env).
   */
  async devLogin(email: string): Promise<void> {
    _loading = true;
    try {
      const result = await authApi.devLogin(email);
      setToken(result.token);
      _profile = result.profile;
      // Redirect based on role
      const redirects: Record<string, string> = {
        admin: '/admin/dashboard',
        teacher: '/teacher/dashboard',
        student: '/student/dashboard',
        company: '/company/dashboard',
      };
      await goto(redirects[result.profile.role] ?? '/');
    } finally {
      _loading = false;
    }
  },

  async logout(): Promise<void> {
    try { await authApi.logout(); } catch { /* ignore */ }
    clearToken();
    _profile = null;
    await goto('/accounts/login');
  },
};
