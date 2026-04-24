package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type GetContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContentLogic {
	return &GetContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetContentLogic) GetContent(in *pb.GetContentRequest) (*pb.GetContentResponse, error) {
	actor, err := actorFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.Services.GetContent.Execute(l.ctx, actor, in)
	return resp, toStatusError(err, "get content failed")
}
