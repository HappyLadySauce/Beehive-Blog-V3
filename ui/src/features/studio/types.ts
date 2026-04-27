import type { AuthUserProfile } from '@/features/auth/types'

export interface StudioUser {
  user_id: string
  username: string
  email: string
  nickname?: string
  avatar_url?: string
  role: string
  status: string
  last_login_at?: number
  created_at: number
  updated_at: number
}

export interface StudioAuditEvent {
  audit_id: string
  user_id?: string
  session_id?: string
  provider?: string
  auth_source?: string
  event_type: string
  result: 'success' | 'failure'
  client_ip?: string
  user_agent?: string
  detail_json?: string
  created_at: number
}

export interface StudioUsersResponse {
  items: StudioUser[]
  total: number
  page: number
  page_size: number
}

export interface StudioAuditsResponse {
  items: StudioAuditEvent[]
  total: number
  page: number
  page_size: number
}

export interface StudioUserListParams {
  keyword?: string
  role?: string
  status?: string
  page?: number
  page_size?: number
}

export interface StudioAuditListParams {
  event_type?: string
  result?: string
  user_id?: string
  started_at?: number
  ended_at?: number
  page?: number
  page_size?: number
}

export interface UserProfileUpdateRequest {
  nickname?: string
  avatar_url?: string
}

export interface UserProfileResponse {
  user: AuthUserProfile
}

export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

export interface StudioMutationResponse {
  ok: boolean
}

export interface StudioUserResponse {
  user: StudioUser
}

export interface StudioUpdateUserRoleRequest {
  role: 'member' | 'admin'
}

export interface StudioUpdateUserStatusRequest {
  status: 'active' | 'disabled' | 'locked'
}

export interface StudioResetPasswordRequest {
  new_password: string
}
