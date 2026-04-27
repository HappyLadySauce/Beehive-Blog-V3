package identity

import (
	"context"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
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

func requireAdminContext(ctx context.Context) (middleware.AuthContext, error) {
	authCtx, ok := middleware.AuthContextFrom(ctx)
	if !ok || authCtx.UserID == "" {
		return middleware.AuthContext{}, errs.New(errs.CodeGatewayAuthorizationRequired, "trusted auth context is missing")
	}
	if !isAdminRole(authCtx.Role) {
		return middleware.AuthContext{}, errs.New(errs.CodeGatewayAccessForbidden, "admin role is required")
	}
	return authCtx, nil
}

func isAdminRole(role string) bool {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "admin", "role_admin":
		return true
	default:
		return false
	}
}
