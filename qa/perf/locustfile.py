"""
Locust entrypoint that exposes the initial smoke load users.
暴露第一批轻量压测用户场景的 Locust 入口。
"""

from qa.perf.scenarios.auth import AuthenticatedUser
from qa.perf.scenarios.health import HealthzUser

__all__ = ["AuthenticatedUser", "HealthzUser"]

