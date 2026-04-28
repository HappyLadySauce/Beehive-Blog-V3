package file

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
)

func RPCContext(ctx context.Context, rpcConf config.FileRPCConf) (context.Context, string, error) {
	authCtx, ok := middleware.AuthContextFrom(ctx)
	if !ok || authCtx.UserID == "" {
		return nil, "", errs.New(errs.CodeGatewayAuthorizationRequired, "trusted auth context is missing")
	}
	requestMeta, ok := middleware.RequestMetaFrom(ctx)
	if !ok {
		requestMeta = ctxmeta.RequestMeta{}
	}
	rpcCtx := ctxmeta.BuildInternalOutgoingContext(ctx, requestMeta, ctxmeta.InternalRPCAuth{
		Token:  rpcConf.InternalAuthToken,
		Caller: rpcConf.InternalCallerName,
	}, ctxmeta.AuthClaims{
		UserID:    authCtx.UserID,
		SessionID: authCtx.SessionID,
		Role:      authCtx.Role,
	})
	return rpcCtx, authCtx.UserID, nil
}
