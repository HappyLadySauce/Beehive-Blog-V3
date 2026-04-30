import { createPinia, setActivePinia } from 'pinia'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { authApi } from '@/features/auth/api/authApi'
import { useAuthStore } from '@/features/auth/stores/authStore'
import { tokenStorage } from '@/shared/storage/tokenStorage'

describe('authStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    tokenStorage.clearRefreshToken()
    if ('clearSnapshot' in tokenStorage && typeof tokenStorage.clearSnapshot === 'function') {
      tokenStorage.clearSnapshot()
    }
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('logs in, stores refresh token, restores and logs out', async () => {
    const store = useAuthStore()
    const refreshSpy = vi.spyOn(authApi, 'refresh')
    await store.login({
      login_identifier: 'admin@beehive.local',
      password: 'Admin@123456',
    })

    expect(store.isAuthenticated).toBe(true)
    expect(store.isAdmin).toBe(true)
    expect(store.refreshToken).toContain('mock_refresh_')

    store.accessToken = ''
    store.currentUser = null
    expect(await store.restoreSession()).toBe(true)
    expect(store.currentUser?.role).toBe('admin')
    expect(refreshSpy).not.toHaveBeenCalled()

    await store.logout()
    expect(store.isAuthenticated).toBe(false)
    expect(store.refreshToken).toBeNull()
  })

  it('restores a valid access token snapshot without refreshing', async () => {
    tokenStorage.setSnapshot({
      accessToken: 'access-token',
      refreshToken: 'refresh-token',
      accessTokenExpiresAt: Date.now() + 5 * 60_000,
      sessionId: 'session-1',
      currentUser: {
        user_id: '1',
        username: 'admin',
        email: 'admin@beehive.local',
        nickname: 'Admin',
        avatar_url: '',
        role: 'admin',
        status: 'active',
      },
    })
    const store = useAuthStore()
    const refreshSpy = vi.spyOn(authApi, 'refresh')

    await expect(store.restoreSession()).resolves.toBe(true)

    expect(store.accessToken).toBe('access-token')
    expect(store.currentUser?.user_id).toBe('1')
    expect(refreshSpy).not.toHaveBeenCalled()
  })

  it('loads current user from a valid snapshot before refreshing', async () => {
    tokenStorage.setSnapshot({
      accessToken: 'access-token',
      refreshToken: 'refresh-token',
      accessTokenExpiresAt: Date.now() + 5 * 60_000,
      sessionId: 'session-1',
      currentUser: null,
    })
    const store = useAuthStore()
    const refreshSpy = vi.spyOn(authApi, 'refresh')
    const meSpy = vi.spyOn(authApi, 'me').mockResolvedValue({
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

    await expect(store.restoreSession()).resolves.toBe(true)

    expect(meSpy).toHaveBeenCalledWith({ accessToken: 'access-token' })
    expect(store.currentUser?.user_id).toBe('1')
    expect(refreshSpy).not.toHaveBeenCalled()
  })

  it('refreshes when the access token snapshot is near expiry', async () => {
    tokenStorage.setSnapshot({
      accessToken: 'stale-access',
      refreshToken: 'refresh-token',
      accessTokenExpiresAt: Date.now() + 30_000,
      sessionId: 'session-1',
      currentUser: {
        user_id: '1',
        username: 'admin',
        email: 'admin@beehive.local',
        nickname: 'Admin',
        avatar_url: '',
        role: 'admin',
        status: 'active',
      },
    })
    const store = useAuthStore()
    const user = {
      user_id: '1',
      username: 'admin',
      email: 'admin@beehive.local',
      nickname: 'Admin',
      avatar_url: '',
      role: 'admin',
      status: 'active',
    }
    const refreshSpy = vi.spyOn(authApi, 'refresh').mockResolvedValue({
      access_token: 'fresh-access',
      refresh_token: 'fresh-refresh',
      expires_in: 900,
      token_type: 'Bearer',
      session_id: 'session-2',
      user,
      session: {
        session_id: 'session-2',
        user_id: '1',
        auth_source: 'local',
        client_type: 'web',
        status: 'active',
      },
    })

    await expect(store.restoreSession()).resolves.toBe(true)

    expect(refreshSpy).toHaveBeenCalledTimes(1)
    expect(store.accessToken).toBe('fresh-access')
    expect(store.refreshToken).toBe('fresh-refresh')
  })

  it('clears invalid refresh token on restore failure', async () => {
    tokenStorage.writeRefreshToken('bad-token')
    const store = useAuthStore()

    expect(await store.restoreSession()).toBe(false)
    expect(store.refreshToken).toBeNull()
  })

  it('coalesces concurrent restore refresh calls', async () => {
    tokenStorage.writeRefreshToken('old-refresh')
    const store = useAuthStore()
    const user = {
      user_id: '1',
      username: 'admin',
      email: 'admin@beehive.local',
      nickname: 'Admin',
      avatar_url: '',
      role: 'admin',
      status: 'active',
    }
    let calls = 0
    vi.spyOn(authApi, 'refresh').mockImplementation(async () => {
      calls += 1
      await Promise.resolve()
      return {
        access_token: 'access',
        refresh_token: 'new-refresh',
        expires_in: 900,
        token_type: 'Bearer',
        session_id: 's1',
        user,
        session: {
          session_id: 's1',
          user_id: '1',
          auth_source: 'local',
          client_type: 'web',
          status: 'active',
        },
      }
    })

    const [first, second] = await Promise.all([store.restoreSession(), store.restoreSession()])

    expect(first).toBe(true)
    expect(second).toBe(true)
    expect(calls).toBe(1)
    expect(store.refreshToken).toBe('new-refresh')
    expect(store.currentUser?.user_id).toBe('1')
  })

  it('recovers when a sibling refresh already rotated the token', async () => {
    tokenStorage.writeRefreshToken('old-refresh')
    setActivePinia(createPinia())
    const storeA = useAuthStore()
    setActivePinia(createPinia())
    const storeB = useAuthStore()
    const user = {
      user_id: '1',
      username: 'admin',
      email: 'admin@beehive.local',
      nickname: 'Admin',
      avatar_url: '',
      role: 'admin',
      status: 'active',
    }
    let resolveFirstOldRefresh: ((value: Awaited<ReturnType<typeof authApi.refresh>>) => void) | null = null
    let rejectSecondOldRefresh: ((reason?: unknown) => void) | null = null
    const refreshSpy = vi.spyOn(authApi, 'refresh').mockImplementation((payload) => {
      if (payload.refresh_token === 'old-refresh') {
        if (!resolveFirstOldRefresh) {
          return new Promise((resolve) => {
            resolveFirstOldRefresh = resolve
          })
        }
        return new Promise((_, reject) => {
          rejectSecondOldRefresh = reject
        })
      }

      expect(payload.refresh_token).toBe('new-refresh')
      return Promise.resolve({
        access_token: 'access-2',
        refresh_token: 'latest-refresh',
        expires_in: 900,
        token_type: 'Bearer',
        session_id: 's2',
        user,
        session: {
          session_id: 's2',
          user_id: '1',
          auth_source: 'local',
          client_type: 'web',
          status: 'active',
        },
      })
    })

    const firstRefresh = storeA.refreshSession()
    const secondRefresh = storeB.refreshSession()

    resolveFirstOldRefresh?.({
      access_token: 'access-1',
      refresh_token: 'new-refresh',
      expires_in: 900,
      token_type: 'Bearer',
      session_id: 's1',
      user,
      session: {
        session_id: 's1',
        user_id: '1',
        auth_source: 'local',
        client_type: 'web',
        status: 'active',
      },
    })
    await firstRefresh
    rejectSecondOldRefresh?.(new Error('invalid refresh token'))

    await expect(secondRefresh).resolves.toBe(true)
    expect(refreshSpy).toHaveBeenCalledTimes(3)
    expect(tokenStorage.readRefreshToken()).toBe('latest-refresh')
    expect(storeB.currentUser?.user_id).toBe('1')
  })

  it('clears the session when refresh still fails without a newer token', async () => {
    tokenStorage.writeRefreshToken('stale-refresh')
    const store = useAuthStore()

    vi.spyOn(authApi, 'refresh').mockRejectedValue(new Error('invalid refresh token'))

    await expect(store.refreshSession()).resolves.toBe(false)

    expect(store.isAuthenticated).toBe(false)
    expect(tokenStorage.readRefreshToken()).toBeNull()
    expect(store.refreshToken).toBeNull()
  })

  it('retries refresh recovery at most once before clearing the session', async () => {
    tokenStorage.writeRefreshToken('old-refresh')
    setActivePinia(createPinia())
    const storeA = useAuthStore()
    setActivePinia(createPinia())
    const storeB = useAuthStore()
    const user = {
      user_id: '1',
      username: 'admin',
      email: 'admin@beehive.local',
      nickname: 'Admin',
      avatar_url: '',
      role: 'admin',
      status: 'active',
    }
    let resolveFirstOldRefresh: ((value: Awaited<ReturnType<typeof authApi.refresh>>) => void) | null = null
    let rejectSecondOldRefresh: ((reason?: unknown) => void) | null = null
    const seenRefreshTokens: string[] = []
    vi.spyOn(authApi, 'refresh').mockImplementation((payload) => {
      seenRefreshTokens.push(payload.refresh_token)

      if (payload.refresh_token === 'old-refresh') {
        if (!resolveFirstOldRefresh) {
          return new Promise((resolve) => {
            resolveFirstOldRefresh = resolve
          })
        }
        return new Promise((_, reject) => {
          rejectSecondOldRefresh = reject
        })
      }

      return Promise.reject(new Error('still invalid'))
    })

    const firstRefresh = storeA.refreshSession()
    const secondRefresh = storeB.refreshSession()

    resolveFirstOldRefresh?.({
      access_token: 'access-1',
      refresh_token: 'new-refresh',
      expires_in: 900,
      token_type: 'Bearer',
      session_id: 's1',
      user,
      session: {
        session_id: 's1',
        user_id: '1',
        auth_source: 'local',
        client_type: 'web',
        status: 'active',
      },
    })
    await firstRefresh
    rejectSecondOldRefresh?.(new Error('invalid refresh token'))

    await expect(secondRefresh).resolves.toBe(false)
    expect(seenRefreshTokens).toEqual(['old-refresh', 'old-refresh', 'new-refresh'])
    expect(tokenStorage.readRefreshToken()).toBeNull()
    expect(storeB.refreshToken).toBeNull()
  })
})
