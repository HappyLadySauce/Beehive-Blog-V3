package logic

import (
	"context"

	fileservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type CreateFileCategoryLogic struct {
	baseLogic
}

func NewCreateFileCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFileCategoryLogic {
	return &CreateFileCategoryLogic{baseLogic: newBaseLogic(ctx, svcCtx)}
}

func (l *CreateFileCategoryLogic) CreateFileCategory(in *pb.CreateFileCategoryRequest) (*pb.FileCategoryResponse, error) {
	category, err := l.svcCtx.Services.CreateFileCategory(l.ctx, fileservice.CreateFileCategoryInput{
		CategoryKey:       in.GetCategoryKey(),
		DisplayName:       in.GetDisplayName(),
		Description:       in.GetDescription(),
		Enabled:           in.GetEnabled(),
		IsDefault:         in.GetIsDefault(),
		SortOrder:         in.GetSortOrder(),
		AllowedExtensions: in.GetAllowedExtensions(),
	})
	if err != nil {
		return nil, toStatus(err, "create file category failed")
	}
	return &pb.FileCategoryResponse{Category: toProtoFileCategory(category)}, nil
}
