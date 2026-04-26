"""
Pydantic response models for gateway HTTP regression checks.
用于 gateway HTTP 回归校验的 Pydantic 响应模型。
"""

from __future__ import annotations

from pydantic import BaseModel, ConfigDict, Field


class StrictModel(BaseModel):
    """
    Base response model that rejects unexpected keys by default.
    默认拒绝未知字段的基础响应模型。
    """

    model_config = ConfigDict(extra="forbid")


class ErrorResponse(StrictModel):
    code: int
    message: str
    reference: str
    request_id: str


class AuthUserProfile(StrictModel):
    user_id: str
    username: str
    email: str
    nickname: str | None = None
    avatar_url: str | None = None
    role: str
    status: str


class AuthSessionView(StrictModel):
    session_id: str
    user_id: str
    auth_source: str
    client_type: str = ""
    device_id: str = ""
    device_name: str = ""
    status: str
    last_seen_at: int | None = None
    expires_at: int | None = None


class AuthRegisterResponse(StrictModel):
    access_token: str
    refresh_token: str
    expires_in: int
    token_type: str
    session_id: str
    user: AuthUserProfile
    session: AuthSessionView


class AuthLoginResponse(AuthRegisterResponse):
    pass


class AuthRefreshResponse(StrictModel):
    access_token: str
    refresh_token: str
    expires_in: int
    token_type: str
    session_id: str
    session: AuthSessionView


class AuthLogoutResponse(StrictModel):
    ok: bool


class AuthMeResponse(StrictModel):
    user: AuthUserProfile


class AuthSsoStartResponse(StrictModel):
    provider: str
    auth_url: str
    state: str


class AuthSsoCallbackResponse(AuthRegisterResponse):
    pass


class HealthzResponse(StrictModel):
    status: str


class ReadyzResponse(StrictModel):
    status: str


# Compatibility alias: keep the typo-safe symbol for any historical callers.
# 兼容历史引用的拼写错误符号，避免历史代码受影响。
ReadyzResponse = ReadyzResponse


class ContentTagView(StrictModel):
    """
    Tag object exposed by content and public content APIs.
    内容接口与公开内容接口返回的标签对象。
    """

    tag_id: str
    name: str
    slug: str
    description: str | None = None
    color: str | None = None
    created_at: int
    updated_at: int


class ContentSummaryView(StrictModel):
    """
    Content summary row for list responses.
    内容列表响应中的摘要视图。
    """

    content_id: str
    type: str
    title: str
    slug: str
    summary: str | None = None
    cover_image_url: str | None = None
    status: str
    visibility: str
    ai_access: str
    published_at: int | None = None
    archived_at: int | None = None
    created_at: int
    updated_at: int
    tags: list[ContentTagView] = Field(default_factory=list)


class ContentDetailView(StrictModel):
    """
    Full content detail payload in studio and public endpoints.
    Studio 与公开接口返回的内容详情结构。
    """

    content_id: str
    type: str
    title: str
    slug: str
    summary: str | None = None
    body_markdown: str
    body_json: str | None = None
    cover_image_url: str | None = None
    status: str
    visibility: str
    ai_access: str
    owner_user_id: str
    author_user_id: str
    source_type: str
    current_revision_id: str | None = None
    comment_enabled: bool
    is_featured: bool
    sort_order: int
    published_at: int | None = None
    archived_at: int | None = None
    created_at: int
    updated_at: int
    tags: list[ContentTagView] = Field(default_factory=list)


class ContentDetailResp(StrictModel):
    """
    Standard response wrapper for content detail.
    内容详情标准返回体。
    """

    content: ContentDetailView


class ContentListResp(StrictModel):
    """
    Standard response wrapper for content list responses.
    内容列表标准返回体。
    """

    items: list[ContentSummaryView]
    total: int
    page: int
    page_size: int


class ContentArchiveResp(StrictModel):
    """
    Standard response wrapper for archive operations.
    内容归档操作返回体。
    """

    ok: bool


class ContentRevisionSummaryView(StrictModel):
    """
    Revision summary in revision list responses.
    版本列表中的摘要结构。
    """

    revision_id: str
    content_id: str
    revision_no: int
    editor_type: str
    change_summary: str | None = None
    source_type: str
    created_at: int


class ContentRevisionDetailView(StrictModel):
    """
    Detailed revision payload.
    版本详情结构体。
    """

    revision_id: str
    content_id: str
    revision_no: int
    title_snapshot: str
    summary_snapshot: str | None = None
    body_markdown: str
    body_json: str | None = None
    editor_type: str
    editor_user_id: str | None = None
    editor_agent_client_id: str | None = None
    change_summary: str | None = None
    source_type: str
    created_at: int


class ContentRevisionDetailResp(StrictModel):
    """
    Standard response wrapper for content revision detail.
    内容版本详情标准返回体。
    """

    revision: ContentRevisionDetailView


class ContentRevisionListResp(StrictModel):
    """
    Standard response wrapper for content revision list.
    内容版本列表标准返回体。
    """

    items: list[ContentRevisionSummaryView]
    total: int
    page: int
    page_size: int


class ContentRelationView(StrictModel):
    """
    Content relation item for relation APIs.
    内容关系条目。
    """

    relation_id: str
    from_content_id: str
    to_content_id: str
    relation_type: str
    weight: int
    sort_order: int
    metadata_json: str | None = None
    created_at: int
    updated_at: int


class ContentRelationResp(StrictModel):
    """
    Standard response wrapper for relation operations that return one relation.
    返回单个关系的标准返回体。
    """

    relation: ContentRelationView


class ContentRelationListResp(StrictModel):
    """
    Standard response wrapper for relation list responses.
    返回关系列表的标准返回体。
    """

    items: list[ContentRelationView]
    total: int
    page: int
    page_size: int


class ContentRelationDeleteResp(StrictModel):
    """
    Standard response wrapper for relation delete operations.
    内容关系删除的标准返回体。
    """

    ok: bool


class ContentTagResp(StrictModel):
    """
    Standard response wrapper for tag create/update APIs.
    标签创建/更新操作的标准返回体。
    """

    tag: ContentTagView


class ContentTagListResp(StrictModel):
    """
    Standard response wrapper for tag list APIs.
    标签列表标准返回体。
    """

    items: list[ContentTagView]
    total: int
    page: int
    page_size: int


class ContentTagDeleteResp(StrictModel):
    """
    Standard response wrapper for tag delete APIs.
    标签删除的标准返回体。
    """

    ok: bool


class PublicContentListResp(ContentListResp):
    """
    Public content list response uses the same item schema as studio list.
    公开内容列表与 studio 列表使用同一 item 模式。
    """


class PublicContentGetResp(ContentDetailResp):
    """
    Public content by slug response wrapper.
    按 slug 查询公开内容的返回体。
    """

