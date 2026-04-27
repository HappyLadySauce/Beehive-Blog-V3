import { appConfig, type ApiMode } from '@/shared/config/env'

import { createLiveStudioApi } from './liveStudioApi'
import { createMockStudioApi } from './mockStudioApi'
import type {
  ChangePasswordRequest,
  StudioAuditsResponse,
  StudioAuditListParams,
  StudioMutationResponse,
  StudioResetPasswordRequest,
  StudioUpdateUserRoleRequest,
  StudioUpdateUserStatusRequest,
  StudioUserListParams,
  StudioUserResponse,
  StudioUsersResponse,
  UserProfileResponse,
  UserProfileUpdateRequest,
} from '../types'

export interface StudioRequestOptions {
  accessToken?: string
}

export interface StudioApi {
  listUsers(params?: StudioUserListParams, options?: StudioRequestOptions): Promise<StudioUsersResponse>
  listAudits(params?: StudioAuditListParams, options?: StudioRequestOptions): Promise<StudioAuditsResponse>
  updateUserRole(userId: string, payload: StudioUpdateUserRoleRequest, options?: StudioRequestOptions): Promise<StudioUserResponse>
  updateUserStatus(userId: string, payload: StudioUpdateUserStatusRequest, options?: StudioRequestOptions): Promise<StudioUserResponse>
  resetUserPassword(userId: string, payload: StudioResetPasswordRequest, options?: StudioRequestOptions): Promise<StudioMutationResponse>
  updateProfile(payload: UserProfileUpdateRequest, options?: StudioRequestOptions): Promise<UserProfileResponse>
  changePassword(payload: ChangePasswordRequest, options?: StudioRequestOptions): Promise<StudioMutationResponse>
}

export function createStudioApi(mode: ApiMode = appConfig.apiMode): StudioApi {
  return mode === 'live' ? createLiveStudioApi() : createMockStudioApi()
}

export const studioApi = createStudioApi()
