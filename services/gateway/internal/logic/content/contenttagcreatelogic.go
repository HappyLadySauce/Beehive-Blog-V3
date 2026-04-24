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

type ContentTagCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Create content tag
func NewContentTagCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentTagCreateLogic {
	return &ContentTagCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentTagCreateLogic) ContentTagCreate(req *types.ContentTagCreateReq) (resp *types.ContentTagResp, err error) {
	rpcCtx, err := studioRPCContext(l.ctx, l.svcCtx.Config.ContentRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.CreateTag(rpcCtx, &contentpb.CreateTagRequest{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Color:       req.Color,
	})
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "content_tag_create", "/api/v3/studio/content/tags", rpcErr)
	}
	return contentadapter.ToTagResp(rpcResp.GetTag()), nil
}
