package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type CreateContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateContentLogic {
	return &CreateContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateContentLogic) CreateContent(in *pb.CreateContentRequest) (*pb.CreateContentResponse, error) {
	actor, err := actorFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.Services.CreateContent.Execute(l.ctx, actor, in)
	return resp, toStatusError(err, "create content failed")
}
