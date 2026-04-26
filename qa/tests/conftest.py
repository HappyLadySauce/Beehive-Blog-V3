"""
Pytest fixtures for the QA regression project.
QA 回归测试工程的 Pytest 夹具定义。
"""

from __future__ import annotations

import pytest

from qa.clients import GatewayClient
from qa.clients.exceptions import GatewayClientError
from qa.config import QASettings, load_settings
from qa.fixtures import AuthFixtureUser, build_unique_user
from qa.flows import AuthFlowContext, AuthFlows, ContentFlows, SSOFlows


@pytest.fixture(scope="session")
def qa_settings() -> QASettings:
    """
    Load session-wide QA settings.
    加载会话级 QA 配置。
    """

    return load_settings()


@pytest.fixture(scope="session")
def gateway_client(qa_settings: QASettings) -> GatewayClient:
    """
    Provide a shared gateway client for regression tests.
    为回归测试提供共享的 gateway 客户端。
    """

    client = GatewayClient(qa_settings)
    yield client
    client.close()


@pytest.fixture(scope="session", autouse=True)
def ensure_gateway_is_ready(gateway_client: GatewayClient) -> None:
    """
    Fail fast when the target gateway is not reachable or not ready.
    当目标 gateway 不可达或未就绪时快速失败。
    """

    try:
        health = gateway_client.healthz()
        ready = gateway_client.readyz()
    except GatewayClientError as exc:
        pytest.exit(
            f"Gateway environment is not reachable: {exc}. Run `uv run --project qa python -m qa.scripts.check_env` for diagnostics.",
            returncode=2,
        )

    if not health.ok or not ready.ok:
        pytest.exit(
            "Gateway environment is not ready. Run `uv run --project qa python -m qa.scripts.check_env` for diagnostics.",
            returncode=2,
        )


@pytest.fixture
def auth_flows(gateway_client: GatewayClient) -> AuthFlows:
    """
    Provide authentication flow helpers.
    提供认证流程辅助封装。
    """

    return AuthFlows(gateway_client)


@pytest.fixture
def sso_flows(gateway_client: GatewayClient) -> SSOFlows:
    """
    Provide SSO flow helpers.
    提供 SSO 流程辅助封装。
    """

    return SSOFlows(gateway_client)


@pytest.fixture
def content_flows(gateway_client: GatewayClient) -> ContentFlows:
    """
    Provide content flow helpers.
    提供内容链路辅助封装。
    """

    return ContentFlows(gateway_client)


@pytest.fixture
def admin_auth_context(qa_settings: QASettings, gateway_client: GatewayClient) -> AuthFlowContext:
    """
    Build an admin auth context or skip studio tests when unavailable.
    构建 admin 鉴权上下文，不可用则跳过 studio 测试。
    """

    if not qa_settings.enable_content_studio_tests:
        pytest.skip("content studio tests are disabled in this QA environment.")

    if not qa_settings.admin_login_identifier or not qa_settings.admin_password:
        pytest.skip(
            "content studio tests require BEEHIVE_QA_ADMIN_LOGIN_IDENTIFIER and BEEHIVE_QA_ADMIN_PASSWORD."
        )

    login = gateway_client.login(
        {
            "login_identifier": qa_settings.admin_login_identifier,
            "password": qa_settings.admin_password,
            "client_type": "web",
            "device_id": "qa-admin-device",
            "device_name": "QA Admin Device",
            "user_agent": "beehive-qa/1.0",
        }
    )

    if not login.ok or login.data is None:
        pytest.skip(f"admin login failed: status={login.response.status_code} payload={login.payload}")

    admin_profile = gateway_client.me(login.data.access_token, login.data.token_type)
    if not admin_profile.ok or admin_profile.data is None:
        pytest.skip(f"admin /api/v3/auth/me failed: status={admin_profile.response.status_code} payload={admin_profile.payload}")

    if admin_profile.data.user.role.lower() != "admin":
        pytest.skip("admin test user is not role=admin in this environment.")

    if qa_settings.admin_login_email_like and qa_settings.admin_login_email_like.lower() not in admin_profile.data.user.email.lower():
        pytest.skip("admin login identity does not match BEEHIVE_QA_ADMIN_LOGIN_EMAIL_LIKE.")

    return AuthFlowContext.from_auth_response(login.data)


@pytest.fixture
def unique_user(qa_settings: QASettings) -> AuthFixtureUser:
    """
    Build a unique user fixture for the current test.
    为当前测试生成唯一用户样例。
    """

    return build_unique_user(qa_settings)
