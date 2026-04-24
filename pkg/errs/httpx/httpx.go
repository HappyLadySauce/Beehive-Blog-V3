package httpx

import (
	"context"
	"net/http"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	gozerohttpx "github.com/zeromicro/go-zero/rest/httpx"
)

// ErrorResponse defines the stable HTTP error payload.
// ErrorResponse 定义稳定的 HTTP 错误响应结构。
type ErrorResponse struct {
	Code      errs.Code `json:"code"`
	Message   string    `json:"message"`
	Reference string    `json:"reference"`
	RequestID string    `json:"request_id"`
}

// CodeToHTTP maps a business code to an HTTP status code.
// CodeToHTTP 将业务错误码映射为 HTTP 状态码。
func CodeToHTTP(code errs.Code) int {
	switch code {
	case errs.CodeGatewayNotReady, errs.CodeGatewayAuthServiceUnavailable, errs.CodeGatewayContentServiceUnavailable:
		return http.StatusServiceUnavailable
	case errs.CodeGatewayUpstreamTimeout:
		return http.StatusGatewayTimeout
	case errs.CodeContentNotFound, errs.CodeContentTagNotFound, errs.CodeContentRevisionNotFound:
		return http.StatusNotFound
	case errs.CodeContentTagInUse:
		return http.StatusPreconditionFailed
	case errs.CodeIdentityAccountPending, errs.CodeIdentityAccountDisabled, errs.CodeIdentityAccountLocked,
		errs.CodeIdentitySSOProviderDisabled, errs.CodeIdentitySSOProviderNotReady, errs.CodeIdentitySSOStateInvalid:
		return http.StatusPreconditionFailed
	}

	category := (int(code) / 100) % 100
	switch category {
	case 1:
		return http.StatusBadRequest
	case 2:
		return http.StatusUnauthorized
	case 3:
		return http.StatusForbidden
	case 4:
		return http.StatusPreconditionFailed
	case 5:
		if int(code)%100 == 1 {
			return http.StatusNotFound
		}
		return http.StatusConflict
	case 6:
		return http.StatusServiceUnavailable
	case 99:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// BuildResponse builds a stable HTTP error response payload.
// BuildResponse 构建稳定的 HTTP 错误响应载荷。
func BuildResponse(err error, requestID string) (int, ErrorResponse) {
	parsed := errs.Parse(err)
	if parsed == nil {
		parsed = &errs.Error{
			Code:    errs.CodeGatewayInternal,
			Message: "internal server error",
		}
	}

	return CodeToHTTP(parsed.Code), ErrorResponse{
		Code:      parsed.Code,
		Message:   parsed.Message,
		Reference: parsed.Reference,
		RequestID: requestID,
	}
}

// WriteError writes a stable HTTP error response.
// WriteError 写入稳定的 HTTP 错误响应。
func WriteError(ctx context.Context, w http.ResponseWriter, err error, requestID string) {
	statusCode, payload := BuildResponse(err, requestID)
	gozerohttpx.WriteJsonCtx(ctx, w, statusCode, payload)
}
