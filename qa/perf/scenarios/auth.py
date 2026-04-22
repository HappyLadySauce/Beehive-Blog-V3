"""
Authenticated Locust users for future protected-endpoint load tests.
面向受保护接口后续压测的已认证 Locust 用户场景。
"""

from __future__ import annotations

from uuid import uuid4

from locust import HttpUser, between, task

from qa.config import load_settings


class AuthenticatedUser(HttpUser):
    """
    Register-and-login user that exercises `/api/v3/auth/me`.
    先注册再登录、随后访问 `/api/v3/auth/me` 的用户场景。
    """

    wait_time = between(1, 3)
    settings = load_settings()
    host = settings.normalized_base_url

    def on_start(self) -> None:
        suffix = uuid4().hex[:10]
        username = f"{self.settings.test_username_prefix}_{suffix}"
        password = self.settings.default_password

        register_response = self.client.post(
            "/api/v3/auth/register",
            json={
                "username": username,
                "email": f"{username}@{self.settings.test_email_domain}",
                "password": password,
                "nickname": f"Locust {suffix[:6]}",
            },
            name="POST /api/v3/auth/register",
        )
        register_response.raise_for_status()
        payload = register_response.json()
        token_type = payload.get("token_type", "Bearer")
        access_token = payload.get("access_token", "")
        self._auth_headers = {"Authorization": f"{token_type} {access_token}"}

    @task
    def me(self) -> None:
        self.client.get("/api/v3/auth/me", headers=self._auth_headers, name="GET /api/v3/auth/me")

