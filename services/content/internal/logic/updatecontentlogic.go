package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type UpdateContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateContentLogic {
	return &UpdateContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateContentLogic) UpdateContent(in *pb.UpdateContentRequest) (*pb.UpdateContentResponse, error) {
	actor, err := actorFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.Services.UpdateContent.Execute(l.ctx, actor, in)
	return resp, toStatusError(err, "update content failed")
}
