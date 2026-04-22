"""
SSO flow helpers reserved for future QA callback scenarios.
为后续 QA 回调场景预留的 SSO 流程辅助封装。
"""

from __future__ import annotations

from qa.clients import GatewayClient


class SSOFlows:
    """
    Encapsulates SSO entrypoints without requiring external provider execution today.
    封装 SSO 入口能力，但当前不依赖真实第三方执行。
    """

    def __init__(self, client: GatewayClient) -> None:
        self.client = client

    def start(self, provider: str, redirect_uri: str, state: str | None = None):
        """
        Trigger the SSO start endpoint for the requested provider.
        调用指定 provider 的 SSO 启动接口。
        """

        payload: dict[str, str] = {
            "provider": provider,
            "redirect_uri": redirect_uri,
        }
        if state:
            payload["state"] = state
        return self.client.sso_start(payload)

    def callback(self, payload: dict[str, str]):
        """
        Trigger the SSO callback endpoint with a prepared payload.
        使用准备好的载荷调用 SSO 回调接口。
        """

        return self.client.sso_callback(payload)

