"""
Content access control regression tests.
内容鉴权回归测试。
"""

from __future__ import annotations

from uuid import uuid4

from qa.clients import GatewayClient
from qa.clients.models import ErrorResponse
from qa.config import QASettings
from qa.fixtures import AuthFixtureUser
from qa.fixtures.content import build_unique_content, build_unique_content_slug


def assert_error_code(error: ErrorResponse | None, expected: set[int]) -> None:
    """
    Assert that the returned business code is in expected set.
    断言返回业务码在预期范围内。
    """

    assert error is not None
    assert error.code in expected


def test_studio_endpoints_require_authorization(gateway_client: GatewayClient, qa_settings: QASettings) -> None:
    unique_suffix = uuid4().hex[:6]
    sample_content = build_unique_content(qa_settings)
    placeholder_id = build_unique_content_slug(qa_settings, prefix=f"placeholder-{unique_suffix}")

    for method, path, body in [
        ("GET", "/api/v3/studio/content/items", None),
        (
            "POST",
            "/api/v3/studio/content/items",
            sample_content.as_create_payload(),
        ),
        ("GET", f"/api/v3/studio/content/items/{placeholder_id}", None),
        (
            "PUT",
            f"/api/v3/studio/content/items/{placeholder_id}",
            sample_content.as_update_payload(status="draft", visibility="private", ai_access="denied", change_summary="unauthorized smoke"),
        ),
        ("DELETE", f"/api/v3/studio/content/items/{placeholder_id}", None),
        ("GET", f"/api/v3/studio/content/items/{placeholder_id}/revisions", None),
        ("GET", f"/api/v3/studio/content/items/{placeholder_id}/revisions/{uuid4().hex[:10]}", None),
        ("GET", f"/api/v3/studio/content/items/{placeholder_id}/relations", None),
        (
            "POST",
            f"/api/v3/studio/content/items/{placeholder_id}/relations",
            {"to_content_id": placeholder_id, "relation_type": "related_to"},
        ),
        ("DELETE", f"/api/v3/studio/content/items/{placeholder_id}/relations/{uuid4().hex[:10]}", None),
        ("GET", "/api/v3/studio/content/tags", None),
        (
            "POST",
            "/api/v3/studio/content/tags",
            {
                "name": f"qa-tag-{unique_suffix}",
                "slug": f"qa-tag-{unique_suffix}",
                "description": "QA placeholder tag",
                "color": "#00ADD8",
            },
        ),
        ("PUT", f"/api/v3/studio/content/tags/{uuid4().hex[:10]}", {"name": "qa", "slug": f"qa-{uuid4().hex[:6]}"}),
        ("DELETE", f"/api/v3/studio/content/tags/{uuid4().hex[:10]}", None),
    ]:
        if method == "GET":
            result = gateway_client.send(method, path)
        elif method == "DELETE":
            result = gateway_client.send(method, path)
        else:
            result = gateway_client.send(method, path, json=body or {})

        assert not result.ok
        assert result.response.status_code in {401, 403}
        assert_error_code(result.error, {100201, 100202, 100203, 100204})


def test_member_token_is_forbidden_on_studio_endpoints(
    gateway_client: GatewayClient,
    unique_user: AuthFixtureUser,
    qa_settings: QASettings,
) -> None:
    register = gateway_client.register(unique_user.as_register_payload())
    assert register.ok

    login = gateway_client.login(unique_user.as_login_payload())
    assert login.ok
    assert login.data is not None
    token = login.data.access_token

    forbidden_codes = {120301}

    list_result = gateway_client.list_studio_content_items(token)
    assert not list_result.ok
    assert_error_code(list_result.error, forbidden_codes)

    create_payload = build_unique_content(qa_settings).as_create_payload()
    create_result = gateway_client.create_studio_content(token, create_payload)
    assert not create_result.ok
    assert_error_code(create_result.error, forbidden_codes)
