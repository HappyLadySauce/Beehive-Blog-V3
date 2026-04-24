package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type ListPublicContentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPublicContentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPublicContentsLogic {
	return &ListPublicContentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPublicContentsLogic) ListPublicContents(in *pb.ListPublicContentsRequest) (*pb.ListPublicContentsResponse, error) {
	resp, err := l.svcCtx.Services.ListPublicContents.Execute(l.ctx, in)
	return resp, toStatusError(err, "list public contents failed")
}
