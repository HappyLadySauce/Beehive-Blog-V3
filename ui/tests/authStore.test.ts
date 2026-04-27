import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

import { useAuthStore } from '@/features/auth/stores/authStore'
import { tokenStorage } from '@/shared/storage/tokenStorage'

describe('authStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
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
})
