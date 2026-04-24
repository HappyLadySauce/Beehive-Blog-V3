package content

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
)

func StudioRPCContext(ctx context.Context, rpcConf config.ContentRPCConf) (context.Context, error) {
	authCtx, ok := middleware.AuthContextFrom(ctx)
	if !ok || authCtx.UserID == "" {
		return nil, errs.New(errs.CodeGatewayAuthorizationRequired, "trusted auth context is missing")
	}
	return internalRPCContext(ctx, rpcConf, ctxmeta.AuthClaims{
		UserID:    authCtx.UserID,
		SessionID: authCtx.SessionID,
		Role:      authCtx.Role,
	}), nil
}

func PublicRPCContext(ctx context.Context, rpcConf config.ContentRPCConf) context.Context {
	return internalRPCContext(ctx, rpcConf, ctxmeta.AuthClaims{})
}

func internalRPCContext(ctx context.Context, rpcConf config.ContentRPCConf, claims ctxmeta.AuthClaims) context.Context {
	requestMeta, ok := middleware.RequestMetaFrom(ctx)
	if !ok {
		requestMeta = ctxmeta.RequestMeta{}
	}
	return ctxmeta.BuildInternalOutgoingContext(ctx, requestMeta, ctxmeta.InternalRPCAuth{
		Token:  rpcConf.InternalAuthToken,
		Caller: rpcConf.InternalCallerName,
	}, claims)
}
