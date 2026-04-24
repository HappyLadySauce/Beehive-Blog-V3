package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type DeleteContentRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteContentRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteContentRelationLogic {
	return &DeleteContentRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteContentRelationLogic) DeleteContentRelation(in *pb.DeleteContentRelationRequest) (*pb.DeleteContentRelationResponse, error) {
	actor, err := actorFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.Services.DeleteContentRelation.Execute(l.ctx, actor, in)
	return resp, toStatusError(err, "delete content relation failed")
}
