"""
Shared performance configuration for Locust scenarios.
Locust 场景使用的共享性能配置。
"""

from __future__ import annotations

from qa.config import QASettings, load_settings


def load_perf_settings() -> QASettings:
    """
    Load shared QA settings for performance scenarios.
    为性能场景加载共享的 QA 配置。
    """

    return load_settings()

