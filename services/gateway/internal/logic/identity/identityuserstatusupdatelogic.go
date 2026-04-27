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

type IdentityUserStatusUpdateLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIdentityUserStatusUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IdentityUserStatusUpdateLogic {
	return &IdentityUserStatusUpdateLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IdentityUserStatusUpdateLogic) IdentityUserStatusUpdate(req *types.AdminUpdateUserStatusReq) (resp *types.AdminUserResp, err error) {
	authCtx, err := requireAdminContext(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcResp, rpcErr := l.svcCtx.IdentityClient.UpdateUserStatus(
		rpcContextWithMeta(l.ctx, l.svcCtx.Config.IdentityRPC),
		identityadapter.BuildUpdateUserStatusRequest(authCtx.UserID, req),
	)
	if rpcErr != nil {
		return nil, identityerrors.MapUpstreamError(l.ctx, "identity_user_status_update", "/api/v3/studio/users/:user_id/status", rpcErr)
	}

	return identityadapter.ToAdminUserResponse(rpcResp.GetUser()), nil
}
