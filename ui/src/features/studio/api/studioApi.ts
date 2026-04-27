import { appConfig, type ApiMode } from '@/shared/config/env'

import { createLiveStudioApi } from './liveStudioApi'
import { createMockStudioApi } from './mockStudioApi'
import type {
  ChangePasswordRequest,
  ContentDetailResponse,
  ContentListParams,
  ContentListResponse,
  ContentMutationResponse,
  ContentRelationListParams,
  ContentRelationListResponse,
  ContentRelationResponse,
  ContentRelationWriteRequest,
  ContentRevisionListResponse,
  ContentRevisionResponse,
  ContentTagListParams,
  ContentTagListResponse,
  ContentTagResponse,
  ContentTagWriteRequest,
  ContentWriteRequest,
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
  deleteUser(userId: string, options?: StudioRequestOptions): Promise<StudioMutationResponse>
  updateUserRole(userId: string, payload: StudioUpdateUserRoleRequest, options?: StudioRequestOptions): Promise<StudioUserResponse>
  updateUserStatus(userId: string, payload: StudioUpdateUserStatusRequest, options?: StudioRequestOptions): Promise<StudioUserResponse>
  resetUserPassword(userId: string, payload: StudioResetPasswordRequest, options?: StudioRequestOptions): Promise<StudioMutationResponse>
  updateProfile(payload: UserProfileUpdateRequest, options?: StudioRequestOptions): Promise<UserProfileResponse>
  changePassword(payload: ChangePasswordRequest, options?: StudioRequestOptions): Promise<StudioMutationResponse>
  listContents(params?: ContentListParams, options?: StudioRequestOptions): Promise<ContentListResponse>
  createContent(payload: ContentWriteRequest, options?: StudioRequestOptions): Promise<ContentDetailResponse>
  getContent(contentId: string, options?: StudioRequestOptions): Promise<ContentDetailResponse>
  updateContent(contentId: string, payload: ContentWriteRequest, options?: StudioRequestOptions): Promise<ContentDetailResponse>
  archiveContent(contentId: string, options?: StudioRequestOptions): Promise<ContentMutationResponse>
  listTags(params?: ContentTagListParams, options?: StudioRequestOptions): Promise<ContentTagListResponse>
  createTag(payload: ContentTagWriteRequest, options?: StudioRequestOptions): Promise<ContentTagResponse>
  updateTag(tagId: string, payload: ContentTagWriteRequest, options?: StudioRequestOptions): Promise<ContentTagResponse>
  deleteTag(tagId: string, options?: StudioRequestOptions): Promise<ContentMutationResponse>
  listRelations(contentId: string, params?: ContentRelationListParams, options?: StudioRequestOptions): Promise<ContentRelationListResponse>
  createRelation(contentId: string, payload: ContentRelationWriteRequest, options?: StudioRequestOptions): Promise<ContentRelationResponse>
  deleteRelation(contentId: string, relationId: string, options?: StudioRequestOptions): Promise<ContentMutationResponse>
  listRevisions(contentId: string, params?: { page?: number; page_size?: number }, options?: StudioRequestOptions): Promise<ContentRevisionListResponse>
  getRevision(contentId: string, revisionId: string, options?: StudioRequestOptions): Promise<ContentRevisionResponse>
}

export function createStudioApi(mode: ApiMode = appConfig.apiMode): StudioApi {
  return mode === 'live' ? createLiveStudioApi() : createMockStudioApi()
}

export const studioApi = createStudioApi()
