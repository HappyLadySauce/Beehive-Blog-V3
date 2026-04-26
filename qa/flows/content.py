"""
Content flow helpers for QA scenarios.
内容场景回归辅助封装。
"""

from __future__ import annotations

from qa.clients import GatewayClient


class ContentFlows:
    """
    Encapsulates common studio content workflows used by QA tests.
    封装 QA 场景里复用的内容链路。
    """

    def __init__(self, client: GatewayClient) -> None:
        self.client = client

    def list_items(
        self,
        access_token: str,
        *,
        page: int = 1,
        page_size: int = 20,
        type: str | None = None,
        status: str | None = None,
        visibility: str | None = None,
        keyword: str | None = None,
        token_type: str = "Bearer",
    ):
        """
        List studio content items.
        列举 studio 内容条目。
        """

        return self.client.list_studio_content_items(
            access_token,
            page=page,
            page_size=page_size,
            type=type,
            status=status,
            visibility=visibility,
            keyword=keyword,
            token_type=token_type,
        )

    def create_item(
        self,
        access_token: str,
        payload: dict[str, object],
        token_type: str = "Bearer",
    ):
        """
        Create one studio content item.
        创建一条 studio 内容。
        """

        return self.client.create_studio_content(access_token, payload, token_type=token_type)

    def get_item(
        self,
        access_token: str,
        content_id: str,
        token_type: str = "Bearer",
    ):
        """
        Get one studio content item.
        按 id 获取 studio 内容。
        """

        return self.client.get_studio_content(access_token, content_id, token_type=token_type)

    def update_item(
        self,
        access_token: str,
        content_id: str,
        payload: dict[str, object],
        token_type: str = "Bearer",
    ):
        """
        Update one studio content item.
        更新 studio 内容。
        """

        return self.client.update_studio_content(access_token, content_id, payload, token_type=token_type)

    def archive_item(
        self,
        access_token: str,
        content_id: str,
        token_type: str = "Bearer",
    ):
        """
        Archive one studio content item.
        归档 studio 内容。
        """

        return self.client.archive_studio_content(access_token, content_id, token_type=token_type)

    def list_revisions(
        self,
        access_token: str,
        content_id: str,
        *,
        page: int = 1,
        page_size: int = 20,
        token_type: str = "Bearer",
    ):
        """
        List revisions for a studio content item.
        列出内容版本。
        """

        return self.client.list_content_revisions(
            access_token,
            content_id,
            page=page,
            page_size=page_size,
            token_type=token_type,
        )

    def get_revision(
        self,
        access_token: str,
        content_id: str,
        revision_id: str,
        token_type: str = "Bearer",
    ):
        """
        Get one revision by revision id.
        按 revision_id 读取版本详情。
        """

        return self.client.get_content_revision(access_token, content_id, revision_id, token_type=token_type)

    def list_relations(
        self,
        access_token: str,
        content_id: str,
        *,
        page: int = 1,
        page_size: int = 20,
        relation_type: str | None = None,
        token_type: str = "Bearer",
    ):
        """
        List relations for a studio content item.
        列举内容关系。
        """

        return self.client.list_content_relations(
            access_token,
            content_id,
            page=page,
            page_size=page_size,
            relation_type=relation_type,
            token_type=token_type,
        )

    def create_relation(
        self,
        access_token: str,
        content_id: str,
        payload: dict[str, object],
        token_type: str = "Bearer",
    ):
        """
        Create one relation for a content item.
        新建内容关系。
        """

        return self.client.create_content_relation(access_token, content_id, payload, token_type=token_type)

    def delete_relation(
        self,
        access_token: str,
        content_id: str,
        relation_id: str,
        token_type: str = "Bearer",
    ):
        """
        Delete one relation by relation id.
        删除内容关系。
        """

        return self.client.delete_content_relation(access_token, content_id, relation_id, token_type=token_type)

    def list_tags(
        self,
        access_token: str,
        *,
        page: int = 1,
        page_size: int = 20,
        keyword: str | None = None,
        token_type: str = "Bearer",
    ):
        """
        List studio tags.
        列举 studio 标签。
        """

        return self.client.list_content_tags(
            access_token,
            page=page,
            page_size=page_size,
            keyword=keyword,
            token_type=token_type,
        )

    def create_tag(
        self,
        access_token: str,
        payload: dict[str, object],
        token_type: str = "Bearer",
    ):
        """
        Create one studio tag.
        创建 studio 标签。
        """

        return self.client.create_content_tag(access_token, payload, token_type=token_type)

    def update_tag(
        self,
        access_token: str,
        tag_id: str,
        payload: dict[str, object],
        token_type: str = "Bearer",
    ):
        """
        Update one studio tag.
        更新 studio 标签。
        """

        return self.client.update_content_tag(access_token, tag_id, payload, token_type=token_type)

    def delete_tag(
        self,
        access_token: str,
        tag_id: str,
        token_type: str = "Bearer",
    ):
        """
        Delete one studio tag.
        删除 studio 标签。
        """

        return self.client.delete_content_tag(access_token, tag_id, token_type=token_type)

    def list_public_items(
        self,
        *,
        page: int = 1,
        page_size: int = 20,
        type: str | None = None,
        keyword: str | None = None,
    ):
        """
        List public content items without auth.
        无鉴权查询公开内容。
        """

        return self.client.list_public_content_items(page=page, page_size=page_size, type=type, keyword=keyword)

    def get_public_item(self, slug: str):
        """
        Get public content by slug.
        按 slug 查询公开内容。
        """

        return self.client.get_public_content_by_slug(slug)
