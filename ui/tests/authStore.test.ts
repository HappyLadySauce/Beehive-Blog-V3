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
});
