import { createPinia, setActivePinia } from 'pinia';
import { beforeEach, describe, expect, it } from 'vitest';

import { router } from '@/app/router';
import { useAuthStore } from '@/features/auth/stores/authStore';

describe('router auth guard', () => {
  beforeEach(async () => {
    setActivePinia(createPinia());
    await router.push('/');
    await router.isReady();
  });

  it('redirects unauthenticated studio visits to login with redirect', async () => {
    await router.push('/studio/content');

    expect(router.currentRoute.value.name).toBe('auth-login');
    expect(router.currentRoute.value.query.redirect).toBe('/studio/content');
  });

  it('allows authenticated studio visits', async () => {
    const authStore = useAuthStore();
    await authStore.login({
      login_identifier: 'creator@beehive.local',
      password: 'Demo@123456',
    });

    await router.push('/studio/content');

    expect(router.currentRoute.value.name).toBe('studio-content');
  });
});
