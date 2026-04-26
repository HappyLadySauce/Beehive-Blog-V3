"""
Content fixture builders for repository-managed QA tests.
仓库内置 QA 测试使用的内容样例构造器。
"""

from __future__ import annotations

from dataclasses import dataclass, field
from uuid import uuid4

from qa.config import QASettings


def _build_slug(prefix: str, settings: QASettings, max_len: int = 48) -> str:
    """
    Build a globally unique slug.
    构建全局唯一 slug。
    """

    candidate = f"{settings.test_username_prefix}-{prefix}-{uuid4().hex[:10]}"
    return candidate[:max_len].lower()


@dataclass(slots=True)
class ContentTagFixture:
    """
    Tag fixture used by content tag test cases.
    内容标签样例。
    """

    name: str
    slug: str
    description: str
    color: str

    def as_create_payload(self) -> dict[str, str]:
        """
        Convert the fixture to a create payload.
        构造标签创建载荷。
        """

        return {
            "name": self.name,
            "slug": self.slug,
            "description": self.description,
            "color": self.color,
        }


@dataclass(slots=True)
class ContentFixture:
    """
    Content fixture for studio content flow tests.
    内容创建样例。
    """

    type: str
    title: str
    slug: str
    summary: str
    body_markdown: str
    visibility: str = "private"
    ai_access: str = "denied"
    source_type: str = "manual"
    comment_enabled: bool = True
    is_featured: bool = False
    sort_order: int = 0
    tag_ids: list[str] = field(default_factory=list)
    change_summary: str = "initial draft"

    def as_create_payload(self) -> dict[str, object]:
        """
        Convert the fixture to a create payload.
        构造内容创建载荷。
        """

        return {
            "type": self.type,
            "title": self.title,
            "slug": self.slug,
            "summary": self.summary,
            "body_markdown": self.body_markdown,
            "visibility": self.visibility,
            "ai_access": self.ai_access,
            "source_type": self.source_type,
            "comment_enabled": self.comment_enabled,
            "is_featured": self.is_featured,
            "sort_order": self.sort_order,
            "tag_ids": self.tag_ids,
            "change_summary": self.change_summary,
        }

    def as_update_payload(
        self,
        *,
        status: str,
        visibility: str,
        ai_access: str,
        change_summary: str,
    ) -> dict[str, object]:
        """
        Convert the fixture to an update payload with required transitions.
        构造包含必填更新字段的更新载荷。
        """

        payload = self.as_create_payload()
        payload["status"] = status
        payload["visibility"] = visibility
        payload["ai_access"] = ai_access
        payload["change_summary"] = change_summary
        return payload


def build_unique_content_slug(settings: QASettings, *, prefix: str = "content") -> str:
    """
    Build a unique content slug.
    构建唯一 content slug。
    """

    return _build_slug(prefix=prefix, settings=settings)


def build_unique_tag(settings: QASettings) -> ContentTagFixture:
    """
    Build a unique content tag fixture.
    构建唯一标签样例。
    """

    suffix = uuid4().hex[:6]
    return ContentTagFixture(
        name=f"qa_tag_{suffix}",
        slug=_build_slug(prefix="tag", settings=settings),
        description="QA tag fixture",
        color="#00ADD8",
    )


def build_unique_content(settings: QASettings, *, tag_ids: list[str] | None = None) -> ContentFixture:
    """
    Build a unique content fixture.
    构建唯一内容样例。
    """

    suffix = uuid4().hex[:8]
    return ContentFixture(
        type="article",
        title=f"QA Content {suffix}",
        slug=_build_slug(prefix="content", settings=settings),
        summary="QA content summary",
        body_markdown="# QA content\n\nThis content is created for regression smoke tests.",
        tag_ids=list(tag_ids or []),
    )
