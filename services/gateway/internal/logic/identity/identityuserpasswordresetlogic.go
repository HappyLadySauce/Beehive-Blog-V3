// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package identity

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	identityadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/identity"
	identityerrors "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/identity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type IdentityUserPasswordResetLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIdentityUserPasswordResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IdentityUserPasswordResetLogic {
	return &IdentityUserPasswordResetLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IdentityUserPasswordResetLogic) IdentityUserPasswordReset(req *types.AdminResetUserPasswordReq) (resp *types.AdminResetUserPasswordResp, err error) {
	authCtx, err := requireAdminContext(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcResp, rpcErr := l.svcCtx.IdentityClient.ResetUserPassword(
		rpcContextWithMeta(l.ctx, l.svcCtx.Config.IdentityRPC),
		identityadapter.BuildResetUserPasswordRequest(authCtx.UserID, req),
	)
	if rpcErr != nil {
		return nil, identityerrors.MapUpstreamError(l.ctx, "identity_user_password_reset", "/api/v3/studio/users/:user_id/password/reset", rpcErr)
	}

	return identityadapter.ToResetUserPasswordResponse(rpcResp), nil
}
