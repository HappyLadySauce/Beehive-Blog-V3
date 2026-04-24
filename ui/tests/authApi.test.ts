import { describe, expect, it, vi } from 'vitest';

import { createAuthApi } from '@/shared/api/authApi';

describe('authApi', () => {
  it('uses mock adapter without backend dependency', async () => {
    const api = createAuthApi('mock');

    const response = await api.login({
      login_identifier: 'tester@beehive.local',
      password: 'Demo@123456',
    });

    expect(response.access_token).toContain('mock_access');
    expect(response.user.email).toBe('tester@beehive.local');
  });

  it('uses live adapter when requested', async () => {
    const fetchMock = vi.fn(async () =>
      new Response(
        JSON.stringify({
          user: {
            user_id: 'user_live_001',
            username: 'live',
            email: 'live@beehive.local',
            nickname: 'Live',
            avatar_url: '',
            role: 'member',
            status: 'active',
          },
        }),
        { status: 200, headers: { 'Content-Type': 'application/json' } },
      ),
    );
    vi.stubGlobal('fetch', fetchMock);

    const api = createAuthApi('live');
    const response = await api.me({ accessToken: 'access-token' });

    expect(response.user.user_id).toBe('user_live_001');
    expect(fetchMock).toHaveBeenCalledWith(
      '/api/v3/auth/me',
      expect.objectContaining({
        method: 'GET',
      }),
    );
    vi.unstubAllGlobals();
  });
});
