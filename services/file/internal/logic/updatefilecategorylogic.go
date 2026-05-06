package logic

import (
	"context"

	fileservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type UpdateFileCategoryLogic struct {
	baseLogic
}

func NewUpdateFileCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFileCategoryLogic {
	return &UpdateFileCategoryLogic{baseLogic: newBaseLogic(ctx, svcCtx)}
}

func (l *UpdateFileCategoryLogic) UpdateFileCategory(in *pb.UpdateFileCategoryRequest) (*pb.FileCategoryResponse, error) {
	category, err := l.svcCtx.Services.UpdateFileCategory(l.ctx, fileservice.UpdateFileCategoryInput{
		CategoryKey: in.GetCategoryKey(),
		DisplayName: in.GetDisplayName(),
		Description: in.GetDescription(),
		Enabled:     in.GetEnabled(),
		SortOrder:   in.GetSortOrder(),
	})
	if err != nil {
		return nil, toStatus(err, "update file category failed")
	}
	return &pb.FileCategoryResponse{Category: toProtoFileCategory(category)}, nil
}
