// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package content

import (
	"context"

	contentadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/content"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type ContentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// List studio content items
func NewContentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentListLogic {
	return &ContentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentListLogic) ContentList(req *types.ContentListReq) (resp *types.ContentListResp, err error) {
	rpcCtx, err := studioRPCContext(l.ctx, l.svcCtx.Config.ContentRPC)
	if err != nil {
		return nil, err
	}
	rpcReq, err := contentadapter.BuildListRequest(req)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.ListStudioContents(rpcCtx, rpcReq)
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "content_list", "/api/v3/studio/content/items", rpcErr)
	}
	return contentadapter.ToContentListResp(rpcResp.GetItems(), rpcResp.GetTotal(), rpcResp.GetPage(), rpcResp.GetPageSize()), nil
}
