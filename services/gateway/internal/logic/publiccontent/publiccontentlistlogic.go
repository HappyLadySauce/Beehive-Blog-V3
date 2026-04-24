// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package publiccontent

import (
	"context"

	contentadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/content"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type PublicContentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// List public content items
func NewPublicContentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicContentListLogic {
	return &PublicContentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicContentListLogic) PublicContentList(req *types.PublicContentListReq) (resp *types.ContentListResp, err error) {
	rpcReq, err := contentadapter.BuildPublicListRequest(req)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.ListPublicContents(publicRPCContext(l.ctx, l.svcCtx.Config.ContentRPC), rpcReq)
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "public_content_list", "/api/v3/public/content/items", rpcErr)
	}
	return contentadapter.ToContentListResp(rpcResp.GetItems(), rpcResp.GetTotal(), rpcResp.GetPage(), rpcResp.GetPageSize()), nil
}
