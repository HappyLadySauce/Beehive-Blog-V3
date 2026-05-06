package logic

import (
	"context"

	fileservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type ListFileCategoriesLogic struct {
	baseLogic
}

func NewListFileCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFileCategoriesLogic {
	return &ListFileCategoriesLogic{baseLogic: newBaseLogic(ctx, svcCtx)}
}

func (l *ListFileCategoriesLogic) ListFileCategories(in *pb.ListFileCategoriesRequest) (*pb.ListFileCategoriesResponse, error) {
	items, err := l.svcCtx.Services.ListFileCategories(l.ctx, fileservice.ListFileCategoriesInput{
		IncludeDisabled: in.GetIncludeDisabled(),
	})
	if err != nil {
		return nil, toStatus(err, "list file categories failed")
	}

	resp := &pb.ListFileCategoriesResponse{Items: make([]*pb.FileCategory, 0, len(items))}
	for _, item := range items {
		resp.Items = append(resp.Items, toProtoFileCategory(item))
	}
	return resp, nil
}
