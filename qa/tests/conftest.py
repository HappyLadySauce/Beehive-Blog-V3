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
from qa.flows import AuthFlows, SSOFlows


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
def unique_user(qa_settings: QASettings) -> AuthFixtureUser:
    """
    Build a unique user fixture for the current test.
    为当前测试生成唯一用户样例。
    """

    return build_unique_user(qa_settings)

