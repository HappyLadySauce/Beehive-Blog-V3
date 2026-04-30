import type { AuthUserProfile } from '@/features/auth/types'

import type { StudioApi } from './studioApi'
import type {
  ContentDetail,
  ContentRelation,
  ContentRevisionSummary,
  ContentTag,
  ContentWriteRequest,
  StudioAuditEvent,
  StudioUser,
} from '../types'

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
  {
    user_id: 'user_mock_writer',
    username: 'writer',
    email: 'writer@beehive.local',
    nickname: 'Writer',
    role: 'member',
    status: 'active',
    created_at: 1776901200,
    updated_at: 1777260000,
    last_login_at: 1777260000,
  },
  {
    user_id: 'user_mock_ops',
    username: 'ops',
    email: 'ops@beehive.local',
    nickname: 'Ops',
    role: 'admin',
    status: 'active',
    created_at: 1776912000,
    updated_at: 1777256400,
    last_login_at: 1777256400,
  },
  {
    user_id: 'user_mock_support',
    username: 'support',
    email: 'support@beehive.local',
    nickname: 'Support',
    role: 'member',
    status: 'active',
    created_at: 1776922800,
    updated_at: 1777249200,
    last_login_at: 1777249200,
  },
  {
    user_id: 'user_mock_qa',
    username: 'qa',
    email: 'qa@beehive.local',
    nickname: 'QA',
    role: 'member',
    status: 'disabled',
    created_at: 1776933600,
    updated_at: 1777242000,
    last_login_at: 1777242000,
  },
  {
    user_id: 'user_mock_growth',
    username: 'growth',
    email: 'growth@beehive.local',
    nickname: 'Growth',
    role: 'member',
    status: 'active',
    created_at: 1776944400,
    updated_at: 1777234800,
    last_login_at: 1777234800,
  },
  {
    user_id: 'user_mock_design',
    username: 'design',
    email: 'design@beehive.local',
    nickname: 'Design',
    role: 'member',
    status: 'active',
    created_at: 1776955200,
    updated_at: 1777227600,
    last_login_at: 1777227600,
  },
  {
    user_id: 'user_mock_data',
    username: 'data',
    email: 'data@beehive.local',
    nickname: 'Data',
    role: 'member',
    status: 'locked',
    created_at: 1776966000,
    updated_at: 1777220400,
  },
  {
    user_id: 'user_mock_release',
    username: 'release',
    email: 'release@beehive.local',
    nickname: 'Release',
    role: 'admin',
    status: 'active',
    created_at: 1776976800,
    updated_at: 1777213200,
    last_login_at: 1777213200,
  },
  {
    user_id: 'user_mock_docs',
    username: 'docs',
    email: 'docs@beehive.local',
    nickname: 'Docs',
    role: 'member',
    status: 'active',
    created_at: 1776987600,
    updated_at: 1777206000,
    last_login_at: 1777206000,
  },
]

const mockAudits: StudioAuditEvent[] = [
  {
    audit_id: 'audit_1001',
    user_id: 'user_mock_admin',
    username: 'admin',
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
    username: 'admin',
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
    username: 'member',
    auth_source: 'local',
    event_type: 'studio_access',
    result: 'failure',
    client_ip: '127.0.0.1',
    detail_json: '{"reason":"forbidden"}',
    created_at: 1777187580,
  },
  {
    audit_id: 'audit_1004',
    user_id: 'user_mock_editor',
    username: 'editor',
    auth_source: 'local',
    event_type: 'refresh_session_token',
    result: 'success',
    client_ip: '127.0.0.1',
    detail_json: '-',
    created_at: 1777186500,
  },
  {
    audit_id: 'audit_1005',
    user_id: 'user_mock_admin',
    username: 'admin',
    auth_source: 'local',
    event_type: 'refresh_session_token',
    result: 'success',
    client_ip: '127.0.0.1',
    detail_json: '-',
    created_at: 1777186200,
  },
  {
    audit_id: 'audit_1006',
    user_id: 'user_mock_admin',
    username: 'admin',
    auth_source: 'local',
    event_type: 'refresh_session_token',
    result: 'success',
    client_ip: '127.0.0.1',
    detail_json: '-',
    created_at: 1777185900,
  },
  {
    audit_id: 'audit_1007',
    user_id: 'user_mock_support',
    username: 'support',
    auth_source: 'local',
    event_type: 'password_reset',
    result: 'success',
    client_ip: '127.0.0.1',
    detail_json: '{"target_user_id":"user_mock_member"}',
    created_at: 1777185600,
  },
  {
    audit_id: 'audit_1008',
    user_id: 'user_mock_qa',
    username: 'qa',
    auth_source: 'local',
    event_type: 'login',
    result: 'failure',
    client_ip: '127.0.0.1',
    detail_json: '{"reason":"invalid_password"}',
    created_at: 1777185300,
  },
  {
    audit_id: 'audit_1009',
    user_id: 'user_mock_ops',
    username: 'ops',
    auth_source: 'local',
    event_type: 'login',
    result: 'success',
    client_ip: '127.0.0.1',
    detail_json: '{"surface":"studio"}',
    created_at: 1777185000,
  },
  {
    audit_id: 'audit_1010',
    user_id: 'user_mock_release',
    username: 'release',
    auth_source: 'local',
    event_type: 'refresh_session_token',
    result: 'success',
    client_ip: '127.0.0.1',
    detail_json: '-',
    created_at: 1777184700,
  },
  {
    audit_id: 'audit_1011',
    user_id: 'user_mock_docs',
    username: 'docs',
    auth_source: 'local',
    event_type: 'role_changed',
    result: 'success',
    client_ip: '127.0.0.1',
    detail_json: '{"new_role":"member"}',
    created_at: 1777184400,
  },
  {
    audit_id: 'audit_1012',
    user_id: 'user_mock_admin',
    username: 'admin',
    auth_source: 'local',
    event_type: 'content_archive',
    result: 'success',
    client_ip: '127.0.0.1',
    detail_json: '{"content_id":"1001"}',
    created_at: 1777184100,
  },
]

const mockTags: ContentTag[] = [
  { tag_id: 'tag_1001', name: 'Gateway', slug: 'gateway', description: 'Gateway notes', color: '#0f8f83', created_at: 1776676320, updated_at: 1776676320 },
  { tag_id: 'tag_1002', name: 'Identity', slug: 'identity', description: 'Identity flows', color: '#e45f35', created_at: 1776676320, updated_at: 1776676320 },
  { tag_id: 'tag_1003', name: 'Release', slug: 'release', description: 'Release process', color: '#5a67d8', created_at: 1776676320, updated_at: 1776676320 },
  { tag_id: 'tag_1004', name: 'Design', slug: 'design', description: 'Design system', color: '#d97706', created_at: 1776676320, updated_at: 1776676320 },
  { tag_id: 'tag_1005', name: 'Platform', slug: 'platform', description: 'Platform changes', color: '#0ea5e9', created_at: 1776676320, updated_at: 1776676320 },
  { tag_id: 'tag_1006', name: 'Infra', slug: 'infra', description: 'Infrastructure work', color: '#64748b', created_at: 1776676320, updated_at: 1776676320 },
  { tag_id: 'tag_1007', name: 'Docs', slug: 'docs', description: 'Documentation', color: '#84cc16', created_at: 1776676320, updated_at: 1776676320 },
  { tag_id: 'tag_1008', name: 'Search', slug: 'search', description: 'Search relevance', color: '#ec4899', created_at: 1776676320, updated_at: 1776676320 },
  { tag_id: 'tag_1009', name: 'Analytics', slug: 'analytics', description: 'Analytics notes', color: '#8b5cf6', created_at: 1776676320, updated_at: 1776676320 },
  { tag_id: 'tag_1010', name: 'Observability', slug: 'observability', description: 'Tracing and logs', color: '#14b8a6', created_at: 1776676320, updated_at: 1776676320 },
  { tag_id: 'tag_1011', name: 'Auth', slug: 'auth', description: 'Authentication guides', color: '#ef4444', created_at: 1776676320, updated_at: 1776676320 },
]

const mockContents: ContentDetail[] = [
  {
    content_id: '1001',
    type: 'article',
    title: 'v3 frontend integration notes',
    slug: 'v3-frontend-integration-notes',
    summary: 'Gateway-first Studio integration notes.',
    body_markdown: 'Gateway-first Studio integration notes.',
    body_json: '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Gateway-first Studio integration notes."}]}]}',
    cover_image_url: '',
    status: 'draft',
    visibility: 'private',
    ai_access: 'denied',
    owner_user_id: 'user_mock_admin',
    author_user_id: 'user_mock_admin',
    source_type: 'manual',
    current_revision_id: 'rev_1001',
    comment_enabled: true,
    is_featured: false,
    sort_order: 0,
    created_at: 1777187580,
    updated_at: 1777273980,
    tags: [mockTags[0]],
  },
]

const mockRelations: ContentRelation[] = []

const mockRevisions: ContentRevisionSummary[] = [
  {
    revision_id: 'rev_1001',
    content_id: '1001',
    revision_no: 1,
    editor_type: 'human',
    change_summary: 'Initial draft',
    source_type: 'manual',
    created_at: 1777273980,
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
  if (password.length < 8) {
    throw new Error('New password must be at least 8 characters.')
  }
}

function normalizePage(page?: number): number {
  if (!Number.isInteger(page) || (page ?? 0) < 1) {
    return 1
  }
  return page as number
}

function normalizePageSize(pageSize?: number): number {
  if (!Number.isInteger(pageSize) || (pageSize ?? 0) < 1) {
    return 20
  }
  return pageSize as number
}

function paginateItems<T>(items: T[], page?: number, pageSize?: number) {
  const normalizedPage = normalizePage(page)
  const normalizedPageSize = normalizePageSize(pageSize)
  const start = (normalizedPage - 1) * normalizedPageSize
  const end = start + normalizedPageSize
  return {
    items: items.slice(start, end),
    page: normalizedPage,
    page_size: normalizedPageSize,
    total: items.length,
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
        const matchesDeleted = params?.include_deleted || user.status !== 'deleted'
        return matchesKeyword && matchesRole && matchesStatus && matchesDeleted
      })
      return paginateItems([...filtered], params?.page, params?.page_size)
    },
    async listAudits(params) {
      const filtered = mockAudits.filter((event) => {
        const matchesEventType = !params?.event_type || event.event_type.includes(params.event_type)
        const matchesResult = !params?.result || event.result === params.result
        const matchesUser = !params?.user_id || event.user_id === params.user_id
        return matchesEventType && matchesResult && matchesUser
      })
      return paginateItems([...filtered], params?.page, params?.page_size)
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
    async updateUserProfile(userId, payload) {
      const user = mockUsers.find((item) => item.user_id === userId)
      if (!user) {
        throw new Error('User not found.')
      }
      if (payload.username !== undefined) {
        const username = payload.username.trim()
        if (!/^[a-zA-Z0-9_]{3,32}$/.test(username)) {
          throw new Error('Username must be 3-32 characters and contain only letters, digits, or underscores.')
        }
        if (mockUsers.some((item) => item.user_id !== userId && item.username === username)) {
          throw new Error('Username already exists.')
        }
        user.username = username
      }
      if (payload.email !== undefined) {
        const email = payload.email.trim().toLowerCase()
        if (email.length > 0 && !email.includes('@')) {
          throw new Error('Email is invalid.')
        }
        if (email.length > 0 && mockUsers.some((item) => item.user_id !== userId && item.email.toLowerCase() === email)) {
          throw new Error('Email already exists.')
        }
        user.email = email
      }
      if (payload.nickname !== undefined) {
        user.nickname = payload.nickname.trim()
      }
      if (payload.avatar_url !== undefined) {
        user.avatar_url = payload.avatar_url.trim()
      }
      user.updated_at = Math.floor(Date.now() / 1000)
      if (user.user_id === mockProfile.user_id) {
        mockProfile = {
          ...mockProfile,
          username: user.username,
          email: user.email,
          nickname: user.nickname,
          avatar_url: user.avatar_url ?? '',
        }
      }
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
    async deleteUser(userId) {
      const user = mockUsers.find((item) => item.user_id === userId)
      if (!user) {
        throw new Error('User not found.')
      }
      user.status = 'deleted'
      user.deleted_at = Math.floor(Date.now() / 1000)
      user.updated_at = user.deleted_at
      return { ok: true }
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
    async listContents(params) {
      const keyword = params?.keyword?.toLowerCase().trim()
      const filtered = mockContents.filter((item) => {
        const matchesKeyword = !keyword || item.title.toLowerCase().includes(keyword) || item.slug.toLowerCase().includes(keyword)
        const matchesType = !params?.type || item.type === params.type
        const matchesStatus = !params?.status || item.status === params.status
        const matchesVisibility = !params?.visibility || item.visibility === params.visibility
        return matchesKeyword && matchesType && matchesStatus && matchesVisibility
      })
      return paginateItems(filtered.map(toSummary), params?.page, params?.page_size)
    },
    async createContent(payload) {
      const now = Math.floor(Date.now() / 1000)
      const content = fromContentWrite(`mock_${now}`, payload, now)
      mockContents.unshift(content)
      return { content: { ...content } }
    },
    async getContent(contentId) {
      const content = findContent(contentId)
      return { content: { ...content, tags: [...content.tags] } }
    },
    async updateContent(contentId, payload) {
      const content = findContent(contentId)
      const next = fromContentWrite(contentId, payload, Math.floor(Date.now() / 1000), content)
      Object.assign(content, next)
      return { content: { ...content, tags: [...content.tags] } }
    },
    async archiveContent(contentId) {
      const content = findContent(contentId)
      const now = Math.floor(Date.now() / 1000)
      content.status = 'archived'
      content.archived_at = now
      content.updated_at = now
      return { ok: true }
    },
    async listTags(params) {
      const keyword = params?.keyword?.toLowerCase().trim()
      const filtered = mockTags.filter((tag) => !keyword || tag.name.toLowerCase().includes(keyword) || tag.slug.toLowerCase().includes(keyword))
      return paginateItems([...filtered], params?.page, params?.page_size)
    },
    async createTag(payload) {
      const now = Math.floor(Date.now() / 1000)
      const tag = { tag_id: `tag_${now}`, name: payload.name, slug: payload.slug, description: payload.description, color: payload.color, created_at: now, updated_at: now }
      mockTags.push(tag)
      return { tag: { ...tag } }
    },
    async updateTag(tagId, payload) {
      const tag = findTag(tagId)
      Object.assign(tag, { ...payload, updated_at: Math.floor(Date.now() / 1000) })
      return { tag: { ...tag } }
    },
    async deleteTag(tagId) {
      const index = mockTags.findIndex((tag) => tag.tag_id === tagId)
      if (index < 0) {
        throw new Error('Tag not found.')
      }
      mockTags.splice(index, 1)
      return { ok: true }
    },
    async listRelations(contentId, params) {
      const filtered = mockRelations.filter((relation) => {
        const matchesContent = relation.from_content_id === contentId
        const matchesType = !params?.relation_type || relation.relation_type === params.relation_type
        return matchesContent && matchesType
      })
      return { items: [...filtered], total: filtered.length, page: params?.page ?? 1, page_size: params?.page_size ?? 20 }
    },
    async createRelation(contentId, payload) {
      const now = Math.floor(Date.now() / 1000)
      const relation = {
        relation_id: `rel_${now}`,
        from_content_id: contentId,
        to_content_id: payload.to_content_id,
        relation_type: payload.relation_type,
        weight: payload.weight ?? 0,
        sort_order: payload.sort_order ?? 0,
        metadata_json: payload.metadata_json,
        created_at: now,
        updated_at: now,
      }
      mockRelations.push(relation)
      return { relation: { ...relation } }
    },
    async deleteRelation(_contentId, relationId) {
      const index = mockRelations.findIndex((relation) => relation.relation_id === relationId)
      if (index < 0) {
        throw new Error('Relation not found.')
      }
      mockRelations.splice(index, 1)
      return { ok: true }
    },
    async listRevisions(contentId, params) {
      const filtered = mockRevisions.filter((revision) => revision.content_id === contentId)
      return { items: [...filtered], total: filtered.length, page: params?.page ?? 1, page_size: params?.page_size ?? 20 }
    },
    async getRevision(contentId, revisionId) {
      const revision = mockRevisions.find((item) => item.content_id === contentId && item.revision_id === revisionId)
      if (!revision) {
        throw new Error('Revision not found.')
      }
      const content = findContent(contentId)
      return {
        revision: {
          ...revision,
          title_snapshot: content.title,
          summary_snapshot: content.summary,
          body_markdown: content.body_markdown,
          body_json: content.body_json,
          editor_user_id: 'user_mock_admin',
        },
      }
    },
  }
}

function toSummary(content: ContentDetail) {
  const { body_markdown: _bodyMarkdown, body_json: _bodyJson, owner_user_id: _ownerUserId, author_user_id: _authorUserId, source_type: _sourceType, current_revision_id: _currentRevisionId, comment_enabled: _commentEnabled, is_featured: _isFeatured, sort_order: _sortOrder, ...summary } = content
  return summary
}

function findContent(contentId: string): ContentDetail {
  const content = mockContents.find((item) => item.content_id === contentId)
  if (!content) {
    throw new Error('Content not found.')
  }
  return content
}

function findTag(tagId: string): ContentTag {
  const tag = mockTags.find((item) => item.tag_id === tagId)
  if (!tag) {
    throw new Error('Tag not found.')
  }
  return tag
}

function fromContentWrite(contentId: string, payload: ContentWriteRequest, now: number, existing?: ContentDetail): ContentDetail {
  const tags = mockTags.filter((tag) => payload.tag_ids?.includes(tag.tag_id))
  return {
    content_id: contentId,
    type: payload.type,
    title: payload.title,
    slug: payload.slug,
    summary: payload.summary,
    body_markdown: payload.body_markdown,
    body_json: payload.body_json,
    cover_image_url: payload.cover_image_url,
    status: payload.status ?? existing?.status ?? 'draft',
    visibility: payload.visibility,
    ai_access: payload.ai_access,
    owner_user_id: existing?.owner_user_id ?? 'user_mock_admin',
    author_user_id: existing?.author_user_id ?? 'user_mock_admin',
    source_type: payload.source_type ?? 'manual',
    current_revision_id: existing?.current_revision_id ?? `rev_${contentId}`,
    comment_enabled: payload.comment_enabled ?? true,
    is_featured: payload.is_featured ?? false,
    sort_order: payload.sort_order ?? 0,
    created_at: existing?.created_at ?? now,
    updated_at: now,
    tags,
  }
}
