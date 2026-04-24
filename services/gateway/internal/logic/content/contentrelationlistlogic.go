// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package content

import (
	"context"

	contentadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/content"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type ContentRelationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// List content relations
func NewContentRelationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentRelationListLogic {
	return &ContentRelationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentRelationListLogic) ContentRelationList(req *types.ContentRelationListReq) (resp *types.ContentRelationListResp, err error) {
	rpcCtx, err := studioRPCContext(l.ctx, l.svcCtx.Config.ContentRPC)
	if err != nil {
		return nil, err
	}
	rpcReq, err := contentadapter.BuildListRelationsRequest(req)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.ListContentRelations(rpcCtx, rpcReq)
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "content_relation_list", "/api/v3/studio/content/items/:content_id/relations", rpcErr)
	}
	return contentadapter.ToRelationListResp(rpcResp.GetItems(), rpcResp.GetTotal(), rpcResp.GetPage(), rpcResp.GetPageSize()), nil
}
