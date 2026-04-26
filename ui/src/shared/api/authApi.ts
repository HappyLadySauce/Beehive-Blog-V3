import { requestJson } from '@/shared/api/httpClient';
import type {
  AuthLoginRequest,
  AuthLoginResponse,
  AuthLogoutRequest,
  AuthLogoutResponse,
  AuthMeResponse,
  AuthRefreshRequest,
  AuthRefreshResponse,
  AuthRegisterRequest,
  AuthRegisterResponse,
  AuthSessionView,
  AuthUserProfile,
} from '@/shared/api/types';
import type { ApiMode } from '@/shared/config/env';
import { appConfig } from '@/shared/config/env';

export interface AuthRequestOptions {
  accessToken?: string;
}

export interface AuthApi {
  register(payload: AuthRegisterRequest): Promise<AuthRegisterResponse>;
  login(payload: AuthLoginRequest): Promise<AuthLoginResponse>;
  refresh(payload: AuthRefreshRequest): Promise<AuthRefreshResponse>;
  me(options?: AuthRequestOptions): Promise<AuthMeResponse>;
  logout(payload?: AuthLogoutRequest, options?: AuthRequestOptions): Promise<AuthLogoutResponse>;
}

function buildMockUser(seed: string): AuthUserProfile {
  const username = seed.includes('@') ? seed.split('@')[0] ?? 'demo_user' : seed;
  return {
    user_id: 'user_mock_001',
    username,
    email: seed.includes('@') ? seed : `${username}@beehive.local`,
    nickname: username,
    avatar_url: '',
    role: 'member',
    status: 'active',
  };
}

function buildMockSession(user: AuthUserProfile): AuthSessionView {
  const now = Math.floor(Date.now() / 1000);
  return {
    session_id: 'sess_mock_001',
    user_id: user.user_id,
    auth_source: 'password',
    client_type: 'web',
    device_id: 'mock-browser',
    device_name: 'Mock Browser',
    status: 'active',
    last_seen_at: now,
    expires_at: now + 30 * 24 * 60 * 60,
  };
}

function buildAuthResponse(user: AuthUserProfile): AuthRegisterResponse {
  return {
    access_token: `mock_access_${user.user_id}`,
    refresh_token: `mock_refresh_${user.user_id}`,
    expires_in: 900,
    token_type: 'Bearer',
    session_id: 'sess_mock_001',
    user,
    session: buildMockSession(user),
  };
}

function buildRefreshResponse(user: AuthUserProfile): AuthRefreshResponse {
  return {
    access_token: `mock_access_${user.user_id}`,
    refresh_token: `mock_refresh_${user.user_id}`,
    expires_in: 900,
    token_type: 'Bearer',
    session_id: 'sess_mock_001',
    session: buildMockSession(user),
  };
}

function createMockAuthApi(): AuthApi {
  let currentUser = buildMockUser('demo@beehive.local');

  return {
    async register(payload) {
      currentUser = buildMockUser(payload.email);
      return buildAuthResponse({
        ...currentUser,
        username: payload.username,
        nickname: payload.nickname ?? payload.username,
      });
    },
    async login(payload) {
      currentUser = buildMockUser(payload.login_identifier);
      return buildAuthResponse(currentUser);
    },
    async refresh(payload) {
      if (!payload.refresh_token.startsWith('mock_refresh_')) {
        throw new Error('Invalid refresh token');
      }
      return buildRefreshResponse(currentUser);
    },
    async me() {
      return { user: currentUser };
    },
    async logout() {
      return { ok: true };
    },
  };
}

function createLiveAuthApi(): AuthApi {
  return {
    register(payload) {
      return requestJson<AuthRegisterResponse>('/api/v3/auth/register', {
        method: 'POST',
        body: JSON.stringify(payload),
      });
    },
    login(payload) {
      return requestJson<AuthLoginResponse>('/api/v3/auth/login', {
        method: 'POST',
        body: JSON.stringify(payload),
      });
    },
    refresh(payload) {
      return requestJson<AuthRefreshResponse>('/api/v3/auth/refresh', {
        method: 'POST',
        body: JSON.stringify(payload),
      });
    },
    me(options) {
      const requestOptions: { method: 'GET'; accessToken?: string } = {
        method: 'GET',
      };
      if (options?.accessToken) {
        requestOptions.accessToken = options.accessToken;
      }
      return requestJson<AuthMeResponse>('/api/v3/auth/me', requestOptions);
    },
    logout(payload, options) {
      const requestOptions: { method: 'POST'; body: string; accessToken?: string } = {
        method: 'POST',
        body: JSON.stringify(payload ?? {}),
      };
      if (options?.accessToken) {
        requestOptions.accessToken = options.accessToken;
      }
      return requestJson<AuthLogoutResponse>('/api/v3/auth/logout', requestOptions);
    },
  };
}

export function createAuthApi(mode: ApiMode = appConfig.apiMode): AuthApi {
  return mode === 'live' ? createLiveAuthApi() : createMockAuthApi();
}

export const authApi = createAuthApi();
