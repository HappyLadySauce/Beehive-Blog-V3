"""
Configuration loading for the repository-managed QA project.
仓库内置 QA 测试工程的配置加载逻辑。
"""

from __future__ import annotations

from functools import lru_cache
from pathlib import Path
import os

from dotenv import load_dotenv
from pydantic import BaseModel, ConfigDict, Field


QA_ROOT = Path(__file__).resolve().parents[1]
ENV_FILE = QA_ROOT / ".env"


class QASettings(BaseModel):
    """
    Stable runtime configuration for QA execution.
    QA 执行阶段使用的稳定运行时配置。
    """

    model_config = ConfigDict(frozen=True, extra="ignore")

    base_url: str = Field(default="http://127.0.0.1:8888")
    timeout_seconds: float = Field(default=5.0)
    verify_ssl: bool = Field(default=False)
    default_password: str = Field(default="Str0ngP@ssw0rd!")
    test_username_prefix: str = Field(default="qa_user")
    test_email_domain: str = Field(default="example.test")
    enable_sso_tests: bool = Field(default=False)
    github_enabled: bool = Field(default=False)
    qq_enabled: bool = Field(default=False)
    wechat_enabled: bool = Field(default=False)
    enable_content_studio_tests: bool = Field(default=False)
    admin_login_identifier: str | None = Field(default=None)
    admin_password: str | None = Field(default=None)
    admin_login_email_like: str | None = Field(default=None)

    @property
    def normalized_base_url(self) -> str:
        """
        Normalize the base URL to avoid duplicated slashes.
        规范化基础地址，避免出现重复斜杠。
        """

        return self.base_url.rstrip("/")


@lru_cache(maxsize=1)
def load_settings() -> QASettings:
    """
    Load QA settings from `.env` and process environment variables.
    从 `.env` 与进程环境变量中加载 QA 配置。
    """

    load_dotenv(ENV_FILE, override=False)

    bool_env = lambda value: str(value or "").strip().lower() in {"1", "true", "yes", "on"}

    return QASettings(
        base_url=os.getenv("BEEHIVE_QA_BASE_URL", "http://127.0.0.1:8888"),
        timeout_seconds=float(os.getenv("BEEHIVE_QA_TIMEOUT_SECONDS", "5")),
        verify_ssl=bool_env(os.getenv("BEEHIVE_QA_VERIFY_SSL", "false")),
        default_password=os.getenv("BEEHIVE_QA_DEFAULT_PASSWORD", "Str0ngP@ssw0rd!"),
        test_username_prefix=os.getenv("BEEHIVE_QA_TEST_USERNAME_PREFIX", "qa_user"),
        test_email_domain=os.getenv("BEEHIVE_QA_TEST_EMAIL_DOMAIN", "example.test"),
        enable_sso_tests=bool_env(os.getenv("BEEHIVE_QA_ENABLE_SSO_TESTS", "false")),
        github_enabled=bool_env(os.getenv("BEEHIVE_QA_GITHUB_ENABLED", "false")),
        qq_enabled=bool_env(os.getenv("BEEHIVE_QA_QQ_ENABLED", "false")),
        wechat_enabled=bool_env(os.getenv("BEEHIVE_QA_WECHAT_ENABLED", "false")),
        enable_content_studio_tests=bool_env(os.getenv("BEEHIVE_QA_ENABLE_CONTENT_STUDIO_TESTS", "false")),
        admin_login_identifier=(os.getenv("BEEHIVE_QA_ADMIN_LOGIN_IDENTIFIER") or None),
        admin_password=(os.getenv("BEEHIVE_QA_ADMIN_PASSWORD") or None),
        admin_login_email_like=(os.getenv("BEEHIVE_QA_ADMIN_LOGIN_EMAIL_LIKE") or None),
    )
