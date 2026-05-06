package logic

import (
	"context"

	fileservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type UpdateFileCategoryExtensionsLogic struct {
	baseLogic
}

func NewUpdateFileCategoryExtensionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFileCategoryExtensionsLogic {
	return &UpdateFileCategoryExtensionsLogic{baseLogic: newBaseLogic(ctx, svcCtx)}
}

func (l *UpdateFileCategoryExtensionsLogic) UpdateFileCategoryExtensions(in *pb.UpdateFileCategoryExtensionsRequest) (*pb.FileCategoryResponse, error) {
	category, err := l.svcCtx.Services.UpdateFileCategoryExtensions(l.ctx, fileservice.UpdateFileCategoryExtensionsInput{
		CategoryKey:       in.GetCategoryKey(),
		AllowedExtensions: in.GetAllowedExtensions(),
	})
	if err != nil {
		return nil, toStatus(err, "update file category extensions failed")
	}
	return &pb.FileCategoryResponse{Category: toProtoFileCategory(category)}, nil
}
