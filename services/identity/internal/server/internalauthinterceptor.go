package server

import (
	"context"
	"crypto/subtle"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	identityServicePrefix = "/identity.Identity/"
)

// InternalAuthInterceptor authenticates all business RPC calls from trusted internal callers.
// InternalAuthInterceptor 对所有业务 RPC 调用执行受信内部调用方认证。
type InternalAuthInterceptor struct {
	token          string
	allowedCallers map[string]struct{}
}

// NewInternalAuthInterceptor creates an internal caller authentication interceptor.
// NewInternalAuthInterceptor 创建内部调用方认证拦截器。
func NewInternalAuthInterceptor(conf config.Config) *InternalAuthInterceptor {
	allowed := make(map[string]struct{}, len(conf.AllowedCallers))
	for _, caller := range conf.AllowedCallers {
		trimmed := strings.TrimSpace(caller)
		if trimmed == "" {
			continue
		}
		allowed[trimmed] = struct{}{}
	}

	return &InternalAuthInterceptor{
		token:          strings.TrimSpace(conf.InternalAuthToken),
		allowedCallers: allowed,
	}
}

// Unary returns the unary server interceptor.
// Unary 返回 unary server interceptor。
func (i *InternalAuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if info == nil || !strings.HasPrefix(info.FullMethod, identityServicePrefix) {
			return handler(ctx, req)
		}

		caller, authErr := i.authenticate(ctx)
		if authErr != nil {
			requestID := ctxmeta.GetRequestIDFromIncomingContext(ctx)
			logs.Ctx(ctx).Warn(
				"identity_internal_auth_failed",
				logs.String("caller", caller),
				logs.RequestID(requestID),
				logs.Err(authErr),
			)
			return nil, errgrpcx.ToStatus(authErr, "internal caller authentication failed")
		}

		logs.Ctx(ctx).Debug(
			"identity_internal_auth",
			logs.String("caller", caller),
			logs.RequestID(ctxmeta.GetRequestIDFromIncomingContext(ctx)),
		)

		return handler(ctxmeta.WithTrustedInternalCaller(ctx, caller), req)
	}
}

func (i *InternalAuthInterceptor) authenticate(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errs.New(errs.CodeIdentityInternalCallerUnauthorized, "internal caller authentication failed")
	}

	token := strings.TrimSpace(firstMetadataValue(md, ctxmeta.MetadataKeyInternalAuthToken))
	caller := strings.TrimSpace(firstMetadataValue(md, ctxmeta.MetadataKeyInternalCaller))
	if token == "" || caller == "" {
		return caller, errs.New(errs.CodeIdentityInternalCallerUnauthorized, "internal caller authentication failed")
	}
	if subtle.ConstantTimeCompare([]byte(token), []byte(i.token)) != 1 {
		return caller, errs.New(errs.CodeIdentityInternalCallerUnauthorized, "internal caller authentication failed")
	}
	if _, ok := i.allowedCallers[caller]; !ok {
		return caller, errs.New(errs.CodeIdentityInternalCallerUnauthorized, "internal caller authentication failed")
	}

	return caller, nil
}

func firstMetadataValue(md metadata.MD, key string) string {
	values := md.Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}
