// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	identityadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/identity"
	identityerrors "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/identity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type AuthUpdateProfileLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthUpdateProfileLogic {
	return &AuthUpdateProfileLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthUpdateProfileLogic) AuthUpdateProfile(req *types.AuthUpdateProfileReq) (resp *types.AuthMeResp, err error) {
	authCtx, ok := middleware.AuthContextFrom(l.ctx)
	if !ok || authCtx.UserID == "" {
		return nil, errs.New(errs.CodeGatewayAuthorizationRequired, "trusted auth context is missing")
	}

	rpcResp, rpcErr := l.svcCtx.IdentityClient.UpdateOwnProfile(
		rpcContextWithMeta(l.ctx, l.svcCtx.Config.IdentityRPC),
		identityadapter.BuildUpdateProfileRequest(authCtx.UserID, req),
	)
	if rpcErr != nil {
		return nil, identityerrors.MapUpstreamError(l.ctx, "auth_update_profile", "/api/v3/auth/me/profile", rpcErr)
	}

	return identityadapter.ToUpdateProfileResponse(rpcResp), nil
}
