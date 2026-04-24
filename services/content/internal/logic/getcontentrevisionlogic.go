package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type GetContentRevisionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetContentRevisionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContentRevisionLogic {
	return &GetContentRevisionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetContentRevisionLogic) GetContentRevision(in *pb.GetContentRevisionRequest) (*pb.GetContentRevisionResponse, error) {
	actor, err := actorFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.Services.GetContentRevision.Execute(l.ctx, actor, in)
	return resp, toStatusError(err, "get content revision failed")
}
