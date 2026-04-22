"""
Custom exceptions for QA HTTP interactions.
QA HTTP 交互过程中使用的自定义异常。
"""


class GatewayClientError(RuntimeError):
    """
    Base error for gateway QA client failures.
    gateway QA 客户端失败的基础异常。
    """


class GatewayTransportError(GatewayClientError):
    """
    Raised when the gateway cannot be reached over HTTP.
    当 gateway 无法通过 HTTP 访问时抛出。
    """


class GatewayResponseDecodeError(GatewayClientError):
    """
    Raised when a response cannot be decoded as the expected JSON shape.
    当响应无法按预期 JSON 结构解码时抛出。
    """

