"""
Gateway HTTP client used by repository-managed QA tests.
仓库内置 QA 测试使用的 gateway HTTP 客户端。
"""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Generic, TypeVar

import httpx
from pydantic import BaseModel, ValidationError

from qa.clients.exceptions import GatewayResponseDecodeError, GatewayTransportError
from qa.clients.models import (
    AuthLoginResponse,
    AuthLogoutResponse,
    AuthMeResponse,
    AuthRefreshResponse,
    AuthRegisterResponse,
    AuthSsoCallbackResponse,
    AuthSsoStartResponse,
    ContentArchiveResp,
    ContentDetailResp,
    ContentListResp,
    ContentRelationDeleteResp,
    ContentRelationListResp,
    ContentRelationResp,
    ContentRevisionDetailResp,
    ContentRevisionListResp,
    ContentTagDeleteResp,
    ContentTagListResp,
    ContentTagResp,
    ErrorResponse,
    HealthzResponse,
    ReadyzResponse,
    PublicContentGetResp,
    PublicContentListResp,
)
from qa.config import QASettings, load_settings

ModelT = TypeVar("ModelT", bound=BaseModel)


@dataclass(slots=True)
class EndpointResult(Generic[ModelT]):
    """
    Combined raw HTTP response and structured payload.
    同时包含原始 HTTP 响应与结构化载荷的结果对象。
    """

    response: httpx.Response
    payload: Any
    data: ModelT | None = None
    error: ErrorResponse | None = None

    @property
    def ok(self) -> bool:
        """
        Return whether the response is successful.
        返回响应是否成功。
        """

        return self.response.is_success


class GatewayClient:
    """
    Thin wrapper around gateway HTTP APIs for QA execution.
    QA 执行阶段使用的 gateway HTTP API 薄封装。
    """

    def __init__(self, settings: QASettings | None = None) -> None:
        self.settings = settings or load_settings()
        self._client = httpx.Client(
            base_url=self.settings.normalized_base_url,
            timeout=self.settings.timeout_seconds,
            verify=self.settings.verify_ssl,
            follow_redirects=False,
        )

    def close(self) -> None:
        """
        Close the underlying HTTP client.
        关闭底层 HTTP 客户端。
        """

        self._client.close()

    def __enter__(self) -> "GatewayClient":
        return self

    def __exit__(self, *_: object) -> None:
        self.close()

    def send(
        self,
        method: str,
        path: str,
        *,
        success_model: type[ModelT] | None = None,
        headers: dict[str, str] | None = None,
        params: dict[str, Any] | None = None,
        json: dict[str, Any] | None = None,
    ) -> EndpointResult[ModelT]:
        """
        Send a request and parse it into success or error models.
        发送请求，并将其解析为成功或错误模型。
        """

        try:
            response = self._client.request(method, path, headers=headers, params=params, json=json)
        except httpx.HTTPError as exc:
            raise GatewayTransportError(f"failed to call gateway {method} {path}: {exc}") from exc

        if response.is_success:
            payload = self._load_json(response)
            if success_model is None:
                return EndpointResult(response=response, payload=payload)
            return EndpointResult(
                response=response,
                payload=payload,
                data=self._validate_model(success_model, payload, response),
            )

        payload = self._load_error_payload(response)
        error = self._try_validate_error(payload)
        return EndpointResult(response=response, payload=payload, error=error)

    def healthz(self) -> EndpointResult[HealthzResponse]:
        return self.send("GET", "/healthz", success_model=HealthzResponse)

    def readyz(self) -> EndpointResult[ReadyzResponse]:
        return self.send("GET", "/readyz", success_model=ReadyzResponse)

    def register(self, payload: dict[str, Any]) -> EndpointResult[AuthRegisterResponse]:
        return self.send("POST", "/api/v3/auth/register", json=payload, success_model=AuthRegisterResponse)

    def login(self, payload: dict[str, Any]) -> EndpointResult[AuthLoginResponse]:
        return self.send("POST", "/api/v3/auth/login", json=payload, success_model=AuthLoginResponse)

    def refresh(self, payload: dict[str, Any]) -> EndpointResult[AuthRefreshResponse]:
        return self.send("POST", "/api/v3/auth/refresh", json=payload, success_model=AuthRefreshResponse)

    def logout(self, access_token: str, refresh_token: str | None = None, token_type: str = "Bearer") -> EndpointResult[AuthLogoutResponse]:
        body: dict[str, Any] = {}
        if refresh_token:
            body["refresh_token"] = refresh_token
        return self.send(
            "POST",
            "/api/v3/auth/logout",
            json=body,
            headers=self.build_auth_headers(access_token, token_type),
            success_model=AuthLogoutResponse,
        )

    def me(self, access_token: str, token_type: str = "Bearer") -> EndpointResult[AuthMeResponse]:
        return self.send(
            "GET",
            "/api/v3/auth/me",
            headers=self.build_auth_headers(access_token, token_type),
            success_model=AuthMeResponse,
        )

    def sso_start(self, payload: dict[str, Any]) -> EndpointResult[AuthSsoStartResponse]:
        return self.send("POST", "/api/v3/auth/sso/start", json=payload, success_model=AuthSsoStartResponse)

    def sso_callback(self, payload: dict[str, Any]) -> EndpointResult[AuthSsoCallbackResponse]:
        return self.send("POST", "/api/v3/auth/sso/callback", json=payload, success_model=AuthSsoCallbackResponse)

    def list_studio_content_items(
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
    ) -> EndpointResult[ContentListResp]:
        """
        List studio content items with optional filters.
        按条件分页查询 studio 内容条目。
        """

        return self.send(
            "GET",
            "/api/v3/studio/content/items",
            headers=self.build_auth_headers(access_token, token_type),
            params=self._build_query(
                {
                    "page": page,
                    "page_size": page_size,
                    "type": type,
                    "status": status,
                    "visibility": visibility,
                    "keyword": keyword,
                }
            ),
            success_model=ContentListResp,
        )

    def create_studio_content(
        self,
        access_token: str,
        payload: dict[str, Any],
        token_type: str = "Bearer",
    ) -> EndpointResult[ContentDetailResp]:
        """
        Create one content item in studio.
        在 studio 中创建内容。
        """

        return self.send(
            "POST",
            "/api/v3/studio/content/items",
            headers=self.build_auth_headers(access_token, token_type),
            json=payload,
            success_model=ContentDetailResp,
        )

    def get_studio_content(
        self,
        access_token: str,
        content_id: str,
        token_type: str = "Bearer",
    ) -> EndpointResult[ContentDetailResp]:
        """
        Get one content item by content id.
        按 content_id 读取单个内容。
        """

        return self.send(
            "GET",
            f"/api/v3/studio/content/items/{content_id}",
            headers=self.build_auth_headers(access_token, token_type),
            success_model=ContentDetailResp,
        )

    def update_studio_content(
        self,
        access_token: str,
        content_id: str,
        payload: dict[str, Any],
        token_type: str = "Bearer",
    ) -> EndpointResult[ContentDetailResp]:
        """
        Update a content item by content id.
        更新单个内容条目。
        """

        return self.send(
            "PUT",
            f"/api/v3/studio/content/items/{content_id}",
            headers=self.build_auth_headers(access_token, token_type),
            json=payload,
            success_model=ContentDetailResp,
        )

    def archive_studio_content(
        self,
        access_token: str,
        content_id: str,
        token_type: str = "Bearer",
    ) -> EndpointResult[ContentArchiveResp]:
        """
        Archive a content item by content id.
        归档内容。
        """

        return self.send(
            "DELETE",
            f"/api/v3/studio/content/items/{content_id}",
            headers=self.build_auth_headers(access_token, token_type),
            success_model=ContentArchiveResp,
        )

    def list_content_revisions(
        self,
        access_token: str,
        content_id: str,
        *,
        page: int = 1,
        page_size: int = 20,
        token_type: str = "Bearer",
    ) -> EndpointResult[ContentRevisionListResp]:
        """
        List content revisions.
        列出内容版本。
        """

        return self.send(
            "GET",
            f"/api/v3/studio/content/items/{content_id}/revisions",
            headers=self.build_auth_headers(access_token, token_type),
            params=self._build_query({"page": page, "page_size": page_size}),
            success_model=ContentRevisionListResp,
        )

    def get_content_revision(
        self,
        access_token: str,
        content_id: str,
        revision_id: str,
        token_type: str = "Bearer",
    ) -> EndpointResult[ContentRevisionDetailResp]:
        """
        Get content revision detail.
        获取内容版本详情。
        """

        return self.send(
            "GET",
            f"/api/v3/studio/content/items/{content_id}/revisions/{revision_id}",
            headers=self.build_auth_headers(access_token, token_type),
            success_model=ContentRevisionDetailResp,
        )

    def list_content_relations(
        self,
        access_token: str,
        content_id: str,
        *,
        page: int = 1,
        page_size: int = 20,
        relation_type: str | None = None,
        token_type: str = "Bearer",
    ) -> EndpointResult[ContentRelationListResp]:
        """
        List content relations.
        列出内容关系。
        """

        return self.send(
            "GET",
            f"/api/v3/studio/content/items/{content_id}/relations",
            headers=self.build_auth_headers(access_token, token_type),
            params=self._build_query(
                {
                    "page": page,
                    "page_size": page_size,
                    "relation_type": relation_type,
                }
            ),
            success_model=ContentRelationListResp,
        )

    def create_content_relation(
        self,
        access_token: str,
        content_id: str,
        payload: dict[str, Any],
        token_type: str = "Bearer",
    ) -> EndpointResult[ContentRelationResp]:
        """
        Create one relation for a content item.
        创建一条内容关系。
        """

        return self.send(
            "POST",
            f"/api/v3/studio/content/items/{content_id}/relations",
            headers=self.build_auth_headers(access_token, token_type),
            json=payload,
            success_model=ContentRelationResp,
        )

    def delete_content_relation(
        self,
        access_token: str,
        content_id: str,
        relation_id: str,
        token_type: str = "Bearer",
    ) -> EndpointResult[ContentRelationDeleteResp]:
        """
        Delete one relation by relation id.
        按 relation id 删除一条关系。
        """

        return self.send(
            "DELETE",
            f"/api/v3/studio/content/items/{content_id}/relations/{relation_id}",
            headers=self.build_auth_headers(access_token, token_type),
            success_model=ContentRelationDeleteResp,
        )

    def list_content_tags(
        self,
        access_token: str,
        *,
        page: int = 1,
        page_size: int = 20,
        keyword: str | None = None,
        token_type: str = "Bearer",
    ) -> EndpointResult[ContentTagListResp]:
        """
        List content tags.
        查询内容标签。
        """

        return self.send(
            "GET",
            "/api/v3/studio/content/tags",
            headers=self.build_auth_headers(access_token, token_type),
            params=self._build_query(
                {
                    "page": page,
                    "page_size": page_size,
                    "keyword": keyword,
                }
            ),
            success_model=ContentTagListResp,
        )

    def create_content_tag(
        self,
        access_token: str,
        payload: dict[str, Any],
        token_type: str = "Bearer",
    ) -> EndpointResult[ContentTagResp]:
        """
        Create one content tag.
        新建内容标签。
        """

        return self.send(
            "POST",
            "/api/v3/studio/content/tags",
            headers=self.build_auth_headers(access_token, token_type),
            json=payload,
            success_model=ContentTagResp,
        )

    def update_content_tag(
        self,
        access_token: str,
        tag_id: str,
        payload: dict[str, Any],
        token_type: str = "Bearer",
    ) -> EndpointResult[ContentTagResp]:
        """
        Update one content tag.
        更新内容标签。
        """

        return self.send(
            "PUT",
            f"/api/v3/studio/content/tags/{tag_id}",
            headers=self.build_auth_headers(access_token, token_type),
            json=payload,
            success_model=ContentTagResp,
        )

    def delete_content_tag(
        self,
        access_token: str,
        tag_id: str,
        token_type: str = "Bearer",
    ) -> EndpointResult[ContentTagDeleteResp]:
        """
        Delete one content tag.
        删除内容标签。
        """

        return self.send(
            "DELETE",
            f"/api/v3/studio/content/tags/{tag_id}",
            headers=self.build_auth_headers(access_token, token_type),
            success_model=ContentTagDeleteResp,
        )

    def list_public_content_items(
        self,
        *,
        page: int = 1,
        page_size: int = 20,
        type: str | None = None,
        keyword: str | None = None,
    ) -> EndpointResult[PublicContentListResp]:
        """
        List public content items.
        查询公开内容。
        """

        return self.send(
            "GET",
            "/api/v3/public/content/items",
            params=self._build_query(
                {
                    "page": page,
                    "page_size": page_size,
                    "type": type,
                    "keyword": keyword,
                }
            ),
            success_model=PublicContentListResp,
        )

    def get_public_content_by_slug(self, slug: str) -> EndpointResult[PublicContentGetResp]:
        """
        Get public content by slug.
        按 slug 获取公开内容。
        """

        return self.send(
            "GET",
            f"/api/v3/public/content/items/{slug}",
            success_model=PublicContentGetResp,
        )

    @staticmethod
    def build_auth_headers(access_token: str, token_type: str = "Bearer") -> dict[str, str]:
        """
        Build the Authorization header for protected endpoints.
        为受保护接口构建 Authorization 请求头。
        """

        return {"Authorization": f"{token_type} {access_token}"}

    @staticmethod
    def _build_query(values: dict[str, Any | None]) -> dict[str, Any]:
        """
        Remove null values from query parameters.
        过滤空查询参数。
        """

        return {key: value for key, value in values.items() if value is not None}

    @staticmethod
    def _load_json(response: httpx.Response) -> Any:
        try:
            return response.json()
        except ValueError as exc:
            preview = response.text[:256]
            raise GatewayResponseDecodeError(
                f"gateway returned non-JSON response for {response.request.method} {response.request.url}: {preview}"
            ) from exc

    @staticmethod
    def _load_error_payload(response: httpx.Response) -> Any:
        """
        Best-effort payload parsing for non-success responses.
        为非成功响应执行尽力而为的载荷解析。
        """

        if not response.content:
            return None

        try:
            return response.json()
        except ValueError:
            return response.text[:256]

    @staticmethod
    def _validate_model(model: type[ModelT], payload: Any, response: httpx.Response) -> ModelT:
        try:
            return model.model_validate(payload)
        except ValidationError as exc:
            raise GatewayResponseDecodeError(
                f"gateway returned unexpected payload for {response.request.method} {response.request.url}: {exc}"
            ) from exc

    @staticmethod
    def _try_validate_error(payload: Any) -> ErrorResponse | None:
        if not isinstance(payload, dict):
            return None
        try:
            return ErrorResponse.model_validate(payload)
        except ValidationError:
            return None
