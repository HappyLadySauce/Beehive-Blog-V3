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

type ContentRevisionListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// List content revisions
func NewContentRevisionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentRevisionListLogic {
	return &ContentRevisionListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentRevisionListLogic) ContentRevisionList(req *types.ContentRevisionListReq) (resp *types.ContentRevisionListResp, err error) {
	rpcCtx, err := studioRPCContext(l.ctx, l.svcCtx.Config.ContentRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.ListContentRevisions(rpcCtx, &contentpb.ListContentRevisionsRequest{
		ContentId: req.ContentId,
		Page:      int32(req.Page),
		PageSize:  int32(req.PageSize),
	})
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "content_revision_list", "/api/v3/studio/content/items/:content_id/revisions", rpcErr)
	}
	return contentadapter.ToRevisionListResp(rpcResp.GetItems(), rpcResp.GetTotal(), rpcResp.GetPage(), rpcResp.GetPageSize()), nil
}
