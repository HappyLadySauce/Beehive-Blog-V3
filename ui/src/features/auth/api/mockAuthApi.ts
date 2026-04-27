import type { AuthApi } from './authApi'
import type { AuthSessionView, AuthUserProfile } from '../types'

const mockAdminEmail = 'admin@beehive.local'
const mockRefreshPrefix = 'mock_refresh_'
const mockRefreshRoleSeparator = '__role_'

function buildMockUser(seed: string, role = 'member'): AuthUserProfile {
  const username = seed.includes('@') ? (seed.split('@')[0] ?? 'demo') : seed
  return {
    user_id: role === 'admin' ? 'user_mock_admin' : 'user_mock_member',
    username,
    email: seed.includes('@') ? seed : `${username}@beehive.local`,
    nickname: role === 'admin' ? 'Admin' : username,
    avatar_url: '',
    role,
    status: 'active',
  }
}

function buildMockSession(user: AuthUserProfile): AuthSessionView {
  const now = Math.floor(Date.now() / 1000)
  return {
    session_id: `sess_${user.user_id}`,
    user_id: user.user_id,
    auth_source: 'local',
    client_type: 'web',
    device_id: 'mock-browser',
    device_name: 'Mock Browser',
    status: 'active',
    last_seen_at: now,
    expires_at: now + 30 * 24 * 60 * 60,
  }
}

function buildMockRefreshToken(user: AuthUserProfile): string {
  return `${mockRefreshPrefix}${encodeURIComponent(user.email)}${mockRefreshRoleSeparator}${encodeURIComponent(user.role)}`
}

function parseMockRefreshToken(refreshToken: string): AuthUserProfile | null {
  if (!refreshToken.startsWith(mockRefreshPrefix)) {
    return null
  }
  const payload = refreshToken.slice(mockRefreshPrefix.length)
  const separatorIndex = payload.indexOf(mockRefreshRoleSeparator)
  if (separatorIndex < 1) {
    return null
  }
  const email = decodeURIComponent(payload.slice(0, separatorIndex))
  const role = decodeURIComponent(payload.slice(separatorIndex + mockRefreshRoleSeparator.length)) || 'member'
  return buildMockUser(email, role)
}

function buildAuthResponse(user: AuthUserProfile) {
  return {
    access_token: `mock_access_${user.user_id}_${Date.now()}`,
    refresh_token: buildMockRefreshToken(user),
    expires_in: 900,
    token_type: 'Bearer',
    session_id: `sess_${user.user_id}`,
    user,
    session: buildMockSession(user),
  }
}

export function createMockAuthApi(): AuthApi {
  let currentUser = buildMockUser(mockAdminEmail, 'admin')

  return {
    async register(payload) {
      currentUser = {
        ...buildMockUser(payload.email, 'member'),
        username: payload.username,
        nickname: payload.nickname ?? payload.username,
      }
      return buildAuthResponse(currentUser)
    },
    async login(payload) {
      const role = payload.login_identifier.toLowerCase() === mockAdminEmail ? 'admin' : 'member'
      currentUser = buildMockUser(payload.login_identifier, role)
      return buildAuthResponse(currentUser)
    },
    async refresh(payload) {
      const restoredUser = parseMockRefreshToken(payload.refresh_token)
      if (!restoredUser) {
        throw new Error('Invalid refresh token')
      }
      currentUser = restoredUser
      const response = buildAuthResponse(currentUser)
      return {
        access_token: response.access_token,
        refresh_token: response.refresh_token,
        expires_in: response.expires_in,
        token_type: response.token_type,
        session_id: response.session_id,
        session: response.session,
      }
    },
    async me() {
      return { user: currentUser }
    },
    async logout() {
      return { ok: true }
    },
    async startSso(payload) {
      return {
        provider: payload.provider,
        auth_url: `https://example.com/oauth/${payload.provider}?state=${encodeURIComponent(payload.state ?? 'mock_state')}`,
        state: payload.state ?? 'mock_state',
      }
    },
    async finishSso(payload) {
      currentUser = buildMockUser(`${payload.provider}_user@beehive.local`, 'member')
      return buildAuthResponse(currentUser)
    },
  }
}
