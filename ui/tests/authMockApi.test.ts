import { describe, expect, it } from 'vitest'
import { createMockAuthApi } from '../src/features/auth/api/mockAuthApi'

describe('mockAuthApi', () => {
  it('supports register and login token flows', async () => {
    const api = createMockAuthApi()

    const registered = await api.register({
      username: 'new_user',
      email: 'new@example.com',
      password: 'Str0ngP@ssw0rd!',
    })
    const loggedIn = await api.login({
      login_identifier: 'new@example.com',
      password: 'Str0ngP@ssw0rd!',
    })

    expect(registered.user.username).toBe('new_user')
    expect(registered.refresh_token).toMatch(/^mock_refresh_/)
    expect(loggedIn.user.email).toBe('new@example.com')
    expect(loggedIn.access_token).toMatch(/^mock_access_/)
  })

  it('refreshes without embedding a user profile', async () => {
    const api = createMockAuthApi()
    const loggedIn = await api.login({
      login_identifier: 'member@beehive.local',
      password: 'Str0ngP@ssw0rd!',
    })

    const refreshed = await api.refresh({ refresh_token: loggedIn.refresh_token })

    expect(refreshed.access_token).toMatch(/^mock_access_/)
    expect(refreshed.user).toBeUndefined()
  })
})
