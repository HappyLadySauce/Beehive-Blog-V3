package middleware

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
)

type contextKey string

const (
	authContextKey        contextKey = "gateway.auth.context"
	requestMetaContextKey contextKey = "gateway.request.meta"
)

// AuthContext stores trusted identity claims.
// AuthContext 保存可信身份声明。
type AuthContext struct {
	UserID        string
	SessionID     string
	Role          string
	AccountStatus string
	AuthSource    string
	AccessToken   string
}

// WithAuthContext stores trusted auth context.
// WithAuthContext 写入可信认证上下文。
func WithAuthContext(ctx context.Context, authCtx AuthContext) context.Context {
	return context.WithValue(ctx, authContextKey, authCtx)
}

// AuthContextFrom reads trusted auth context.
// AuthContextFrom 读取可信认证上下文。
func AuthContextFrom(ctx context.Context) (AuthContext, bool) {
	authCtx, ok := ctx.Value(authContextKey).(AuthContext)
	return authCtx, ok
}

// WithRequestMeta stores extracted request metadata.
// WithRequestMeta 写入提取后的请求元数据。
func WithRequestMeta(ctx context.Context, meta ctxmeta.RequestMeta) context.Context {
	return context.WithValue(ctx, requestMetaContextKey, meta)
}

// RequestMetaFrom reads request metadata.
// RequestMetaFrom 读取请求元数据。
func RequestMetaFrom(ctx context.Context) (ctxmeta.RequestMeta, bool) {
	meta, ok := ctx.Value(requestMetaContextKey).(ctxmeta.RequestMeta)
	return meta, ok
}
