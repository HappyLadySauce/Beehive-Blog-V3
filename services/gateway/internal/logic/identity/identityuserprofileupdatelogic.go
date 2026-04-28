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

type IdentityUserProfileUpdateLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIdentityUserProfileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IdentityUserProfileUpdateLogic {
	return &IdentityUserProfileUpdateLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// IdentityUserProfileUpdate updates basic user profile fields through identity RPC.
// IdentityUserProfileUpdate 通过 identity RPC 更新用户基础资料字段。
func (l *IdentityUserProfileUpdateLogic) IdentityUserProfileUpdate(req *types.AdminUpdateUserProfileReq) (resp *types.AdminUserResp, err error) {
	authCtx, err := requireAdminContext(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcResp, rpcErr := l.svcCtx.IdentityClient.UpdateUserProfile(
		rpcContextWithMeta(l.ctx, l.svcCtx.Config.IdentityRPC),
		identityadapter.BuildUpdateUserProfileRequest(authCtx.UserID, req),
	)
	if rpcErr != nil {
		return nil, identityerrors.MapUpstreamError(l.ctx, "identity_user_profile_update", "/api/v3/studio/users/:user_id/profile", rpcErr)
	}

	return identityadapter.ToAdminUserResponse(rpcResp.GetUser()), nil
}
