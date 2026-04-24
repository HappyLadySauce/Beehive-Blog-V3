// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package content

import (
	"context"

	contentpb "github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type ContentRelationDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete content relation
func NewContentRelationDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentRelationDeleteLogic {
	return &ContentRelationDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentRelationDeleteLogic) ContentRelationDelete(req *types.ContentRelationIdReq) (resp *types.ContentRelationDeleteResp, err error) {
	rpcCtx, err := studioRPCContext(l.ctx, l.svcCtx.Config.ContentRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.DeleteContentRelation(rpcCtx, &contentpb.DeleteContentRelationRequest{
		ContentId:  req.ContentId,
		RelationId: req.RelationId,
	})
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "content_relation_delete", "/api/v3/studio/content/items/:content_id/relations/:relation_id", rpcErr)
	}
	return &types.ContentRelationDeleteResp{Ok: rpcResp.GetOk()}, nil
}
