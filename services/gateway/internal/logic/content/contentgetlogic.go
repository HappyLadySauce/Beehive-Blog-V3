// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package content

import (
	"context"

	contentpb "github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	contentadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/content"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type ContentGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get studio content item
func NewContentGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentGetLogic {
	return &ContentGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentGetLogic) ContentGet(req *types.ContentIdReq) (resp *types.ContentDetailResp, err error) {
	rpcCtx, err := studioRPCContext(l.ctx, l.svcCtx.Config.ContentRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.GetContent(rpcCtx, &contentpb.GetContentRequest{ContentId: req.ContentId})
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "content_get", "/api/v3/studio/content/items/:content_id", rpcErr)
	}
	return contentadapter.ToContentDetailResp(rpcResp.GetContent()), nil
}
