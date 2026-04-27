import { afterEach, describe, expect, it, vi } from 'vitest'

import { createStudioApi } from '@/features/studio'

describe('studioApi', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('returns mock users and audit events', async () => {
    const api = createStudioApi('mock')

    const users = await api.listUsers()
    const audits = await api.listAudits()

    expect(users.items.map((user) => user.email)).toContain('admin@beehive.local')
    expect(audits.items.map((event) => event.event_type)).toContain('studio_access')
  })

  it('sends live studio requests through gateway paths', async () => {
    const fetchMock = vi.fn(async (input: RequestInfo | URL) => {
      const url = String(input)
      if (url.endsWith('/studio/users')) {
        return Response.json({ items: [], total: 0, page: 1, page_size: 20 })
      }
      if (url.endsWith('/studio/audits')) {
        return Response.json({ items: [], total: 0, page: 1, page_size: 20 })
      }
      if (url.endsWith('/auth/me/profile')) {
        return Response.json({
          user: {
            user_id: 'user_mock_admin',
            username: 'admin',
            email: 'admin@beehive.local',
            nickname: 'Admin',
            avatar_url: '',
            role: 'admin',
            status: 'active',
          },
        })
      }
      return Response.json({ ok: true })
    })
    vi.stubGlobal('fetch', fetchMock)

    const api = createStudioApi('live')
    await api.listUsers(undefined, { accessToken: 'access' })
    await api.listAudits(undefined, { accessToken: 'access' })
    await api.updateProfile({ avatar_url: '' }, { accessToken: 'access' })
    await api.changePassword({ old_password: 'Admin@123456', new_password: 'Admin@123456789' }, { accessToken: 'access' })

    expect(fetchMock).toHaveBeenNthCalledWith(
      1,
      '/api/v3/studio/users',
      expect.objectContaining({ method: 'GET' }),
    )
    expect(fetchMock).toHaveBeenNthCalledWith(
      2,
      '/api/v3/studio/audits',
      expect.objectContaining({ method: 'GET' }),
    )
    expect(fetchMock).toHaveBeenNthCalledWith(
      3,
      '/api/v3/auth/me/profile',
      expect.objectContaining({ method: 'PATCH' }),
    )
    expect(fetchMock).toHaveBeenNthCalledWith(
      4,
      '/api/v3/auth/me/password',
      expect.objectContaining({ method: 'POST' }),
    )

    const headers = fetchMock.mock.calls[0]?.[1]?.headers as Headers
    expect(headers.get('Authorization')).toBe('Bearer access')
    expect(fetchMock.mock.calls[2]?.[1]?.body).toBe(JSON.stringify({ avatar_url: '' }))
  })

  it('accepts the same minimum password length as the live identity contract', async () => {
    const api = createStudioApi('mock')

    await expect(api.changePassword({
      old_password: 'Admin@123456',
      new_password: '12345678',
    })).resolves.toEqual({ ok: true })
  })

  it('rejects weak mock password changes', async () => {
    const api = createStudioApi('mock')

    await expect(api.changePassword({
      old_password: 'Admin@123456',
      new_password: 'short',
    })).rejects.toThrow('New password')
  })
})
