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

type AuthMeLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthMeLogic {
	return &AuthMeLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthMeLogic) AuthMe(req *types.AuthMeReq) (resp *types.AuthMeResp, err error) {
	authCtx, ok := middleware.AuthContextFrom(l.ctx)
	if !ok || authCtx.UserID == "" {
		return nil, errs.New(errs.CodeGatewayAuthorizationRequired, "trusted auth context is missing")
	}

	rpcResp, rpcErr := l.svcCtx.IdentityClient.GetCurrentUser(
		rpcContextWithMeta(l.ctx),
		identityadapter.BuildMeRequest(authCtx.UserID),
	)
	if rpcErr != nil {
		return nil, identityerrors.MapUpstreamError(l.ctx, "auth_me", "/api/v3/auth/me", rpcErr)
	}

	return identityadapter.ToMeResponse(rpcResp), nil
}
