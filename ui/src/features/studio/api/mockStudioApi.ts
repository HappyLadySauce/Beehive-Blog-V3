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

const mockTags: ContentTag[] = [
  { tag_id: 'tag_1001', name: 'Gateway', slug: 'gateway', description: 'Gateway notes', color: '#0f8f83', created_at: 1776676320, updated_at: 1776676320 },
  { tag_id: 'tag_1002', name: 'Identity', slug: 'identity', description: 'Identity flows', color: '#e45f35', created_at: 1776676320, updated_at: 1776676320 },
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
      return { items: filtered.map(toSummary), total: filtered.length, page: params?.page ?? 1, page_size: params?.page_size ?? 20 }
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
      return { items: [...filtered], total: filtered.length, page: params?.page ?? 1, page_size: params?.page_size ?? 50 }
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
