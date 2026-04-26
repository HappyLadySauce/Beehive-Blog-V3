import { createPinia, setActivePinia } from 'pinia';
import { beforeEach, describe, expect, it } from 'vitest';

import { useAuthStore } from '@/features/auth/stores/authStore';

describe('authStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('stores access token in memory and refresh token in storage', async () => {
    const store = useAuthStore();

    await store.login({
      login_identifier: 'creator@beehive.local',
      password: 'Demo@123456',
    });

    expect(store.isAuthenticated).toBe(true);
    expect(store.accessToken).toContain('mock_access');
    expect(window.localStorage.getItem('beehive.auth.refreshToken')).toContain('mock_refresh');
  });

  it('restores an in-memory session from refresh token storage', async () => {
    const initialStore = useAuthStore();
    await initialStore.login({
      login_identifier: 'creator@beehive.local',
      password: 'Demo@123456',
    });
    initialStore.$reset();

    const restored = await initialStore.restoreSession();

    expect(restored).toBe(true);
    expect(initialStore.isAuthenticated).toBe(true);
    expect(initialStore.accessToken).toContain('mock_access');
    expect(initialStore.currentUser?.email).toBe('creator@beehive.local');
  });

  it('clears invalid refresh tokens when restore fails', async () => {
    window.localStorage.setItem('beehive.auth.refreshToken', 'bad-refresh-token');
    const store = useAuthStore();

    const restored = await store.restoreSession();

    expect(restored).toBe(false);
    expect(store.isAuthenticated).toBe(false);
    expect(window.localStorage.getItem('beehive.auth.refreshToken')).toBeNull();
  });
});
