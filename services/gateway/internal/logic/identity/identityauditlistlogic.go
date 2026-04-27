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

type IdentityAuditListLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIdentityAuditListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IdentityAuditListLogic {
	return &IdentityAuditListLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IdentityAuditListLogic) IdentityAuditList(req *types.IdentityAuditListReq) (resp *types.IdentityAuditListResp, err error) {
	authCtx, err := requireAdminContext(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcResp, rpcErr := l.svcCtx.IdentityClient.ListIdentityAudits(
		rpcContextWithMeta(l.ctx, l.svcCtx.Config.IdentityRPC),
		identityadapter.BuildListAuditsRequest(authCtx.UserID, req),
	)
	if rpcErr != nil {
		return nil, identityerrors.MapUpstreamError(l.ctx, "identity_audit_list", "/api/v3/studio/audits", rpcErr)
	}

	return identityadapter.ToAuditListResponse(rpcResp), nil
}
