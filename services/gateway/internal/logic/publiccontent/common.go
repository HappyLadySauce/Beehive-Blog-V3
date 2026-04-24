package publiccontent

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	contentadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/content"
)

func publicRPCContext(ctx context.Context, rpcConf config.ContentRPCConf) context.Context {
	return contentadapter.PublicRPCContext(ctx, rpcConf)
}

func mapContentError(ctx context.Context, action string, route string, err error) error {
	return contentadapter.MapUpstreamError(ctx, action, route, err)
}
