package content

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	contentadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/content"
)

func studioRPCContext(ctx context.Context, rpcConf config.ContentRPCConf) (context.Context, error) {
	return contentadapter.StudioRPCContext(ctx, rpcConf)
}

func mapContentError(ctx context.Context, action string, route string, err error) error {
	return contentadapter.MapUpstreamError(ctx, action, route, err)
}
