package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

type GetPublicContentBySlugLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPublicContentBySlugLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPublicContentBySlugLogic {
	return &GetPublicContentBySlugLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPublicContentBySlugLogic) GetPublicContentBySlug(in *pb.GetPublicContentBySlugRequest) (*pb.GetPublicContentBySlugResponse, error) {
	resp, err := l.svcCtx.Services.GetPublicContentBySlug.Execute(l.ctx, in)
	return resp, toStatusError(err, "get public content failed")
}
