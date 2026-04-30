import { describe, expect, it } from 'vitest'

import { tokenStorage } from '@/shared/storage/tokenStorage'

describe('tokenStorage', () => {
  it('reads, writes and clears auth snapshots', () => {
    tokenStorage.setSnapshot({
      accessToken: 'access-token',
      refreshToken: 'refresh-token',
      accessTokenExpiresAt: Date.now() + 60_000,
      sessionId: 'session-id',
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

    expect(tokenStorage.getSnapshot()).toMatchObject({
      accessToken: 'access-token',
      refreshToken: 'refresh-token',
      sessionId: 'session-id',
    })
    expect(tokenStorage.readRefreshToken()).toBe('refresh-token')

    tokenStorage.clearSnapshot()
    expect(tokenStorage.getSnapshot()).toBeNull()
    expect(tokenStorage.readRefreshToken()).toBeNull()
  })

  it('returns null for malformed auth snapshots', () => {
    window.localStorage.setItem('beehive.v3.auth.session_snapshot', '{broken-json')

    expect(tokenStorage.getSnapshot()).toBeNull()
  })

  it('reads, writes and clears refresh tokens', () => {
    tokenStorage.writeRefreshToken('refresh-token')
    expect(tokenStorage.readRefreshToken()).toBe('refresh-token')

    tokenStorage.clearRefreshToken()
    expect(tokenStorage.readRefreshToken()).toBeNull()
  })
})
