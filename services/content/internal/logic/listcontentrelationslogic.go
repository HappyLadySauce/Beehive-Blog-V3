package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type ListContentRelationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListContentRelationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListContentRelationsLogic {
	return &ListContentRelationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListContentRelationsLogic) ListContentRelations(in *pb.ListContentRelationsRequest) (*pb.ListContentRelationsResponse, error) {
	actor, err := actorFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.Services.ListContentRelations.Execute(l.ctx, actor, in)
	return resp, toStatusError(err, "list content relations failed")
}
