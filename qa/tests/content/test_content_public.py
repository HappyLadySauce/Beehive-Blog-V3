"""
Content public content smoke checks.
公开内容链路冒烟校验。
"""

from __future__ import annotations

from uuid import uuid4

from qa.clients import GatewayClient


def test_public_content_list_is_accessible(gateway_client: GatewayClient) -> None:
    result = gateway_client.list_public_content_items(page=1, page_size=20)

    assert result.ok
    assert result.response.status_code == 200
    assert result.data is not None
    assert isinstance(result.data.items, list)


def test_public_content_get_missing_slug_returns_not_found(gateway_client: GatewayClient) -> None:
    slug = f"qa-public-missing-content-{uuid4().hex[:10]}"
    result = gateway_client.get_public_content_by_slug(slug)

    assert not result.ok
    assert result.response.status_code in {404, 400}
    assert result.error is not None
    assert result.error.code == 120501
