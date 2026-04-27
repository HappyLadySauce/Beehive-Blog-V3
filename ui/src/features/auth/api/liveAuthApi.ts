import { requestJson } from '@/shared/api/httpClient'

import type { AuthApi } from './authApi'
import type {
  AuthLoginResponse,
  AuthLogoutResponse,
  AuthMeResponse,
  AuthRefreshResponse,
  AuthRegisterResponse,
  AuthSsoCallbackResponse,
  AuthSsoStartResponse,
} from '../types'

const authBasePath = '/api/v3/auth'

export function createLiveAuthApi(): AuthApi {
  return {
    register(payload) {
      return requestJson<AuthRegisterResponse>(`${authBasePath}/register`, {
        method: 'POST',
        body: JSON.stringify(payload),
      })
    },
    login(payload) {
      return requestJson<AuthLoginResponse>(`${authBasePath}/login`, {
        method: 'POST',
        body: JSON.stringify(payload),
      })
    },
    refresh(payload) {
      return requestJson<AuthRefreshResponse>(`${authBasePath}/refresh`, {
        method: 'POST',
        body: JSON.stringify(payload),
      })
    },
    me(options) {
      return requestJson<AuthMeResponse>(`${authBasePath}/me`, {
        method: 'GET',
        accessToken: options?.accessToken,
      })
    },
    logout(payload, options) {
      return requestJson<AuthLogoutResponse>(`${authBasePath}/logout`, {
        method: 'POST',
        body: JSON.stringify(payload ?? {}),
        accessToken: options?.accessToken,
      })
    },
    startSso(payload) {
      return requestJson<AuthSsoStartResponse>(`${authBasePath}/sso/start`, {
        method: 'POST',
        body: JSON.stringify(payload),
      })
    },
    finishSso(payload) {
      return requestJson<AuthSsoCallbackResponse>(`${authBasePath}/sso/callback`, {
        method: 'POST',
        body: JSON.stringify(payload),
      })
    },
  }
}
