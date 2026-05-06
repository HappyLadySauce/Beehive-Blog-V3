package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type SetDefaultFileCategoryLogic struct {
	baseLogic
}

func NewSetDefaultFileCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDefaultFileCategoryLogic {
	return &SetDefaultFileCategoryLogic{baseLogic: newBaseLogic(ctx, svcCtx)}
}

func (l *SetDefaultFileCategoryLogic) SetDefaultFileCategory(in *pb.SetDefaultFileCategoryRequest) (*pb.FileCategoryResponse, error) {
	category, err := l.svcCtx.Services.SetDefaultFileCategory(l.ctx, in.GetCategoryKey())
	if err != nil {
		return nil, toStatus(err, "set default file category failed")
	}
	return &pb.FileCategoryResponse{Category: toProtoFileCategory(category)}, nil
}
