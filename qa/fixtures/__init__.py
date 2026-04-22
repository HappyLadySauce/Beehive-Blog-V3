"""
Reusable test fixtures for the QA project.
QA 测试工程的可复用测试样例入口。
"""

from .auth import AuthFixtureUser, build_unique_user

__all__ = ["AuthFixtureUser", "build_unique_user"]

