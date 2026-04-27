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
  deleted_at?: number
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
  include_deleted?: boolean
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

export type ContentType =
  | 'article'
  | 'note'
  | 'project'
  | 'experience'
  | 'timeline_event'
  | 'insight'
  | 'portfolio'
  | 'page'

export type ContentStatus = 'draft' | 'review' | 'published' | 'archived'
export type ContentVisibility = 'public' | 'member' | 'private'
export type ContentAIAccess = 'allowed' | 'denied'
export type ContentRelationType =
  | 'belongs_to'
  | 'related_to'
  | 'derived_from'
  | 'references'
  | 'part_of'
  | 'depends_on'
  | 'timeline_of'

export interface ContentTag {
  tag_id: string
  name: string
  slug: string
  description?: string
  color?: string
  created_at: number
  updated_at: number
}

export interface ContentSummary {
  content_id: string
  type: ContentType
  title: string
  slug: string
  summary?: string
  cover_image_url?: string
  status: ContentStatus
  visibility: ContentVisibility
  ai_access: ContentAIAccess
  published_at?: number
  archived_at?: number
  created_at: number
  updated_at: number
  tags: ContentTag[]
}

export interface ContentDetail extends ContentSummary {
  body_markdown: string
  body_json?: string
  owner_user_id: string
  author_user_id: string
  source_type: string
  current_revision_id?: string
  comment_enabled: boolean
  is_featured: boolean
  sort_order: number
}

export interface ContentListParams {
  page?: number
  page_size?: number
  type?: string
  status?: string
  visibility?: string
  keyword?: string
}

export interface ContentListResponse {
  items: ContentSummary[]
  total: number
  page: number
  page_size: number
}

export interface ContentDetailResponse {
  content: ContentDetail
}

export interface ContentMutationResponse {
  ok: boolean
}

export interface ContentWriteRequest {
  type: ContentType
  title: string
  slug: string
  summary?: string
  body_markdown: string
  body_json?: string
  cover_image_url?: string
  status?: ContentStatus
  visibility: ContentVisibility
  ai_access: ContentAIAccess
  source_type?: string
  comment_enabled?: boolean
  is_featured?: boolean
  sort_order?: number
  tag_ids?: string[]
  change_summary?: string
}

export interface ContentTagListParams {
  page?: number
  page_size?: number
  keyword?: string
}

export interface ContentTagListResponse {
  items: ContentTag[]
  total: number
  page: number
  page_size: number
}

export interface ContentTagResponse {
  tag: ContentTag
}

export interface ContentTagWriteRequest {
  name: string
  slug: string
  description?: string
  color?: string
}

export interface ContentRelation {
  relation_id: string
  from_content_id: string
  to_content_id: string
  relation_type: ContentRelationType
  weight: number
  sort_order: number
  metadata_json?: string
  created_at: number
  updated_at: number
}

export interface ContentRelationListParams {
  page?: number
  page_size?: number
  relation_type?: string
}

export interface ContentRelationListResponse {
  items: ContentRelation[]
  total: number
  page: number
  page_size: number
}

export interface ContentRelationResponse {
  relation: ContentRelation
}

export interface ContentRelationWriteRequest {
  to_content_id: string
  relation_type: ContentRelationType
  weight?: number
  sort_order?: number
  metadata_json?: string
}

export interface ContentRevisionSummary {
  revision_id: string
  content_id: string
  revision_no: number
  editor_type: string
  change_summary?: string
  source_type: string
  created_at: number
}

export interface ContentRevisionDetail extends ContentRevisionSummary {
  title_snapshot: string
  summary_snapshot?: string
  body_markdown: string
  body_json?: string
  editor_user_id?: string
  editor_agent_client_id?: string
}

export interface ContentRevisionListResponse {
  items: ContentRevisionSummary[]
  total: number
  page: number
  page_size: number
}

export interface ContentRevisionResponse {
  revision: ContentRevisionDetail
}
