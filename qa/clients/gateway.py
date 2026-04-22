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
    ErrorResponse,
    HealthzResponse,
    ReadyzResponse,
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
        json: dict[str, Any] | None = None,
    ) -> EndpointResult[ModelT]:
        """
        Send a request and parse it into success or error models.
        发送请求，并将其解析为成功或错误模型。
        """

        try:
            response = self._client.request(method, path, headers=headers, json=json)
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

    @staticmethod
    def build_auth_headers(access_token: str, token_type: str = "Bearer") -> dict[str, str]:
        """
        Build the Authorization header for protected endpoints.
        为受保护接口构建 Authorization 请求头。
        """

        return {"Authorization": f"{token_type} {access_token}"}

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
