package auth

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
)

func rpcContextWithMeta(ctx context.Context) context.Context {
	requestMeta, ok := middleware.RequestMetaFrom(ctx)
	if !ok {
		return ctx
	}
	return ctxmeta.OutgoingContextWithRequestMeta(ctx, requestMeta)
}
