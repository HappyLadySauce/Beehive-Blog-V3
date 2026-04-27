import { describe, expect, it, vi } from 'vitest'

import { createAuthApi } from '@/features/auth/api/authApi'

describe('authApi', () => {
  it('uses mock login and refresh without gateway calls', async () => {
    const api = createAuthApi('mock')
    const login = await api.login({
      login_identifier: 'admin@beehive.local',
      password: 'Admin@123456',
    })

    expect(login.user.role).toBe('admin')
    expect(login.refresh_token).toContain('mock_refresh_')

    const refresh = await api.refresh({ refresh_token: login.refresh_token })
    expect(refresh.refresh_token).toContain('mock_refresh_')
    expect(refresh.session_id).toBe(login.session_id)
  })

  it('sends live login and me requests through gateway paths', async () => {
    const fetchMock = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const url = String(input)
      if (url.endsWith('/login')) {
        return Response.json({
          access_token: 'access',
          refresh_token: 'refresh',
          expires_in: 900,
          token_type: 'Bearer',
          session_id: 'session',
          user: {
            user_id: '1',
            username: 'admin',
            email: 'admin@beehive.local',
            nickname: 'Admin',
            avatar_url: '',
            role: 'admin',
            status: 'active',
          },
          session: {
            session_id: 'session',
            user_id: '1',
            auth_source: 'local',
            client_type: 'web',
            device_id: 'browser',
            device_name: 'Browser',
            status: 'active',
            last_seen_at: 1,
            expires_at: 2,
          },
        })
      }
      expect((init?.headers as Headers).get('Authorization')).toBe('Bearer access')
      return Response.json({
        user: {
          user_id: '1',
          username: 'admin',
          email: 'admin@beehive.local',
          nickname: 'Admin',
          avatar_url: '',
          role: 'admin',
          status: 'active',
        },
      })
    })
    vi.stubGlobal('fetch', fetchMock)

    const api = createAuthApi('live')
    await api.login({ login_identifier: 'admin@beehive.local', password: 'Admin@123456' })
    await api.me({ accessToken: 'access' })

    expect(fetchMock).toHaveBeenNthCalledWith(
      1,
      '/api/v3/auth/login',
      expect.objectContaining({ method: 'POST' }),
    )
    expect(fetchMock).toHaveBeenNthCalledWith(
      2,
      '/api/v3/auth/me',
      expect.objectContaining({ method: 'GET' }),
    )
  })
})
