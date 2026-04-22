"""
Pydantic response models for gateway HTTP regression checks.
用于 gateway HTTP 回归校验的 Pydantic 响应模型。
"""

from __future__ import annotations

from pydantic import BaseModel, ConfigDict


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

