import { requestJson } from '@/shared/api/httpClient'

import type { StudioApi } from './studioApi'
import type {
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
  StudioAuditListParams,
  StudioAuditsResponse,
  StudioMutationResponse,
  StudioUserListParams,
  StudioUserResponse,
  StudioUsersResponse,
  UserProfileResponse,
} from '../types'

const studioBasePath = '/api/v3/studio'
const studioContentBasePath = '/api/v3/studio/content'
const authMeBasePath = '/api/v3/auth/me'

function withQuery(path: string, params?: Record<string, string | number | boolean | undefined>): string {
  const searchParams = new URLSearchParams()
  for (const [key, value] of Object.entries(params ?? {})) {
    if (value !== undefined && value !== '') {
      searchParams.set(key, String(value))
    }
  }
  const query = searchParams.toString()
  return query.length > 0 ? `${path}?${query}` : path
}

export function createLiveStudioApi(): StudioApi {
  return {
    listUsers(params, options) {
      return requestJson<StudioUsersResponse>(withQuery(`${studioBasePath}/users`, params), {
        method: 'GET',
        accessToken: options?.accessToken,
      })
    },
    listAudits(params, options) {
      return requestJson<StudioAuditsResponse>(withQuery(`${studioBasePath}/audits`, params), {
        method: 'GET',
        accessToken: options?.accessToken,
      })
    },
    deleteUser(userId, options) {
      return requestJson<StudioMutationResponse>(`${studioBasePath}/users/${encodeURIComponent(userId)}`, {
        method: 'DELETE',
        accessToken: options?.accessToken,
      })
    },
    updateUserRole(userId, payload, options) {
      return requestJson<StudioUserResponse>(`${studioBasePath}/users/${encodeURIComponent(userId)}/role`, {
        method: 'PATCH',
        body: JSON.stringify(payload),
        accessToken: options?.accessToken,
      })
    },
    updateUserStatus(userId, payload, options) {
      return requestJson<StudioUserResponse>(`${studioBasePath}/users/${encodeURIComponent(userId)}/status`, {
        method: 'PATCH',
        body: JSON.stringify(payload),
        accessToken: options?.accessToken,
      })
    },
    resetUserPassword(userId, payload, options) {
      return requestJson<StudioMutationResponse>(`${studioBasePath}/users/${encodeURIComponent(userId)}/password/reset`, {
        method: 'POST',
        body: JSON.stringify(payload),
        accessToken: options?.accessToken,
      })
    },
    updateProfile(payload, options) {
      return requestJson<UserProfileResponse>(`${authMeBasePath}/profile`, {
        method: 'PATCH',
        body: JSON.stringify(payload),
        accessToken: options?.accessToken,
      })
    },
    changePassword(payload, options) {
      return requestJson<StudioMutationResponse>(`${authMeBasePath}/password`, {
        method: 'POST',
        body: JSON.stringify(payload),
        accessToken: options?.accessToken,
      })
    },
    listContents(params: ContentListParams = {}, options) {
      return requestJson<ContentListResponse>(withQuery(`${studioContentBasePath}/items`, params), {
        method: 'GET',
        accessToken: options?.accessToken,
      })
    },
    createContent(payload: ContentWriteRequest, options) {
      return requestJson<ContentDetailResponse>(`${studioContentBasePath}/items`, {
        method: 'POST',
        body: JSON.stringify(payload),
        accessToken: options?.accessToken,
      })
    },
    getContent(contentId, options) {
      return requestJson<ContentDetailResponse>(`${studioContentBasePath}/items/${encodeURIComponent(contentId)}`, {
        method: 'GET',
        accessToken: options?.accessToken,
      })
    },
    updateContent(contentId, payload: ContentWriteRequest, options) {
      return requestJson<ContentDetailResponse>(`${studioContentBasePath}/items/${encodeURIComponent(contentId)}`, {
        method: 'PUT',
        body: JSON.stringify(payload),
        accessToken: options?.accessToken,
      })
    },
    archiveContent(contentId, options) {
      return requestJson<ContentMutationResponse>(`${studioContentBasePath}/items/${encodeURIComponent(contentId)}`, {
        method: 'DELETE',
        accessToken: options?.accessToken,
      })
    },
    listTags(params: ContentTagListParams = {}, options) {
      return requestJson<ContentTagListResponse>(withQuery(`${studioContentBasePath}/tags`, params), {
        method: 'GET',
        accessToken: options?.accessToken,
      })
    },
    createTag(payload: ContentTagWriteRequest, options) {
      return requestJson<ContentTagResponse>(`${studioContentBasePath}/tags`, {
        method: 'POST',
        body: JSON.stringify(payload),
        accessToken: options?.accessToken,
      })
    },
    updateTag(tagId, payload: ContentTagWriteRequest, options) {
      return requestJson<ContentTagResponse>(`${studioContentBasePath}/tags/${encodeURIComponent(tagId)}`, {
        method: 'PUT',
        body: JSON.stringify(payload),
        accessToken: options?.accessToken,
      })
    },
    deleteTag(tagId, options) {
      return requestJson<ContentMutationResponse>(`${studioContentBasePath}/tags/${encodeURIComponent(tagId)}`, {
        method: 'DELETE',
        accessToken: options?.accessToken,
      })
    },
    listRelations(contentId, params: ContentRelationListParams = {}, options) {
      return requestJson<ContentRelationListResponse>(withQuery(`${studioContentBasePath}/items/${encodeURIComponent(contentId)}/relations`, params), {
        method: 'GET',
        accessToken: options?.accessToken,
      })
    },
    createRelation(contentId, payload: ContentRelationWriteRequest, options) {
      return requestJson<ContentRelationResponse>(`${studioContentBasePath}/items/${encodeURIComponent(contentId)}/relations`, {
        method: 'POST',
        body: JSON.stringify(payload),
        accessToken: options?.accessToken,
      })
    },
    deleteRelation(contentId, relationId, options) {
      return requestJson<ContentMutationResponse>(`${studioContentBasePath}/items/${encodeURIComponent(contentId)}/relations/${encodeURIComponent(relationId)}`, {
        method: 'DELETE',
        accessToken: options?.accessToken,
      })
    },
    listRevisions(contentId, params = {}, options) {
      return requestJson<ContentRevisionListResponse>(withQuery(`${studioContentBasePath}/items/${encodeURIComponent(contentId)}/revisions`, params), {
        method: 'GET',
        accessToken: options?.accessToken,
      })
    },
    getRevision(contentId, revisionId, options) {
      return requestJson<ContentRevisionResponse>(`${studioContentBasePath}/items/${encodeURIComponent(contentId)}/revisions/${encodeURIComponent(revisionId)}`, {
        method: 'GET',
        accessToken: options?.accessToken,
      })
    },
  }
}
