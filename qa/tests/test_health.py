"""
Health endpoint regression tests.
健康检查接口的回归测试。
"""

from __future__ import annotations

from qa.clients import GatewayClient


def test_healthz_returns_ok(gateway_client: GatewayClient) -> None:
    result = gateway_client.healthz()

    assert result.ok
    assert result.response.status_code == 200
    assert result.response.headers["content-type"].startswith("application/json")
    assert result.data is not None
    assert result.data.status == "ok"


def test_readyz_returns_ready(gateway_client: GatewayClient) -> None:
    result = gateway_client.readyz()

    assert result.ok
    assert result.response.status_code == 200
    assert result.response.headers["content-type"].startswith("application/json")
    assert result.data is not None
    assert result.data.status == "ready"

