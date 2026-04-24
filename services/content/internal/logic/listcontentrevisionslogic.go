package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type ListContentRevisionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListContentRevisionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListContentRevisionsLogic {
	return &ListContentRevisionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListContentRevisionsLogic) ListContentRevisions(in *pb.ListContentRevisionsRequest) (*pb.ListContentRevisionsResponse, error) {
	actor, err := actorFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.Services.ListContentRevisions.Execute(l.ctx, actor, in)
	return resp, toStatusError(err, "list content revisions failed")
}
