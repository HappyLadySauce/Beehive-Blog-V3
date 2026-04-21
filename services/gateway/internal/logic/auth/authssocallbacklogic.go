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

type AuthSsoCallbackLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthSsoCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthSsoCallbackLogic {
	return &AuthSsoCallbackLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthSsoCallbackLogic) AuthSsoCallback(req *types.AuthSsoCallbackReq) (resp *types.AuthSsoCallbackResp, err error) {
	rpcResp, rpcErr := l.svcCtx.IdentityClient.FinishSsoLogin(
		rpcContextWithMeta(l.ctx, l.svcCtx.Config.IdentityRPC),
		identityadapter.BuildSsoCallbackRequest(req),
	)
	if rpcErr != nil {
		return nil, identityerrors.MapUpstreamError(l.ctx, "auth_sso_callback", "/api/v3/auth/sso/callback", rpcErr)
	}

	return identityadapter.ToSsoCallbackResponse(rpcResp), nil
}
