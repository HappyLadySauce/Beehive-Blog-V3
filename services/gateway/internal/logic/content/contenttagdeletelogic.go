// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package content

import (
	"context"

	contentpb "github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type ContentTagDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete content tag
func NewContentTagDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentTagDeleteLogic {
	return &ContentTagDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentTagDeleteLogic) ContentTagDelete(req *types.ContentTagIdReq) (resp *types.ContentTagDeleteResp, err error) {
	rpcCtx, err := studioRPCContext(l.ctx, l.svcCtx.Config.ContentRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.DeleteTag(rpcCtx, &contentpb.DeleteTagRequest{TagId: req.TagId})
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "content_tag_delete", "/api/v3/studio/content/tags/:tag_id", rpcErr)
	}
	return &types.ContentTagDeleteResp{Ok: rpcResp.GetOk()}, nil
}
