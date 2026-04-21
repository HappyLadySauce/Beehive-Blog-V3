// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	identityadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/identity"
	identityerrors "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/identity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type AuthRefreshLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthRefreshLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthRefreshLogic {
	return &AuthRefreshLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthRefreshLogic) AuthRefresh(req *types.AuthRefreshReq) (resp *types.AuthRefreshResp, err error) {
	rpcResp, rpcErr := l.svcCtx.IdentityClient.RefreshSessionToken(
		rpcContextWithMeta(l.ctx, l.svcCtx.Config.IdentityRPC),
		identityadapter.BuildRefreshRequest(req),
	)
	if rpcErr != nil {
		return nil, identityerrors.MapUpstreamError(l.ctx, "auth_refresh", "/api/v3/auth/refresh", rpcErr)
	}

	return identityadapter.ToRefreshResponse(rpcResp), nil
}
