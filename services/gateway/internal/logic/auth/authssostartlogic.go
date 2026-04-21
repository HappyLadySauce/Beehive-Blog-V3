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

type AuthSsoStartLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthSsoStartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthSsoStartLogic {
	return &AuthSsoStartLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthSsoStartLogic) AuthSsoStart(req *types.AuthSsoStartReq) (resp *types.AuthSsoStartResp, err error) {
	rpcResp, rpcErr := l.svcCtx.IdentityClient.StartSsoLogin(
		rpcContextWithMeta(l.ctx, l.svcCtx.Config.IdentityRPC),
		identityadapter.BuildSsoStartRequest(req),
	)
	if rpcErr != nil {
		return nil, identityerrors.MapUpstreamError(l.ctx, "auth_sso_start", "/api/v3/auth/sso/start", rpcErr)
	}

	return identityadapter.ToSsoStartResponse(rpcResp), nil
}
