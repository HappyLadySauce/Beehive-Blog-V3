package file

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
)

func rpcContextWithActor(ctx context.Context, rpcConf config.FileRPCConf) (context.Context, string, error) {
	return fileadapter.RPCContext(ctx, rpcConf)
}

func mapFileError(ctx context.Context, action string, route string, err error) error {
	return fileadapter.MapUpstreamError(ctx, action, route, err)
}
