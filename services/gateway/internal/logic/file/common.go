package file

import (
	"context"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
)

func rpcContextWithActor(ctx context.Context, rpcConf config.FileRPCConf) (context.Context, string, error) {
	return fileadapter.RPCContext(ctx, rpcConf)
}

func rpcContextWithAdminActor(ctx context.Context, rpcConf config.FileRPCConf) (context.Context, string, error) {
	authCtx, ok := middleware.AuthContextFrom(ctx)
	if !ok || authCtx.UserID == "" {
		return nil, "", errs.New(errs.CodeGatewayAuthorizationRequired, "trusted auth context is missing")
	}
	if !isAdminRole(authCtx.Role) {
		return nil, "", errs.New(errs.CodeGatewayAccessForbidden, "admin role is required")
	}
	return fileadapter.RPCContext(ctx, rpcConf)
}

func isAdminRole(role string) bool {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "admin", "role_admin":
		return true
	default:
		return false
	}
}

func mapFileError(ctx context.Context, action string, route string, err error) error {
	return fileadapter.MapUpstreamError(ctx, action, route, err)
}
