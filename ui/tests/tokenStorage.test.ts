import { describe, expect, it } from 'vitest'

import { tokenStorage } from '@/shared/storage/tokenStorage'

describe('tokenStorage', () => {
  it('reads, writes and clears refresh tokens', () => {
    tokenStorage.writeRefreshToken('refresh-token')
    expect(tokenStorage.readRefreshToken()).toBe('refresh-token')

    tokenStorage.clearRefreshToken()
    expect(tokenStorage.readRefreshToken()).toBeNull()
  })
})
