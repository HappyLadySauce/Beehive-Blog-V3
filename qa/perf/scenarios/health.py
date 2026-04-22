"""
Health-check Locust users for smoke load tests.
用于轻量压测的健康检查 Locust 用户场景。
"""

from __future__ import annotations

from locust import HttpUser, between, task

from qa.perf.config import load_perf_settings


class HealthzUser(HttpUser):
    """
    Minimal probe user that repeatedly checks `/healthz`.
    持续探测 `/healthz` 的最小用户场景。
    """

    wait_time = between(1, 2)
    host = load_perf_settings().normalized_base_url

    @task
    def healthz(self) -> None:
        self.client.get("/healthz", name="GET /healthz")

