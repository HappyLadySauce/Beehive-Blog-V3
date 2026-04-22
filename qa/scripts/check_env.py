"""
Environment diagnostics for repository-managed QA execution.
仓库内置 QA 执行环境的探测脚本。
"""

from __future__ import annotations

import sys

from qa.clients import GatewayClient
from qa.clients.exceptions import GatewayClientError
from qa.config import load_settings


def main() -> int:
    """
    Check whether the gateway target environment is ready for QA execution.
    检查 gateway 目标环境是否已准备好执行 QA 测试。
    """

    settings = load_settings()
    print(f"[qa] target base url: {settings.normalized_base_url}")

    try:
        with GatewayClient(settings) as client:
            health = client.healthz()
            print(f"[qa] healthz -> {health.response.status_code} {health.payload}")

            ready = client.readyz()
            print(f"[qa] readyz -> {ready.response.status_code} {ready.payload}")

            if not health.ok:
                print("[qa] gateway healthz failed")
                return 1
            if not ready.ok:
                print("[qa] gateway readyz failed")
                return 1
    except GatewayClientError as exc:
        print(f"[qa] environment check failed: {exc}")
        return 1

    print("[qa] environment is ready")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

