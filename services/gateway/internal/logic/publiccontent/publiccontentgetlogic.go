// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package publiccontent

import (
	"context"

	contentpb "github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	contentadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/content"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type PublicContentGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get public content item by slug
func NewPublicContentGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicContentGetLogic {
	return &PublicContentGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublicContentGetLogic) PublicContentGet(req *types.PublicContentSlugReq) (resp *types.ContentDetailResp, err error) {
	rpcResp, rpcErr := l.svcCtx.ContentClient.GetPublicContentBySlug(publicRPCContext(l.ctx, l.svcCtx.Config.ContentRPC), &contentpb.GetPublicContentBySlugRequest{Slug: req.Slug})
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "public_content_get", "/api/v3/public/content/items/:slug", rpcErr)
	}
	return contentadapter.ToContentDetailResp(rpcResp.GetContent()), nil
}
