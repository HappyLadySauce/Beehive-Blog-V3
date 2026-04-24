// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package content

import (
	"context"

	contentpb "github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type ContentArchiveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Archive content item
func NewContentArchiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ContentArchiveLogic {
	return &ContentArchiveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ContentArchiveLogic) ContentArchive(req *types.ContentIdReq) (resp *types.ContentArchiveResp, err error) {
	rpcCtx, err := studioRPCContext(l.ctx, l.svcCtx.Config.ContentRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.ContentClient.ArchiveContent(rpcCtx, &contentpb.ArchiveContentRequest{ContentId: req.ContentId})
	if rpcErr != nil {
		return nil, mapContentError(l.ctx, "content_archive", "/api/v3/studio/content/items/:content_id", rpcErr)
	}
	return &types.ContentArchiveResp{Ok: rpcResp.GetOk()}, nil
}
