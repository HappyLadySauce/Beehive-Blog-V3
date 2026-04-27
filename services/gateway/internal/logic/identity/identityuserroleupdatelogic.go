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

type IdentityUserRoleUpdateLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIdentityUserRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IdentityUserRoleUpdateLogic {
	return &IdentityUserRoleUpdateLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IdentityUserRoleUpdateLogic) IdentityUserRoleUpdate(req *types.AdminUpdateUserRoleReq) (resp *types.AdminUserResp, err error) {
	authCtx, err := requireAdminContext(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcResp, rpcErr := l.svcCtx.IdentityClient.UpdateUserRole(
		rpcContextWithMeta(l.ctx, l.svcCtx.Config.IdentityRPC),
		identityadapter.BuildUpdateUserRoleRequest(authCtx.UserID, req),
	)
	if rpcErr != nil {
		return nil, identityerrors.MapUpstreamError(l.ctx, "identity_user_role_update", "/api/v3/studio/users/:user_id/role", rpcErr)
	}

	return identityadapter.ToAdminUserResponse(rpcResp.GetUser()), nil
}
