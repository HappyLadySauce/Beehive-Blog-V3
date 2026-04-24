package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type ListTagsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTagsLogic {
	return &ListTagsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTagsLogic) ListTags(in *pb.ListTagsRequest) (*pb.ListTagsResponse, error) {
	actor, err := actorFromContext(l.ctx)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.Services.ListTags.Execute(l.ctx, actor, in)
	return resp, toStatusError(err, "list tags failed")
}
