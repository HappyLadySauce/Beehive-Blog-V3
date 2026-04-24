package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type ArchiveContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArchiveContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArchiveContentLogic {
	return &ArchiveContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArchiveContentLogic) ArchiveContent(in *pb.ArchiveContentRequest) (*pb.ArchiveContentResponse, error) {
	actor, err := actorFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.Services.ArchiveContent.Execute(l.ctx, actor, in)
	return resp, toStatusError(err, "archive content failed")
}
