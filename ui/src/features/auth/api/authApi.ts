import { appConfig, type ApiMode } from '@/shared/config/env'

import { createLiveAuthApi } from './liveAuthApi'
import { createMockAuthApi } from './mockAuthApi'
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
  AuthSsoCallbackRequest,
  AuthSsoCallbackResponse,
  AuthSsoStartRequest,
  AuthSsoStartResponse,
} from '../types'

export interface AuthRequestOptions {
  accessToken?: string
}

export interface AuthApi {
  register(payload: AuthRegisterRequest): Promise<AuthRegisterResponse>
  login(payload: AuthLoginRequest): Promise<AuthLoginResponse>
  refresh(payload: AuthRefreshRequest): Promise<AuthRefreshResponse>
  me(options?: AuthRequestOptions): Promise<AuthMeResponse>
  logout(payload?: AuthLogoutRequest, options?: AuthRequestOptions): Promise<AuthLogoutResponse>
  startSso(payload: AuthSsoStartRequest): Promise<AuthSsoStartResponse>
  finishSso(payload: AuthSsoCallbackRequest): Promise<AuthSsoCallbackResponse>
}

export function createAuthApi(mode: ApiMode = appConfig.apiMode): AuthApi {
  return mode === 'live' ? createLiveAuthApi() : createMockAuthApi()
}

export const authApi = createAuthApi()
