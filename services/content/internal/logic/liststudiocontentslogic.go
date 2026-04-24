package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type ListStudioContentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListStudioContentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListStudioContentsLogic {
	return &ListStudioContentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListStudioContentsLogic) ListStudioContents(in *pb.ListStudioContentsRequest) (*pb.ListStudioContentsResponse, error) {
	actor, err := actorFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.Services.ListStudioContents.Execute(l.ctx, actor, in)
	return resp, toStatusError(err, "list studio contents failed")
}
