import type { AuthUserProfile } from '@/features/auth/types'

import type { StudioApi } from './studioApi'
import type { StudioAuditEvent, StudioUser } from '../types'

const mockUsers: StudioUser[] = [
  {
    user_id: 'user_mock_admin',
    username: 'admin',
    email: 'admin@beehive.local',
    nickname: 'Admin',
    role: 'admin',
    status: 'active',
    created_at: 1776676320,
    updated_at: 1777276200,
    last_login_at: 1777276200,
  },
  {
    user_id: 'user_mock_editor',
    username: 'editor',
    email: 'editor@beehive.local',
    nickname: 'Editor',
    role: 'member',
    status: 'active',
    created_at: 1776781080,
    updated_at: 1777219320,
    last_login_at: 1777219320,
  },
  {
    user_id: 'user_mock_member',
    username: 'member',
    email: 'member@beehive.local',
    nickname: 'Member',
    role: 'member',
    status: 'disabled',
    created_at: 1776858060,
    updated_at: 1776858060,
  },
]

const mockAudits: StudioAuditEvent[] = [
  {
    audit_id: 'audit_1001',
    user_id: 'user_mock_admin',
    auth_source: 'local',
    event_type: 'login',
    result: 'success',
    client_ip: '127.0.0.1',
    detail_json: '{"surface":"studio"}',
    created_at: 1777275840,
  },
  {
    audit_id: 'audit_1002',
    user_id: 'user_mock_admin',
    auth_source: 'local',
    event_type: 'admin_update_user_status',
    result: 'success',
    client_ip: '127.0.0.1',
    detail_json: '{"target_user_id":"user_mock_member","old_status":"active","new_status":"disabled"}',
    created_at: 1777218000,
  },
  {
    audit_id: 'audit_1003',
    user_id: 'user_mock_member',
    auth_source: 'local',
    event_type: 'studio_access',
    result: 'failure',
    client_ip: '127.0.0.1',
    detail_json: '{"reason":"forbidden"}',
    created_at: 1777187580,
  },
]

let mockProfile: AuthUserProfile = {
  user_id: 'user_mock_admin',
  username: 'admin',
  email: 'admin@beehive.local',
  nickname: 'Admin',
  avatar_url: '',
  role: 'admin',
  status: 'active',
}

function rejectWeakPassword(password: string): void {
  if (password.length < 12 || !/[A-Z]/.test(password) || !/[0-9]/.test(password)) {
    throw new Error('New password must be at least 12 characters and include uppercase letters and numbers.')
  }
}

export function createMockStudioApi(): StudioApi {
  return {
    async listUsers(params) {
      const keyword = params?.keyword?.toLowerCase().trim()
      const filtered = mockUsers.filter((user) => {
        const matchesKeyword = !keyword
          || user.username.toLowerCase().includes(keyword)
          || user.email.toLowerCase().includes(keyword)
          || user.nickname?.toLowerCase().includes(keyword)
        const matchesRole = !params?.role || user.role === params.role
        const matchesStatus = !params?.status || user.status === params.status
        return matchesKeyword && matchesRole && matchesStatus
      })
      return { items: [...filtered], total: filtered.length, page: params?.page ?? 1, page_size: params?.page_size ?? 20 }
    },
    async listAudits(params) {
      const filtered = mockAudits.filter((event) => {
        const matchesEventType = !params?.event_type || event.event_type.includes(params.event_type)
        const matchesResult = !params?.result || event.result === params.result
        const matchesUser = !params?.user_id || event.user_id === params.user_id
        return matchesEventType && matchesResult && matchesUser
      })
      return { items: [...filtered], total: filtered.length, page: params?.page ?? 1, page_size: params?.page_size ?? 20 }
    },
    async updateUserRole(userId, payload) {
      const user = mockUsers.find((item) => item.user_id === userId)
      if (!user) {
        throw new Error('User not found.')
      }
      user.role = payload.role
      user.updated_at = Math.floor(Date.now() / 1000)
      return { user: { ...user } }
    },
    async updateUserStatus(userId, payload) {
      const user = mockUsers.find((item) => item.user_id === userId)
      if (!user) {
        throw new Error('User not found.')
      }
      user.status = payload.status
      user.updated_at = Math.floor(Date.now() / 1000)
      return { user: { ...user } }
    },
    async resetUserPassword(userId, payload) {
      if (!mockUsers.some((item) => item.user_id === userId)) {
        throw new Error('User not found.')
      }
      rejectWeakPassword(payload.new_password)
      return { ok: true }
    },
    async updateProfile(payload) {
      mockProfile = {
        ...mockProfile,
        nickname: payload.nickname,
        avatar_url: payload.avatar_url ?? '',
      }
      return { user: mockProfile }
    },
    async changePassword(payload) {
      if (payload.old_password.length === 0) {
        throw new Error('Current password is required.')
      }
      rejectWeakPassword(payload.new_password)
      return { ok: true }
    },
  }
}
