package auth

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
)

func rpcContextWithMeta(ctx context.Context, rpcConf config.IdentityRPCConf) context.Context {
	requestMeta, ok := middleware.RequestMetaFrom(ctx)
	if !ok {
		return ctxmeta.BuildIdentityOutgoingContext(ctx, ctxmeta.RequestMeta{}, ctxmeta.InternalRPCAuth{
			Token:  rpcConf.InternalAuthToken,
			Caller: rpcConf.InternalCallerName,
		})
	}
	return ctxmeta.BuildIdentityOutgoingContext(ctx, requestMeta, ctxmeta.InternalRPCAuth{
		Token:  rpcConf.InternalAuthToken,
		Caller: rpcConf.InternalCallerName,
	})
}
