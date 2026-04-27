export type AuthProvider = 'github' | 'qq' | 'wechat'

export interface AuthUserProfile {
  user_id: string
  username: string
  email: string
  nickname?: string
  avatar_url?: string
  role: string
  status: string
}

export interface AuthSessionView {
  session_id: string
  user_id: string
  auth_source: string
  client_type?: string
  device_id?: string
  device_name?: string
  status: string
  last_seen_at?: number
  expires_at?: number
}

export interface AuthRegisterRequest {
  username: string
  email: string
  password: string
  nickname?: string
}

export interface AuthRegisterResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  token_type: string
  session_id: string
  user: AuthUserProfile
  session: AuthSessionView
}

export interface AuthLoginRequest {
  login_identifier: string
  password: string
  client_type?: string
  device_id?: string
  device_name?: string
  user_agent?: string
}

export type AuthLoginResponse = AuthRegisterResponse

export interface AuthRefreshRequest {
  refresh_token: string
  user_agent?: string
}

export interface AuthRefreshResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  token_type: string
  session_id: string
  session: AuthSessionView
  user?: AuthUserProfile
}

export interface AuthMeResponse {
  user: AuthUserProfile
}

export interface AuthLogoutRequest {
  refresh_token?: string
}

export interface AuthLogoutResponse {
  ok: boolean
}

export interface AuthSsoStartRequest {
  provider: AuthProvider
  redirect_uri: string
  state?: string
}

export interface AuthSsoStartResponse {
  provider: AuthProvider
  auth_url: string
  state: string
}

export interface AuthSsoCallbackRequest {
  provider: AuthProvider
  code: string
  state: string
  redirect_uri: string
  client_type?: string
  device_id?: string
  device_name?: string
  user_agent?: string
}

export type AuthSsoCallbackResponse = AuthRegisterResponse
