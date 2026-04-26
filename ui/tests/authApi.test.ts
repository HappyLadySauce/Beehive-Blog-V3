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

  it('refreshes mock sessions from a stored refresh token', async () => {
    const api = createAuthApi('mock');
    const login = await api.login({
      login_identifier: 'tester@beehive.local',
      password: 'Demo@123456',
    });

    const response = await api.refresh({ refresh_token: login.refresh_token });

    expect(response.access_token).toContain('mock_access');
    expect(response.refresh_token).toContain('mock_refresh');
    expect('user' in response).toBe(false);
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

  it('posts refresh tokens through the live adapter', async () => {
    const fetchMock = vi.fn(async () =>
      new Response(
        JSON.stringify({
          access_token: 'live_access',
          refresh_token: 'live_refresh',
          expires_in: 900,
          token_type: 'Bearer',
          session_id: 'session_live_001',
          session: {
            session_id: 'session_live_001',
            user_id: 'user_live_001',
            auth_source: 'password',
            client_type: 'web',
            device_id: 'browser',
            device_name: 'Browser',
            status: 'active',
            last_seen_at: 1713772000,
            expires_at: 1713772900,
          },
        }),
        { status: 200, headers: { 'Content-Type': 'application/json' } },
      ),
    );
    vi.stubGlobal('fetch', fetchMock);

    const api = createAuthApi('live');
    const response = await api.refresh({ refresh_token: 'refresh-token' });

    expect(response.access_token).toBe('live_access');
    expect(fetchMock).toHaveBeenCalledWith(
      '/api/v3/auth/refresh',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ refresh_token: 'refresh-token' }),
      }),
    );
    vi.unstubAllGlobals();
  });
});
