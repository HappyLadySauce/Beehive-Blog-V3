import { createPinia, setActivePinia } from 'pinia'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { authApi } from '@/features/auth/api/authApi'
import { useAuthStore } from '@/features/auth/stores/authStore'
import { tokenStorage } from '@/shared/storage/tokenStorage'

describe('authStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('logs in, stores refresh token, restores and logs out', async () => {
    const store = useAuthStore()
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

    await store.logout()
    expect(store.isAuthenticated).toBe(false)
    expect(store.refreshToken).toBeNull()
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
})
