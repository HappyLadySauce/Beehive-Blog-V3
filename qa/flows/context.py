"""
Runtime flow context shared across QA multi-step scenarios.
QA 多步骤场景共享的运行时上下文。
"""

from __future__ import annotations

from dataclasses import dataclass

from qa.clients.models import AuthLoginResponse, AuthRegisterResponse, AuthSsoCallbackResponse


@dataclass(slots=True)
class AuthFlowContext:
    """
    Captures the current authenticated state for chained requests.
    保存链式请求所需的当前认证状态。
    """

    access_token: str
    refresh_token: str
    token_type: str
    session_id: str
    user_id: str
    username: str
    email: str

    @classmethod
    def from_auth_response(
        cls,
        payload: AuthRegisterResponse | AuthLoginResponse | AuthSsoCallbackResponse,
    ) -> "AuthFlowContext":
        """
        Build flow context from a successful auth response.
        从成功的认证响应中构建流程上下文。
        """

        return cls(
            access_token=payload.access_token,
            refresh_token=payload.refresh_token,
            token_type=payload.token_type,
            session_id=payload.session_id,
            user_id=payload.user.user_id,
            username=payload.user.username,
            email=payload.user.email,
        )

