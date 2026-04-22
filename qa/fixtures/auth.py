"""
Auth fixture builders for repository-managed QA tests.
仓库内置 QA 测试使用的认证样例构造器。
"""

from __future__ import annotations

from dataclasses import dataclass
from uuid import uuid4

from qa.config import QASettings


@dataclass(slots=True)
class AuthFixtureUser:
    """
    Generated credential bundle used by auth regression flows.
    认证回归流程使用的生成型凭据集合。
    """

    username: str
    email: str
    password: str
    nickname: str

    def as_register_payload(self) -> dict[str, str]:
        """
        Convert the fixture to a register payload.
        将当前样例转换为注册请求载荷。
        """

        return {
            "username": self.username,
            "email": self.email,
            "password": self.password,
            "nickname": self.nickname,
        }

    def as_login_payload(self) -> dict[str, str]:
        """
        Convert the fixture to a login payload.
        将当前样例转换为登录请求载荷。
        """

        return {
            "login_identifier": self.email,
            "password": self.password,
            "client_type": "web",
            "device_id": "qa-web-device",
            "device_name": "QA Chrome",
            "user_agent": "beehive-qa/1.0",
        }


def build_unique_user(settings: QASettings) -> AuthFixtureUser:
    """
    Build a unique auth fixture to avoid clashes in persistent test environments.
    构建唯一认证样例，避免在持久化测试环境中发生冲突。
    """

    suffix = uuid4().hex[:10]
    username = f"{settings.test_username_prefix}_{suffix}"
    return AuthFixtureUser(
        username=username,
        email=f"{username}@{settings.test_email_domain}",
        password=settings.default_password,
        nickname=f"QA {suffix[:6]}",
    )

