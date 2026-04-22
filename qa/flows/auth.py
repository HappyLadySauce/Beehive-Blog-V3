"""
Authentication flow helpers for multi-step QA scenarios.
多步骤 QA 场景使用的认证流程辅助封装。
"""

from __future__ import annotations

from qa.clients import GatewayClient
from qa.clients.models import AuthRefreshResponse
from qa.fixtures import AuthFixtureUser
from qa.flows.context import AuthFlowContext


class AuthFlows:
    """
    Encapsulates common authentication chains used by regression tests.
    封装回归测试中复用的常见认证链路。
    """

    def __init__(self, client: GatewayClient) -> None:
        self.client = client

    def register_and_capture(self, user: AuthFixtureUser) -> AuthFlowContext:
        """
        Register a new user and capture the issued tokens.
        注册新用户并捕获签发的令牌。
        """

        result = self.client.register(user.as_register_payload())
        if not result.ok or result.data is None:
            raise AssertionError(f"register_and_capture failed: status={result.response.status_code} payload={result.payload}")
        return AuthFlowContext.from_auth_response(result.data)

    def login_and_capture(self, user: AuthFixtureUser) -> AuthFlowContext:
        """
        Log in with an existing fixture user and capture the issued tokens.
        使用已有样例用户登录并捕获签发的令牌。
        """

        result = self.client.login(user.as_login_payload())
        if not result.ok or result.data is None:
            raise AssertionError(f"login_and_capture failed: status={result.response.status_code} payload={result.payload}")
        return AuthFlowContext.from_auth_response(result.data)

    def refresh_and_capture(self, context: AuthFlowContext) -> AuthFlowContext:
        """
        Refresh tokens and return the next authenticated context.
        刷新令牌并返回下一轮认证上下文。
        """

        result = self.client.refresh(
            {
                "refresh_token": context.refresh_token,
                "user_agent": "beehive-qa/1.0",
            }
        )
        if not result.ok or result.data is None:
            raise AssertionError(f"refresh_and_capture failed: status={result.response.status_code} payload={result.payload}")

        refreshed: AuthRefreshResponse = result.data
        return AuthFlowContext(
            access_token=refreshed.access_token,
            refresh_token=refreshed.refresh_token,
            token_type=refreshed.token_type,
            session_id=refreshed.session_id,
            user_id=context.user_id,
            username=context.username,
            email=context.email,
        )

    def me(self, context: AuthFlowContext):
        """
        Load the current user for the provided auth context.
        为给定认证上下文读取当前用户。
        """

        return self.client.me(context.access_token, context.token_type)

    def logout(self, context: AuthFlowContext):
        """
        Revoke the current session and refresh chain.
        吊销当前会话及其刷新链。
        """

        return self.client.logout(
            access_token=context.access_token,
            refresh_token=context.refresh_token,
            token_type=context.token_type,
        )

