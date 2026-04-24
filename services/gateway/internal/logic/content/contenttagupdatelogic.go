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

type ContentTagUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Update content tag
func NewContentTagUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentTagUpdateLogic {
	return &ContentTagUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentTagUpdateLogic) ContentTagUpdate(req *types.ContentTagUpdateReq) (resp *types.ContentTagResp, err error) {
	rpcCtx, err := studioRPCContext(l.ctx, l.svcCtx.Config.ContentRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.UpdateTag(rpcCtx, &contentpb.UpdateTagRequest{
		TagId:       req.TagId,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Color:       req.Color,
	})
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "content_tag_update", "/api/v3/studio/content/tags/:tag_id", rpcErr)
	}
	return contentadapter.ToTagResp(rpcResp.GetTag()), nil
}
