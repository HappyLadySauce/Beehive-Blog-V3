export interface GatewayErrorResponse {
  code: number;
  message: string;
  reference: string;
  request_id: string;
}

export interface AuthUserProfile {
  user_id: string;
  username: string;
  email: string;
  nickname: string;
  avatar_url: string;
  role: string;
  status: string;
}

export interface AuthSessionView {
  session_id: string;
  user_id: string;
  auth_source: string;
  client_type: string;
  device_id: string;
  device_name: string;
  status: string;
  last_seen_at: number;
  expires_at: number;
}

export interface AuthRegisterRequest {
  username: string;
  email: string;
  password: string;
  nickname?: string;
}

export interface AuthRegisterResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
  session_id: string;
  user: AuthUserProfile;
  session: AuthSessionView;
}

export interface AuthLoginRequest {
  login_identifier: string;
  password: string;
  client_type?: string;
  device_id?: string;
  device_name?: string;
  user_agent?: string;
}

export type AuthLoginResponse = AuthRegisterResponse;

export interface AuthRefreshRequest {
  refresh_token: string;
}

export interface AuthRefreshResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
  session_id: string;
  session: AuthSessionView;
}

export interface AuthMeResponse {
  user: AuthUserProfile;
}

export interface AuthLogoutRequest {
  refresh_token?: string | undefined;
}

export interface AuthLogoutResponse {
  ok: boolean;
}

export interface ContentTagView {
  tag_id: string;
  name: string;
  slug: string;
  description: string;
  color: string;
  created_at: number;
  updated_at: number;
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
  | string;

export type ContentStatus = 'draft' | 'review' | 'published' | 'archived' | string;
export type ContentVisibility = 'public' | 'member' | 'private' | string;
export type ContentAiAccess = 'allowed' | 'denied' | string;
export type ContentSourceType =
  | 'manual'
  | 'import_v1'
  | 'import_markdown'
  | 'agent_generated'
  | 'agent_assisted'
  | string;

export interface ContentSummaryView {
  content_id: string;
  type: ContentType;
  title: string;
  slug: string;
  summary: string;
  cover_image_url: string;
  status: ContentStatus;
  visibility: ContentVisibility;
  ai_access: ContentAiAccess;
  published_at: number;
  archived_at: number;
  created_at: number;
  updated_at: number;
  tags: ContentTagView[];
}

export interface ContentDetailView extends ContentSummaryView {
  body_markdown: string;
  body_json: string;
  owner_user_id: string;
  author_user_id: string;
  source_type: ContentSourceType;
  current_revision_id: string;
  comment_enabled: boolean;
  is_featured: boolean;
  sort_order: number;
}

export interface ContentListResponse {
  items: ContentSummaryView[];
  total: number;
  page: number;
  page_size: number;
}

export interface ContentPublicBySlugResponse {
  content: ContentDetailView;
}

export interface PublicContentQuery {
  page?: number;
  page_size?: number;
  type?: ContentType;
  keyword?: string;
}

export interface StudioContentListQuery extends PublicContentQuery {
  status?: ContentStatus;
  visibility?: ContentVisibility;
}
