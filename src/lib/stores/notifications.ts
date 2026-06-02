import { writable, get } from 'svelte/store';
import { notifications as notifApi } from '$lib/api';

const _unreadCount = writable(0);

let _polling: ReturnType<typeof setInterval> | null = null;

export const notificationStore = {
  /** Svelte store - subscribe for reactive unread count. */
  subscribe: _unreadCount.subscribe,

  get unreadCount() {
    return get(_unreadCount);
  },

  /** Fetch the current unread count from the backend. */
  async refresh() {
    try {
      const count = await notifApi.unreadCount();
      _unreadCount.set(count);
    } catch {

    }
  },

  /** Start polling every `intervalMs` (default 30 s). */
  startPolling(intervalMs = 30_000) {
    this.stopPolling();
    this.refresh();
    _polling = setInterval(() => this.refresh(), intervalMs);
  },

  stopPolling() {
    if (_polling) {
      clearInterval(_polling);
      _polling = null;
    }
  },

  /** Call after marking a single notification as read. */
  decrement() {
    _unreadCount.update((n) => Math.max(0, n - 1));
  },

  /** Call after marking all as read. */
  clear() {
    _unreadCount.set(0);
  },
};
