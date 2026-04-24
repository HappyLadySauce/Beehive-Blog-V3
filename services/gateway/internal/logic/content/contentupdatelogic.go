// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package content

import (
	"context"

	contentadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/content"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type ContentUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Update content item
func NewContentUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentUpdateLogic {
	return &ContentUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentUpdateLogic) ContentUpdate(req *types.ContentUpdateReq) (resp *types.ContentDetailResp, err error) {
	rpcCtx, err := studioRPCContext(l.ctx, l.svcCtx.Config.ContentRPC)
	if err != nil {
		return nil, err
	}
	rpcReq, err := contentadapter.BuildUpdateRequest(req)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.UpdateContent(rpcCtx, rpcReq)
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "content_update", "/api/v3/studio/content/items/:content_id", rpcErr)
	}
	return contentadapter.ToContentDetailResp(rpcResp.GetContent()), nil
}
