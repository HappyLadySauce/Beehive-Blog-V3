package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type CreateContentRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateContentRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateContentRelationLogic {
	return &CreateContentRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateContentRelationLogic) CreateContentRelation(in *pb.CreateContentRelationRequest) (*pb.CreateContentRelationResponse, error) {
	actor, err := actorFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.Services.CreateContentRelation.Execute(l.ctx, actor, in)
	return resp, toStatusError(err, "create content relation failed")
}
