"""
Studio content CRUD and visibility tests.
Studio 内容管理链路的 CRUD 与可见性验证。
"""

from __future__ import annotations

import time

from qa.clients import GatewayClient
from qa.config import QASettings
from qa.flows import AuthFlowContext
from qa.fixtures.content import build_unique_content, build_unique_tag


def test_studio_content_happy_path(
    qa_settings: QASettings,
    gateway_client: GatewayClient,
    admin_auth_context: AuthFlowContext,
    content_flows,
) -> None:
    admin_token = admin_auth_context.access_token
    created_content_ids: list[str] = []
    created_tag_ids: list[str] = []
    created_relation_id: str | None = None
    primary_id: str | None = None

    primary_content = build_unique_content(qa_settings)
    relation_content = build_unique_content(qa_settings)

    try:
        tag = build_unique_tag(qa_settings)
        tag_result = content_flows.create_tag(admin_token, tag.as_create_payload())
        assert tag_result.ok
        assert tag_result.data is not None

        tag_id = tag_result.data.tag.tag_id
        created_tag_ids.append(tag_id)

        tag_list = content_flows.list_tags(admin_token, keyword=tag.name)
        assert tag_list.ok
        assert tag_list.data is not None
        assert any(item.tag_id == tag_id for item in tag_list.data.items)

        primary_content.tag_ids = [tag_id]
        primary_result = content_flows.create_item(admin_token, primary_content.as_create_payload())
        assert primary_result.ok
        assert primary_result.data is not None

        primary_id = primary_result.data.content.content_id
        created_content_ids.append(primary_id)

        primary_get = content_flows.get_item(admin_token, primary_id)
        assert primary_get.ok
        assert primary_get.data is not None
        assert primary_get.data.content.content_id == primary_id

        primary_list = content_flows.list_items(admin_token, keyword=primary_content.slug)
        assert primary_list.ok
        assert primary_list.data is not None
        assert any(item.content_id == primary_id for item in primary_list.data.items)

        updated = content_flows.update_item(
            admin_token,
            primary_id,
            primary_content.as_update_payload(
                status="published",
                visibility="public",
                ai_access="allowed",
                change_summary="publish content in qa",
            ),
        )
        assert updated.ok
        assert updated.data is not None
        assert updated.data.content.status == "published"
        assert updated.data.content.visibility == "public"

        revisions = content_flows.list_revisions(admin_token, primary_id, page=1, page_size=20)
        assert revisions.ok
        assert revisions.data is not None
        assert revisions.data.total >= 1
        assert revisions.data.items

        revision_id = revisions.data.items[0].revision_id
        revision_detail = content_flows.get_revision(admin_token, primary_id, revision_id)
        assert revision_detail.ok
        assert revision_detail.data is not None
        assert revision_detail.data.revision.revision_id == revision_id

        relation_content.tag_ids = [tag_id]
        relation_content_result = content_flows.create_item(admin_token, relation_content.as_create_payload())
        assert relation_content_result.ok
        assert relation_content_result.data is not None

        relation_content_id = relation_content_result.data.content.content_id
        created_content_ids.append(relation_content_id)

        relation_result = content_flows.create_relation(
            admin_token,
            primary_id,
            {
                "to_content_id": relation_content_id,
                "relation_type": "related_to",
                "weight": 10,
                "sort_order": 0,
            },
        )
        assert relation_result.ok
        assert relation_result.data is not None

        relation_id = relation_result.data.relation.relation_id
        created_relation_id = relation_id

        relation_list = content_flows.list_relations(admin_token, primary_id, relation_type="related_to")
        assert relation_list.ok
        assert relation_list.data is not None
        assert any(item.relation_id == relation_id for item in relation_list.data.items)

        removed = content_flows.delete_relation(admin_token, primary_id, relation_id)
        assert removed.ok
        created_relation_id = None

        archive_primary = content_flows.archive_item(admin_token, primary_id)
        assert archive_primary.ok
        assert archive_primary.data is not None
        assert archive_primary.data.ok
        created_content_ids.remove(primary_id)

        archive_related = content_flows.archive_item(admin_token, relation_content_id)
        assert archive_related.ok
        created_content_ids.remove(relation_content_id)

    finally:
        if created_relation_id and primary_id:
            content_flows.delete_relation(admin_token, primary_id, created_relation_id)

        for content_id in created_content_ids:
            content_flows.archive_item(admin_token, content_id)

        for tag_id in created_tag_ids:
            delete_tag = content_flows.delete_tag(admin_token, tag_id)
            if not delete_tag.ok and delete_tag.error is not None:
                assert delete_tag.error.code in {120503, 120506}


def test_studio_content_visibility_check(
    qa_settings: QASettings,
    gateway_client: GatewayClient,
    admin_auth_context: AuthFlowContext,
    content_flows,
) -> None:
    admin_token = admin_auth_context.access_token
    content_id: str | None = None
    tag_id: str | None = None
    content_archived = False

    tag = build_unique_tag(qa_settings)
    tag_result = content_flows.create_tag(admin_token, tag.as_create_payload())
    assert tag_result.ok
    assert tag_result.data is not None
    tag_id = tag_result.data.tag.tag_id

    content = build_unique_content(qa_settings, tag_ids=[tag_result.data.tag.tag_id])
    try:
        create_result = content_flows.create_item(admin_token, content.as_create_payload())
        assert create_result.ok
        assert create_result.data is not None

        content_id = create_result.data.content.content_id
        slug = create_result.data.content.slug

        update_payload = content.as_update_payload(
            status="published",
            visibility="public",
            ai_access="allowed",
            change_summary="make content public",
        )
        update_result = content_flows.update_item(admin_token, content_id, update_payload)
        assert update_result.ok

        public_before_archive = gateway_client.get_public_content_by_slug(slug)
        assert public_before_archive.ok
        assert public_before_archive.data is not None
        assert public_before_archive.data.content.content_id == content_id

        archived = content_flows.archive_item(admin_token, content_id)
        assert archived.ok
        content_archived = True

        # Give public list/read path a moment to receive transition if async cache is involved.
        deadline = time.time() + 3
        after_archive = gateway_client.get_public_content_by_slug(slug)
        while after_archive.ok and time.time() < deadline:
            time.sleep(0.2)
            after_archive = gateway_client.get_public_content_by_slug(slug)

        assert not after_archive.ok
        assert after_archive.error is not None
        assert after_archive.error.code == 120501
    finally:
        if content_id and not content_archived:
            content_flows.archive_item(admin_token, content_id)
        if tag_id:
            delete_tag = content_flows.delete_tag(admin_token, tag_id)
            if not delete_tag.ok and delete_tag.error is not None:
                assert delete_tag.error.code in {120503, 120506}
