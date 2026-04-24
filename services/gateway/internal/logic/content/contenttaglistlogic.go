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

type ContentTagListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// List content tags
func NewContentTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentTagListLogic {
	return &ContentTagListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentTagListLogic) ContentTagList(req *types.ContentTagListReq) (resp *types.ContentTagListResp, err error) {
	rpcCtx, err := studioRPCContext(l.ctx, l.svcCtx.Config.ContentRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.ListTags(rpcCtx, &contentpb.ListTagsRequest{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
		Keyword:  req.Keyword,
	})
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "content_tag_list", "/api/v3/studio/content/tags", rpcErr)
	}
	return contentadapter.ToTagListResp(rpcResp.GetItems(), rpcResp.GetTotal(), rpcResp.GetPage(), rpcResp.GetPageSize()), nil
}
