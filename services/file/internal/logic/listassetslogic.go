package logic

import (
	"context"

	fileservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type ListAssetsLogic struct {
	baseLogic
}

func NewListAssetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAssetsLogic {
	return &ListAssetsLogic{baseLogic: newBaseLogic(ctx, svcCtx)}
}

func (l *ListAssetsLogic) ListAssets(in *pb.ListAssetsRequest) (*pb.ListAssetsResponse, error) {
	result, err := l.svcCtx.Services.ListAssets(l.ctx, fileservice.ListAssetsInput{
		ActorUserID: in.GetActorUserId(),
		CategoryKey: in.GetCategoryKey(),
		Status:      toOptionalServiceStatus(in.GetStatus()),
		Visibility:  toOptionalServiceVisibility(in.GetVisibility()),
		OwnerUserID: in.GetOwnerUserId(),
		Keyword:     in.GetKeyword(),
		Page:        int(in.GetPage()),
		PageSize:    int(in.GetPageSize()),
	})
	if err != nil {
		return nil, toStatus(err, "list file assets failed")
	}

	items := make([]*pb.Asset, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toProtoAsset(item))
	}
	return &pb.ListAssetsResponse{
		Items:    items,
		Total:    result.Total,
		Page:     int32(result.Page),
		PageSize: int32(result.PageSize),
	}, nil
}
