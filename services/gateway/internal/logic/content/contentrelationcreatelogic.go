// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package content

import (
	"context"

	contentadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/content"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type ContentRelationCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Create content relation
func NewContentRelationCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentRelationCreateLogic {
	return &ContentRelationCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentRelationCreateLogic) ContentRelationCreate(req *types.ContentRelationCreateReq) (resp *types.ContentRelationResp, err error) {
	rpcCtx, err := studioRPCContext(l.ctx, l.svcCtx.Config.ContentRPC)
	if err != nil {
		return nil, err
	}
	rpcReq, err := contentadapter.BuildCreateRelationRequest(req)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.CreateContentRelation(rpcCtx, rpcReq)
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "content_relation_create", "/api/v3/studio/content/items/:content_id/relations", rpcErr)
	}
	return contentadapter.ToRelationResp(rpcResp.GetRelation()), nil
}
