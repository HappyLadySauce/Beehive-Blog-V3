"""
SSO regression placeholders controlled by explicit provider flags.
由显式 provider 开关控制的 SSO 回归预留测试。
"""

from __future__ import annotations

import pytest

from qa.config import QASettings
from qa.flows import SSOFlows


@pytest.mark.skipif(True, reason="SSO live-provider regression is disabled by default in phase one.")
def test_sso_start_placeholder() -> None:
    """
    Placeholder to keep the SSO test package visible until provider stubs are added.
    占位测试，用于在补充 provider stub 前保留 SSO 测试入口。
    """


def test_sso_start_is_skipped_when_disabled(qa_settings: QASettings, sso_flows: SSOFlows) -> None:
    if not qa_settings.enable_sso_tests:
        pytest.skip("SSO regression is disabled in the current QA environment.")

    result = sso_flows.start(
        provider="github",
        redirect_uri="https://app.example.com/auth/callback/github",
        state="qa-sso-state",
    )

    assert result.ok
    assert result.data is not None
    assert result.data.provider == "github"

