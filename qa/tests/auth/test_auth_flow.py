"""
Gateway authentication regression flows.
gateway 认证链路的回归测试。
"""

from __future__ import annotations

from qa.clients import GatewayClient
from qa.clients.models import ErrorResponse
from qa.fixtures import AuthFixtureUser
from qa.flows import AuthFlowContext, AuthFlows


def assert_error_code(error: ErrorResponse | None, expected_codes: set[int]) -> None:
    """
    Assert that the returned error code is one of the expected business codes.
    断言返回的错误码属于预期业务错误码集合。
    """

    assert error is not None
    assert error.code in expected_codes


def test_register_success(gateway_client: GatewayClient, unique_user: AuthFixtureUser) -> None:
    result = gateway_client.register(unique_user.as_register_payload())

    assert result.ok
    assert result.response.status_code == 200
    assert result.data is not None
    assert result.data.token_type == "Bearer"
    assert result.data.access_token
    assert result.data.refresh_token
    assert result.data.user.username == unique_user.username
    assert result.data.user.email == unique_user.email


def test_duplicate_register_is_rejected(gateway_client: GatewayClient, unique_user: AuthFixtureUser) -> None:
    first = gateway_client.register(unique_user.as_register_payload())
    assert first.ok

    second = gateway_client.register(unique_user.as_register_payload())

    assert not second.ok
    assert second.response.status_code in {400, 409}
    assert_error_code(second.error, {110502, 110503})


def test_login_success(gateway_client: GatewayClient, unique_user: AuthFixtureUser) -> None:
    gateway_client.register(unique_user.as_register_payload())

    result = gateway_client.login(unique_user.as_login_payload())

    assert result.ok
    assert result.response.status_code == 200
    assert result.data is not None
    assert result.data.token_type == "Bearer"
    assert result.data.user.username == unique_user.username


def test_login_with_wrong_password_fails(gateway_client: GatewayClient, unique_user: AuthFixtureUser) -> None:
    gateway_client.register(unique_user.as_register_payload())
    payload = unique_user.as_login_payload()
    payload["password"] = "wrong-password"

    result = gateway_client.login(payload)

    assert not result.ok
    assert result.response.status_code == 401
    assert_error_code(result.error, {110201})


def test_me_requires_authorization(gateway_client: GatewayClient) -> None:
    result = gateway_client.send("GET", "/api/v3/auth/me")

    assert not result.ok
    assert result.response.status_code == 401
    assert_error_code(result.error, {100201, 100202, 100203, 100204})


def test_register_me_logout_flow(auth_flows: AuthFlows, unique_user: AuthFixtureUser) -> None:
    context = auth_flows.register_and_capture(unique_user)
    current = auth_flows.me(context)

    assert current.ok
    assert current.data is not None
    assert current.data.user.user_id == context.user_id
    assert current.data.user.email == context.email

    logout = auth_flows.logout(context)
    assert logout.ok
    assert logout.data is not None
    assert logout.data.ok is True

    revoked_me = auth_flows.me(context)
    assert not revoked_me.ok
    assert revoked_me.response.status_code == 401
    assert_error_code(revoked_me.error, {100201, 100202, 100203, 100204})


def test_login_refresh_me_flow(auth_flows: AuthFlows, gateway_client: GatewayClient, unique_user: AuthFixtureUser) -> None:
    gateway_client.register(unique_user.as_register_payload())
    context = auth_flows.login_and_capture(unique_user)
    refreshed = auth_flows.refresh_and_capture(context)
    current = auth_flows.me(refreshed)

    assert current.ok
    assert current.data is not None
    assert current.data.user.username == unique_user.username
    assert refreshed.access_token
    assert refreshed.refresh_token
    assert refreshed.session_id == context.session_id


def test_logout_revokes_refresh_chain(auth_flows: AuthFlows, gateway_client: GatewayClient, unique_user: AuthFixtureUser) -> None:
    gateway_client.register(unique_user.as_register_payload())
    context = auth_flows.login_and_capture(unique_user)
    logout = auth_flows.logout(context)

    assert logout.ok

    refreshed = gateway_client.refresh(
        {
            "refresh_token": context.refresh_token,
            "user_agent": "beehive-qa/1.0",
        }
    )

    assert not refreshed.ok
    assert refreshed.response.status_code == 401
    assert_error_code(refreshed.error, {110202, 110203, 110204})

