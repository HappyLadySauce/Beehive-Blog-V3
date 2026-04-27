import { requestJson } from '@/shared/api/httpClient'

import type { StudioApi } from './studioApi'
import type {
  StudioAuditListParams,
  StudioAuditsResponse,
  StudioMutationResponse,
  StudioUserListParams,
  StudioUserResponse,
  StudioUsersResponse,
  UserProfileResponse,
} from '../types'

const studioBasePath = '/api/v3/studio'
const authMeBasePath = '/api/v3/auth/me'

function withQuery(path: string, params?: StudioUserListParams | StudioAuditListParams): string {
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
  }
}
