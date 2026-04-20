package middleware

import (
	"net/http"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/google/uuid"
)

// RequestMetaMiddleware extracts trusted request metadata once per request.
// RequestMetaMiddleware 每个请求只提取一次可信元数据。
type RequestMetaMiddleware struct {
	securityConf config.GatewaySecurityConf
}

// NewRequestMetaMiddleware creates a request metadata middleware.
// NewRequestMetaMiddleware 创建请求元数据中间件。
func NewRequestMetaMiddleware(securityConf config.GatewaySecurityConf) *RequestMetaMiddleware {
	return &RequestMetaMiddleware{
		securityConf: securityConf,
	}
}

// Handle extracts request metadata and stores it in context.
// Handle 提取请求元数据并写入上下文。
func (m *RequestMetaMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestMeta := ctxmeta.BuildRequestMetaFromHTTP(r, m.securityConf.TrustedProxyHeaders)
		if requestMeta.RequestID == "" {
			requestMeta.RequestID = uuid.NewString()
		}

		ctx := logs.WithRequestID(r.Context(), requestMeta.RequestID)
		ctx = WithRequestMeta(ctx, requestMeta)
		r = r.WithContext(ctx)
		logs.Ctx(ctx).Info(
			"request_meta_extract",
			logs.String("route", r.URL.Path),
			logs.String("client_ip", requestMeta.ClientIP),
		)

		next(w, r)
	}
}
