// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package content

import (
	"context"

	contentadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/content"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type ContentCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Create content item
func NewContentCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentCreateLogic {
	return &ContentCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentCreateLogic) ContentCreate(req *types.ContentCreateReq) (resp *types.ContentDetailResp, err error) {
	rpcCtx, err := studioRPCContext(l.ctx, l.svcCtx.Config.ContentRPC)
	if err != nil {
		return nil, err
	}
	rpcReq, err := contentadapter.BuildCreateRequest(req)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.CreateContent(rpcCtx, rpcReq)
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "content_create", "/api/v3/studio/content/items", rpcErr)
	}
	return contentadapter.ToContentDetailResp(rpcResp.GetContent()), nil
}
