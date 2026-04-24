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

type ContentRevisionGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get content revision
func NewContentRevisionGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentRevisionGetLogic {
	return &ContentRevisionGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentRevisionGetLogic) ContentRevisionGet(req *types.ContentRevisionIdReq) (resp *types.ContentRevisionDetailResp, err error) {
	rpcCtx, err := studioRPCContext(l.ctx, l.svcCtx.Config.ContentRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.GetContentRevision(rpcCtx, &contentpb.GetContentRevisionRequest{
		ContentId:  req.ContentId,
		RevisionId: req.RevisionId,
	})
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "content_revision_get", "/api/v3/studio/content/items/:content_id/revisions/:revision_id", rpcErr)
	}
	return contentadapter.ToRevisionDetailResp(rpcResp.GetRevision()), nil
}
